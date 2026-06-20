package payrolls

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/aeschyllus/sweldo-rest/internal/pkg/auth"
	"github.com/aeschyllus/sweldo-rest/internal/pkg/json"
	"github.com/aeschyllus/sweldo-rest/internal/pkg/money"
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
			r.Patch("/finalize", h.FinalizePayrollRun)
			r.Post("/details", h.CreatePayrollDetail)
			r.Get("/details", h.ListPayrollDetailsByRunID)
		})
	})

	r.Route("/payroll-details/{detailID}/deductions", func(r chi.Router) {
		r.Post("/", h.CreateDeduction)
		r.Get("/", h.ListDeductions)
		r.Delete("/{deductionID}", h.DeleteDeduction)
	})
}

func (h *handler) CreatePayrollRun(w http.ResponseWriter, r *http.Request) {
	var req createPayrollRunRequest
	if err := json.Read(w, r, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.RunDate == "" {
		json.WriteError(w, http.StatusBadRequest, "run_date is required")
		return
	}
	if req.TotalPay == "" {
		json.WriteError(w, http.StatusBadRequest, "total_pay is required")
		return
	}

	totalPay, err := money.ParseCents(req.TotalPay)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid total_pay")
		return
	}

	companyIDStr := r.URL.Query().Get("company_id")
	companyID, err := strconv.ParseInt(companyIDStr, 10, 64)
	if err != nil || companyID == 0 {
		json.WriteError(w, http.StatusBadRequest, "company_id is required")
		return
	}

	userID := auth.UserIDFromContext(r.Context())

	response, err := h.service.CreatePayrollRun(r.Context(), CreatePayrollRunParams{
		CompanyID:      companyID,
		RunDate:        req.RunDate,
		TotalEmployees: req.TotalEmployees,
		TotalPay:       totalPay,
		CreatedBy:      userID,
	})
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to create payroll run", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	json.Write(w, http.StatusCreated, response)
}

func (h *handler) ListPayrollRuns(w http.ResponseWriter, r *http.Request) {
	query, err := parseListPayrollRunsQuery(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	responses, err := h.service.ListPayrollRunsByCompanyID(r.Context(), query.CompanyID)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to list payroll runs", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "internal server error")
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
	if err := json.Read(w, r, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	totalPay, err := money.ParseCents(req.TotalPay)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid total_pay")
		return
	}

	userID := auth.UserIDFromContext(r.Context())

	response, err := h.service.UpdatePayrollRunByID(r.Context(), id, UpdatePayrollRunParams{
		TotalEmployees: req.TotalEmployees,
		TotalPay:       totalPay,
		UpdatedBy:      userID,
	})
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to update payroll run", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	json.Write(w, http.StatusOK, response)
}

func (h *handler) FinalizePayrollRun(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "runID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid run ID")
		return
	}

	userID := auth.UserIDFromContext(r.Context())

	response, err := h.service.FinalizePayrollRun(r.Context(), id, userID)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to finalize payroll run", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "internal server error")
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
	if err := json.Read(w, r, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.EmployeeID <= 0 {
		json.WriteError(w, http.StatusBadRequest, "employee_id is required")
		return
	}

	grossPay, err := money.ParseCents(req.GrossPay)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid gross_pay")
		return
	}

	taxDeduction, err := money.ParseCents(req.TaxDeduction)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid tax_deduction")
		return
	}

	var hourlyRate, hoursWorked int64
	if req.HourlyRate != nil {
		hourlyRate, err = money.ParseCents(*req.HourlyRate)
		if err != nil {
			json.WriteError(w, http.StatusBadRequest, "invalid hourly_rate")
			return
		}
	}
	if req.HoursWorked != nil {
		hoursWorked, err = money.ParseCents(*req.HoursWorked)
		if err != nil {
			json.WriteError(w, http.StatusBadRequest, "invalid hours_worked")
			return
		}
	}

	userID := auth.UserIDFromContext(r.Context())

	response, err := h.service.CreatePayrollDetail(r.Context(), runID, CreatePayrollDetailParams{
		EmployeeID:   req.EmployeeID,
		GrossPay:     grossPay,
		TaxDeduction: taxDeduction,
		HourlyRate:   hourlyRate,
		HoursWorked:  hoursWorked,
		CreatedBy:    userID,
	})
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to create payroll detail", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "internal server error")
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
		slog.ErrorContext(r.Context(), "failed to list payroll details by run", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	json.Write(w, http.StatusOK, responses)
}

func (h *handler) ListPayrollDetailsByEmployeeID(w http.ResponseWriter, r *http.Request) {
	query, err := parseListPayrollDetailsQuery(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	responses, err := h.service.ListAllPayrollDetailsByEmployeeID(r.Context(), query.EmployeeID)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to list payroll details by employee", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	json.Write(w, http.StatusOK, responses)
}

func (h *handler) CreateDeduction(w http.ResponseWriter, r *http.Request) {
	detailIDStr := chi.URLParam(r, "detailID")
	detailID, err := strconv.ParseInt(detailIDStr, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid detail ID")
		return
	}

	var req createDeductionRequest
	if err := json.Read(w, r, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.DeductionType == "" {
		json.WriteError(w, http.StatusBadRequest, "deduction_type is required")
		return
	}

	amount, err := money.ParseCents(req.Amount)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid amount")
		return
	}

	userID := auth.UserIDFromContext(r.Context())

	response, err := h.service.CreateDeduction(r.Context(), CreateDeductionParams{
		PayrollDetailID: detailID,
		DeductionType:   req.DeductionType,
		Amount:          amount,
		CreatedBy:       userID,
	})
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to create deduction", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	json.Write(w, http.StatusCreated, response)
}

func (h *handler) ListDeductions(w http.ResponseWriter, r *http.Request) {
	detailIDStr := chi.URLParam(r, "detailID")
	detailID, err := strconv.ParseInt(detailIDStr, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid detail ID")
		return
	}

	deductions, err := h.service.ListDeductionsByDetailID(r.Context(), detailID)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to list deductions", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	json.Write(w, http.StatusOK, deductions)
}

func (h *handler) DeleteDeduction(w http.ResponseWriter, r *http.Request) {
	deductionIDStr := chi.URLParam(r, "deductionID")
	deductionID, err := strconv.ParseInt(deductionIDStr, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid deduction ID")
		return
	}

	if err := h.service.DeleteDeduction(r.Context(), deductionID); err != nil {
		slog.ErrorContext(r.Context(), "failed to delete deduction", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
