package pgconvert

import "github.com/jackc/pgx/v5/pgtype"

func ToText(s *string) pgtype.Text {
	var name pgtype.Text

	if s != nil {
		name = pgtype.Text{
			String: *s,
			Valid:  true,
		}
	}

	return name
}
