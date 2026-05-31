package pgconvert

import (
	"fmt"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
)

func ToNumeric(cents int64) (pgtype.Numeric, error) {
	val := big.NewRat(cents, 100)
	var n pgtype.Numeric
	if err := n.Scan(val.FloatString(2)); err != nil {
		return pgtype.Numeric{}, err
	}
	return n, nil
}

func FromNumeric(n pgtype.Numeric) (int64, error) {
	if !n.Valid {
		return 0, fmt.Errorf("numeric is null")
	}

	f, err := n.Float64Value()
	if err != nil {
		return 0, err
	}

	// Convert to cents with rounding
	cents := int64(f.Float64*100 + 0.5)
	return cents, nil
}
