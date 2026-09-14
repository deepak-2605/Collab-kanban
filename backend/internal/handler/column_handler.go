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

type ColumnHandler struct {
	service *service.ColumnService
}

func NewColumnHandler(s *service.ColumnService) *ColumnHandler {
	return &ColumnHandler{service: s}
}

// Create - it creates a newColumn POST:- /boards/{boardID}/column
func (h *ColumnHandler) Create(w http.ResponseWriter, r *http.Request) {

	boardID, err := bson.ObjectIDFromHex(chi.URLParam(r, "boardID"))

	if err != nil {
		http.Error(w, "invalid board id", http.StatusBadRequest)
		return
	}

	var req struct {
		Name string `json:"name"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
	}

	userID, err := bson.ObjectIDFromHex(middleware.UserIDFromContext(r.Context()))
	if err != nil {
		http.Error(w, "invalid user", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	column, err := h.service.Create(ctx, boardID, userID, req.Name)
	if err != nil {
		mapBoardError(w, err) // reuses 404/403/500 mapping
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(column)
}

// List :- list all columns for a board GET: /boards/{boardID}/columns
func (h *ColumnHandler) List(w http.ResponseWriter, r *http.Request) {

	boardID, err := bson.ObjectIDFromHex(chi.URLParam(r, "boardID"))

	if err != nil {
		http.Error(w, "invalid board id", http.StatusBadRequest)
		return
	}

	userID, err := bson.ObjectIDFromHex(middleware.UserIDFromContext(r.Context()))
	if err != nil {
		http.Error(w, "invalid user", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	cols, err := h.service.List(ctx, boardID, userID)
	if err != nil {
		mapBoardError(w, err) // reuses 404/403/500 mapping
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cols)
}

// Rename :- rename a col PATCH /columns/{columnID}
func (h *ColumnHandler) Rename(w http.ResponseWriter, r *http.Request) {

	columnID, err := bson.ObjectIDFromHex(chi.URLParam(r, "columnID"))

	if err != nil {
		http.Error(w, "invalid column id", http.StatusBadRequest)
		return
	}

	var req struct {
		Name string `json:"name"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	userID, err := bson.ObjectIDFromHex(middleware.UserIDFromContext(r.Context()))
	if err != nil {
		http.Error(w, "invalid user", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err = h.service.Rename(ctx, columnID, userID, req.Name)

	if err != nil {
		mapBoardError(w, err) // reuses 404/403/500 mapping
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// Delete - delete the column
// DELETE /columns/{columnID}
func (h *ColumnHandler) Delete(w http.ResponseWriter, r *http.Request) {

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

	err = h.service.Delete(ctx, columnID, userID)

	if err != nil {
		mapBoardError(w, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
