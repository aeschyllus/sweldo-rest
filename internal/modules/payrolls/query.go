package payrolls

import (
	"fmt"
	"net/http"
	"strconv"
)

func parseListPayrollRunsQuery(r *http.Request) (listPayrollRunsQuery, error) {
	companyIDStr := r.URL.Query().Get("company_id")
	if companyIDStr == "" {
		return listPayrollRunsQuery{}, fmt.Errorf("company_id is required")
	}

	companyID, err := strconv.ParseInt(companyIDStr, 10, 64)
	if err != nil {
		return listPayrollRunsQuery{}, err
	}

	return listPayrollRunsQuery{CompanyID: companyID}, nil
}

func parseListPayrollDetailsQuery(r *http.Request) (listPayrollDetailsQuery, error) {
	employeeIDStr := r.URL.Query().Get("employee_id")
	if employeeIDStr == "" {
		return listPayrollDetailsQuery{}, fmt.Errorf("employee_id is required")
	}

	employeeID, err := strconv.ParseInt(employeeIDStr, 10, 64)
	if err != nil {
		return listPayrollDetailsQuery{}, err
	}

	return listPayrollDetailsQuery{EmployeeID: employeeID}, nil
}
