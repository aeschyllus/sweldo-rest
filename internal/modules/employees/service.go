package employees

import (
	"context"

	"github.com/aeschyllus/sweldo-rest/internal/adapters/postgresql/sqlc"
	"github.com/aeschyllus/sweldo-rest/internal/pkg/money"
	"github.com/aeschyllus/sweldo-rest/internal/pkg/pgconvert"
)

func NewService(repo sqlc.Querier) Service {
	return &service{repo}
}

func (s *service) CreateEmployee(ctx context.Context, params CreateEmployeeParams) (EmployeeResponse, error) {
	baseSalary, err := pgconvert.ToNumeric(params.BaseSalary)
	if err != nil {
		return EmployeeResponse{}, err
	}

	employee, err := s.repo.CreateEmployee(ctx, sqlc.CreateEmployeeParams{
		CompanyID:      params.CompanyID,
		FirstName:      params.FirstName,
		LastName:       params.LastName,
		EmploymentType: params.EmploymentType,
		SalaryType:     params.SalaryType,
		BaseSalary:     baseSalary,
	})
	if err != nil {
		return EmployeeResponse{}, err
	}

	return toEmployeeResponse(employee), nil
}

func (s *service) ListEmployeesByCompanyID(ctx context.Context, params ListEmployeesParams) ([]EmployeeResponse, error) {
	employees, err := s.repo.ListEmployeesByCompanyID(ctx, sqlc.ListEmployeesByCompanyIDParams{
		CompanyID: params.CompanyID,
		Name:      pgconvert.ToText(params.Name),
	})
	if err != nil {
		return nil, err
	}

	responses := make([]EmployeeResponse, len(employees))
	for i, emp := range employees {
		responses[i] = toEmployeeResponse(emp)
	}

	return responses, nil
}

func (s *service) FindEmployeeByID(ctx context.Context, params FindEmployeeParams) (EmployeeResponse, error) {
	employee, err := s.repo.FindEmployeeByID(ctx, sqlc.FindEmployeeByIDParams{
		ID:        params.ID,
		CompanyID: params.CompanyID,
	})
	if err != nil {
		return EmployeeResponse{}, err
	}

	return toEmployeeResponse(employee), nil
}

func (s *service) UpdateEmployeeByID(ctx context.Context, params UpdateEmployeeParams) (EmployeeResponse, error) {
	baseSalary, err := pgconvert.ToNumeric(params.BaseSalary)
	if err != nil {
		return EmployeeResponse{}, err
	}

	employee, err := s.repo.UpdateEmployeeByID(ctx, sqlc.UpdateEmployeeByIDParams{
		ID:             params.ID,
		FirstName:      params.FirstName,
		LastName:       params.LastName,
		EmploymentType: params.EmploymentType,
		SalaryType:     params.SalaryType,
		BaseSalary:     baseSalary,
	})
	if err != nil {
		return EmployeeResponse{}, err
	}

	return toEmployeeResponse(employee), nil
}

func toEmployeeResponse(e sqlc.Employee) EmployeeResponse {
	cents, err := pgconvert.FromNumeric(e.BaseSalary)
	if err != nil {
		cents = 0
	}

	return EmployeeResponse{
		ID:             e.ID,
		CompanyID:      e.CompanyID,
		FirstName:      e.FirstName,
		LastName:       e.LastName,
		EmploymentType: e.EmploymentType,
		SalaryType:     e.SalaryType,
		BaseSalary:     money.FormatCents(cents),
		CreatedAt:      e.CreatedAt.Time,
		UpdatedAt:      e.UpdatedAt.Time,
	}
}
