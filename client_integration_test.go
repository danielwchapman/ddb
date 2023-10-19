package ddb

import (
    "context"
    "errors"
    "os"
    "testing"
    "time"

    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/danielwchapman/grpcerrors"
    "github.com/google/uuid"

    "github.com/google/go-cmp/cmp"
)

var uut = func() *StdClient {
    if os.Getenv("INTEGRATION") == "" {
        return nil
    }

    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        panic(err)
    }

    table := os.Getenv("TEST_TABLE")

    if table == "" {
        panic("QUESTIONS_TABLE env var not set")
    }

    return &StdClient{
        Ddb:   dynamodb.NewFromConfig(cfg),
        Table: table,
    }
}()

type testRow struct {
    PK         string
    SK         string
    TestString string
    TestInt    int
    TestFloat  float64
    TestBool   bool
    TestTime   time.Time
    TestSlice  []string
    TestMap    map[string]string
    // TODO add pointer types
    // TODO add embedded struct
}

func makeRandomTestRow(suffix string) testRow {
    return testRow{
        PK:         "PK" + suffix,
        SK:         "SK" + suffix,
        TestString: "test string",
        TestInt:    123,
        TestFloat:  123.456,
        TestBool:   true,
        TestTime:   time.Now(),
        TestSlice:  []string{"a", "b", "c"},
        TestMap:    map[string]string{"a": "b", "c": "d"},
    }
}

func TestIntegrationGet(t *testing.T) {
    t.Parallel()

    if os.Getenv("INTEGRATION") == "" {
        t.Skip("skipping integration tests")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    t.Run("Get non-existent row returns Not Found error", func(t *testing.T) {
        var got testRow
        if err := uut.Get(ctx, "non-existent", "non-existent", &got); !errors.Is(err, grpcerrors.ErrNotFound) {
            t.Errorf("expected NotFound error, got: %v", err)
        }
    })
}

func TestIntegrationPut(t *testing.T) {
    t.Parallel()

    if os.Getenv("INTEGRATION") == "" {
        t.Skip("skipping integration tests")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    t.Run("Put then get same row", func(t *testing.T) {
        want := makeRandomTestRow(t.Name())

        t.Cleanup(func() {
            if err := uut.Delete(ctx, want.PK, want.SK); err != nil {
                t.Errorf("unexpected error: %v", err)
            }
        })

        if err := uut.Put(ctx, nil, want); err != nil {
            t.Errorf("unexpected error: %v", err)
        }

        var got testRow
        if err := uut.Get(ctx, want.PK, want.SK, &got); err != nil {
            t.Errorf("unexpected error: %v", err)
        }

        if diff := cmp.Diff(want, got); diff != "" {
            t.Errorf("unexpected diff: %s", diff)
        }
    })

    t.Run("Put condition already exists returns correct error type", func(t *testing.T) {
        row := makeRandomTestRow(t.Name())

        t.Cleanup(func() {
            if err := uut.Delete(ctx, row.PK, row.SK); err != nil {
                t.Errorf("unexpected error: %v", err)
            }
        })

        // First Put should succeed
        if err := uut.Put(ctx, &DoesNotExist, row); err != nil {
            t.Errorf("unexpected error: %v", err)
        }

        // Second Put should fail
        if err := uut.Put(ctx, &DoesNotExist, row); err == nil {
            t.Errorf("expected error because row alrady exists, but got nil")
        }
    })
}

func TestIntegrationTransactPuts(t *testing.T) {
    t.Parallel()

    if os.Getenv("INTEGRATION") == "" {
        t.Skip("skipping integration tests")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    t.Run("Put many then get and compare them all", func(t *testing.T) {
        var (
            want1 = makeRandomTestRow(t.Name() + "1")
            want2 = makeRandomTestRow(t.Name() + "2")
            want3 = makeRandomTestRow(t.Name() + "3")
        )

        t.Cleanup(func() {
            if err := uut.Delete(ctx, want1.PK, want1.SK); err != nil {
                t.Errorf("unexpected error: %v", err)
            }
            if err := uut.Delete(ctx, want2.PK, want2.SK); err != nil {
                t.Errorf("unexpected error: %v", err)
            }
            if err := uut.Delete(ctx, want3.PK, want3.SK); err != nil {
                t.Errorf("unexpected error: %v", err)
            }
        })

        rows := []PutRow{
            {
                Condition: nil,
                Row:       want1,
            },
            {
                Condition: nil,
                Row:       want2,
            },
            {
                Condition: nil,
                Row:       want3,
            },
        }

        token := uuid.New().String()
        if err := uut.TransactPuts(ctx, token, rows...); err != nil {
            t.Fatalf("unexpected error: %v", err)
        }

        var got1, got2, got3 testRow

        if err := uut.Get(ctx, want1.PK, want1.SK, &got1); err != nil {
            t.Errorf("unexpected error: %v", err)
        }

        if err := uut.Get(ctx, want2.PK, want2.SK, &got2); err != nil {
            t.Errorf("unexpected error: %v", err)
        }

        if err := uut.Get(ctx, want3.PK, want3.SK, &got3); err != nil {
            t.Errorf("unexpected error: %v", err)
        }

        if diff := cmp.Diff(want1, got1); diff != "" {
            t.Errorf("unexpected diff: %s", diff)
        }

        if diff := cmp.Diff(want2, got2); diff != "" {
            t.Errorf("unexpected diff: %s", diff)
        }

        if diff := cmp.Diff(want3, got3); diff != "" {
            t.Errorf("unexpected diff: %s", diff)
        }
    })
}