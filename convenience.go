package ddb

import (
    "errors"
    "fmt"
    "strings"

    "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func marshalMapList(items []any) ([]map[string]types.AttributeValue, error) {
    if len(items) == 0 {
        return nil, nil
    }

    var (
        out = make([]map[string]types.AttributeValue, len(items))
        err error
    )

    for i := range items {
        out[i], err = attributevalue.MarshalMap(items[i])
        if err != nil {
            return nil, fmt.Errorf("MarshalMapList: MarshalMap: %w", err)
        }
    }

    return out, nil
}

func makePuts(table string, rows ...PutRow) ([]types.Put, error) {
    items := make([]types.Put, len(rows))
    for i := range rows {
        item, err := attributevalue.MarshalMap(rows[i].Row)
        if err != nil {
            return nil, fmt.Errorf("TransactionPuts: MarshalMap: %w", err)
        }

        items[i] = types.Put{
            Item:                item,
            ConditionExpression: rows[i].Condition,
            TableName:           &table,
        }
    }

    return items, nil
}

func makeDeletes(table string, rows ...DeleteRow) ([]types.Delete, error) {
    items := make([]types.Delete, len(rows))
    for i := range rows {
        items[i] = types.Delete{
            Key: map[string]types.AttributeValue{
                "PK": &types.AttributeValueMemberS{Value: rows[i].PK},
                "SK": &types.AttributeValueMemberS{Value: rows[i].SK},
            },
            ConditionExpression: rows[i].Condition,
            TableName:           &table,
        }
    }
    return items, nil
}

func makeTransactionWriteItems(puts []types.Put) []types.TransactWriteItem {
    out := make([]types.TransactWriteItem, len(puts))
    for i := range puts {
        out[i] = types.TransactWriteItem{
            Put: &puts[i],
        }
    }
    return out
}

func makeTransactionWriteItems2(
    puts []types.Put,
    deletes []types.Delete,
    updates []types.Update,
) []types.TransactWriteItem {
    out := make([]types.TransactWriteItem, len(puts)+len(deletes)+len(updates))
    i := 0
    for j := range puts {
        out[i] = types.TransactWriteItem{
            Put: &puts[j],
        }
    }
    for j := range deletes {
        out[i] = types.TransactWriteItem{
            Delete: &deletes[j],
        }
    }
    for j := range updates {
        out[i] = types.TransactWriteItem{
            Update: &updates[j],
        }
    }
    return out
}

// IsTransactionCanceled checks if the error is a TransactionCanceledException. If it is,
// it returns improved error message and true.
func IsTransactionCanceled(err error) (error, bool) {
    var e *types.TransactionCanceledException
    if !errors.As(err, &e) {
        return nil, false
    }
    var builder strings.Builder
    builder.WriteString("Transaction canceled: ")
    if e.Message != nil {
        builder.WriteString(*e.Message)
    }
    for _, reason := range e.CancellationReasons {
        if reason.Code != nil {
            _, _ = fmt.Fprintf(&builder, "; code: %s", *reason.Code)
        }
        if reason.Message != nil {
            _, _ = fmt.Fprintf(&builder, "; message: %s", *reason.Message)
        }
        if len(reason.Item) != 0 {
            // don't print the whole item for data privacy reasons.
            _, _ = fmt.Fprintf(&builder, "; item: PK: %s; SK: %s", reason.Item["PK"], reason.Item["SK"])
        }
    }
    return errors.New(builder.String()), true
}
