package pgconvert

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func FromTimestamptz(ts pgtype.Timestamptz) time.Time {
	if !ts.Valid {
		return time.Time{}
	}
	return ts.Time
}
