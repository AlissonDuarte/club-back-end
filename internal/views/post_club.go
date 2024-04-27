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
	db := database.NewDb()

	allowed, err := models.IsUserIDInClub(db, uint(userID), uint(clubID))

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Error to check if you're member of this club", http.StatusInternalServerError)
		return
	}

	if !allowed {
		http.Error(w, "You're not member of this club", http.StatusUnauthorized)
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

	upload := models.UserUploadPost{
		UserID:   uint(userID),
		FilePath: filePath,
		FileSize: fileHeader.Size,
	}

	if err := db.Create(&upload).Error; err != nil {
		http.Error(w, fmt.Sprintf("Error while saving file to database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

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
	userID, err := functions.UserIdFromToken(app)

	if err != nil {
		http.Error(w, fmt.Sprintf("User not found: %s", err.Error()), http.StatusForbidden)
		return
	}

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
	allowed, err := models.IsUserIDInClub(db, uint(userID), uint(clubID))

	if err != nil {
		http.Error(w, "Error to check if you're member of this club", http.StatusInternalServerError)
		return
	}

	if !allowed {
		http.Error(w, "You're not member of this club", http.StatusUnauthorized)
		return
	}

	post, err := models.PostClubGetByID(db, uint(postID), uint(clubID))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)

}

func PostClubDelete(w http.ResponseWriter, app *http.Request) {
	type UserIDRequest struct {
		UserID int `json:"userID"`
	}

	var userIDRequest UserIDRequest
	postIDStr := chi.URLParam(app, "postID")
	clubIDStr := chi.URLParam(app, "id")

	postID, err := strconv.Atoi(postIDStr)

	if err != nil {
		http.Error(w, "Invalid post id format", http.StatusBadRequest)
		return
	}

	clubID, err := strconv.Atoi(clubIDStr)

	if err != nil {
		http.Error(w, "Invalid club id format", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(app.Body).Decode(&userIDRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := userIDRequest.UserID
	db := database.NewDb()

	// verificar se o user é criador do post ou dono do club
	post, err := models.PostClubGetByID(db, uint(postID), uint(clubID))

	if err != nil {
		http.Error(w, "Error to get club, try later", http.StatusInternalServerError)
		return
	}

	if userID != int(post.UserID) || clubID != int(post.ClubID) {
		http.Error(w, "You're not the owner of this post nor the owner of the club", http.StatusUnauthorized)
		return
	}

	// verificar se o usuario ainda faz parte do clube
	allowed, err := models.IsUserIDInClub(db, uint(userID), uint(clubID))

	if err != nil {
		http.Error(w, "Error to verify if you're a member of this club", http.StatusInternalServerError)
		return
	}

	if !allowed {
		http.Error(w, "You're not allowed to delete this post", http.StatusInternalServerError)
	}
	if err := db.Delete(&post).Error; err != nil {
		http.Error(w, fmt.Sprintf("Error to delete this post: %s", err.Error()), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	response.Code = http.StatusNoContent
	response.Message = "Post deleted!"
	response.Status = "ok"

	jsonData, err := json.Marshal(response)

	if err != nil {
		http.Error(w, "Error to retrive data", http.StatusInternalServerError)

	}
	w.Write(jsonData)
}
