package payrolls

import (
	"net/http"
	"strconv"

	"github.com/aeschyllus/sweldo-rest/internal/pkg/json"
	"github.com/go-chi/chi/v5"
)

func NewHandler(service Service) *handler {
	return &handler{service}
}

func (h *handler) RegisterRoutes(r chi.Router) {
	r.Route("/payroll-runs", func(r chi.Router) {
		r.Post("/", h.CreatePayrollRun)
		r.Get("/", h.ListPayrollRuns)
		r.Get("/details", h.ListPayrollDetailsByEmployeeID)

		r.Route("/{runID}", func(r chi.Router) {
			r.Get("/", h.FindPayrollRunByID)
			r.Put("/", h.UpdatePayrollRunByID)
			r.Post("/details", h.CreatePayrollDetail)
			r.Get("/details", h.ListPayrollDetailsByRunID)
		})

	})
}

func (h *handler) CreatePayrollRun(w http.ResponseWriter, r *http.Request) {
	var req createPayrollRunRequest
	if err := json.Read(r, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.service.CreatePayrollRun(r.Context(), CreatePayrollRunParams{
		CompanyID:      req.CompanyID,
		RunDate:        req.RunDate,
		TotalEmployees: req.TotalEmployees,
		TotalPay:       req.TotalPay,
	})
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusCreated, response)
}

func (h *handler) ListPayrollRuns(w http.ResponseWriter, r *http.Request) {
	query, err := parseListPayrollRunsQuery(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid company_id")
		return
	}

	if query.CompanyID == 0 {
		json.WriteError(w, http.StatusBadRequest, "company_id is required")
		return
	}

	responses, err := h.service.ListPayrollRunsByCompanyID(r.Context(), query.CompanyID)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, responses)
}

func (h *handler) FindPayrollRunByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "runID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid run ID")
		return
	}

	response, err := h.service.FindPayrollRunByID(r.Context(), id)
	if err != nil {
		json.WriteError(w, http.StatusNotFound, "payroll run not found")
		return
	}

	json.Write(w, http.StatusOK, response)
}

func (h *handler) UpdatePayrollRunByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "runID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid run ID")
		return
	}

	var req updatePayrollRunRequest
	if err := json.Read(r, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.service.UpdatePayrollRunByID(r.Context(), id, UpdatePayrollRunParams{
		TotalEmployees: req.TotalEmployees,
		TotalPay:       req.TotalPay,
	})
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, response)
}

func (h *handler) CreatePayrollDetail(w http.ResponseWriter, r *http.Request) {
	runIDStr := chi.URLParam(r, "runID")
	runID, err := strconv.ParseInt(runIDStr, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid run ID")
		return
	}

	var req createPayrollDetailRequest
	if err := json.Read(r, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.service.CreatePayrollDetail(r.Context(), runID, CreatePayrollDetailParams{
		EmployeeID:   req.EmployeeID,
		GrossPay:     req.GrossPay,
		TaxDeduction: req.TaxDeduction,
		NetPay:       req.NetPay,
	})
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusCreated, response)
}

func (h *handler) ListPayrollDetailsByRunID(w http.ResponseWriter, r *http.Request) {
	runIDStr := chi.URLParam(r, "runID")
	runID, err := strconv.ParseInt(runIDStr, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid run ID")
		return
	}

	responses, err := h.service.ListAllPayrollDetailsByRunID(r.Context(), runID)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, responses)
}

func (h *handler) ListPayrollDetailsByEmployeeID(w http.ResponseWriter, r *http.Request) {
	query, err := parseListPayrollDetailsQuery(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid employee_id")
		return
	}

	if query.EmployeeID == 0 {
		json.WriteError(w, http.StatusBadRequest, "employee_id is required")
		return
	}

	responses, err := h.service.ListAllPayrollDetailsByEmployeeID(r.Context(), query.EmployeeID)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, responses)
}
