package usermodel

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/konstantin-suspitsyn/datacomrade/configs"
	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/validator"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int64     `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type password struct {
	plaintext *string
	hash      []byte
}

type UserModel struct {
	DB *sql.DB
}

///////////////////////////////////////////////////////
////                  VALIDATION                   ////
///////////////////////////////////////////////////////

func ValidateUser(v *validator.Validator, user *User) {
	nameMaxChars := 50
	minPasswordLength := 1
	maxPasswordLength := 50

	// Check Name
	v.Check(len(user.Name) <= nameMaxChars, "name", fmt.Sprintf("Name is longer than %d", nameMaxChars))
	v.Check(user.Name != "", "name", "Name must be provided")

	// Check email
	ValidateUserEmail(v, user.Email)
	// Check password
	v.Check(*user.Password.plaintext != "", "email", "Email must be provided")
	v.Check(v.Matches(user.Email, validator.EmailRX), "email", "It's not an email")
	v.Check(len(*user.Password.plaintext) >= minPasswordLength, "password", fmt.Sprintf("Password should be longer than %d", minPasswordLength))
	v.Check(len(*user.Password.plaintext) <= maxPasswordLength, "password", fmt.Sprintf("Password should be smaller than %d", minPasswordLength))

	if user.Password.hash == nil {
		panic("Missing hash for password")
	}

}

func ValidateUserEmail(v *validator.Validator, email string) {
	emailMaxChars := 100
	v.Check(len(email) <= emailMaxChars, "email", fmt.Sprintf("Email is longer than %d", emailMaxChars))
	v.Check(email != "", "email", "email must be provided")
	v.Check(v.Matches(email, validator.EmailRX), "email", "It's not an email")

}

///////////////////////////////////////////////////////
////                  CHECK PASSWORD               ////
///////////////////////////////////////////////////////

// Create hashed password
func (p *password) Set(plaintextPassword string) error {

	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)

	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

// Compare hasp and pass
func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil

}

///////////////////////////////////////////////////////
////                     CRUD                      ////
///////////////////////////////////////////////////////

func (m UserModel) Insert(ctx context.Context, user *User) error {
	sqlQuery := `INSERT INTO users.user (email, "name", password_hash, activated, created_at, updated_at) 
	VALUES($1, $2, $3, false, now(), now())
	RETURNING id, created_at, updated_at;`

	args := []any{user.Email, user.Name, user.Password.hash}

	ctx, cancel := context.WithTimeout(ctx, configs.QueryTimeoutShort)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, sqlQuery, args...).Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "user_email_unique"):
			return ErrDuplicateEmail
		case strings.Contains(err.Error(), "user_name_unique"):
			return ErrDuplicateName
		default:
			return err

		}
	}
	return nil
}

func (m UserModel) GetByEmail(ctx context.Context, email string) (*User, error) {
	sqlQuery := `SELECT id, email, "name", password_hash, activated, created_at, updated_at FROM users.user
	where email = $1
	and is_deleted = false;`

	ctx, cancel := context.WithTimeout(ctx, configs.QueryTimeoutShort)
	defer cancel()

	var user User

	err := m.DB.QueryRowContext(ctx, sqlQuery, email).Scan(
		&user.Id,
		&user.Email,
		&user.Name,
		&user.Password.hash,
		&user.Activated,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, fmt.Errorf("%w. Searched email: %s", ErrUserRecordNotFound, email)
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) GetActiveByEmail(ctx context.Context, email string) (*User, error) {
	sqlQuery := `SELECT id, email, "name", password_hash, activated, created_at, updated_at FROM users.user
	where email = $1
	and activated = true
	and is_deleted = false;`

	ctx, cancel := context.WithTimeout(ctx, configs.QueryTimeoutShort)
	defer cancel()

	var user User

	err := m.DB.QueryRowContext(ctx, sqlQuery, email).Scan(
		&user.Id,
		&user.Email,
		&user.Name,
		&user.Password.hash,
		&user.Activated,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, fmt.Errorf("%w. Searched email: %s", ErrUserRecordNotFound, email)
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) ActivateUserById(ctx context.Context, id int64) error {
	sqlQuery := `UPDATE users.user SET activated=true, updated_at=now() 
	WHERE id = $1;`

	ctx, cancel := context.WithTimeout(ctx, configs.QueryTimeoutShort)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, sqlQuery, id)

	if err != nil {
		return err
	}

	return nil

}

func (m UserModel) GetById(ctx context.Context, userId int64) (*User, error) {
	user := User{}

	query := `SELECT id, email, "name", password_hash, activated, created_at, updated_at FROM users.user
	WHERE id = $1;`

	ctx, cancel := context.WithTimeout(ctx, configs.QueryTimeoutShort)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, userId).Scan(
		&user.Id,
		&user.Email,
		&user.Name,
		&user.Password.hash,
		&user.Activated,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrUserRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil

}

func (m *UserModel) UpdatePassword(ctx context.Context, userId int64, plainPassword string) error {

	var pass password

	pass.Set(plainPassword)

	query := `UPDATE users.user
	SET password_hash=$1
	WHERE id=$2;`

	ctx, cancel := context.WithTimeout(ctx, configs.QueryTimeoutShort)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, pass.hash, userId)

	if err != nil {
		return err
	}

	return nil

}
