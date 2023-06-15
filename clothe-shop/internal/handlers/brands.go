package handlers

import (
	"clothing-store/internal/client"
	"clothing-store/internal/data"
	brand "clothing-store/pkg/pb/brand"
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net/http"
)

func (app *Application) createBrandHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name        string `json:"name"`
		Country     string `json:"country"`
		Description string `json:"description"`
		ImageURL    string `json:"image_url"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	res, err := client.GetBrandClient().CreateBrand(context.Background(), &brand.Brand{
		Name:        input.Name,
		Country:     input.Country,
		Description: input.Description,
		ImageUrl:    input.ImageURL,
	})

	if err != nil {
		app.rpcErrorResponse(w, r, err.Error())
		return
	}
	response := data.Brand{
		ID:          res.Id,
		Name:        res.Name,
		Country:     input.Country,
		Description: input.Description,
		ImageURL:    res.ImageUrl,
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/brands/%d", response.ID))
	err = app.writeJSON(w, http.StatusCreated, response, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) showBrandHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	req := &brand.ShowBrandRequest{Id: id}
	reqBytes, err := proto.Marshal(req)
	if err != nil {
		app.errorResponse(w, r, 500, "Error while marshalling")
	}

	resByte, err := GetRabbitResponse("brand", "brand_queue", reqBytes)
	if err != nil {
		app.errorResponse(w, r, 500, "Error while getting response on RabbitMQ")
	}
	response := &brand.Brand{}

	err = proto.Unmarshal(resByte, response)
	failOnError(err, "Failed to convert body to Response")

	if response.Name == "" {
		app.writeJSON(w, http.StatusNotFound, envelope{"message": "Brand not found"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, response, nil)
}

func (app *Application) updateBrandHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Name        *string `json:"name"`
		Country     *string `json:"country"`
		Description *string `json:"description"`
		ImageURL    *string `json:"image_url"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	updateBrand := &brand.Brand{}

	if input.Name != nil {
		updateBrand.Name = *input.Name
	}
	if input.Country != nil {
		updateBrand.Country = *input.Country
	}
	if input.Description != nil {
		updateBrand.Description = *input.Description
	}
	if input.ImageURL != nil {
		updateBrand.ImageUrl = *input.ImageURL
	}

	res, err := client.GetBrandClient().UpdateBrand(context.Background(),
		&brand.UpdateBrandRequest{Id: id, Brand: updateBrand})

	if err != nil {
		app.rpcErrorResponse(w, r, err.Error())
		return
	}

	err = app.writeJSON(w, http.StatusOK, res, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *Application) deleteBrandHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	_, err = client.GetBrandClient().DeleteBrand(context.Background(), &brand.DeleteBrandRequest{Id: id})

	if err != nil {
		app.rpcErrorResponse(w, r, err.Error())
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "brand successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) listBrandsHandler(w http.ResponseWriter, r *http.Request) {
	res, err := client.GetBrandClient().ListBrand(context.Background(), &brand.ListBrandRequest{})
	if err != nil {
		app.rpcErrorResponse(w, r, err.Error())
		return
	}

	err = app.writeJSON(w, http.StatusOK, res.BrandList, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
