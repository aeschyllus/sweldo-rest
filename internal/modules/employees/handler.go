package employees

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
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	employee, err := h.service.CreateEmployee(r.Context(), CreateEmployeeParams{
		CompanyID:      req.CompanyID,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		EmploymentType: req.EmploymentType,
		SalaryType:     req.SalaryType,
		BaseSalary:     req.BaseSalary,
	})
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to create employee", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	json.Write(w, http.StatusCreated, employee)
}

// Supports searching by first/last name via query parameters.
//
// NOTE: It is mandatory to pass in the company_id query parameter.
// We can inject the company_id in the future via middleware that extracts the company_id from the JWT
//
// TODO: add pagination support
//
// e.g.:
//   - /employees?company_id=1&name=juan
func (h *handler) ListEmployeesByCompanyID(w http.ResponseWriter, r *http.Request) {
	params, err := parseListEmployeesQuery(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	employee, err := h.service.ListEmployeesByCompanyID(r.Context(), params)
	if err != nil {
		json.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	json.Write(w, http.StatusOK, employee)
}

// TODO: refactor to prevent companies from searching employees of other companies
func (h *handler) FindEmployeeByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "employeeID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req findEmployeeRequest
	if err := json.Read(w, r, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	employee, err := h.service.FindEmployeeByID(r.Context(), FindEmployeeParams{
		ID:        id,
		CompanyID: req.CompanyID,
	})
	if err != nil {
		json.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	json.Write(w, http.StatusOK, employee)
}

func (h *handler) UpdateEmployeeByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "employeeID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req updateEmployeeRequest
	if err := json.Read(w, r, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	employee, err := h.service.UpdateEmployeeByID(r.Context(), UpdateEmployeeParams{
		ID:             id,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		EmploymentType: req.EmploymentType,
		SalaryType:     req.SalaryType,
		BaseSalary:     req.BaseSalary,
	})
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to update employee", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	json.Write(w, http.StatusOK, employee)
}
