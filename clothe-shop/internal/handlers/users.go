package handlers

import (
	"clothing-store/internal/client"
	"clothing-store/internal/data"
	auth "clothing-store/pkg/pb/auth"
	"context"
	"errors"
	"net/http"
)

func (app *Application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var response struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		Activated bool   `json:"activated"`
		Money     int64  `json:"money"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	res, err := client.GetAuthClient().Register(context.Background(), &auth.RegisterRequest{
		Name:     input.Name,
		Password: input.Password,
		Email:    input.Email,
	})

	if err != nil {
		app.rpcErrorResponse(w, r, err.Error())
		return
	}

	response.Name = res.Name
	response.Email = res.Email
	response.Activated = res.Activated
	response.Money = res.Money

	err = app.writeJSON(w, http.StatusAccepted, response, nil)
}

func (app *Application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var response struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		Activated bool   `json:"activated"`
		Money     int64  `json:"money"`
	}

	qs := r.URL.Query()

	token := app.readString(qs, "token", "")

	res, err := client.GetAuthClient().Activate(context.Background(), &auth.ActivateRequest{Token: token})

	if err != nil {
		app.rpcErrorResponse(w, r, err.Error())
		return
	}

	response.Name = res.Name
	response.Email = res.Email
	response.Activated = res.Activated
	response.Money = res.Money

	err = app.writeJSON(w, http.StatusAccepted, response, nil)

}

func (app *Application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.Models.Users.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "user successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
