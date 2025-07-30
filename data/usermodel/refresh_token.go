package usermodel

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const ScopeRefresh = "refresh_token"

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	Id           string    `json:"id"`
	UserId       int64     `json:"user_id"`
	Expire       time.Time `json:"expire"`
	CreatedAt    time.Time `json:"-"`
	IsActive     bool      `json:"is_active"`
	UpdatedAt    time.Time `json:"updated_at"`
	RefreshToken string    `json:"refresh_token"`
}

type RefreshTokenModel struct {
	DB *sql.DB
}

// /////////////////////////////////////////////////////
// //                     CRUD                      ////
// /////////////////////////////////////////////////////
func (m *RefreshTokenModel) Insert(token *RefreshToken) error {

	query := `INSERT INTO users.refresh_token (user_id, expire, created_at, is_active, updated_at, refresh_token, id) VALUES($1, $2, now(), true, now(), $3, $4);`

	args := []any{token.UserId, token.Expire, token.RefreshToken, token.Id}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

// Sends all refresh tokens to black list
func (rt *RefreshTokenModel) DeactivateRefreshTokensForUserId(userId int64) error {
	query := `UPDATE users.refresh_token
	SET is_active = false, updated_at = now()
	WHERE user_id = $1;`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rt.DB.ExecContext(ctx, query, userId)
	if err != nil {
		return err
	}

	return nil

}

func (rt *RefreshTokenModel) GetById(id string) (*RefreshToken, error) {

	query := `SELECT user_id, expire, created_at, is_active, updated_at, refresh_token, id FROM users.refresh_token
	WHERE id = $1;`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var refreshToken RefreshToken

	err := rt.DB.QueryRowContext(ctx, query, id).Scan(
		&refreshToken.UserId,
		&refreshToken.Expire,
		&refreshToken.CreatedAt,
		&refreshToken.IsActive,
		&refreshToken.UpdatedAt,
		&refreshToken.RefreshToken,
		&refreshToken.Id,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, fmt.Errorf("%w. ID: %s", ErrNoRefreshTokenRecord, id)
		default:
			return nil, err
		}
	}
	return &refreshToken, nil

}
