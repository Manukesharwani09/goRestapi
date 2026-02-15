package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/Manukesharwani09/goRestapi/internal/types"
	"github.com/Manukesharwani09/goRestapi/internal/utils/response"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return

		}
		slog.Info("creating a request")
		response.WriteJson(w, http.StatusCreated, map[string]string{"message": "student created successfully"})
	}
}
