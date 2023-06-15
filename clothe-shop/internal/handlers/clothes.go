package handlers

import (
	"clothing-store/internal/client"
	"clothing-store/internal/data"
	"clothing-store/internal/validator"
	brand "clothing-store/pkg/pb/brand"
	clothes "clothing-store/pkg/pb/clothe"
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"net/http"
)

func (app *Application) createClotheHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name     string   `json:"name"`
		Price    int64    `json:"price"`
		Brand    string   `json:"brand"`
		Color    string   `json:"color"`
		Sizes    []string `json:"sizes"`
		Sex      string   `json:"sex"`
		Type     string   `json:"type"`
		ImageURL string   `json:"image_url"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	res, err := client.GetClotheClient().CreateClothe(context.Background(), &clothes.Clothe{
		Name:     input.Name,
		Price:    input.Price,
		Brand:    input.Brand,
		Color:    input.Color,
		Sizes:    input.Sizes,
		Sex:      input.Sex,
		Type:     input.Type,
		ImageUrl: input.ImageURL,
	})

	if err != nil {
		app.rpcErrorResponse(w, r, err.Error())
		return
	}
	response := data.Clothe{
		ID:       res.Id,
		Name:     res.Name,
		Price:    res.Price,
		Brand:    res.Brand,
		Color:    res.Color,
		Sizes:    res.Sizes,
		Sex:      res.Sex,
		Type:     res.Type,
		ImageURL: res.ImageUrl,
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/clothes/%d", response.ID))
	err = app.writeJSON(w, http.StatusCreated, response, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *Application) showClotheHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	req := &clothes.ShowClotheRequest{Id: id}
	reqBytes, err := proto.Marshal(req)
	if err != nil {
		app.errorResponse(w, r, 500, "Error while marshalling")
	}

	resByte, err := GetRabbitResponse("clothe", "clothe_queue", reqBytes)
	if err != nil {
		app.errorResponse(w, r, 500, "Error while getting response on RabbitMQ")
	}
	response := &clothes.Clothe{}

	err = proto.Unmarshal(resByte, response)
	failOnError(err, "Failed to convert body to Response")

	if response.Name == "" {
		app.writeJSON(w, http.StatusNotFound, envelope{"message": "Clothe not found"}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, response, nil)
}

func (app *Application) updateClotheHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Name     *string  `json:"name"`
		Price    *int64   `json:"price"`
		Brand    *string  `json:"brand"`
		Color    *string  `json:"color"`
		Sizes    []string `json:"sizes"`
		Sex      *string  `json:"sex"`
		Type     *string  `json:"type"`
		ImageURL *string  `json:"image_url"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	updateClothe := &clothes.Clothe{}

	if input.Name != nil {
		updateClothe.Name = *input.Name
	}
	if input.Price != nil {
		updateClothe.Price = *input.Price
	}
	if input.Brand != nil {
		updateClothe.Brand = *input.Brand
	}
	if input.Color != nil {
		updateClothe.Color = *input.Color
	}
	if input.Sizes != nil {
		updateClothe.Sizes = input.Sizes
	}
	if input.Sex != nil {
		updateClothe.Sex = *input.Sex
	}
	if input.Type != nil {
		updateClothe.Type = *input.Type
	}
	if input.ImageURL != nil {
		updateClothe.ImageUrl = *input.ImageURL
	}

	res, err := client.GetClotheClient().UpdateClothe(context.Background(),
		&clothes.UpdateClotheRequest{Id: id, Clothe: updateClothe})

	if err != nil {
		app.rpcErrorResponse(w, r, err.Error())
		return
	}

	err = app.writeJSON(w, http.StatusOK, res, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) deleteClotheHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	_, err = client.GetClotheClient().DeleteClothe(context.Background(), &clothes.DeleteClotheRequest{Id: id})

	if err != nil {
		app.rpcErrorResponse(w, r, err.Error())
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "clothe successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) listClothesHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name     string
		Brand    string
		PriceMax int64
		PriceMin int64
		Sizes    []string
		Color    string
		Type     string
		Sex      string
		data.Filters
	}
	v := validator.New()
	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.Brand = app.readString(qs, "brand", "")
	input.PriceMax = app.readInt(qs, "price_max", 10000000, v)
	input.PriceMin = app.readInt(qs, "price_min", 0, v)
	input.Name = app.readString(qs, "name", "")
	input.Sizes = app.readCSV(qs, "sizes", []string{})
	input.Sex = app.readString(qs, "sex", "")
	input.Sex = app.readString(qs, "sex", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = 20
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "price", "sex", "brand", "-id", "-name", "-price", "-sex", "-brand"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	var allBrands []string
	resBrands, err := client.GetBrandClient().ListBrand(context.Background(), &brand.ListBrandRequest{})

	for i := 0; i < len(resBrands.GetBrandList()); i++ {
		allBrands = append(allBrands, resBrands.GetBrandList()[i].Name)
	}
	allBrands = append(allBrands, "")

	keys := data.Keys{
		PriceMax:       input.PriceMax,
		PriceMin:       input.PriceMin,
		Brand:          input.Brand,
		Sizes:          input.Sizes,
		SizesSafelist:  []string{"XS", "S", "M", "L", "XL", ""},
		BrandsSafelist: allBrands,
	}

	if data.ValidateKeys(v, keys); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	res, err := client.GetClotheClient().ListClothe(context.Background(), &clothes.ListClotheRequest{
		Name:     input.Name,
		Brand:    input.Brand,
		PriceMax: input.PriceMax,
		PriceMin: input.PriceMin,
		Sizes:    input.Sizes,
		Color:    input.Color,
		Type:     input.Type,
		Sex:      input.Sex,
		Filter: &clothes.Filter{
			Page:         input.Filters.Page,
			PageSize:     input.Filters.PageSize,
			Sort:         input.Filters.Sort,
			SortSafeList: input.Filters.SortSafelist,
		},
	})

	if err != nil {
		app.rpcErrorResponse(w, r, err.Error())
		return
	}
	err = app.writeJSON(w, http.StatusOK, res.GetClotheList(), nil)

}
