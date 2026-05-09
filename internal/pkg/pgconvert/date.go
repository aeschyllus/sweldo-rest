package pgconvert

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ToDate(t time.Time) pgtype.Date {
	return pgtype.Date{
		Time:  t,
		Valid: true,
	}
}

func FromDate(d pgtype.Date) string {
	if !d.Valid {
		return ""
	}
	return d.Time.Format("2026-01-02")
}
