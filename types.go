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

type RowType string

// assume composite primary key where columns are named 'PK' and 'SK'

type RowHeader struct {
    PK      string
    SK      string
    RowType string
}
