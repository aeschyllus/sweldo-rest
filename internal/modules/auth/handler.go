package auth

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aeschyllus/sweldo-rest/internal/pkg/json"
	"github.com/go-chi/chi/v5"
)

func NewHandler(service Service) *handler {
	return &handler{service}
}

func (h *handler) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
	})
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.Read(w, r, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" {
		json.WriteError(w, http.StatusBadRequest, "email is required")
		return
	}
	if req.Password == "" {
		json.WriteError(w, http.StatusBadRequest, "password is required")
		return
	}
	if req.CompanyID == 0 {
		json.WriteError(w, http.StatusBadRequest, "company_id is required")
		return
	}

	response, err := h.service.Register(r.Context(), RegisterParams{
		CompanyID: req.CompanyID,
		Email:     req.Email,
		Password:  req.Password,
	})
	if err != nil {
		slog.ErrorContext(r.Context(), "registration failed", "error", err)
		json.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("registration failed: %s", err.Error()))
		return
	}

	json.Write(w, http.StatusCreated, response)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.Read(w, r, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" {
		json.WriteError(w, http.StatusBadRequest, "email is required")
		return
	}
	if req.Password == "" {
		json.WriteError(w, http.StatusBadRequest, "password is required")
		return
	}

	response, err := h.service.Login(r.Context(), LoginParams{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		json.WriteError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	json.Write(w, http.StatusOK, response)
}
