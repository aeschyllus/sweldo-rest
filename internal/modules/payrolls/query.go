package payrolls

import (
	"net/http"
	"strconv"
)

func parseListPayrollRunsQuery(r *http.Request) (listPayrollRunsQuery, error) {
	companyIDStr := r.URL.Query().Get("company_id")
	if companyIDStr == "" {
		return listPayrollRunsQuery{}, nil
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
		return listPayrollDetailsQuery{}, nil
	}

	employeeID, err := strconv.ParseInt(employeeIDStr, 10, 64)
	if err != nil {
		return listPayrollDetailsQuery{}, err
	}

	return listPayrollDetailsQuery{EmployeeID: employeeID}, nil
}
