package employees

import (
	"fmt"
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
	r.Route("/employees", func(r chi.Router) {
		r.Post("/", h.CreateEmployee)
		r.Get("/", h.ListEmployeesByCompanyID)
		r.Get("/{employeeID}", h.FindEmployeeByID)
		r.Put("/{employeeID}", h.UpdateEmployeeByID)
	})
}

func (h *handler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var req createEmployeeRequest
	if err := json.Read(w, r, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.FirstName == "" || req.LastName == "" {
		json.WriteError(w, http.StatusBadRequest, "first_name and last_name are required")
		return
	}
	if req.EmploymentType == "" || !EmploymentType(req.EmploymentType).Valid() {
		json.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid employment_type: %s", req.EmploymentType))
		return
	}
	if req.SalaryType == "" || !SalaryType(req.SalaryType).Valid() {
		json.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid salary_type: %s", req.SalaryType))
		return
	}

	baseSalary, err := money.ParseCents(req.BaseSalary)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid base_salary: %s", req.BaseSalary))
		return
	}

	companyIDStr := r.URL.Query().Get("company_id")
	companyID, err := strconv.ParseInt(companyIDStr, 10, 64)
	if err != nil || companyID == 0 {
		json.WriteError(w, http.StatusBadRequest, "company_id is required")
		return
	}

	userID := auth.UserIDFromContext(r.Context())

	employee, err := h.service.CreateEmployee(r.Context(), CreateEmployeeParams{
		CompanyID:      companyID,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		EmploymentType: req.EmploymentType,
		SalaryType:     req.SalaryType,
		BaseSalary:     baseSalary,
		CreatedBy:      userID,
	})
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to create employee", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	json.Write(w, http.StatusCreated, employee)
}

func (h *handler) ListEmployeesByCompanyID(w http.ResponseWriter, r *http.Request) {
	params, err := parseListEmployeesQuery(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	employees, err := h.service.ListEmployeesByCompanyID(r.Context(), params)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to list employees", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if employees == nil {
		employees = []EmployeeResponse{}
	}

	json.Write(w, http.StatusOK, employees)
}

func (h *handler) FindEmployeeByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "employeeID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid employee ID")
		return
	}

	companyIDStr := r.URL.Query().Get("company_id")
	companyID, err := strconv.ParseInt(companyIDStr, 10, 64)
	if err != nil || companyID == 0 {
		json.WriteError(w, http.StatusBadRequest, "company_id is required")
		return
	}

	employee, err := h.service.FindEmployeeByID(r.Context(), FindEmployeeParams{
		ID:        id,
		CompanyID: companyID,
	})
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to find employee", "error", err)
		json.WriteError(w, http.StatusNotFound, "employee not found")
		return
	}

	json.Write(w, http.StatusOK, employee)
}

func (h *handler) UpdateEmployeeByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "employeeID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid employee ID")
		return
	}

	var req updateEmployeeRequest
	if err := json.Read(w, r, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.FirstName == "" || req.LastName == "" {
		json.WriteError(w, http.StatusBadRequest, "first_name and last_name are required")
		return
	}
	if req.EmploymentType == "" || !EmploymentType(req.EmploymentType).Valid() {
		json.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid employment_type: %s", req.EmploymentType))
		return
	}
	if req.SalaryType == "" || !SalaryType(req.SalaryType).Valid() {
		json.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid salary_type: %s", req.SalaryType))
		return
	}

	baseSalary, err := money.ParseCents(req.BaseSalary)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid base_salary: %s", req.BaseSalary))
		return
	}

	companyIDStr := r.URL.Query().Get("company_id")
	companyID, err := strconv.ParseInt(companyIDStr, 10, 64)
	if err != nil || companyID == 0 {
		json.WriteError(w, http.StatusBadRequest, "company_id is required")
		return
	}

	userID := auth.UserIDFromContext(r.Context())

	employee, err := h.service.UpdateEmployeeByID(r.Context(), UpdateEmployeeParams{
		ID:             id,
		CompanyID:      companyID,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		EmploymentType: req.EmploymentType,
		SalaryType:     req.SalaryType,
		BaseSalary:     baseSalary,
		UpdatedBy:      userID,
	})
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to update employee", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	json.Write(w, http.StatusOK, employee)
}
