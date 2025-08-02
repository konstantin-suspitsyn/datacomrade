package usermodel

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"errors"
	"log/slog"
	"time"

	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/validator"
)

const ScopeActivation = "activation"
const ScopeAuthentication = "authentication"
const ScopeForgotPassword = "forgotPassword"

type Token struct {
	PlainText   string    `json:"token"`
	Hash        []byte    `json:"-"`
	UserId      int64     `json:"-"`
	Expire      time.Time `json:"expire"`
	Scope       string    `json:"-"`
	ActivatedAt time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	IsActive    bool      `json:"-"`
}

type TokenModel struct {
	DB *sql.DB
}

type TokenInput struct {
	TokenPlainText string `json:"token"`
}

///////////////////////////////////////////////////////
////                  VALIDATION                   ////
///////////////////////////////////////////////////////

func ValidateTokenPlainText(v *validator.Validator, tokenPlainText string) {
	v.Check(tokenPlainText != "", "Token", "Must be provided")
	v.Check(len(tokenPlainText) == 26, "Token", "Must be 26 chars long")
}

///////////////////////////////////////////////////////
////                     CRUD                      ////
///////////////////////////////////////////////////////

func (m *TokenModel) New(userId int64, expirationDuration time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userId, expirationDuration, scope)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)

	return token, err

}

func (m *TokenModel) DeactivateTokensForUsers(scope string, userId int64) error {
	query := `UPDATE users.mail_token
	SET is_active = false, updated_at=now()
	where scope = $1 and user_id = $2 and is_active = true`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, scope, userId)

	return err
}

func (m *TokenModel) Insert(token *Token) error {
	query := `INSERT INTO users.mail_token (hash, user_id, expire, "scope", created_at, updated_at, is_active) 
	VALUES($1, $2, $3, $4, now(), now(), true);`

	args := []any{token.Hash, token.UserId, token.Expire, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	slog.Info("USER ID: ", "USER ID: ", token.UserId)
	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m *TokenModel) GetByPlainText(plainText, scope string) (*Token, error) {
	query := `SELECT hash, user_id, expire, "scope", created_at, updated_at FROM users.mail_token
	WHERE hash = $1 and scope = $2 and is_active = true;`

	tokenHash := sha256.Sum256([]byte(plainText))
	token := Token{}

	args := []any{tokenHash[:], scope}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&token.Hash,
		&token.UserId,
		&token.Expire,
		&token.Scope,
		&token.ActivatedAt,
		&token.UpdatedAt,
	)
	if err != nil {

		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrTokenRecordNotFound
		default:
			return nil, err
		}
	}

	token.IsActive = true

	return &token, nil
}

func (m *TokenModel) GetTokenByUserId(userId int64, scope string) (*Token, error) {
	query := `SELECT hash, user_id, expire, "scope", created_at, updated_at FROM users.mail_token
	WHERE id = $1 and scope = $2 and is_active = true;`

	token := Token{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, userId, scope).Scan(
		&token.Hash,
		&token.UserId,
		&token.Expire,
		&token.Scope,
		&token.ActivatedAt,
		&token.UpdatedAt,
	)
	if err != nil {

		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrTokenRecordNotFound
		default:
			return nil, err
		}
	}

	token.IsActive = true

	return &token, nil
}

///////////////////////////////////////////////////////
////                  UTILITIES                    ////
///////////////////////////////////////////////////////

func generateToken(userId int64, expirationDuration time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserId: userId,
		Expire: time.Now().Add(expirationDuration),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)

	if err != nil {
		return nil, err
	}

	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.PlainText))

	token.Hash = hash[:]

	return token, nil
}
