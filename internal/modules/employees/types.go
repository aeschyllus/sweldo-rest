package employees

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
	CreateEmployee(ctx context.Context, params CreateEmployeeParams) (EmployeeResponse, error)
	ListEmployeesByCompanyID(ctx context.Context, params ListEmployeesParams) ([]EmployeeResponse, error)
	FindEmployeeByID(ctx context.Context, params FindEmployeeParams) (EmployeeResponse, error)
	UpdateEmployeeByID(ctx context.Context, params UpdateEmployeeParams) (EmployeeResponse, error)
}

type createEmployeeRequest struct {
	CompanyID      int64  `json:"company_id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	EmploymentType string `json:"employment_type"`
	SalaryType     string `json:"salary_type"`
	BaseSalary     string `json:"base_salary"`
}

type updateEmployeeRequest struct {
	ID             int64  `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	EmploymentType string `json:"employment_type"`
	SalaryType     string `json:"salary_type"`
	BaseSalary     string `json:"base_salary"`
}

type CreateEmployeeParams struct {
	CompanyID      int64
	FirstName      string
	LastName       string
	EmploymentType string
	SalaryType     string
	BaseSalary     int64
}

type ListEmployeesParams struct {
	CompanyID int64
	Name      *string
}

type FindEmployeeParams struct {
	ID        int64
	CompanyID int64
}

type UpdateEmployeeParams struct {
	ID             int64
	FirstName      string
	LastName       string
	EmploymentType string
	SalaryType     string
	BaseSalary     int64
}

type EmployeeResponse struct {
	ID             int64     `json:"id"`
	CompanyID      int64     `json:"company_id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	EmploymentType string    `json:"employment_type"`
	SalaryType     string    `json:"salary_type"`
	BaseSalary     string    `json:"base_salary"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
