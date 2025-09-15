package shareddata

import (
	"net/http"

	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/custresponse"
)

func (sds *SharedDataService) GetAllDomains(w http.ResponseWriter, r *http.Request) {

	dto, err := sds.getDataWithPager(r)

	if err != nil {
		custresponse.BadRequestResponse(w, r, err)
		return
	}

	custresponse.WriteJSON(w, http.StatusOK, dto, nil)

}

func (sds *SharedDataService) CreateDomain(w http.ResponseWriter, r *http.Request) {
}

func (sds *SharedDataService) DeleteDomain(w http.ResponseWriter, r *http.Request) {
}
func (sds *SharedDataService) UpdateDomain(w http.ResponseWriter, r *http.Request) {
}
