package handlers

import (
	auth "clothing-store/pkg/pb/auth"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
)

func (app *Application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	req := &auth.LoginRequest{Email: input.Email, Password: input.Password}
	reqBytes, err := proto.Marshal(req)

	if err != nil {
		app.errorResponse(w, r, 500, "Error while marshalling")
	}

	resByte, err := GetRabbitResponse("token", "auth_queue", reqBytes)

	if err != nil {
		app.errorResponse(w, r, 500, "Error while getting response on RabbitMQ")
	}
	response := &auth.LoginResponse{}

	err = proto.Unmarshal(resByte, response)
	failOnError(err, "Failed to convert body to Response")

	if response.Token == "" {
		app.writeJSON(w, http.StatusNotFound, envelope{"message": "Incorrect credentials or user doesn't exist"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, response, nil)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
