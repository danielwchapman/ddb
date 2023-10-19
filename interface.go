package ddb

import "context"

// TODO break this into individual methods and create a new package out of it

//go:generate mockgen -source=interface.go -destination=./mocks/mocks.go -package=mock

type StdHelper interface {
    Delete(ctx context.Context, pk, sk string) error
    Get(ctx context.Context, pk, sk string, out any) error
    Put(ctx context.Context, condition *string, row any) error
    TransactPuts(ctx context.Context, token string, rows ...PutRow) error
}

var _ StdHelper = (*StdClient)(nil)

type PutRow struct {
    Row       any
    Condition *string
}

type DeleteRow struct {
    PK        string
    SK        string
    Condition *string
}
