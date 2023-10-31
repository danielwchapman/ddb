// Package ddb is a collection of functions that make AWS DynamoDB easier to work with.
package ddb

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Client provides convenience methods for working with a DynamoDB table following Single Table Design.
type Client struct {
	Ddb   *dynamodb.Client
	Table string
}

var _ ClientInterface = (*Client)(nil)

func (c *Client) Delete(ctx context.Context, pk, sk string) error {
	_, err := c.Ddb.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &c.Table,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	})

	if err != nil {
		return fmt.Errorf("Delete: DeleteItem: %w", err)
	}

	return nil
}

func (c *Client) Get(ctx context.Context, pk, sk string, out any) error {
	req := dynamodb.GetItemInput{
		TableName: &c.Table,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	}

	resp, err := c.Ddb.GetItem(ctx, &req)
	if err != nil {
		return fmt.Errorf("Get: GetItem: %w", err)
	}

	if len(resp.Item) == 0 {
		return ErrNotFound
	}

	if err := attributevalue.UnmarshalMap(resp.Item, out); err != nil {
		return fmt.Errorf("Get: UnmarshalMap: %w", err)
	}

	return nil
}

func (c *Client) Query(ctx context.Context, pk, skBeginsWith string, out any, opts ...Option) error {
	if pk == "" {
		return &InvalidArgumentError{errors.New("pk cannot be empty")}
	}

	var queryOptions options
	for _, opt := range opts {
		if err := opt(&queryOptions); err != nil {
			return fmt.Errorf("Query: %w", err)
		}
	}

	keyCondition := expression.Key(queryOptions.indexName + "PK").Equal(expression.Value(pk))
	if skBeginsWith != "" {
		keyCondition = keyCondition.And(expression.Key(queryOptions.indexName + "SK").BeginsWith(skBeginsWith))
	}

	expr, err := expression.
		NewBuilder().
		WithKeyCondition(keyCondition).
		Build()

	scanForward := !queryOptions.scanBackwards

	req := dynamodb.QueryInput{
		ExclusiveStartKey:         queryOptions.startKey,
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeValues: expr.Values(),
		ExpressionAttributeNames:  expr.Names(),
		Limit:                     queryOptions.pageSize,
		ScanIndexForward:          &scanForward,
		TableName:                 &c.Table,
	}

	result, err := c.Ddb.Query(ctx, &req)
	if err != nil {
		return fmt.Errorf("Query: %w", err)
	}

	if len(result.Items) == 0 {
		return ErrNotFound
	}

	if err = attributevalue.UnmarshalListOfMaps(result.Items, &out); err != nil {
		return &InternalError{err: fmt.Errorf("Query: UnmarshalListOfMaps: %w", err)}
	}

	if queryOptions.pageOut != nil && len(result.LastEvaluatedKey) > 0 {
		lastEvaluatedKey, err := SerializeExclusiveStartKey(result.LastEvaluatedKey)
		if err != nil {
			return &InternalError{err: fmt.Errorf("Query: SerializeExclusiveStartKey: %w", err)}
		}
		*queryOptions.pageOut = lastEvaluatedKey
	}

	return nil
}

func (c *Client) Put(ctx context.Context, row any, opts ...Option) error {
	var putOptions options
	for _, opt := range opts {
		err := opt(&putOptions)
		if err != nil {
			return err
		}
	}

	if putOptions.updatesCount > 0 {
		return &InvalidArgumentError{errors.New("put cannot update items with options")}
	}

	var (
		expressionAttributeValues map[string]types.AttributeValue
		expressionAttributeNames  map[string]string
		condition                 *string
	)

	if putOptions.conditionsCount > 0 {
		expr, err := expression.NewBuilder().
			WithCondition(putOptions.conditions).
			Build()

		if err != nil {
			return fmt.Errorf("Put: expression builder: %w", err)
		}

		expressionAttributeValues = expr.Values()
		expressionAttributeNames = expr.Names()
		condition = expr.Condition()
	}

	item, err := attributevalue.MarshalMap(row)
	if err != nil {
		return fmt.Errorf("Put: MarshalMap: %w", err)
	}

	req := dynamodb.PutItemInput{
		TableName:                 &c.Table,
		Item:                      item,
		ConditionExpression:       condition,
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
		ReturnValues:              putOptions.returnValues,
	}

	out, err := c.Ddb.PutItem(ctx, &req)
	if err != nil {
		// TODO check for conditional errors
		return fmt.Errorf("Put: PutItem: %w", err)
	}

	if putOptions.returnValues != "" {
		if err := attributevalue.UnmarshalMap(out.Attributes, putOptions.returnValuesOut); err != nil {
			return fmt.Errorf("Put: UnmarshalMap: %w", err)
		}
	}

	return nil
}

// TransactDeletes uses a DynamoDB transaction to delete multiple items in one atomic request.
func (c *Client) TransactDeletes(ctx context.Context, token string, rows ...DeleteRow) error {
	if len(rows) > 100 {
		return &InvalidArgumentError{errors.New("cannot exceed 100 rows")}
	}

	items := makeDeletes(c.Table, rows...)

	req := dynamodb.TransactWriteItemsInput{
		TransactItems:      makeTransactionWriteItems(nil, items, nil),
		ClientRequestToken: &token,
	}

	if _, err := c.Ddb.TransactWriteItems(ctx, &req); err != nil {
		var condFailedErr *types.ConditionalCheckFailedException
		if errors.As(err, &condFailedErr) {
			return fmt.Errorf("TransactDeletes: TransactWriteItems: Condition failed %w", condFailedErr)
		}

		// TODO tidy up
		if canceledErr, ok := IsTransactionCanceled(err); ok {
			return fmt.Errorf("TransactDeletes: TransactWriteItems: Canceled Transaction: %w", canceledErr)
		}

		return fmt.Errorf("TransactDeletes: TransactWriteItems: %w", err)
	}

	return nil
}

