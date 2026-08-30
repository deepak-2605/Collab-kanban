package service

import (
	"context"
	"errors"
	"time"

	"github.com/deepak-2605/collab-kanban/backend/internal/models"
	"github.com/deepak-2605/collab-kanban/backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      *repository.UserRepository
	jwtSecret string
}

func NewAuthService(repo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret}
}

// ErrInvalidCredentials is a business error the handler will map to a 401.
var ErrInvalidCredentials = errors.New("invalid credentials")

func (s *AuthService) Register(ctx context.Context, name, email, password string) (*models.User, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:      name,
		Email:     email,
		Password:  string(hash),
		CreatedAt: time.Now(),
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// generateToken builds a signed JWT for a user.
func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID.Hex(),
		"email": user.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return "", ErrInvalidCredentials
	}
	return s.generateToken(user)
}

// GetProfile fetches a user by their id string (from the JWT).
func (s *AuthService) GetProfile(ctx context.Context, userID string) (*models.User, error) {
	objID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, objID)
}
