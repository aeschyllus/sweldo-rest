package payrolls

import (
	"context"
	"fmt"
	"time"

	"github.com/aeschyllus/sweldo-rest/internal/adapters/postgresql/sqlc"
	"github.com/aeschyllus/sweldo-rest/internal/pkg/money"
	"github.com/jackc/pgx/v5/pgtype"
)

func NewService(repo sqlc.Querier) Service {
	return &service{repo}
}

func (s *service) CreatePayrollRun(ctx context.Context, params CreatePayrollRunParams) (PayrollRunResponse, error) {
	runDate, err := time.Parse("2006-01-02", params.RunDate)
	if err != nil {
		return PayrollRunResponse{}, err
	}

	totalPay, err := money.ToNumeric(params.TotalPay)
	if err != nil {
		return PayrollRunResponse{}, err
	}

	sqlcParams := sqlc.CreatePayrollRunParams{
		CompanyID:      params.CompanyID,
		RunDate:        toDate(runDate),
		TotalEmployees: params.TotalEmployees,
		TotalPay:       totalPay,
		CreatedBy:      pgtype.Int8{Int64: params.CreatedBy, Valid: true},
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
	totalPay, err := money.ToNumeric(params.TotalPay)
	if err != nil {
		return PayrollRunResponse{}, err
	}

	sqlcParams := sqlc.UpdatePayrollRunByIDParams{
		TotalEmployees: params.TotalEmployees,
		TotalPay:       totalPay,
		ID:             id,
		UpdatedBy:      pgtype.Int8{Int64: params.UpdatedBy, Valid: true},
	}

	run, err := s.repo.UpdatePayrollRunByID(ctx, sqlcParams)
	if err != nil {
		return PayrollRunResponse{}, err
	}

	return toPayrollRunResponse(run), nil
}

func (s *service) FinalizePayrollRun(ctx context.Context, id int64, userID int64) (PayrollRunResponse, error) {
	run, err := s.repo.FinalizePayrollRunByID(ctx, sqlc.FinalizePayrollRunByIDParams{
		ID:        id,
		UpdatedBy: pgtype.Int8{Int64: userID, Valid: true},
	})
	if err != nil {
		return PayrollRunResponse{}, fmt.Errorf("finalize payroll run: %w", err)
	}

	return toPayrollRunResponse(run), nil
}

func (s *service) CreatePayrollDetail(ctx context.Context, runID int64, params CreatePayrollDetailParams) (PayrollDetailResponse, error) {
	run, err := s.repo.FindPayrollRunByID(ctx, runID)
	if err != nil {
		return PayrollDetailResponse{}, fmt.Errorf("payroll run not found: %w", err)
	}

	if run.Status == "FINALIZED" {
		return PayrollDetailResponse{}, fmt.Errorf("cannot add details to a finalized payroll run")
	}

	grossPay, err := money.ToNumeric(params.GrossPay)
	if err != nil {
		return PayrollDetailResponse{}, err
	}

	taxDeduction, err := money.ToNumeric(params.TaxDeduction)
	if err != nil {
		return PayrollDetailResponse{}, err
	}

	netPayCents := params.GrossPay - params.TaxDeduction
	netPay, err := money.ToNumeric(netPayCents)
	if err != nil {
		return PayrollDetailResponse{}, err
	}

	hourlyRate, err := money.ToNumeric(params.HourlyRate)
	if err != nil {
		return PayrollDetailResponse{}, err
	}

	hoursWorked, err := money.ToNumeric(params.HoursWorked)
	if err != nil {
		return PayrollDetailResponse{}, err
	}

	sqlcParams := sqlc.CreatePayrollDetailParams{
		PayrollRunID: runID,
		EmployeeID:   params.EmployeeID,
		GrossPay:     grossPay,
		TaxDeduction: taxDeduction,
		NetPay:       netPay,
		HourlyRate:   hourlyRate,
		HoursWorked:  hoursWorked,
		CreatedBy:    pgtype.Int8{Int64: params.CreatedBy, Valid: true},
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

func (s *service) CreateDeduction(ctx context.Context, params CreateDeductionParams) (DeductionResponse, error) {
	amount, err := money.ToNumeric(params.Amount)
	if err != nil {
		return DeductionResponse{}, err
	}

	deduction, err := s.repo.CreateDeduction(ctx, sqlc.CreateDeductionParams{
		PayrollDetailID: params.PayrollDetailID,
		DeductionType:   params.DeductionType,
		Amount:          amount,
		CreatedBy:       pgtype.Int8{Int64: params.CreatedBy, Valid: true},
	})
	if err != nil {
		return DeductionResponse{}, err
	}

	return toDeductionResponse(deduction), nil
}

func (s *service) ListDeductionsByDetailID(ctx context.Context, detailID int64) ([]DeductionResponse, error) {
	deductions, err := s.repo.ListDeductionsByPayrollDetailID(ctx, detailID)
	if err != nil {
		return nil, err
	}

	responses := make([]DeductionResponse, len(deductions))
	for i, d := range deductions {
		responses[i] = toDeductionResponse(d)
	}

	return responses, nil
}

func (s *service) DeleteDeduction(ctx context.Context, id int64) error {
	_, err := s.repo.DeleteDeduction(ctx, id)
	return err
}

func toPayrollRunResponse(run sqlc.PayrollRun) PayrollRunResponse {
	cents, err := money.FromNumeric(run.TotalPay)
	if err != nil {
		cents = 0
	}

	return PayrollRunResponse{
		ID:             run.ID,
		CompanyID:      run.CompanyID,
		RunDate:        fromDate(run.RunDate),
		TotalEmployees: run.TotalEmployees,
		TotalPay:       money.FormatCents(cents),
		Status:         run.Status,
		CreatedAt:      fromTimestamptz(run.CreatedAt),
	}
}

func toPayrollDetailResponse(detail sqlc.PayrollDetail) PayrollDetailResponse {
	gross, _ := money.FromNumeric(detail.GrossPay)
	tax, _ := money.FromNumeric(detail.TaxDeduction)
	net, _ := money.FromNumeric(detail.NetPay)
	rate, _ := money.FromNumeric(detail.HourlyRate)
	hours, _ := money.FromNumeric(detail.HoursWorked)

	return PayrollDetailResponse{
		ID:           detail.ID,
		PayrollRunID: detail.PayrollRunID,
		EmployeeID:   detail.EmployeeID,
		GrossPay:     money.FormatCents(gross),
		TaxDeduction: money.FormatCents(tax),
		NetPay:       money.FormatCents(net),
		HourlyRate:   money.FormatCents(rate),
		HoursWorked:  money.FormatCents(hours),
		CreatedAt:    fromTimestamptz(detail.CreatedAt),
	}
}

func toDeductionResponse(d sqlc.Deduction) DeductionResponse {
	amt, _ := money.FromNumeric(d.Amount)

	return DeductionResponse{
		ID:              d.ID,
		PayrollDetailID: d.PayrollDetailID,
		DeductionType:   d.DeductionType,
		Amount:          money.FormatCents(amt),
		CreatedAt:       fromTimestamptz(d.CreatedAt),
	}
}

func toDate(t time.Time) pgtype.Date {
	return pgtype.Date{
		Time:  t,
		Valid: true,
	}
}

func fromDate(d pgtype.Date) string {
	if !d.Valid {
		return ""
	}
	return d.Time.Format("2006-01-02")
}

func fromTimestamptz(ts pgtype.Timestamptz) time.Time {
	if !ts.Valid {
		return time.Time{}
	}
	return ts.Time
}
