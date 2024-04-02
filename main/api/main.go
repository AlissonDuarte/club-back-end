package main

import (
    "net/http"
    "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type user struct {
	Name string
	BirthDate string
	Cep string
}

func main() {
	// Call the function
	app := chi.NewRouter()
	app.Use(JWTMiddleware)
	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)
	app.Use(middleware.Throttle(1000))
	
	app.Get("/", func(w http.ResponseWriter, app *http.Request) {
		w.Write([]byte("Hello World"))
	})
	
	app.Post("/users", func(w http.ResponseWriter, app *http.Request) {
		var userData user
		render.DecodeJSON(app.Body, &userData)
		render.JSON(w, app, userData)
	})

	http.ListenAndServe(":3000", app)
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, app *http.Request) {
		// before
		next.ServeHTTP(w, app)
		// after
	})
}