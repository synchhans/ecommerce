package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

type fakeRepo struct {
	createFn func(ctx context.Context, email, passwordHash, name string) (*User, error)
	getByE   func(ctx context.Context, email string) (*User, string, error)
	getByID  func(ctx context.Context, userID string) (*User, error)
}

func (f fakeRepo) CreateUser(ctx context.Context, email, passwordHash, name string) (*User, error) {
	return f.createFn(ctx, email, passwordHash, name)
}
func (f fakeRepo) GetUserByEmail(ctx context.Context, email string) (*User, string, error) {
	return f.getByE(ctx, email)
}
func (f fakeRepo) GetUserByID(ctx context.Context, userID string) (*User, error) {
	return f.getByID(ctx, userID)
}

func TestRegister_201(t *testing.T) {
	repo := fakeRepo{
		createFn: func(ctx context.Context, email, passwordHash, name string) (*User, error) {
			require.Equal(t, "a@b.com", email)
			require.NotEmpty(t, passwordHash)
			return &User{ID: "u1", Email: email, Name: name, Status: "active"}, nil
		},
		getByE:  func(ctx context.Context, email string) (*User, string, error) { return nil, "", ErrNotFound },
		getByID: func(ctx context.Context, userID string) (*User, error) { return &User{ID: "u1"}, nil },
	}

	svc := NewService(repo, "secret")
	h := NewHandler(svc, "secret")

	r := chi.NewRouter()
	h.Routes(r)

	reqBody := []byte(`{"email":"A@B.com","password":"pass123","name":"A"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	var out AuthResult
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &out))
	require.Equal(t, "u1", out.User.ID)
	require.NotEmpty(t, out.Token)
}

func TestLogin_401(t *testing.T) {
	repo := fakeRepo{
		createFn: func(ctx context.Context, email, passwordHash, name string) (*User, error) { return nil, nil },
		getByE: func(ctx context.Context, email string) (*User, string, error) {
			return nil, "", ErrNotFound
		},
		getByID: func(ctx context.Context, userID string) (*User, error) { return nil, ErrNotFound },
	}
	svc := NewService(repo, "secret")
	h := NewHandler(svc, "secret")

	r := chi.NewRouter()
	h.Routes(r)

	reqBody := []byte(`{"email":"a@b.com","password":"wrong"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}
