package payrolls

import (
	"context"
	"time"

	"github.com/aeschyllus/sweldo-rest/internal/adapters/postgresql/sqlc"
	"github.com/aeschyllus/sweldo-rest/internal/pkg/money"
	"github.com/aeschyllus/sweldo-rest/internal/pkg/pgconvert"
)

func NewService(repo sqlc.Querier) Service {
	return &service{repo}
}

// Payroll Run operations

func (s *service) CreatePayrollRun(ctx context.Context, params CreatePayrollRunParams) (PayrollRunResponse, error) {
	runDate, err := time.Parse("2006-01-02", params.RunDate)
	if err != nil {
		return PayrollRunResponse{}, err
	}

	totalPay, err := pgconvert.ToNumeric(params.TotalPay)
	if err != nil {
		return PayrollRunResponse{}, err
	}

	sqlcParams := sqlc.CreatePayrollRunParams{
		CompanyID:      params.CompanyID,
		RunDate:        pgconvert.ToDate(runDate),
		TotalEmployees: params.TotalEmployees,
		TotalPay:       totalPay,
	}

	run, err := s.repo.CreatePayrollRun(ctx, sqlcParams)
	if err != nil {
		return PayrollRunResponse{}, err
	}

	return toPayrollRunResponse(run), nil
}

func (s *service) ListPayrollRunsByCompanyID(ctx context.Context, companyID int64) ([]PayrollRunResponse, error) {
	runs, err := s.repo.ListPayrollRunsByCompanyID(ctx, companyID)
	if err != nil {
		return nil, err
	}

	responses := make([]PayrollRunResponse, len(runs))
	for i, run := range runs {
		responses[i] = toPayrollRunResponse(run)
	}

	return responses, nil
}

func (s *service) FindPayrollRunByID(ctx context.Context, id int64) (PayrollRunResponse, error) {
	run, err := s.repo.FindPayrollRunByID(ctx, id)
	if err != nil {
		return PayrollRunResponse{}, err
	}

	return toPayrollRunResponse(run), nil
}

func (s *service) UpdatePayrollRunByID(ctx context.Context, id int64, params UpdatePayrollRunParams) (PayrollRunResponse, error) {
	totalPay, err := pgconvert.ToNumeric(params.TotalPay)
	if err != nil {
		return PayrollRunResponse{}, err
	}

	sqlcParams := sqlc.UpdatePayrollRunByIDParams{
		TotalEmployees: params.TotalEmployees,
		TotalPay:       totalPay,
		ID:             id,
	}

	run, err := s.repo.UpdatePayrollRunByID(ctx, sqlcParams)
	if err != nil {
		return PayrollRunResponse{}, err
	}

	return toPayrollRunResponse(run), nil
}

// Payroll Detail operations

func (s *service) CreatePayrollDetail(ctx context.Context, runID int64, params CreatePayrollDetailParams) (PayrollDetailResponse, error) {
	grossPay, err := pgconvert.ToNumeric(params.GrossPay)
	if err != nil {
		return PayrollDetailResponse{}, err
	}

	taxDeduction, err := pgconvert.ToNumeric(params.TaxDeduction)
	if err != nil {
		return PayrollDetailResponse{}, err
	}

	netPay, err := pgconvert.ToNumeric(params.NetPay)
	if err != nil {
		return PayrollDetailResponse{}, err
	}

	sqlcParams := sqlc.CreatePayrollDetailParams{
		PayrollRunID: runID,
		EmployeeID:   params.EmployeeID,
		GrossPay:     grossPay,
		TaxDeduction: taxDeduction,
		NetPay:       netPay,
	}

	detail, err := s.repo.CreatePayrollDetail(ctx, sqlcParams)
	if err != nil {
		return PayrollDetailResponse{}, err
	}

	return toPayrollDetailResponse(detail), nil
}

func (s *service) ListAllPayrollDetailsByRunID(ctx context.Context, runID int64) ([]PayrollDetailResponse, error) {
	details, err := s.repo.ListAllPayrollDetailsByRunID(ctx, runID)
	if err != nil {
		return nil, err
	}

	responses := make([]PayrollDetailResponse, len(details))
	for i, detail := range details {
		responses[i] = toPayrollDetailResponse(detail)
	}

	return responses, nil
}

func (s *service) ListAllPayrollDetailsByEmployeeID(ctx context.Context, employeeID int64) ([]PayrollDetailResponse, error) {
	details, err := s.repo.ListAllPayrollDetailsByEmployeeID(ctx, employeeID)
	if err != nil {
		return nil, err
	}

	responses := make([]PayrollDetailResponse, len(details))
	for i, detail := range details {
		responses[i] = toPayrollDetailResponse(detail)
	}

	return responses, nil
}

// Helper functions (package-level, not methods)

func toPayrollRunResponse(run sqlc.PayrollRun) PayrollRunResponse {
	cents, err := pgconvert.FromNumeric(run.TotalPay)
	if err != nil {
		cents = 0
	}

	return PayrollRunResponse{
		ID:             run.ID,
		CompanyID:      run.CompanyID,
		RunDate:        pgconvert.FromDate(run.RunDate),
		TotalEmployees: run.TotalEmployees,
		TotalPay:       money.FormatCents(cents),
		CreatedAt:      pgconvert.FromTimestamptz(run.CreatedAt),
	}
}

func toPayrollDetailResponse(detail sqlc.PayrollDetail) PayrollDetailResponse {
	gross, _ := pgconvert.FromNumeric(detail.GrossPay)
	tax, _ := pgconvert.FromNumeric(detail.TaxDeduction)
	net, _ := pgconvert.FromNumeric(detail.NetPay)

	return PayrollDetailResponse{
		ID:           detail.ID,
		PayrollRunID: detail.PayrollRunID,
		EmployeeID:   detail.EmployeeID,
		GrossPay:     money.FormatCents(gross),
		TaxDeduction: money.FormatCents(tax),
		NetPay:       money.FormatCents(net),
		CreatedAt:    pgconvert.FromTimestamptz(detail.CreatedAt),
	}
}
