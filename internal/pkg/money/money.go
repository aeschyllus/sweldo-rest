package money

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

func ParseCents(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty money string")
	}

	negative := s[0] == '-'
	if negative {
		s = s[1:]
	}

	parts := strings.SplitN(s, ".", 2)
	var whole, frac int64

	whole, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid money value: %s", s)
	}

	if len(parts) == 2 {
		if len(parts[1]) > 2 {
			frac, err = strconv.ParseInt(parts[1][:2], 10, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid money value: %s", s)
			}
		} else {
			frac, err = strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid money value: %s", s)
			}
			// Pad to 2 decimal places
			frac *= int64(math.Pow10(2 - len(parts[1])))
		}
	}

	cents := whole*100 + frac
	if negative {
		cents = -cents
	}

	return cents, nil
}

func FormatCents(c int64) string {
	sign := ""
	if c < 0 {
		sign = "-"
		c = -c
	}
	whole := c / 100
	frac := c % 100
	return fmt.Sprintf("%s%d.%02d", sign, whole, frac)
}

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

	cents := int64(f.Float64*100 + 0.5)
	return cents, nil
}
