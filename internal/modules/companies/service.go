package companies

import (
	"context"

	"github.com/aeschyllus/sweldo-rest/internal/adapters/postgresql/sqlc"
	"github.com/aeschyllus/sweldo-rest/internal/pkg/pgconvert"
)

func NewService(repo sqlc.Querier) Service {
	return &service{repo}
}

func (s *service) CreateCompany(ctx context.Context, params CreateCompanyParams) (sqlc.Company, error) {
	return s.repo.CreateCompany(ctx, sqlc.CreateCompanyParams{
		Name:  params.Name,
		TaxID: params.TaxID,
	})
}

func (s *service) ListCompanies(ctx context.Context, params ListCompaniesParams) ([]sqlc.Company, error) {

	return s.repo.ListCompanies(ctx, sqlc.ListCompaniesParams{
		Name:       pgconvert.ToText(params.Name),
		PageLimit:  params.PageLimit,
		PageOffset: params.PageOffset,
	})
}

func (s *service) FindCompanyByID(ctx context.Context, id int64) (sqlc.Company, error) {
	return s.repo.FindCompanyByID(ctx, id)
}

func (s *service) UpdateCompanyByID(ctx context.Context, params UpdateCompanyParams) (sqlc.Company, error) {
	return s.repo.UpdateCompanyByID(ctx, sqlc.UpdateCompanyByIDParams{
		ID:    params.ID,
		Name:  params.Name,
		TaxID: params.TaxID,
	})
}
