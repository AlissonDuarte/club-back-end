package views

import (
	"clube/infraestructure/database"
	"clube/infraestructure/models"
	"clube/internal/functions"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func PostCreate(w http.ResponseWriter, app *http.Request) {
	postContent := app.FormValue("content")
	postTitle := app.FormValue("title")
	userIDStr := app.FormValue("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	userIDUint, err := strconv.ParseUint(userIDStr, 10, 64)

	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	conn := database.NewDb()

	err = app.ParseMultipartForm(10 << 20)
	var errorMessage string

	if err != nil {
		errorMessage = "Error to read file"
		errorMessage = fmt.Sprintf("Error to read file: %s", err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	file, fileData, err := app.FormFile("file")
	if err != nil {
		errorMessage = "Error to get file"
		errorMessage = fmt.Sprintf("Error to get file: %s", err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}
	defer file.Close()

	uploadDir := "./uploads/post"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, fmt.Sprintf("Error to create upload directory: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	filePath, err := os.OpenFile(filepath.Join(uploadDir, functions.GenerateKeys(16)+".png"), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		errorMessage = "Error to create file"
		errorMessage = fmt.Sprintf("Error to create file: %s", err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}
	defer filePath.Close()

	db := database.NewDb()

	upload := models.UserUploadPost{
		UserID:   uint(userID),
		FilePath: filePath.Name(),
		FileSize: fileData.Size,
	}

	if err := db.Create(&upload).Error; err != nil {
		http.Error(w, fmt.Sprintf("Error while saving file to database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	newPost := models.NewPost(
		postTitle,
		postContent,
		uint(userIDUint),
		upload.ID,
		conn,
	)

	if err := conn.Create(&newPost).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	response.Code = http.StatusOK
	response.Message = "Post created!"
	response.Status = "ok"

	jsonData, err := json.Marshal(response)

	if err != nil {
		http.Error(w, "Error to retrive data", http.StatusInternalServerError)

	}
	w.Write(jsonData)
}

// leitura de post pelo ID do post
// para abrir comentários, likes e etc
func PostRead(w http.ResponseWriter, app *http.Request) {
	postIDStr := chi.URLParam(app, "id")
	postID, err := strconv.Atoi(postIDStr)

	if err != nil {
		http.Error(w, "Invalid post ID format", http.StatusBadRequest)
		return
	}

	db := database.NewDb()

	post, err := models.PostGetByID(db, postID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
