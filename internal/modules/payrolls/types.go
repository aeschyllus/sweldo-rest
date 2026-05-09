package payrolls

import (
	"context"
	"time"

	"github.com/aeschyllus/sweldo-rest/internal/adapters/postgresql/sqlc"
)

type handler struct {
	service Service
}

type service struct {
	repo sqlc.Querier
}

type Service interface {
	CreatePayrollRun(ctx context.Context, params CreatePayrollRunParams) (PayrollRunResponse, error)
	ListPayrollRunsByCompanyID(ctx context.Context, companyID int64) ([]PayrollRunResponse, error)
	FindPayrollRunByID(ctx context.Context, id int64) (PayrollRunResponse, error)
	UpdatePayrollRunByID(ctx context.Context, id int64, params UpdatePayrollRunParams) (PayrollRunResponse, error)
	CreatePayrollDetail(ctx context.Context, runID int64, params CreatePayrollDetailParams) (PayrollDetailResponse, error)
	ListAllPayrollDetailsByRunID(ctx context.Context, runID int64) ([]PayrollDetailResponse, error)
	ListAllPayrollDetailsByEmployeeID(ctx context.Context, employeeID int64) ([]PayrollDetailResponse, error)
}

type createPayrollRunRequest struct {
	CompanyID      int64   `json:"company_id"`
	RunDate        string  `json:"run_date"`
	TotalEmployees int32   `json:"total_employees"`
	TotalPay       float64 `json:"total_pay"`
}

type updatePayrollRunRequest struct {
	TotalEmployees int32   `json:"total_employees"`
	TotalPay       float64 `json:"total_pay"`
}

type createPayrollDetailRequest struct {
	EmployeeID   int64   `json:"employee_id"`
	GrossPay     float64 `json:"gross_pay"`
	TaxDeduction float64 `json:"tax_deduction"`
	NetPay       float64 `json:"net_pay"`
}

type CreatePayrollRunParams struct {
	CompanyID      int64
	RunDate        string
	TotalEmployees int32
	TotalPay       float64
}

type UpdatePayrollRunParams struct {
	TotalEmployees int32
	TotalPay       float64
}

type CreatePayrollDetailParams struct {
	EmployeeID   int64
	GrossPay     float64
	TaxDeduction float64
	NetPay       float64
}

type listPayrollRunsQuery struct {
	CompanyID int64
}

type listPayrollDetailsQuery struct {
	EmployeeID int64
}

type PayrollRunResponse struct {
	ID             int64     `json:"id"`
	CompanyID      int64     `json:"company_id"`
	RunDate        string    `json:"run_date"`
	TotalEmployees int32     `json:"total_employees"`
	TotalPay       float64   `json:"total_pay"`
	CreatedAt      time.Time `json:"created_at"`
}

type PayrollDetailResponse struct {
	ID           int64     `json:"id"`
	PayrollRunID int64     `json:"payroll_run_id"`
	EmployeeID   int64     `json:"employee_id"`
	GrossPay     float64   `json:"gross_pay"`
	TaxDeduction float64   `json:"tax_deduction"`
	NetPay       float64   `json:"net_pay"`
	CreatedAt    time.Time `json:"created_at"`
}
