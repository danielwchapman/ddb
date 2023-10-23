package ddb

type DeleteRow struct {
    PK        string
    SK        string
    Condition *string
}

type PutRow struct {
    Row       any
    Condition *string
}

// RowHeader are fields that must exist in every database row. It enforces a composite primary key where
// columns are named 'PK' and 'SK'. It also enforces a RowType column for identification.
type RowHeader struct {
    PK      string
    SK      string
    RowType string
}
