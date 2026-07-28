package tablesapiv1

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	tablesconv "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter/tables"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesvalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/tables"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GetGroupLevelById отдаёт активную строку dc.group_levels по id.
func (t *TablesApiV1) GetGroupLevelById(ctx context.Context, req *tablesv1.GetGroupLevelByIdRequest) (*tablesv1.GetGroupLevelByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetGroupLevelById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetGroupLevelByIdResponse{GroupLevel: tablesconv.GroupLevelToProto(row)}, nil
}

// GetGroupLevels отдаёт все активные строки dc.group_levels.
func (t *TablesApiV1) GetGroupLevels(ctx context.Context, req *tablesv1.GetGroupLevelsRequest) (*tablesv1.GetGroupLevelsResponse, error) {
	rows, err := t.services.TablesService.GetGroupLevels(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetGroupLevelsResponse{GroupLevels: tablesconv.GroupLevelsToProto(rows)}, nil
}

// GetDeletedGroupLevelById отдаёт мягко удалённую строку dc.group_levels по id.
func (t *TablesApiV1) GetDeletedGroupLevelById(ctx context.Context, req *tablesv1.GetDeletedGroupLevelByIdRequest) (*tablesv1.GetDeletedGroupLevelByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetDeletedGroupLevelById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedGroupLevelByIdResponse{GroupLevel: tablesconv.GroupLevelToProto(row)}, nil
}

// GetDeletedGroupLevels отдаёт все мягко удалённые строки dc.group_levels.
func (t *TablesApiV1) GetDeletedGroupLevels(ctx context.Context, req *tablesv1.GetDeletedGroupLevelsRequest) (*tablesv1.GetDeletedGroupLevelsResponse, error) {
	rows, err := t.services.TablesService.GetDeletedGroupLevels(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetDeletedGroupLevelsResponse{GroupLevels: tablesconv.GroupLevelsToProto(rows)}, nil
}

// CreateGroupLevel вставляет строку dc.group_levels и отдаёт её целиком.
func (t *TablesApiV1) CreateGroupLevel(ctx context.Context, req *tablesv1.CreateGroupLevelRequest) (*tablesv1.CreateGroupLevelResponse, error) {
	if err := tablesvalidation.ValidateCreateGroupLevel(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.CreateGroupLevel(ctx, tablesconv.ToCreateGroupLevelParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.CreateGroupLevelResponse{GroupLevel: tablesconv.GroupLevelToProto(row)}, nil
}

// UpdateGroupLevelById обновляет активную строку dc.group_levels и отдаёт её целиком.
func (t *TablesApiV1) UpdateGroupLevelById(ctx context.Context, req *tablesv1.UpdateGroupLevelByIdRequest) (*tablesv1.UpdateGroupLevelByIdResponse, error) {
	if err := tablesvalidation.ValidateUpdateGroupLevelById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.UpdateGroupLevelById(ctx, tablesconv.ToUpdateGroupLevelByIdParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UpdateGroupLevelByIdResponse{GroupLevel: tablesconv.GroupLevelToProto(row)}, nil
}

// DeleteGroupLevelById мягко удаляет строку dc.group_levels.
func (t *TablesApiV1) DeleteGroupLevelById(ctx context.Context, req *tablesv1.DeleteGroupLevelByIdRequest) (*tablesv1.DeleteGroupLevelByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.DeleteGroupLevelById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.DeleteGroupLevelByIdResponse{Empty: &emptypb.Empty{}}, nil
}

// UndeleteGroupLevelById восстанавливает мягко удалённую строку dc.group_levels.
func (t *TablesApiV1) UndeleteGroupLevelById(ctx context.Context, req *tablesv1.UndeleteGroupLevelByIdRequest) (*tablesv1.UndeleteGroupLevelByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := t.services.TablesService.UndeleteGroupLevelById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.UndeleteGroupLevelByIdResponse{Empty: &emptypb.Empty{}}, nil
}
