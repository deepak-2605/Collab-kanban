package service

import (
	"context"
	"time"

	"github.com/deepak-2605/collab-kanban/backend/internal/models"
	"github.com/deepak-2605/collab-kanban/backend/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CardService struct {
	repo       *repository.CardRepository
	columnRepo *repository.ColumnRepository
	boardRepo  *repository.BoardRepository
}

func NewCardService(repo *repository.CardRepository, columnRepo *repository.ColumnRepository, boardRepo *repository.BoardRepository) *CardService {
	return &CardService{repo: repo, columnRepo: columnRepo, boardRepo: boardRepo}
}

// assertBoardAccess — owner-or-member rule (mirrors ColumnService).
func (s *CardService) assertBoardAccess(ctx context.Context, boardID, userID bson.ObjectID) error {
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

// Create — a card goes into a column; position = current count in that column.
func (s *CardService) Create(ctx context.Context, columnID, userID bson.ObjectID, title, description string) (*models.Card, error) {
	col, err := s.columnRepo.FindByID(ctx, columnID)
	if err != nil {
		return nil, ErrNotFound
	}
	if err := s.assertBoardAccess(ctx, col.BoardID, userID); err != nil {
		return nil, err
	}

	existing, err := s.repo.FindByColumn(ctx, columnID)
	if err != nil {
		return nil, err
	}

	card := &models.Card{
		BoardID:     col.BoardID, // denormalized from the column
		ColumnID:    columnID,
		Title:       title,
		Description: description,
		Position:    len(existing),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.repo.Create(ctx, card); err != nil {
		return nil, err
	}
	return card, nil
}

// ListByColumn — all cards in a column (position order), if the user can access the board.
func (s *CardService) ListByColumn(ctx context.Context, columnID, userID bson.ObjectID) ([]models.Card, error) {
	col, err := s.columnRepo.FindByID(ctx, columnID)
	if err != nil {
		return nil, ErrNotFound
	}
	if err := s.assertBoardAccess(ctx, col.BoardID, userID); err != nil {
		return nil, err
	}
	return s.repo.FindByColumn(ctx, columnID)
}

// Update — edit a card's title/description.
func (s *CardService) Update(ctx context.Context, cardID, userID bson.ObjectID, title, description string) error {
	card, err := s.repo.FindByID(ctx, cardID)
	if err != nil {
		return ErrNotFound
	}
	if err := s.assertBoardAccess(ctx, card.BoardID, userID); err != nil {
		return err
	}
	card.Title = title
	card.Description = description
	card.UpdatedAt = time.Now()
	return s.repo.Update(ctx, card)
}

// Move — the heart of Kanban: move a card to another column + position.
func (s *CardService) Move(ctx context.Context, cardID, targetColumnID bson.ObjectID, newPosition int, userID bson.ObjectID) error {
	card, err := s.repo.FindByID(ctx, cardID)
	if err != nil {
		return ErrNotFound
	}
	if err := s.assertBoardAccess(ctx, card.BoardID, userID); err != nil {
		return err
	}

	// the target column must exist AND be on the same board (no cross-board moves)
	targetCol, err := s.columnRepo.FindByID(ctx, targetColumnID)
	if err != nil {
		return ErrNotFound
	}
	if targetCol.BoardID != card.BoardID {
		return ErrForbidden
	}

	card.ColumnID = targetColumnID
	card.Position = newPosition
	card.UpdatedAt = time.Now()
	return s.repo.Update(ctx, card)
}

// Delete — remove a card.
func (s *CardService) Delete(ctx context.Context, cardID, userID bson.ObjectID) error {
	card, err := s.repo.FindByID(ctx, cardID)
	if err != nil {
		return ErrNotFound
	}
	if err := s.assertBoardAccess(ctx, card.BoardID, userID); err != nil {
		return err
	}
	return s.repo.Delete(ctx, cardID)
}
