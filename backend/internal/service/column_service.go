package service

import (
	"context"
	"time"

	"github.com/deepak-2605/collab-kanban/backend/internal/models"
	"github.com/deepak-2605/collab-kanban/backend/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ColumnService struct {
	repo      *repository.ColumnRepository
	boardRepo *repository.BoardRepository
}

func NewColumnService(repo *repository.ColumnRepository, boardRepo *repository.BoardRepository) *ColumnService {
	return &ColumnService{repo: repo, boardRepo: boardRepo}
}

// assertBoardAccess — owner OR member may manage this board's columns.
// (Reuses ErrNotFound / ErrForbidden already defined in board_service.go — same package.)
func (s *ColumnService) assertBoardAccess(ctx context.Context, boardID, userID bson.ObjectID) error {
	board, err := s.boardRepo.FindByID(ctx, boardID)
	if err != nil {
		return ErrNotFound
	}
	if board.OwnerID == userID {
		return nil
	}
	for _, m := range board.Members {
		if m == userID {
			return nil
		}
	}
	return ErrForbidden
}

// Create - the column inside board
func (s *ColumnService) Create(ctx context.Context, boardID, userID bson.ObjectID, name string) (*models.Column, error) {
	if err := s.assertBoardAccess(ctx, boardID, userID); err != nil {
		return nil, err
	}

	existing, err := s.repo.FindByBoard(ctx, boardID)
	if err != nil {
		return nil, err
	}

	column := &models.Column{
		BoardID:   boardID,
		Name:      name,
		Position:  len(existing), // append to the end
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.repo.Create(ctx, column); err != nil {
		return nil, err
	}
	return column, nil
}

// List - get all columns for a particular board
func (s *ColumnService) List(ctx context.Context, boardID, userID bson.ObjectID) ([]models.Column, error) {

	if err := s.assertBoardAccess(ctx, boardID, userID); err != nil {
		return nil, err
	}

	return s.repo.FindByBoard(ctx, boardID)
}

// Rename - rename the Column
func (s *ColumnService) Rename(ctx context.Context, columnID, userID bson.ObjectID, newName string) error {

	col, err := s.repo.FindByID(ctx, columnID)

	if err != nil {
		return ErrNotFound
	}

	if err = s.assertBoardAccess(ctx, col.BoardID, userID); err != nil {
		return err
	}

	col.Name = newName
	col.UpdatedAt = time.Now()

	return s.repo.Update(ctx, col)
}

// Delete - delete the column
func (s *ColumnService) Delete(ctx context.Context, columnID, userID bson.ObjectID) error {
	col, err := s.repo.FindByID(ctx, columnID)

	if err != nil {
		return ErrNotFound
	}

	if err = s.assertBoardAccess(ctx, col.BoardID, userID); err != nil {
		return err
	}

	return s.repo.Delete(ctx, columnID)
}
