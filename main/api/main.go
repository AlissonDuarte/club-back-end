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

	fileServer := http.FileServer(http.Dir("."))

	app.Handle("/files/*", http.StripPrefix("/files/", fileServer))
	app.Handle("/images/*", http.StripPrefix("/images/", fileServer))

	app.Post("/users", views.UserCreate)
	app.Get("/home", views.Home)
	app.Post("/user/login", views.UserLogin)

	app.Group(func(app chi.Router) {
		// protected views by jwt
		app.Use(middles.AuthMiddleware)

		app.Route("/user/{id}", func(app chi.Router) {
			app.Get("/", views.UserRead)
			app.Patch("/", views.UserUpdate)
			app.Delete("/", views.UserSoftDelete)

			app.Post("/follow", views.UserFollow)
			app.Post("/unfollow", views.UserUnfollow)

			app.Get("/followers", views.UserGetFollowers)
			app.Get("/following", views.UserGetFollowing)

			app.Get("/feed", views.UserFeed)

			app.Post("/images/user", views.UserUploadProfilePicture)
			app.Get("/images/user", views.UserProfilePicture)

			app.Get("/images/posts/{imageID}", views.UserPostsPictures)
		})

		app.Route("/post", func(app chi.Router) {
			app.Get("/{id}", views.PostRead)
			app.Post("/", views.PostCreate)
			app.Patch("/", views.PostUpdate)
			app.Delete("/", views.PostDelete)

		})

		app.Post("/clubs", views.ClubCreate)
		app.Route("/club/{id}", func(app chi.Router) {
			app.Get("/", views.ClubRead)
			app.Put("/", views.ClubUpdate)
			app.Delete("/", views.ClubSoftDelete)

			app.Route("/post", func(app chi.Router) {
				app.Post("/", views.PostClubCreate)
				app.Get("/{postID}", views.PostClubRead)
				app.Delete("/{postID}", views.PostClubDelete)
			})
		})

		app.Route("/comment", func(app chi.Router) {
			app.Post("/", views.CommentCreate)
			app.Patch("/", views.CommentUpdate)
			app.Delete("/", views.CommentDelete)
		})
	})

	err := models.Migrate(conn)
	if err != nil {
		panic(err)
	}

	http.ListenAndServe(":3000", app)
}
