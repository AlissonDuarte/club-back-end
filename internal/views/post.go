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
	userIDStr := chi.URLParam(app, "id")
	userID, err := strconv.Atoi(userIDStr)
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

	uploadDir := "./uploads"
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

	postContent := app.FormValue("content")
	postTitle := app.FormValue("title")

	newPost := models.NewPost(
		postTitle,
		postContent,
		uint(userIDUint),
		conn,
	)

	if err := conn.Create(&newPost).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	upload := models.UserUploadPost{
		UserID:   uint(userID),
		FilePath: filePath.Name(),
		FileSize: fileData.Size,
		PostID:   newPost.ID,
	}

	if err := db.Create(&upload).Error; err != nil {
		http.Error(w, fmt.Sprintf("Error saving file to database: %s", err.Error()), http.StatusInternalServerError)
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
