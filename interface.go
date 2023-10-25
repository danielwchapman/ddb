package ddb

import "context"

//go:generate mockgen -source=interface.go -destination=./mocks/mocks.go -package=mocks

type ClientInterface interface {
    Deleter
    Getter
    Putter
    TransactPutter
    Updater
}

type Deleter interface {
    Delete(ctx context.Context, pk, sk string) error
}

type Getter interface {
    Get(ctx context.Context, pk, sk string, out any) error
}

type Putter interface {
    Put(ctx context.Context, row any, opts ...Option) error
}

type Queryer interface {
    Query(ctx context.Context, pk, skPrefix string, opts ...Option) error
}

type TransactPutter interface {
    TransactPuts(ctx context.Context, token string, rows ...PutRow) error
}

type Updater interface {
    Update(ctx context.Context, pk, sk string, opts ...Option) error
}
