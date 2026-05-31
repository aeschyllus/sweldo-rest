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
	FinalizePayrollRun(ctx context.Context, id int64) (PayrollRunResponse, error)
	CreatePayrollDetail(ctx context.Context, runID int64, params CreatePayrollDetailParams) (PayrollDetailResponse, error)
	ListAllPayrollDetailsByRunID(ctx context.Context, runID int64) ([]PayrollDetailResponse, error)
	ListAllPayrollDetailsByEmployeeID(ctx context.Context, employeeID int64) ([]PayrollDetailResponse, error)
	CreateDeduction(ctx context.Context, params CreateDeductionParams) (DeductionResponse, error)
	ListDeductionsByDetailID(ctx context.Context, detailID int64) ([]DeductionResponse, error)
	DeleteDeduction(ctx context.Context, id int64) error
}

type createPayrollRunRequest struct {
	RunDate        string `json:"run_date"`
	TotalEmployees int32  `json:"total_employees"`
	TotalPay       string `json:"total_pay"`
}

type updatePayrollRunRequest struct {
	TotalEmployees int32  `json:"total_employees"`
	TotalPay       string `json:"total_pay"`
}

type createPayrollDetailRequest struct {
	EmployeeID   int64   `json:"employee_id"`
	GrossPay     string  `json:"gross_pay"`
	TaxDeduction string  `json:"tax_deduction"`
	HourlyRate   *string `json:"hourly_rate,omitempty"`
	HoursWorked  *string `json:"hours_worked,omitempty"`
}

type createDeductionRequest struct {
	DeductionType string `json:"deduction_type"`
	Amount        string `json:"amount"`
}

type CreatePayrollRunParams struct {
	CompanyID      int64
	RunDate        string
	TotalEmployees int32
	TotalPay       int64
}

type UpdatePayrollRunParams struct {
	TotalEmployees int32
	TotalPay       int64
}

type CreatePayrollDetailParams struct {
	EmployeeID   int64
	GrossPay     int64
	TaxDeduction int64
	HourlyRate   int64
	HoursWorked  int64
}

type CreateDeductionParams struct {
	PayrollDetailID int64
	DeductionType   string
	Amount          int64
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
	TotalPay       string    `json:"total_pay"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

type PayrollDetailResponse struct {
	ID           int64     `json:"id"`
	PayrollRunID int64     `json:"payroll_run_id"`
	EmployeeID   int64     `json:"employee_id"`
	GrossPay     string    `json:"gross_pay"`
	TaxDeduction string    `json:"tax_deduction"`
	NetPay       string    `json:"net_pay"`
	HourlyRate   string    `json:"hourly_rate"`
	HoursWorked  string    `json:"hours_worked"`
	CreatedAt    time.Time `json:"created_at"`
}

type DeductionResponse struct {
	ID              int64     `json:"id"`
	PayrollDetailID int64     `json:"payroll_detail_id"`
	DeductionType   string    `json:"deduction_type"`
	Amount          string    `json:"amount"`
	CreatedAt       time.Time `json:"created_at"`
}
