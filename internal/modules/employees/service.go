package employees

import (
	"context"

	"github.com/aeschyllus/sweldo-rest/internal/adapters/postgresql/sqlc"
	"github.com/aeschyllus/sweldo-rest/internal/pkg/pgconvert"
)

func NewService(repo sqlc.Querier) Service {
	return &service{repo}
}

func (s *service) CreateEmployee(ctx context.Context, params CreateEmployeeParams) (sqlc.Employee, error) {
	baseSalary, err := pgconvert.ToNumeric(params.BaseSalary)
	if err != nil {
		return sqlc.Employee{}, err
	}

	return s.repo.CreateEmployee(ctx, sqlc.CreateEmployeeParams{
		CompanyID:      params.CompanyID,
		FirstName:      params.FirstName,
		LastName:       params.LastName,
		EmploymentType: params.EmploymentType,
		SalaryType:     params.SalaryType,
		BaseSalary:     baseSalary,
	})
}

func (s *service) ListEmployeesByCompanyID(ctx context.Context, params ListEmployeesParams) ([]sqlc.Employee, error) {
	return s.repo.ListEmployeesByCompanyID(ctx, sqlc.ListEmployeesByCompanyIDParams{
		CompanyID: params.CompanyID,
		Name:      pgconvert.ToText(params.Name),
	})
}

func (s *service) FindEmployeeByID(ctx context.Context, params FindEmployeeParams) (sqlc.Employee, error) {
	return s.repo.FindEmployeeByID(ctx, sqlc.FindEmployeeByIDParams{
		ID:        params.ID,
		CompanyID: params.CompanyID,
	})
}

func (s *service) UpdateEmployeeByID(ctx context.Context, params UpdateEmployeeParams) (sqlc.Employee, error) {
	baseSalary, err := pgconvert.ToNumeric(params.BaseSalary)
	if err != nil {
		return sqlc.Employee{}, err
	}

	return s.repo.UpdateEmployeeByID(ctx, sqlc.UpdateEmployeeByIDParams{
		ID:             params.ID,
		FirstName:      params.FirstName,
		LastName:       params.LastName,
		EmploymentType: params.EmploymentType,
		SalaryType:     params.SalaryType,
		BaseSalary:     baseSalary,
	})
}
