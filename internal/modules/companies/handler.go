package companies

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/aeschyllus/sweldo-rest/internal/pkg/json"
	"github.com/go-chi/chi/v5"
)

func NewHandler(service Service) *handler {
	return &handler{service}
}

func (h *handler) RegisterRoutes(r chi.Router) {
	r.Route("/companies", func(r chi.Router) {
		r.Post("/", h.CreateCompany)
		r.Get("/", h.ListCompanies)
		r.Get("/{companyID}", h.FindCompanyByID)
		r.Put("/{companyID}", h.UpdateCompanyByID)
	})
}

func (h *handler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var req createCompanyRequest
	if err := json.Read(w, r, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		json.WriteError(w, http.StatusBadRequest, "name is required")
		return
	}
	if req.TaxID == "" {
		json.WriteError(w, http.StatusBadRequest, "tax_id is required")
		return
	}

	company, err := h.service.CreateCompany(r.Context(), CreateCompanyParams{
		Name:  req.Name,
		TaxID: req.TaxID,
	})
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to create company", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	json.Write(w, http.StatusCreated, company)
}

func (h *handler) ListCompanies(w http.ResponseWriter, r *http.Request) {
	params, err := parseListCompaniesQuery(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	companies, err := h.service.ListCompanies(r.Context(), params)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to list companies", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	json.Write(w, http.StatusOK, companies)
}

func (h *handler) FindCompanyByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "companyID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid company ID")
		return
	}

	company, err := h.service.FindCompanyByID(r.Context(), id)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to find company", "error", err)
		json.WriteError(w, http.StatusNotFound, "company not found")
		return
	}

	json.Write(w, http.StatusOK, company)
}

func (h *handler) UpdateCompanyByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "companyID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid company ID")
		return
	}

	var req updateCompanyRequest
	if err := json.Read(w, r, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		json.WriteError(w, http.StatusBadRequest, "name is required")
		return
	}
	if req.TaxID == "" {
		json.WriteError(w, http.StatusBadRequest, "tax_id is required")
		return
	}

	company, err := h.service.UpdateCompanyByID(r.Context(), UpdateCompanyParams{
		ID:    id,
		Name:  req.Name,
		TaxID: req.TaxID,
	})
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to update company", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	json.Write(w, http.StatusOK, company)
}
