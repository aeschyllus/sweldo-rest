package employees

import (
	"fmt"
	"net/http"
	"strconv"
)

func parseListEmployeesQuery(r *http.Request) (ListEmployeesParams, error) {
	query := r.URL.Query()

	companyIDStr := query.Get("company_id")
	if companyIDStr == "" {
		return ListEmployeesParams{}, fmt.Errorf("company_id is required")
	}

	companyID, err := strconv.ParseInt(companyIDStr, 10, 64)
	if err != nil || companyID <= 0 {
		return ListEmployeesParams{}, fmt.Errorf("Invalid company_id")
	}

	limit := int32(10)
	offset := int32(0)

	if l := query.Get("limit"); l != "" {
		val, err := strconv.Atoi(l)
		if err != nil || val <= 0 {
			return ListEmployeesParams{}, fmt.Errorf("Invalid limit")
		}
		limit = int32(val)
	}

	if o := query.Get("offset"); o != "" {
		val, err := strconv.Atoi(o)
		if err != nil || val < 0 {
			return ListEmployeesParams{}, fmt.Errorf("Invalid offset")
		}
		offset = int32(val)
	}

	var name *string
	if n := query.Get("name"); n != "" {
		name = &n
	}

	return ListEmployeesParams{
		CompanyID:  companyID,
		Name:       name,
		PageLimit:  limit,
		PageOffset: offset,
	}, nil
}
