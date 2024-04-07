package main

import (
    "net/http"
	"clube/internal/views"
	"clube/infraestructure/database"
	"clube/infraestructure/models"
    "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)


func main() {
	// Call the function
	conn := database.NewDb()
	app := chi.NewRouter()
	app.Use(JWTMiddleware)
	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)
	app.Use(middleware.Throttle(1000))
	
	fileServer := http.FileServer(http.Dir("./templates/imgs"))
	htmlServer := http.FileServer(http.Dir("./templates"))

    app.Handle("/*", http.StripPrefix("/", fileServer))
	app.Handle("/*", http.StripPrefix("/", htmlServer))

	app.Get("/home", views.Home)
	
	app.Post("/users", views.UserCreate)

	app.Route("/user/{id}", func(app chi.Router) {
		app.Get("/", views.UserRead)
		app.Put("/", views.UserUpdate)
		app.Delete("/", views.UserSoftDelete)
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