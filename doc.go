/*
Package memeduck provides tools to build Spanner SQL queries.


Supported Types

The following types can be used as a SQL expression:

  * If a value implements `ToASTExpr() (*ast.Expr, error)`, memeduck uses this method to convert Go values into SQL expressions.
  * If a value is nil (of any type), it is converted into NULL.
  * If a value is one of string, *string, or spanner.NullString, it is converted into STRING literal.
  * If a value is []byte, it is converted into BYTES literal.
  * If a value is one of int, *int, int64, or *int64, spanner.NullInt64, it is converted into INT64 literal.
  * If a value is one of bool, *bool, or spanner.NullBool, it is converted into BOOL literal.
  * If a value is one of float64, *float64, or spanner.NullFloat64, it is converted into FLOAT64 literal.
  * If a value is one of time.Time, *time.Time, or spanner.NullTime, it is converted into TIMESTAMP literal.
  * If a value is one of civil.Date, *civil.Date, or spanner.NullDate, it is converted into DATE literal.
  * If a value is a slice of the above types, it is converted into ARRAY<T> literal.


Struct Tags

You can add `spanner:"Name"` tag to struct fields to indicate which field in struct corresponds to which column, otherwise memeduck uses field name as column name.
See examples section of Insert function for more details.
*/
package memeduck
