package roles

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/konstantin-suspitsyn/datacomrade/data/rolesmodel"
	"github.com/konstantin-suspitsyn/datacomrade/data/sharedtypes"
	"github.com/konstantin-suspitsyn/datacomrade/data/usermodel"
	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/custresponse"
)

func (rs *RoleService) GetAllRolesHandler(w http.ResponseWriter, r *http.Request) {
	dto, err := rs.getDataWithPager(r)

	if err != nil {
		custresponse.BadRequestResponse(w, r, err)
		return
	}

	custresponse.WriteJSON(w, http.StatusOK, dto, nil)

}

func (rs *RoleService) CreateRoleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var input rolesmodel.RoleInputDTO
	appUser := r.Context().Value(sharedtypes.AuthKey{}).(*usermodel.AppUser)
	err := custresponse.ReadJSON(w, r, &input)
	if err != nil {
		custresponse.BadRequestResponse(w, r, err)
		return
	}

	role, err := rs.CreateRole(ctx, input.NameShort, input.NameLong, input.Description, input.JWTExport, appUser.Id)

	if err != nil {
		custresponse.BadRequestResponse(w, r, err)
		return
	}
	custresponse.WriteJSON(w, http.StatusCreated, role, nil)
}

func (rs *RoleService) DeleteRoleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	roleId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		custresponse.BadRequestResponse(w, r, fmt.Errorf("ERROR readind id of the role, %w", err))
		return

	}
	appUser := r.Context().Value(sharedtypes.AuthKey{}).(*usermodel.AppUser)

	err = rs.DeleteRole(ctx, roleId, appUser.Id)

	if err != nil {
		custresponse.BadRequestResponse(w, r, fmt.Errorf("ERROR deletind a role by id: %d, %w", roleId, err))
		return
	}

	custresponse.WriteJSON(w, http.StatusOK, nil, nil)

}
func (rs *RoleService) UpdateRoleHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	roleId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		custresponse.BadRequestResponse(w, r, fmt.Errorf("ERROR readind id of the role, %w", err))
		return

	}
	var input rolesmodel.RoleInputDTO
	appUser := r.Context().Value(sharedtypes.AuthKey{}).(*usermodel.AppUser)
	err = custresponse.ReadJSON(w, r, &input)
	if err != nil {
		custresponse.BadRequestResponse(w, r, err)
		return
	}

	role, err := rs.UpdateRole(ctx, roleId, input.NameShort, input.NameLong, input.JWTExport, input.Description, appUser.Id)

	if err != nil {
		custresponse.BadRequestResponse(w, r, err)
		return
	}
	custresponse.WriteJSON(w, http.StatusOK, role, nil)
}
