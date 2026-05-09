package pgconvert

import (
	"fmt"
	"math"

	"github.com/jackc/pgx/v5/pgtype"
)

func ToNumeric(f float64) (pgtype.Numeric, error) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return pgtype.Numeric{}, fmt.Errorf("invalid numeric value: %v", f)
	}
	var n pgtype.Numeric
	if err := n.Scan(f); err != nil {
		return pgtype.Numeric{}, err
	}
	return n, nil
}

func FromNumeric(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	f, _ := n.Float64Value()
	return f.Float64
}
