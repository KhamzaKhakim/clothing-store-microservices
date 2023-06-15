package handlers

import (
	"clothing-store/internal/client"
	"clothing-store/internal/data"
	clothes "clothing-store/pkg/pb/clothe"
	"context"
	"net/http"
)

func (app *Application) addToCartHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	//clothe, err := app.models.Clothes.Get(id)
	//if err != nil {
	//	switch {
	//	case errors.Is(err, data.ErrRecordNotFound):
	//		app.notFoundResponse(w, r)
	//	default:
	//		app.serverErrorResponse(w, r, err)
	//	}
	//	return
	//}

	res, err := client.GetClotheClient().ShowClothe(context.Background(), &clothes.ShowClotheRequest{Id: id})

	clothe := data.Clothe{
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
	user := app.contextGetUser(r)
	err = app.Models.Users.UpdateMoney(user, clothe.Price)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	_, err = app.Models.Carts.GetById(user.ID)
	if err != nil {
		app.Models.Carts.CreateCartForUser(user.ID)
	}

	err = app.Models.Carts.AddClotheForCart(user.ID, clothe)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "clothe successfully added to the cart"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) showCartHandler(w http.ResponseWriter, r *http.Request) {

	var response struct {
		Name      string  `json:"name"`
		Money     int64   `json:"money"`
		ClothesID []int64 `json:"clothes_id"`
	}

	user := app.contextGetUser(r)

	clothes, err := app.Models.Carts.GetById(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response.Name = user.Name
	response.Money = user.Money
	response.ClothesID = clothes

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
