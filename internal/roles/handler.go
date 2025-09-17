package roles

import (
	"net/http"

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

func (rs *RoleService) DeleteDomain(w http.ResponseWriter, r *http.Request) {
}
func (rs *RoleService) UpdateDomain(w http.ResponseWriter, r *http.Request) {
}
