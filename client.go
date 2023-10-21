// Package ddb is a collection of functions that make AWS DynamoDB easier to work with.
package ddb

import (
    "context"
    "errors"
    "fmt"

    "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

    "github.com/danielwchapman/grpcerrors"
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
        return grpcerrors.ErrNotFound
    }

    if err := attributevalue.UnmarshalMap(resp.Item, out); err != nil {
        return fmt.Errorf("Get: UnmarshalMap: %w", err)
    }

    return nil
}

func (c *Client) Put(ctx context.Context, condition *string, row any) error {
    item, err := attributevalue.MarshalMap(row)
    if err != nil {
        return fmt.Errorf("Put: MarshalMap: %w", err)
    }

    req := dynamodb.PutItemInput{
        TableName:           &c.Table,
        Item:                item,
        ConditionExpression: condition,
    }

    if _, err := c.Ddb.PutItem(ctx, &req); err != nil {
        // TODO check for conditional errors
        return fmt.Errorf("Put: PutItem: %w", err)
    }

    return nil
}

// TransactDeletes uses a DynamoDB transaction to delete multiple items in one atomic request.
func (c *Client) TransactDeletes(ctx context.Context, token string, rows ...DeleteRow) error {
    if len(rows) > 100 {
        return grpcerrors.MakeInvalidArgumentError("cannot exceed 100 rows")
    }

    items, err := makeDeletes(c.Table, rows...)
    if err != nil {
        return fmt.Errorf("TransactionPuts: %w", err)
    }

    req := dynamodb.TransactWriteItemsInput{
        TransactItems:      makeTransactionWriteItems(nil, items, nil),
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

// TransactPuts uses a DynamoDB transaction to put multiple items in one atomic request.
func (c *Client) TransactPuts(ctx context.Context, token string, rows ...PutRow) error {
    if len(rows) > 100 {
        return grpcerrors.MakeInvalidArgumentError("cannot exceed 100 rows")
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
