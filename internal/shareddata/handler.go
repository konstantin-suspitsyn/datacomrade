package shareddata

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/konstantin-suspitsyn/datacomrade/data/sharedmodels"
	"github.com/konstantin-suspitsyn/datacomrade/data/sharedtypes"
	"github.com/konstantin-suspitsyn/datacomrade/data/usermodel"
	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/custresponse"
	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/urlparams"
)

func (sds *SharedDataService) GetAllDomainsHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	urlPagingParams, err := urlparams.GetPager(r)

	if err != nil {
		custresponse.BadRequestResponse(w, r, err)
	}

	dto, err := sds.getDataWithPager(ctx, urlPagingParams)

	if err != nil {
		custresponse.BadRequestResponse(w, r, err)
		return
	}

	custresponse.WriteJSON(w, http.StatusOK, dto, nil)

}

func (sds *SharedDataService) CreateDomainHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var input sharedmodels.DomainInputDTO
	appUser := r.Context().Value(sharedtypes.AuthKey{}).(*usermodel.AppUser)
	err := custresponse.ReadJSON(w, r, &input)
	if err != nil {
		custresponse.BadRequestResponse(w, r, err)
		return
	}
	domain, err := sds.CreateDomain(ctx, input.Name, input.Description, appUser.Id)

	if err != nil {
		custresponse.BadRequestResponse(w, r, err)
		return
	}

	custresponse.WriteJSON(w, http.StatusCreated, domain, nil)

}

func (sds *SharedDataService) DeleteDomainHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	domainId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		custresponse.BadRequestResponse(w, r, fmt.Errorf("ERROR readind id, %w", err))
		return

	}
	appUser := r.Context().Value(sharedtypes.AuthKey{}).(*usermodel.AppUser)

	err = sds.DeleteDomain(ctx, domainId, appUser.Id)

	if err != nil {
		custresponse.BadRequestResponse(w, r, fmt.Errorf("ERROR deleting domain by id = %d, %w", domainId, err))
		return

	}
}
func (sds *SharedDataService) UpdateDomainHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	domainId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		custresponse.BadRequestResponse(w, r, fmt.Errorf("ERROR readind id, %w", err))
		return

	}
	var input sharedmodels.DomainInputDTO
	appUser := r.Context().Value(sharedtypes.AuthKey{}).(*usermodel.AppUser)
	err = custresponse.ReadJSON(w, r, &input)
	if err != nil {
		custresponse.BadRequestResponse(w, r, err)
		return
	}

	domain, err := sds.UpdateDomainByUser(ctx, domainId, input.Name, input.Description, appUser.Id)

	if err != nil {
		custresponse.BadRequestResponse(w, r, err)
		return
	}

	custresponse.WriteJSON(w, http.StatusCreated, domain, nil)

}
