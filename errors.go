package ddb

import (
    "errors"
    "fmt"
    "strings"

    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// IsTransactionCanceled checks if the error is a TransactionCanceledException and
// returns improved error message and true if it is.
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
