package healthcheck

import (
	"net/http"

	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/custresponse"
)

func ReturnOk(w http.ResponseWriter, r *http.Request) {
	messageOk := map[string]string{"message": "Ok"}
	custresponse.WriteJSON(w, http.StatusOK, messageOk, nil)
}
