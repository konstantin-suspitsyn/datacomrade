package userservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_model"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetUserById возвращает активную строку dc.user по id.
func (s *UserService) GetUserById(ctx context.Context, id int64) (user_model.DcUser, error) {
	row, err := s.UserRepository.GetUserById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user_model.DcUser{}, fmt.Errorf("dc.user id = %d: %w", id, customerrors.ErrNotFound)
		}

		return user_model.DcUser{}, fmt.Errorf("get dc.user id = %d: %w", id, err)
	}

	return row, nil
}

// GetDeletedUserById возвращает мягко удалённую строку dc.user по id.
func (s *UserService) GetDeletedUserById(ctx context.Context, id int64) (user_model.DcUser, error) {
	row, err := s.UserRepository.GetDeletedUserById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user_model.DcUser{}, fmt.Errorf("deleted dc.user id = %d: %w", id, customerrors.ErrNotFound)
		}

		return user_model.DcUser{}, fmt.Errorf("get deleted dc.user id = %d: %w", id, err)
	}

	return row, nil
}

// GetUsers возвращает строки dc.user.
func (s *UserService) GetUsers(ctx context.Context) ([]user_model.DcUser, error) {
	rows, err := s.UserRepository.GetUsers(ctx)

	if err != nil {
		return nil, fmt.Errorf("GetUsers: %w", err)
	}

	return rows, nil
}

// GetDeletedUsers возвращает строки dc.user.
func (s *UserService) GetDeletedUsers(ctx context.Context) ([]user_model.DcUser, error) {
	rows, err := s.UserRepository.GetDeletedUsers(ctx)

	if err != nil {
		return nil, fmt.Errorf("GetDeletedUsers: %w", err)
	}

	return rows, nil
}

// GetUserByExternalId возвращает активную строку dc.user по уникальной колонке external_id.
func (s *UserService) GetUserByExternalId(ctx context.Context, externalID uuid.UUID) (user_model.DcUser, error) {
	row, err := s.UserRepository.GetUserByExternalId(ctx, externalID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user_model.DcUser{}, fmt.Errorf("dc.user external_id = %v: %w", externalID, customerrors.ErrNotFound)
		}

		return user_model.DcUser{}, fmt.Errorf("get dc.user external_id = %v: %w", externalID, err)
	}

	return row, nil
}

// CreateUser вставляет строку dc.user и возвращает её целиком.
func (s *UserService) CreateUser(ctx context.Context, params user_model.CreateUserParams) (user_model.DcUser, error) {
	row, err := s.UserRepository.CreateUser(ctx, params)

	if err != nil {
		return user_model.DcUser{}, fmt.Errorf("%w: dc.user: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateUserById обновляет активную строку dc.user и возвращает её целиком.
//
// Запрос фильтрует по is_deleted = false, поэтому попытка обновить удалённую
// или несуществующую запись даёт sql.ErrNoRows — переводим его в ErrNotFound,
// чтобы api-слой ответил NotFound, а не Internal.
func (s *UserService) UpdateUserById(ctx context.Context, params user_model.UpdateUserByIdParams) (user_model.DcUser, error) {
	row, err := s.UserRepository.UpdateUserById(ctx, params)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user_model.DcUser{}, fmt.Errorf("dc.user id = %d: %w", params.ID, customerrors.ErrNotFound)
		}

		return user_model.DcUser{}, fmt.Errorf("%w: dc.user id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return row, nil
}

// DeleteUserById мягко удаляет строку dc.user.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *UserService) DeleteUserById(ctx context.Context, id int64) error {
	if _, err := s.GetUserById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.UserRepository.DeleteUserById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.user id = %d: %w", customerrors.ErrDelete, id, err)
	}

	return nil
}

// UndeleteUserById восстанавливает мягко удалённую строку dc.user.
//
// Существование удалённой записи проверяется заранее по той же причине,
// что и в DeleteUserById.
func (s *UserService) UndeleteUserById(ctx context.Context, id int64) error {
	if _, err := s.GetDeletedUserById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrUndelete, err)
	}

	if err := s.UserRepository.UndeleteUserById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.user id = %d: %w", customerrors.ErrUndelete, id, err)
	}

	return nil
}
