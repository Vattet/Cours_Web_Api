package handler

import (
	"net/http"
	"encoding/json"

	"rest-api/model"
	"rest-api/service"
)

func UserConnection() http.HandlerFunc {
	return func(writer http.ResponseWriter, reader *http.Request) {
		user := model.User{}
		if err := json.NewDecoder(reader.Body).Decode(&user); err != nil {
			writeError(writer, http.StatusBadRequest, "invalid JSON body")
			return
		}

		token, err := service.Login(user.Username, user.Password)
		if err != nil {
			writeError(writer, http.StatusUnauthorized, err.Error())
			return
		}

		body := map[string]string{"token": token}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(&body)
	}
}
