package service

import (
	"context"
	"errors"
	"time"

	"github.com/deepak-2605/collab-kanban/backend/internal/models"
	"github.com/deepak-2605/collab-kanban/backend/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var ErrNotFound = errors.New("board not found")
var ErrForbidden = errors.New("not allowed")

type BoardService struct {
	repo *repository.BoardRepository
}

func NewBoardService(repo *repository.BoardRepository) *BoardService {
	return &BoardService{repo: repo}
}

func (s *BoardService) getOwned(ctx context.Context, boardID, userID bson.ObjectID) (*models.Board, error) {
	board, err := s.repo.FindByID(ctx, boardID)

	if err != nil {
		return nil, ErrNotFound
	}

	if board.OwnerID != userID {
		return nil, ErrForbidden
	}
	return board, nil
}

func (s *BoardService) Delete(ctx context.Context, boardID, userID bson.ObjectID) error {
	if _, err := s.getOwned(ctx, boardID, userID); err != nil {
		return err
	}
	return s.repo.Delete(ctx, boardID)
}

// Create - build a models.Board
func (s *BoardService) Create(ctx context.Context, name string, userID bson.ObjectID) (*models.Board, error) {

	board := &models.Board{
		Name:      name,
		OwnerID:   userID,
		Members:   []bson.ObjectID{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.repo.Create(ctx, board); err != nil {
		return nil, err
	}
	return board, nil
}

// List the repo
func (s *BoardService) List(ctx context.Context, userID bson.ObjectID) ([]models.Board, error) {
	boards,err:= s.repo.FindByUser(ctx, userID)

	if err!=nil {
		return nil,err
	}

	if boards == nil {
		boards = []models.Board{}
	}
	return boards,nil
}

// Rename the board
func (s *BoardService) Rename(ctx context.Context, boardID, userID bson.ObjectID, newName string) error {

	board, err := s.getOwned(ctx, boardID, userID)
	if err != nil {
		return err
	}
	board.Name = newName
	board.UpdatedAt = time.Now()
	return s.repo.Update(ctx, board)
}
