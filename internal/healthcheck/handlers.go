package healthcheck

import (
	"net/http"

	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/custresponse"
)

func ReturnOk(w http.ResponseWriter, r *http.Request) {
	custresponse.WriteJSON(w, http.StatusOK, nil, nil)
}
