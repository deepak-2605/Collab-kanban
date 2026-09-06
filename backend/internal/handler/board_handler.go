package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/deepak-2605/collab-kanban/backend/internal/middleware"
	"github.com/deepak-2605/collab-kanban/backend/internal/service"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type BoardHandler struct {
	service *service.BoardService
}

func NewBoardHandler(s *service.BoardService) *BoardHandler {
	return &BoardHandler{service: s}
}

type BoardRequest struct {
	Name string `json:"name"`
}

func mapBoardError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrNotFound):
		http.Error(w, "board not found", http.StatusNotFound)
	case errors.Is(err, service.ErrForbidden):
		http.Error(w, "not allowed", http.StatusForbidden)
	default:
		http.Error(w, "server error", http.StatusInternalServerError)
	}
}

func (h *BoardHandler) Create(w http.ResponseWriter, r *http.Request) {

	var req BoardRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	userId, err := bson.ObjectIDFromHex(middleware.UserIDFromContext(r.Context()))

	if err != nil {
		http.Error(w, "invalid user", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	board, err := h.service.Create(ctx, req.Name, userId)

	if err != nil {
		http.Error(w, "could not create board", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(board)
}

// Delete - boardID comes from the URL path
func (h *BoardHandler) Delete(w http.ResponseWriter, r *http.Request) {

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

	if err := h.service.Delete(ctx, boardID, userID); err != nil {
		mapBoardError(w, err) // 404 / 403 / 500
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *BoardHandler) List(w http.ResponseWriter, r *http.Request) {

	userID, err := bson.ObjectIDFromHex(middleware.UserIDFromContext(r.Context()))

	if err != nil {
		http.Error(w, "invalid user", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	boards, err := h.service.List(ctx, userID)

	if err != nil {
		mapBoardError(w, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(boards)

}

func (h *BoardHandler) Rename(w http.ResponseWriter, r *http.Request) {

	var req BoardRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

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

	err = h.service.Rename(ctx, boardID, userID, req.Name)

	if err != nil {
		mapBoardError(w, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusNoContent)

}
