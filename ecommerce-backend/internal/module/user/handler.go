package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	httpx "github.com/synchhans/ecommerce-backend/internal/platform/http"
)

type Handler struct {
	svc       *Service
	jwtSecret []byte
}

func NewHandler(svc *Service, jwtSecret string) *Handler {
	return &Handler{svc: svc, jwtSecret: []byte(jwtSecret)}
}

func (h *Handler) Routes(r chi.Router) {
	r.Post("/auth/register", h.register)
	r.Post("/auth/login", h.login)

	r.Group(func(pr chi.Router) {
		pr.Use(httpx.AuthMiddleware(h.jwtSecret))
		pr.Get("/me", h.me)
	})
}

type registerReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Fail(w, http.StatusBadRequest, "invalid_json", "Invalid JSON payload")
		return
	}

	res, err := h.svc.Register(r.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidPayload):
			httpx.Fail(w, http.StatusBadRequest, "invalid_payload", "Invalid payload")
		case errors.Is(err, ErrEmailTaken):
			httpx.Fail(w, http.StatusConflict, "email_taken", "Email already taken")
		default:
			httpx.Fail(w, http.StatusInternalServerError, "internal_error", "Internal server error")
		}
		return
	}

	httpx.Created(w, res)
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Fail(w, http.StatusBadRequest, "invalid_json", "Invalid JSON payload")
		return
	}

	res, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidPayload) {
			httpx.Fail(w, http.StatusBadRequest, "invalid_payload", "Invalid payload")
			return
		}
		httpx.Fail(w, http.StatusUnauthorized, "invalid_credentials", "Invalid credentials")
		return
	}

	httpx.OK(w, res)
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	userID, ok := httpx.UserIDFromContext(r.Context())
	if !ok {
		httpx.Fail(w, http.StatusUnauthorized, "unauthorized", "Unauthorized")
		return
	}
	u, err := h.svc.Me(r.Context(), userID)
	if err != nil {
		httpx.Fail(w, http.StatusUnauthorized, "unauthorized", "Unauthorized")
		return
	}
	httpx.OK(w, u)
}


