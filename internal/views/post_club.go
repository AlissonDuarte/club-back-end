package views

import (
	"clube/infraestructure/database"
	"clube/infraestructure/models"
	"clube/internal/functions"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func PostClubCreate(w http.ResponseWriter, app *http.Request) {
	postContent := app.FormValue("content")
	postTitle := app.FormValue("title")
	userIDStr := app.FormValue("userID")
	clubIDStr := chi.URLParam(app, "id")

	userID, err := strconv.Atoi(userIDStr)

	if err != nil {
		http.Error(w, "Invalid user id format", http.StatusBadRequest)
		return
	}

	clubID, err := strconv.Atoi(clubIDStr)

	if err != nil {
		http.Error(w, "Invalid club id format", http.StatusBadRequest)
		return
	}

	if err := app.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %s", err.Error()), http.StatusBadRequest)
		return
	}

	file, fileHeader, err := app.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting file: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	uploadDir := "./uploads/clubs/post"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, fmt.Sprintf("Error creating upload directory: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	filePath := filepath.Join(uploadDir, functions.GenerateKeys(16)+".png")
	outFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating file: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	// Copiar o conteúdo do arquivo recebido para o arquivo que estamos criando
	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error copying file: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	db := database.NewDb()

	upload := models.UserUploadPost{
		UserID:   uint(userID),
		FilePath: filePath,
		FileSize: fileHeader.Size,
	}

	if err := db.Create(&upload).Error; err != nil {
		http.Error(w, fmt.Sprintf("Error while saving file to database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	fmt.Println(postTitle, postContent, userID, upload.ID, clubID)
	newPost := models.NewPostClub(
		postTitle,
		postContent,
		uint(userID),
		upload.ID,
		uint(clubID),
		db,
	)

	if err := db.Create(&newPost).Error; err != nil {
		http.Error(w, fmt.Sprintf("Error creating post: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Post created!",
		"status":  "ok",
	})

}

func PostClubRead(w http.ResponseWriter, app *http.Request) {
	postIDStr := chi.URLParam(app, "postID")
	clubIDStr := chi.URLParam(app, "id")

	fmt.Println(postIDStr, clubIDStr)

	postID, err := strconv.Atoi(postIDStr)

	if err != nil {
		http.Error(w, "Invalid post ID format", http.StatusBadRequest)
		return
	}

	clubID, err := strconv.Atoi(clubIDStr)

	if err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	db := database.NewDb()

	post, err := models.PostClubGetByID(db, uint(postID), uint(clubID))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)

}
