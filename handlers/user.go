package handlers

import (
	"encoding/json"
	"net/http"
	"w00k/go/rest-ws/models"
	"w00k/go/rest-ws/repository"
	"w00k/go/rest-ws/server"

	"github.com/segmentio/ksuid"
)

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

// SignUpHandler: endpoint para insertar un user en la base de datos
// los casos son:
// - si el request es inválido, se retorna un response de error con HTTP 400
// - si no se puede generar el id, se retorna un resopnse de error con HTTP 500
// - si hay algún error con el repositorio al insertar el user, se retorna un response de error con HTTP 500
// - si es caso exitoso, se responde con un SignUpResponse con HTTP 200
func SignUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request = SignUpRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := ksuid.NewRandom()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		var user = models.User{
			Email:    request.Email,
			Password: request.Password,
			Id:       id.String(),
		}
		err = repository.InsertUser(r.Context(), &user)
		if err != nil {
			if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
				http.Error(w, "User is in use", http.StatusConflict)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SignUpResponse{
			Id:    user.Id,
			Email: user.Email,
		})
	}
}
