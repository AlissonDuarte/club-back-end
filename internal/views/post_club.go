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
	"time"

	"github.com/go-chi/chi/v5"
)

func PostClubCreate(w http.ResponseWriter, app *http.Request) {
	postContent := app.FormValue("content")
	postTitle := app.FormValue("title")
	userIDStr := app.FormValue("userID")
	clubIDStr := chi.URLParam(app, "id")

	if !functions.PostContentMaxLength(postContent, 516) || !functions.PostContentMaxLength(postTitle, 128) {
		http.Error(w, "Post content is to big", http.StatusBadRequest)
		return
	}

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
		PostID int `json:"postID"`
	}

	var userIDRequest UserIDRequest
	clubIDStr := chi.URLParam(app, "id")

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
	postID := userIDRequest.PostID
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

func PostClubUpdate(w http.ResponseWriter, app *http.Request) {
	postContent := app.FormValue("content")
	postTitle := app.FormValue("title")
	postIDstr := app.FormValue("postID")
	userIDstr := app.FormValue("userID")
	clubIDStr := chi.URLParam(app, "id")

	if !functions.PostContentMaxLength(postContent, 516) || !functions.PostContentMaxLength(postTitle, 128) {
		http.Error(w, "Post content is to big", http.StatusBadRequest)
		return
	}

	userIDUint, err := strconv.ParseUint(userIDstr, 10, 64)

	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	postIDUint, err := strconv.ParseUint(postIDstr, 10, 64)

	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	clubID, err := strconv.Atoi(clubIDStr)

	if err != nil {
		http.Error(w, "Invalid club ID format", http.StatusBadRequest)
		return
	}

	conn := database.NewDb()

	file, fileData, err := app.FormFile("file")

	if err != nil {
		http.Error(w, fmt.Sprintf("Error to get file: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	defer file.Close()

	uploadDir := "./uploads/clubs/post"

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, fmt.Sprintf("Error to find upload directory: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	filePath, err := os.OpenFile(filepath.Join(uploadDir, functions.GenerateKeys(16)+".png"), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error to create file: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer filePath.Close()

	_, err = io.Copy(filePath, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error copying file: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	upload := models.UserUploadPost{
		UserID:   uint(userIDUint),
		FilePath: filePath.Name(),
		FileSize: fileData.Size,
	}

	if err := conn.Create(&upload).Error; err != nil {
		http.Error(w, fmt.Sprintf("Error while saving file to database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	updatedPost, err := models.PostClubGetByID(
		conn,
		uint(postIDUint),
		uint(clubID),
	)

	if err != nil {
		http.Error(w, "Cannot retrive this post data", http.StatusInternalServerError)
	}

	if updatedPost.UserID != uint(userIDUint) {
		http.Error(w, "You cannot update this post", http.StatusUnauthorized)
		return
	}
	oldImagepath := updatedPost.Image.FilePath
	oldImageId := updatedPost.ImageID

	updatedPost.Content = postContent
	updatedPost.Title = postTitle
	updatedPost.ImageID = upload.ID
	updatedPost.Image = upload
	updatedPost.Updated = true
	updatedPost.UpdatedAt = time.Now()

	err = updatedPost.Update(conn)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error to update info %s", err.Error()), http.StatusInternalServerError)
		return
	}

	os.Remove(oldImagepath)
	oldImage, _ := models.GetUserUploadPostByID(conn, oldImageId)

	if err := conn.Delete(&oldImage).Error; err != nil {
		http.Error(w, fmt.Sprintf("Error to delete this post: %s", err.Error()), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)

}
