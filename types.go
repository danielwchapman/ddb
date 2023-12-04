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

type RowGSI1Header struct {
	GSI1PK string
	GSI1SK string
}

type RowGSI2Header struct {
	GSI2PK string
	GSI2SK string
}

type RowGSI3Header struct {
	GSI3PK string
	GSI3SK string
}

type RowGSI4Header struct {
	GSI4PK string
	GSI4SK string
}

type RowGSI5Header struct {
	GSI5PK string
	GSI5SK string
}