// TransactPuts uses a DynamoDB transaction to put multiple items in one atomic request.
func (c *Client) TransactPuts(ctx context.Context, token string, rows ...PutRow) error {
	if len(rows) > 100 {
		return &InvalidArgumentError{errors.New("cannot exceed 100 rows")}
	}

	items, err := makePuts(c.Table, rows...)
	if err != nil {
		return fmt.Errorf("TransactionPuts: %w", err)
	}

	req := dynamodb.TransactWriteItemsInput{
		TransactItems:      makeTransactionWriteItems(items, nil, nil),
		ClientRequestToken: &token,
	}

	if _, err := c.Ddb.TransactWriteItems(ctx, &req); err != nil {
		var condFailedErr *types.ConditionalCheckFailedException
		if errors.As(err, &condFailedErr) {
			return fmt.Errorf("TransactPuts: TransactWriteItems: Condition failed %w", condFailedErr)
		}

		// TODO tidy up
		if canceledErr, ok := IsTransactionCanceled(err); ok {
			return fmt.Errorf("TransactPuts: TransactWriteItems: Canceled Transaction: %w", canceledErr)
		}

		return fmt.Errorf("TransactPuts: TransactWriteItems: %w", err)
	}

	return nil
}

//// TransactWrites uses a DynamoDB transaction to put multiple items in one atomic request.
//func (c *Client) TransactWrites(ctx context.Context, token string, puts []PutRow, deletes []DeleteRow, updates []UpdateRow) error {
//    if len(puts)+len(deletes)+len(updates) > 100 {
//        return grpcerrors.MakeInvalidArgumentError("cannot exceed 100 rows")
//    }
//
//    putItems, err := makePuts(c.Table, puts...)
//    if err != nil {
//        return fmt.Errorf("TransactWrites: %w", err)
//    }
//
//    deleteItems, err := makeDeletes(c.Table, deletes...)
//    if err != nil {
//        return fmt.Errorf("TransactWrites: %w", err)
//    }
//
//    updateItems, err := makeUpdates(c.Table, updates...)
//    if err != nil {
//        return fmt.Errorf("TransactWrites: %w", err)
//    }
//
//    req := dynamodb.TransactWriteItemsInput{
//        TransactItems:      makeTransactionWriteItems2(putItems, deleteItems, updateItems),
//        ClientRequestToken: &token,
//    }
//
//    if _, err := c.Ddb.TransactWriteItems(ctx, &req); err != nil {
//        var condFailedErr *types.ConditionalCheckFailedException
//        if errors.As(err, &condFailedErr) {
//            return fmt.Errorf("TransactPuts: TransactWriteItems: Condition failed %w", condFailedErr)
//        }
//
//        // TODO tidy up
//        if canceledErr, ok := IsTransactionCanceled(err); ok {
//            return fmt.Errorf("TransactPuts: TransactWriteItems: Canceled Transaction: %w", canceledErr)
//        }
//
//        return fmt.Errorf("TransactPuts: TransactWriteItems: %w", err)
//    }
//
//    return nil
//}

// Update updates an item in a table. The row map must contain the updated values for the item. If a key is not
// in the row map, the value will be unchanged. Careful when working with arrays and maps, as the entire value
// will be replaced.
func (c *Client) Update(ctx context.Context, pk, sk string, opts ...Option) error {
	var updateOptions options
	for _, opt := range opts {
		err := opt(&updateOptions)
		if err != nil {
			return err
		}
	}

	var (
		conditionExpression *string
		expr                expression.Expression
		err                 error
	)

	if updateOptions.conditionsCount > 0 {
		expr, err = expression.NewBuilder().
			WithCondition(updateOptions.conditions).
			WithUpdate(updateOptions.updates).
			Build()
	} else {
		expr, err = expression.NewBuilder().
			WithUpdate(updateOptions.updates).
			Build()
	}

	if err != nil {
		return fmt.Errorf("Update: expression builder: %w", err)
	}

	if updateOptions.conditionsCount > 0 {
		conditionExpression = expr.Condition()
	}

	req := dynamodb.UpdateItemInput{
		TableName: &c.Table,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		ConditionExpression:       conditionExpression,
		ExpressionAttributeValues: expr.Values(),
		ExpressionAttributeNames:  expr.Names(),
		UpdateExpression:          expr.Update(),
		ReturnValues:              updateOptions.returnValues,
	}

	if req.UpdateExpression == nil || *req.UpdateExpression == "" {
		return &InvalidArgumentError{err: errors.New("no updates provided")}
	}

	out, err := c.Ddb.UpdateItem(ctx, &req)

	if err != nil {
		// TODO add conditional check failed error and map to
		// grpcerrors.ErrNotFound or grpcerrors.AlreadyExists

		return fmt.Errorf("Update: %w", err)
	}

	if updateOptions.returnValues != "" {
		if err := attributevalue.UnmarshalMap(out.Attributes, updateOptions.returnValuesOut); err != nil {
			return fmt.Errorf("Update: UnmarshalMap: %w", err)
		}
	}

	return nil
}
