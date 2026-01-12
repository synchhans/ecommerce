package cart

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	httpx "github.com/synchhans/ecommerce-backend/internal/platform/http"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Routes(r chi.Router) {
	r.Post("/cart", h.createCart)
	r.Get("/cart/{id}", h.getCart)

	r.Post("/cart/{id}/items", h.upsertItem)
	r.Patch("/cart/{id}/items/{itemId}", h.updateItemQty)
	r.Delete("/cart/{id}/items/{itemId}", h.deleteItem)
}

func (h *Handler) createCart(w http.ResponseWriter, r *http.Request) {
	id, err := h.svc.CreateCart(r.Context())
	if err != nil {
		httpx.Fail(w, http.StatusInternalServerError, "internal_error")
		return
	}
	httpx.Created(w, httpx.Envelope{"id": id})
}

func (h *Handler) getCart(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	c, err := h.svc.GetCart(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			httpx.Fail(w, http.StatusNotFound, "not_found")
			return
		}
		httpx.Fail(w, http.StatusInternalServerError, "internal_error")
		return
	}
	httpx.OK(w, c)
}

type upsertItemReq struct {
	VariantID string `json:"variant_id"`
	Qty       int    `json:"qty"`
}

func (h *Handler) upsertItem(w http.ResponseWriter, r *http.Request) {
	cartID := chi.URLParam(r, "id")

	var req upsertItemReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Fail(w, http.StatusBadRequest, "invalid_json")
		return
	}
	if req.VariantID == "" || req.Qty <= 0 {
		httpx.Fail(w, http.StatusBadRequest, "invalid_payload")
		return
	}

	if err := h.svc.AddOrReplaceItem(r.Context(), cartID, req.VariantID, req.Qty); err != nil {
		if errors.Is(err, ErrInvalidQty) {
			httpx.Fail(w, http.StatusBadRequest, "invalid_qty")
			return
		}
		httpx.Fail(w, http.StatusInternalServerError, "internal_error")
		return
	}

	httpx.OK(w, httpx.Envelope{"ok": true})
}

type updateQtyReq struct {
	Qty int `json:"qty"`
}

func (h *Handler) updateItemQty(w http.ResponseWriter, r *http.Request) {
	cartID := chi.URLParam(r, "id")
	itemID := chi.URLParam(r, "itemId")

	var req updateQtyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Fail(w, http.StatusBadRequest, "invalid_json")
		return
	}
	if req.Qty <= 0 {
		httpx.Fail(w, http.StatusBadRequest, "invalid_payload")
		return
	}

	if err := h.svc.UpdateItemQty(r.Context(), cartID, itemID, req.Qty); err != nil {
		if errors.Is(err, ErrNotFound) {
			httpx.Fail(w, http.StatusNotFound, "not_found")
			return
		}
		if errors.Is(err, ErrInvalidQty) {
			httpx.Fail(w, http.StatusBadRequest, "invalid_qty")
			return
		}
		httpx.Fail(w, http.StatusInternalServerError, "internal_error")
		return
	}

	httpx.OK(w, httpx.Envelope{"ok": true})
}

func (h *Handler) deleteItem(w http.ResponseWriter, r *http.Request) {
	cartID := chi.URLParam(r, "id")
	itemID := chi.URLParam(r, "itemId")

	if err := h.svc.RemoveItem(r.Context(), cartID, itemID); err != nil {
		if errors.Is(err, ErrNotFound) {
			httpx.Fail(w, http.StatusNotFound, "not_found")
			return
		}
		httpx.Fail(w, http.StatusInternalServerError, "internal_error")
		return
	}
	httpx.OK(w, httpx.Envelope{"ok": true})
}
