package ddb

import "context"

//go:generate mockgen -source=interface.go -destination=./mocks/mocks.go -package=mocks

type ClientInterface interface {
    Deleter
    Getter
    Putter
    TransactPutter
}

type Deleter interface {
    Delete(ctx context.Context, pk, sk string) error
}

type Getter interface {
    Get(ctx context.Context, pk, sk string, out any) error
}

type Putter interface {
    Put(ctx context.Context, condition *string, row any) error
}

type TransactPutter interface {
    TransactPuts(ctx context.Context, token string, rows ...PutRow) error
}

type PutRow struct {
    Row       any
    Condition *string
}

type DeleteRow struct {
    PK        string
    SK        string
    Condition *string
}
