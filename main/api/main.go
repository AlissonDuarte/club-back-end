package main

import (
    "net/http"

	"clube/internal/serializer"
	"clube/infraestructure/database"
	"clube/infraestructure/models"
    "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)


func main() {
	// Call the function
	conn := database.NewDb()
	app := chi.NewRouter()
	app.Use(JWTMiddleware)
	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)
	app.Use(middleware.Throttle(1000))
	
	app.Get("/", func(w http.ResponseWriter, app *http.Request) {
		w.Write([]byte("Hello World"))
	})
	
	app.Post("/users", func(w http.ResponseWriter, app *http.Request) {
		var userData serializer.UserSerializer

		render.DecodeJSON(app.Body, &userData)
		render.JSON(w, app, userData)
	})

	err :=	models.Migrate(conn)
	if err != nil {
		panic(err)
	}

	http.ListenAndServe(":3000", app)
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, app *http.Request) {
		// before
		next.ServeHTTP(w, app)
		// after
	})
}