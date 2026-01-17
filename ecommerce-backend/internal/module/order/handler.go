package order

import (
	"encoding/json"
	"errors"
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
	r.Post("/checkout", h.checkout)
	r.Get("/orders/{id}", h.getOrder)
}

type checkoutReq struct {
	CartID  string          `json:"cart_id"`
	Address AddressSnapshot `json:"address"`
}

func (h *Handler) checkout(w http.ResponseWriter, r *http.Request) {
	var req checkoutReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_json"})
		return
	}
	if req.CartID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_payload"})
		return
	}

	orderID, err := h.svc.Checkout(r.Context(), req.CartID, req.Address)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "not_found"})
		case errors.Is(err, ErrEmptyCart):
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "empty_cart"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "internal_error"})
		}
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"order_id": orderID})
}

func (h *Handler) getOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	o, err := h.svc.GetOrder(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "not_found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "internal_error"})
		return
	}
	writeJSON(w, http.StatusOK, o)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
