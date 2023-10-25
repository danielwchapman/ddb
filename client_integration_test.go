package ddb

import (
    "context"
    "errors"
    "fmt"
    "os"
    "testing"
    "time"

    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
    "github.com/danielwchapman/grpcerrors"
    "github.com/google/uuid"

    "github.com/google/go-cmp/cmp"
)

// TODO use a test table instead

var now = time.Now().Truncate(0)

var uut = func() *Client {
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

    return &Client{
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
        PK:         "PK#" + suffix,
        SK:         "SK#" + suffix,
        TestString: "test string",
        TestInt:    123,
        TestFloat:  123.456,
        TestBool:   true,
        TestTime:   now,
        TestSlice:  []string{"a", "b", "c"},
        TestMap:    map[string]string{"a": "a1", "b": "b1"},
    }
}

func makeQueryTestRows(pkSuffix string, count int) []testRow {
    rows := make([]testRow, count)

    for i := 0; i < count; i++ {
        rows[i] = testRow{
            PK:         "PK#" + pkSuffix,
            SK:         fmt.Sprintf("SK#%d", i),
            TestString: "test string",
            TestInt:    123,
            TestFloat:  123.456,
            TestBool:   true,
            TestTime:   now,
            TestSlice:  []string{"a", "b", "c"},
            TestMap:    map[string]string{"a": "a1", "b": "b1"},
        }
    }

    return rows
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

        if err := uut.Put(ctx, want); err != nil {
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

    t.Run("Put ItemExists condition honored", func(t *testing.T) {
        row := makeRandomTestRow(t.Name())

        t.Cleanup(func() {
            if err := uut.Delete(ctx, row.PK, row.SK); err != nil {
                t.Errorf("unexpected error: %v", err)
            }
        })

        // should fail because item does not exist
        if err := uut.Put(ctx, row, WithItemExists()); err == nil {
            t.Errorf("unexpected error: %v", err)
        }
    })

    t.Run("Put return values all old", func(t *testing.T) {
        want := makeRandomTestRow(t.Name())

        t.Cleanup(func() {
            if err := uut.Delete(ctx, want.PK, want.SK); err != nil {
                t.Errorf("unexpected error: %v", err)
            }
        })

        if err := uut.Put(ctx, want); err != nil {
            t.Errorf("unexpected error: %v", err)
        }

        var got testRow
        if err := uut.Put(ctx, want, WithReturnValues(types.ReturnValueAllOld, &got)); err != nil {
            t.Errorf("unexpected error: %v", err)
        }

        if diff := cmp.Diff(want, got); diff != "" {
            t.Errorf("unexpected diff: %s", diff)
        }
    })
}

func TestIntegrationQuery(t *testing.T) {
    t.Parallel()

    if os.Getenv("INTEGRATION") == "" {
        t.Skip("skipping integration tests")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    testRows := makeQueryTestRows(t.Name(), 8)
    testTime, err := time.Parse(time.RFC3339, "2023-10-25T09:17:47.855071-04:00")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    for i := range testRows {
        testRows[i].TestTime = testTime
    }

    //t.Cleanup(func() {
    //    // TODO use batch delete when ready
    //    for i := range testRows {
    //        if err := uut.Delete(ctx, testRows[i].PK, testRows[i].SK); err != nil {
    //            t.Errorf("unexpected error: %v", err)
    //        }
    //    }
    //})

    // TODO use BatchWriteItems when it's ready
    //for i := range testRows {
    //    if err := uut.Put(ctx, testRows[i]); err != nil {
    //        t.Errorf("unexpected error: %v", err)
    //    }
    //}

    t.Run("Basic", func(t *testing.T) {
        var got []testRow
        if err := uut.Query(ctx, testRows[0].PK, "SK", &got); err != nil {
            t.Errorf("unexpected error: %v", err)
        }

        if diff := cmp.Diff(testRows, got); diff != "" {
            t.Errorf("unexpected diff: %s", diff)
        }
    })

    // Empty BeginsWith should return all PK rows
    t.Run("Empty BeginsWith", func(t *testing.T) {
        var got []testRow
        if err := uut.Query(ctx, testRows[0].PK, "", &got); err != nil {
            t.Errorf("unexpected error: %v", err)
        }

        if diff := cmp.Diff(testRows, got); diff != "" {
            t.Errorf("unexpected diff: %s", diff)
        }
    })

    t.Run("NotFound Error", func(t *testing.T) {
        var got []testRow
        err := uut.Query(ctx, "NotInTable", "NotThere", &got)
        if !errors.Is(err, ErrNotFound) {
            t.Errorf("unexpected error: %v", err)
        }
    })

    t.Run("Honor ExclusiveStartKey", func(t *testing.T) {
        const pageSize = 2

        var (
            got1      []testRow
            got2      []testRow
            pageToken string
        )

        if err := uut.Query(
            ctx,
            testRows[0].PK,
            "SK", &got1,
            WithPageSize(pageSize),
            WithPage(pageToken, &pageToken), // empty pageToken the first time
        ); err != nil {
            t.Errorf("unexpected error: %v", err)
        }

        if pageToken == "" {
            t.Errorf("expected pageToken to be set")
        }

        if diff := cmp.Diff(testRows[:pageSize], got1); diff != "" {
            t.Errorf("unexpected diff: %s", diff)
        }

        if err := uut.Query(
            ctx,
            testRows[0].PK,
            "SK", &got2,
            WithPageSize(pageSize),
            WithPage(pageToken, &pageToken),
        ); err != nil {
            t.Errorf("unexpected error: %v", err)
        }

        if pageToken == "" {
            t.Errorf("expected pageToken to be set")
        }

        if diff := cmp.Diff(testRows[pageSize:pageSize+pageSize], got2); diff != "" {
            t.Errorf("unexpected diff: %s", diff)
        }
    })

    t.Run("Honor PageSize", func(t *testing.T) {
        const pageSize = 3
        var got []testRow
        if err := uut.Query(ctx, testRows[0].PK, "", &got, WithPageSize(pageSize)); err != nil {
            t.Errorf("unexpected error: %v", err)
        }

        if diff := cmp.Diff(testRows[:pageSize], got); diff != "" {
            t.Errorf("unexpected diff: %s", diff)
        }
    })

    t.Run("Honor ScanBackwards", func(t *testing.T) {
        var got []testRow
        if err := uut.Query(ctx, testRows[0].PK, "SK", &got, WithScanBackwards()); err != nil {
            t.Errorf("unexpected error: %v", err)
        }

        want := make([]testRow, len(testRows))
        for i := range testRows {
            want[i] = testRows[len(testRows)-1-i]
        }

        if diff := cmp.Diff(want, got); diff != "" {
            t.Errorf("unexpected diff: %s", diff)
        }
    })

    t.Run("GSI", func(t *testing.T) {
        // TODO
        t.Skip("needs implemented")
    })

    t.Run("LSI", func(t *testing.T) {
        // TODO
        t.Skip("needs implemented")
    })
}

func TestIntegrationTransactPuts(t *testing.T) {
    t.Parallel()

    if os.Getenv("INTEGRATION") == "" {
        t.Skip("skipping integration tests")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    t.Run("Basic", func(t *testing.T) {
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

func TestIntegrationUpdate(t *testing.T) {
    t.Parallel()

    if os.Getenv("INTEGRATION") == "" {
        t.Skip("skipping integration tests")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    t.Run("Basic", func(t *testing.T) {
        row := makeRandomTestRow(t.Name())

        t.Cleanup(func() {
            if err := uut.Delete(ctx, row.PK, row.SK); err != nil {
                t.Errorf("unexpected error: %v", err)
            }
        })

        // put a row so we have something to update
        if err := uut.Put(ctx, row, WithItemNotExist()); err != nil {
            t.Errorf("unexpected error: %v", err)
        }

        updates := map[string]any{
            "TestString": "updated string",
        }

        var got testRow
        if err := uut.Update(
            ctx,
            row.PK,
            row.SK,
            WithItemExists(),
            WithFieldUpdates(updates),
            WithReturnValues(types.ReturnValueAllNew, &got),
        ); err != nil {
            t.Errorf("unexpected error: %v", err)
        }

        row.TestString = "updated string"
        if diff := cmp.Diff(row, got); diff != "" {
            t.Errorf("unexpected diff: %s", diff)
        }
    })

    t.Run("Update Embedded Map", func(t *testing.T) {
        row := makeRandomTestRow(t.Name())

        t.Cleanup(func() {
            if err := uut.Delete(ctx, row.PK, row.SK); err != nil {
                t.Errorf("unexpected error: %v", err)
            }
        })

        // put a row so we have something to update
        if err := uut.Put(ctx, row, WithItemNotExist()); err != nil {
            t.Errorf("unexpected error: %v", err)
        }

        updates := map[string]any{
            "TestMap.a": "a2",
        }

        var got testRow
        if err := uut.Update(ctx, row.PK, row.SK, WithFieldUpdates(updates), WithReturnValues(types.ReturnValueAllNew, &got)); err != nil {
            t.Errorf("unexpected error: %v", err)
        }

        row.TestMap["a"] = "a2"
        if diff := cmp.Diff(row, got); diff != "" {
            t.Errorf("unexpected diff: %s", diff)
        }
    })

    t.Run("Update Slice", func(t *testing.T) {
        // TODO
        t.Skip("needs implemented")
    })

    t.Run("Honor WithExists", func(t *testing.T) {
        row := makeRandomTestRow(t.Name())

        updates := map[string]any{
            "TestString": "updated string",
        }

        // item does not exist, so this should fail
        if err := uut.Update(ctx, row.PK, row.SK, WithItemExists(), WithFieldUpdates(updates)); err == nil {
            t.Errorf("unexpected error: %v", err)
        }
    })

    t.Run("Honor WithNotExists", func(t *testing.T) {
        row := makeRandomTestRow(t.Name())

        updates := map[string]any{
            "TestString": "updated string",
        }

        t.Cleanup(func() {
            if err := uut.Delete(ctx, row.PK, row.SK); err != nil {
                t.Errorf("unexpected error: %v", err)
            }
        })

        // put a row so we have something to update
        if err := uut.Put(ctx, row, WithItemNotExist()); err != nil {
            t.Errorf("unexpected error: %v", err)
        }

        // item does not exist, so this should fail
        if err := uut.Update(ctx, row.PK, row.SK, WithItemNotExist(), WithFieldUpdates(updates)); err == nil {
            t.Errorf("unexpected error: %v", err)
        }
    })

    t.Run("Empty Conditions", func(t *testing.T) {
        row := makeRandomTestRow(t.Name())

        t.Cleanup(func() {
            if err := uut.Delete(ctx, row.PK, row.SK); err != nil {
                t.Errorf("unexpected error: %v", err)
            }
        })

        updates := map[string]any{
            "TestString": "updated string",
        }

        if err := uut.Update(
            ctx,
            row.PK,
            row.SK,
            WithFieldUpdates(updates),
        ); err != nil {
            t.Errorf("unexpected error: %v", err)
        }
    })

    t.Run("Multiple Conditions", func(t *testing.T) {
        // TODO
        t.Skip("needs implemented")
    })

    t.Run("Honor WithReturnValues", func(t *testing.T) {
        row := makeRandomTestRow(t.Name())

        t.Cleanup(func() {
            if err := uut.Delete(ctx, row.PK, row.SK); err != nil {
                t.Errorf("unexpected error: %v", err)
            }
        })

        // put a row so we have something to update
        if err := uut.Put(ctx, row, WithItemNotExist()); err != nil {
            t.Errorf("unexpected error: %v", err)
        }

        updates := map[string]any{
            "TestString": "updated string",
        }

        var got testRow
        if err := uut.Update(
            ctx,
            row.PK,
            row.SK,
            WithItemExists(),
            WithFieldUpdates(updates),
            WithReturnValues(types.ReturnValueAllOld, &got),
        ); err != nil {
            t.Errorf("unexpected error: %v", err)
        }

        if diff := cmp.Diff(row, got); diff != "" {
            t.Errorf("unexpected diff: %s", diff)
        }
    })
}
