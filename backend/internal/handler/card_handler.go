package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/deepak-2605/collab-kanban/backend/internal/middleware"
	"github.com/deepak-2605/collab-kanban/backend/internal/service"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CardHandler struct {
	service *service.CardService
}

func NewCardHandler(s *service.CardService) *CardHandler {
	return &CardHandler{service: s}
}

// Create — POST /columns/{columnID}/cards
func (h *CardHandler) Create(w http.ResponseWriter, r *http.Request) {
	columnID, err := bson.ObjectIDFromHex(chi.URLParam(r, "columnID"))
	if err != nil {
		http.Error(w, "invalid column id", http.StatusBadRequest)
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	userID, err := bson.ObjectIDFromHex(middleware.UserIDFromContext(r.Context()))
	if err != nil {
		http.Error(w, "invalid user", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	card, err := h.service.Create(ctx, columnID, userID, req.Title, req.Description)
	if err != nil {
		mapBoardError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(card)
}

// List — GET /columns/{columnID}/cards
func (h *CardHandler) List(w http.ResponseWriter, r *http.Request) {
	columnID, err := bson.ObjectIDFromHex(chi.URLParam(r, "columnID"))
	if err != nil {
		http.Error(w, "invalid column id", http.StatusBadRequest)
		return
	}

	userID, err := bson.ObjectIDFromHex(middleware.UserIDFromContext(r.Context()))
	if err != nil {
		http.Error(w, "invalid user", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	cards, err := h.service.ListByColumn(ctx, columnID, userID)
	if err != nil {
		mapBoardError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cards)
}

// Update — PATCH /cards/{cardID}  (edit title/description)
func (h *CardHandler) Update(w http.ResponseWriter, r *http.Request) {
	cardID, err := bson.ObjectIDFromHex(chi.URLParam(r, "cardID"))
	if err != nil {
		http.Error(w, "invalid card id", http.StatusBadRequest)
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	userID, err := bson.ObjectIDFromHex(middleware.UserIDFromContext(r.Context()))
	if err != nil {
		http.Error(w, "invalid user", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if err := h.service.Update(ctx, cardID, userID, req.Title, req.Description); err != nil {
		mapBoardError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Move — PATCH /cards/{cardID}/move   ← YOU WRITE THIS ONE
//
// Body shape:  {"columnId": "<target column id>", "position": <int>}
//
// Steps:
//  1. cardID  := ObjectIDFromHex(chi.URLParam(r, "cardID"))          -> 400 on error
//  2. decode a request struct:
//     var req struct {
//     ColumnID string `json:"columnId"`
//     Position int    `json:"position"`
//     }
//     -> 400 if decode fails
//  3. targetColumnID := ObjectIDFromHex(req.ColumnID)                -> 400 on error
//  4. userID := ObjectIDFromHex(middleware.UserIDFromContext(...))   -> 401 on error
//  5. ctx with timeout (+ defer cancel)
//  6. err := h.service.Move(ctx, cardID, targetColumnID, req.Position, userID)
//     if err != nil -> mapBoardError(w, err); return
//  7. w.WriteHeader(http.StatusNoContent)
func (h *CardHandler) Move(w http.ResponseWriter, r *http.Request) {
	// TODO: implement per the steps above
	cardID, err := bson.ObjectIDFromHex(chi.URLParam(r, "cardID"))

	if err != nil {
		http.Error(w, "invalid card id", http.StatusBadRequest)
		return
	}
	var req struct {
		ColumnID string `json:"columnId"`
		Position int    `json:"position"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.ColumnID == "" || req.Position < 0 {
		http.Error(w, "columnId and a valid position are required", http.StatusBadRequest)
		return
	}

	targetColumnID, err := bson.ObjectIDFromHex(req.ColumnID)

	if err != nil {
		http.Error(w, "invalid card id", http.StatusBadRequest)
		return
	}

	userID, err := bson.ObjectIDFromHex(middleware.UserIDFromContext(r.Context()))
	if err != nil {
		http.Error(w, "invalid user", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if err := h.service.Move(ctx, cardID, targetColumnID, req.Position, userID); err != nil {
		mapBoardError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

// Delete — DELETE /cards/{cardID}
func (h *CardHandler) Delete(w http.ResponseWriter, r *http.Request) {
	cardID, err := bson.ObjectIDFromHex(chi.URLParam(r, "cardID"))
	if err != nil {
		http.Error(w, "invalid card id", http.StatusBadRequest)
		return
	}

	userID, err := bson.ObjectIDFromHex(middleware.UserIDFromContext(r.Context()))
	if err != nil {
		http.Error(w, "invalid user", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if err := h.service.Delete(ctx, cardID, userID); err != nil {
		mapBoardError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
