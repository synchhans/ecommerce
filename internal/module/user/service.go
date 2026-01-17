package user

import (
	"context"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	httpx "github.com/synchhans/ecommerce-backend/internal/platform/http"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrInvalidPayload = errors.New("invalid payload")

type Service struct {
	repo      Repository
	jwtSecret []byte
	jwtTTL    time.Duration
}

func NewService(repo Repository, jwtSecret string) *Service {
	return &Service{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
		jwtTTL:    24 * time.Hour,
	}
}

func (s *Service) Register(ctx context.Context, email, password, name string) (*AuthResult, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	name = strings.TrimSpace(name)
	if email == "" || password == "" {
		return nil, ErrInvalidPayload
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u, err := s.repo.CreateUser(ctx, email, string(hash), name)
	if err != nil {
		return nil, err
	}

	tok, err := httpx.SignJWT(u.ID, s.jwtSecret, s.jwtTTL)
	if err != nil {
		return nil, err
	}

	return &AuthResult{User: *u, Token: tok}, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || password == "" {
		return nil, ErrInvalidPayload
	}

	u, hash, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	tok, err := httpx.SignJWT(u.ID, s.jwtSecret, s.jwtTTL)
	if err != nil {
		return nil, err
	}

	return &AuthResult{User: *u, Token: tok}, nil
}

func (s *Service) Me(ctx context.Context, userID string) (*User, error) {
	return s.repo.GetUserByID(ctx, userID)
}
