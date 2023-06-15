package handlers

import (
	"clothing-store/internal/client"
	"clothing-store/internal/data"
	pb "clothing-store/pkg/pb/auth"
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"net/http"
)

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *Application) rateLimit(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(2, 4)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			app.rateLimitExceededResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *Application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		res, err := client.GetAuthClient().Authenticate(context.Background(), &pb.AuthenticateRequest{Token: authorizationHeader})

		if err != nil || res == nil {
			app.rpcErrorResponse(w, r, err.Error())
			return
		}

		user := &data.User{ID: res.Id, Activated: res.Activated, Money: res.Money}

		r = app.contextSetUser(r, user)

		next.ServeHTTP(w, r)
	})
}

func (app *Application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}
		if !user.Activated {
			app.inactiveAccountResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})

}

func (app *Application) requireRole(role string, next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		_, err := client.GetAuthClient().Authorize(context.Background(), &pb.AuthorizeRequest{Id: user.ID, Role: role})

		if err != nil {
			app.rpcErrorResponse(w, r, err.Error())
			return
		}
		next.ServeHTTP(w, r)
	}
	return app.requireActivatedUser(fn)
}

func (app *Application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

		next.ServeHTTP(w, r)
	})
}
