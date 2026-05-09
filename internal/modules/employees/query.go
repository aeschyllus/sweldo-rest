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

	var name *string
	if n := query.Get("name"); n != "" {
		name = &n
	}

	return ListEmployeesParams{
		CompanyID: companyID,
		Name:      name,
	}, nil
}
