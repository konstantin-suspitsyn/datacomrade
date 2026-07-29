package userapiv1

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	userconv "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter/user"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	uservalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/user"
	userv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GetUserById отдаёт активную строку dc.user по id.
func (u *UserApiV1) GetUserById(ctx context.Context, req *userv1.GetUserByIdRequest) (*userv1.GetUserByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserService.GetUserById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userv1.GetUserByIdResponse{User: userconv.UserToProto(row)}, nil
}

// GetUsers отдаёт все активные строки dc.user.
func (u *UserApiV1) GetUsers(ctx context.Context, _ *userv1.GetUsersRequest) (*userv1.GetUsersResponse, error) {
	rows, err := u.services.UserService.GetUsers(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userv1.GetUsersResponse{Users: userconv.UsersToProto(rows)}, nil
}

// GetDeletedUserById отдаёт мягко удалённую строку dc.user по id.
func (u *UserApiV1) GetDeletedUserById(ctx context.Context, req *userv1.GetDeletedUserByIdRequest) (*userv1.GetDeletedUserByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserService.GetDeletedUserById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userv1.GetDeletedUserByIdResponse{User: userconv.UserToProto(row)}, nil
}

// GetDeletedUsers отдаёт все мягко удалённые строки dc.user.
func (u *UserApiV1) GetDeletedUsers(ctx context.Context, _ *userv1.GetDeletedUsersRequest) (*userv1.GetDeletedUsersResponse, error) {
	rows, err := u.services.UserService.GetDeletedUsers(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userv1.GetDeletedUsersResponse{Users: userconv.UsersToProto(rows)}, nil
}

// GetUserByExternalId отдаёт активную строку dc.user по уникальной колонке external_id.
func (u *UserApiV1) GetUserByExternalId(ctx context.Context, req *userv1.GetUserByExternalIdRequest) (*userv1.GetUserByExternalIdResponse, error) {
	if err := uservalidation.ValidateGetUserByExternalId(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserService.GetUserByExternalId(ctx, userconv.ToGetUserByExternalIdArg(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userv1.GetUserByExternalIdResponse{User: userconv.UserToProto(row)}, nil
}

// CreateUser вставляет строку dc.user и отдаёт её целиком.
func (u *UserApiV1) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	if err := uservalidation.ValidateCreateUser(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserService.CreateUser(ctx, userconv.ToCreateUserParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userv1.CreateUserResponse{User: userconv.UserToProto(row)}, nil
}

// UpdateUserById обновляет активную строку dc.user и отдаёт её целиком.
func (u *UserApiV1) UpdateUserById(ctx context.Context, req *userv1.UpdateUserByIdRequest) (*userv1.UpdateUserByIdResponse, error) {
	if err := uservalidation.ValidateUpdateUserById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserService.UpdateUserById(ctx, userconv.ToUpdateUserByIdParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userv1.UpdateUserByIdResponse{User: userconv.UserToProto(row)}, nil
}

// DeleteUserById мягко удаляет строку dc.user.
func (u *UserApiV1) DeleteUserById(ctx context.Context, req *userv1.DeleteUserByIdRequest) (*userv1.DeleteUserByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := u.services.UserService.DeleteUserById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userv1.DeleteUserByIdResponse{Empty: &emptypb.Empty{}}, nil
}

// UndeleteUserById восстанавливает мягко удалённую строку dc.user.
func (u *UserApiV1) UndeleteUserById(ctx context.Context, req *userv1.UndeleteUserByIdRequest) (*userv1.UndeleteUserByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := u.services.UserService.UndeleteUserById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userv1.UndeleteUserByIdResponse{Empty: &emptypb.Empty{}}, nil
}
