package payment

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Routes(r chi.Router) {
	r.Post("/payments/initiate", h.initiate)
	r.Post("/payments/webhook/{provider}", h.webhook)
}

type initiateReq struct {
	OrderID  string `json:"order_id"`
	Provider string `json:"provider"`
}

func (h *Handler) initiate(w http.ResponseWriter, r *http.Request) {
	var req initiateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_json"})
		return
	}
	if req.OrderID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_payload"})
		return
	}

	res, err := h.svc.Initiate(r.Context(), req.OrderID, req.Provider)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "not_found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "internal_error"})
		return
	}

	writeJSON(w, http.StatusCreated, res)
}

type webhookReq struct {
	ProviderRef string `json:"provider_ref"`
	Status      string `json:"status"` // pending/paid/failed/expired/refunded
}

func (h *Handler) webhook(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	raw, _ := io.ReadAll(r.Body)

	var req webhookReq
	if err := json.Unmarshal(raw, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_json"})
		return
	}
	if req.ProviderRef == "" || req.Status == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_payload"})
		return
	}

	res, err := h.svc.Webhook(r.Context(), provider, req.ProviderRef, req.Status, raw)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidStatus):
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_status"})
		case errors.Is(err, ErrNotFound):
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "not_found"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "internal_error"})
		}
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
