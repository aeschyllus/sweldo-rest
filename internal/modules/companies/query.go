package companies

import (
	"fmt"
	"net/http"
	"strconv"
)

func parseListCompaniesQuery(r *http.Request) (ListCompaniesParams, error) {
	query := r.URL.Query()

	limit := int32(10)
	offset := int32(0)

	if l := query.Get("limit"); l != "" {
		val, err := strconv.Atoi(l)
		if err != nil || val <= 0 {
			return ListCompaniesParams{}, fmt.Errorf("Invalid limit")
		}
		limit = int32(val)
	}

	if o := query.Get("offset"); o != "" {
		val, err := strconv.Atoi(o)
		if err != nil || val < 0 {
			return ListCompaniesParams{}, fmt.Errorf("Invalid offset")
		}
		offset = int32(val)
	}

	var name *string
	if n := query.Get("name"); n != "" {
		name = &n
	}

	return ListCompaniesParams{
		Name:       name,
		PageLimit:  limit,
		PageOffset: offset,
	}, nil
}
