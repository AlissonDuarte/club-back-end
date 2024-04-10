package main

import (
	"clube/infraestructure/database"
	"clube/infraestructure/models"
	"clube/internal/middles"
	"clube/internal/views"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func main() {
	// Call the function
	conn := database.NewDb()
	app := chi.NewRouter()
	app.Use(cors.AllowAll().Handler)
	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)
	app.Use(middleware.Throttle(1000))

	fileServer := http.FileServer(http.Dir("./templates/imgs"))
	htmlServer := http.FileServer(http.Dir("./templates"))

	app.Handle("/*", http.StripPrefix("/", fileServer))
	app.Handle("/*", http.StripPrefix("/", htmlServer))

	app.Post("/users", views.UserCreate)
	app.Get("/home", views.Home)
	app.Post("/user/login", views.UserLogin)

	app.Group(func(app chi.Router) {
		// protected views by jwt
		app.Use(middles.AuthMiddleware)

		app.Route("/user/{id}", func(app chi.Router) {
			app.Get("/", views.UserRead)
			app.Put("/", views.UserUpdate)
			app.Delete("/", views.UserSoftDelete)
		})

		app.Post("/clubs", views.ClubCreate)
		app.Route("/club/{id}", func(app chi.Router) {
			app.Get("/", views.ClubRead)
		})
	})

	err := models.Migrate(conn)
	if err != nil {
		panic(err)
	}

	http.ListenAndServe(":3000", app)
}
