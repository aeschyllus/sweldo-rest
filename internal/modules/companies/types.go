package companies

import (
	"context"

	"github.com/aeschyllus/sweldo-rest/internal/adapters/postgresql/sqlc"
)

type handler struct {
	service Service
}

type service struct {
	repo sqlc.Querier
}

type Service interface {
	CreateCompany(ctx context.Context, params CreateCompanyParams) (sqlc.Company, error)
	ListCompanies(ctx context.Context, params ListCompaniesParams) ([]sqlc.Company, error)
	FindCompanyByID(ctx context.Context, id int64) (sqlc.Company, error)
	UpdateCompanyByID(ctx context.Context, params UpdateCompanyParams) (sqlc.Company, error)
}

type createCompanyRequest struct {
	Name  string `json:"name"`
	TaxID string `json:"tax_id"`
}

type updateCompanyRequest struct {
	Name  string `json:"name"`
	TaxID string `json:"tax_id"`
}

type CreateCompanyParams struct {
	Name  string
	TaxID string
}

type ListCompaniesParams struct {
	Name       *string
	PageLimit  int32
	PageOffset int32
}

type UpdateCompanyParams struct {
	ID    int64
	Name  string
	TaxID string
}
