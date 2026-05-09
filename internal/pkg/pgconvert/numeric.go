package pgconvert

import "github.com/jackc/pgx/v5/pgtype"

func ToNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	n.Scan(f)
	return n
}

func FromNumeric(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	f, _ := n.Float64Value()
	return f.Float64
}
