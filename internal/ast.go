package internal

import (
	"reflect"
	"strconv"
	"time"

	"cloud.google.com/go/civil"
	"cloud.google.com/go/spanner"
	"github.com/MakeNowJust/memefish/pkg/ast"
	"github.com/pkg/errors"
)

// ASTExpr is a type that can be converted into ast.Expr.
type ASTExpr interface {
	ToASTExpr() (ast.Expr, error)
}

func ToExpr(val interface{}) (ast.Expr, error) {
	switch v := val.(type) {
	case nil:
		return NullLit(), nil
	case string:
		return StringLit(v), nil
	case *string:
		if v == nil {
			return NullLit(), nil
		}
		return StringLit(*v), nil
	case spanner.NullString:
		if !v.Valid {
			return NullLit(), nil
		}
		return StringLit(v.StringVal), nil
	case []byte:
		if v == nil {
			return NullLit(), nil
		}
		return BytesLit(v), nil
	case int:
		return IntLit(int64(v)), nil
	case *int:
		if v == nil {
			return NullLit(), nil
		}
		return IntLit(int64(*v)), nil
	case int64:
		return IntLit(v), nil
	case *int64:
		if v == nil {
			return NullLit(), nil
		}
		return IntLit(*v), nil
	case spanner.NullInt64:
		if !v.Valid {
			return NullLit(), nil
		}
		return IntLit(v.Int64), nil
	case bool:
		return BoolLit(v), nil
	case *bool:
		if v == nil {
			return NullLit(), nil
		}
		return BoolLit(*v), nil
	case spanner.NullBool:
		if !v.Valid {
			return NullLit(), nil
		}
		return BoolLit(v.Bool), nil
	case float64:
		return FloatLit(v), nil
	case *float64:
		if v == nil {
			return NullLit(), nil
		}
		return FloatLit(*v), nil
	case spanner.NullFloat64:
		if !v.Valid {
			return NullLit(), nil
		}
		return FloatLit(v.Float64), nil
	case time.Time:
		return TimeLit(v), nil
	case *time.Time:
		if v == nil {
			return NullLit(), nil
		}
		return TimeLit(*v), nil
	case spanner.NullTime:
		if !v.Valid {
			return NullLit(), nil
		}
		return TimeLit(v.Time), nil
	case civil.Date:
		return DateLit(v), nil
	case *civil.Date:
		if v == nil {
			return NullLit(), nil
		}
		return DateLit(*v), nil
	case spanner.NullDate:
		if !v.Valid {
			return NullLit(), nil
		}
		return DateLit(v.Date), nil
	default:
		if se, ok := val.(ASTExpr); ok {
			return se.ToASTExpr()
		}
		// TODO: support big.Rat
		// Slices
		valV := reflect.ValueOf(val)
		if valV.Type().Kind() == reflect.Slice {
			exprs := make([]ast.Expr, 0, valV.Len())
			for i := 0; i < valV.Len(); i++ {
				vi := valV.Index(i).Interface()
				ei, err := ToExpr(vi)
				if err != nil {
					return nil, errors.WithMessagef(err, "at index %d", i)
				}
				exprs = append(exprs, ei)
			}
			return ArrayLit(exprs), nil
		} else {
			// TODO: support Go structs
			return nil, errors.Errorf("can't convert %T into SQL expr", val)

		}
	}
}

func StringLit(v string) *ast.StringLiteral {
	return &ast.StringLiteral{
		Value: v,
	}
}

func BytesLit(v []byte) *ast.BytesLiteral {
	return &ast.BytesLiteral{
		Value: v,
	}
}

func IntLit(v int64) *ast.IntLiteral {
	return &ast.IntLiteral{
		Base:  10,
		Value: strconv.FormatInt(v, 10),
	}
}

func BoolLit(v bool) *ast.BoolLiteral {
	return &ast.BoolLiteral{
		Value: v,
	}
}

func FloatLit(v float64) *ast.FloatLiteral {
	return &ast.FloatLiteral{
		Value: strconv.FormatFloat(v, 'e', -1, 64),
	}
}

func TimeLit(v time.Time) *ast.TimestampLiteral {
	return &ast.TimestampLiteral{
		Value: &ast.StringLiteral{
			Value: v.Format(time.RFC3339Nano),
		},
	}
}

func DateLit(v civil.Date) *ast.DateLiteral {
	return &ast.DateLiteral{
		Value: &ast.StringLiteral{
			Value: v.String(),
		},
	}
}

func ArrayLit(exprs []ast.Expr) *ast.ArrayLiteral {
	return &ast.ArrayLiteral{
		Values: exprs,
	}
}

func NullLit() *ast.NullLiteral {
	return &ast.NullLiteral{}
}
