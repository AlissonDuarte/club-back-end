package views

import (
	"clube/infraestructure/database"
	"clube/infraestructure/models"
	"clube/internal/functions"
	"clube/internal/serializer"
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

func PostCreate(w http.ResponseWriter, r *http.Request) {
	// Obter os dados do formulário
	postContent := r.FormValue("content")
	postTitle := r.FormValue("title")
	userIDStr := r.FormValue("userID")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// ParseMultipartForm analisa o formulário multipart do pedido.
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Obter o arquivo enviado
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting file: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Criar diretório de upload se não existir
	uploadDir := "./uploads/post"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, fmt.Sprintf("Error creating upload directory: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Criar o arquivo no diretório de upload
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

	// Criar registro no banco de dados para o upload do arquivo
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

	// Criar novo post no banco de dados
	newPost := models.NewPost(
		postTitle,
		postContent,
		uint(userID),
		upload.ID,
		db,
	)
	if err := db.Create(&newPost).Error; err != nil {
		http.Error(w, fmt.Sprintf("Error creating post: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Responder com sucesso
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Post created!",
		"status":  "ok",
	})
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

	post, err := models.PostGetByID(db, uint(postID))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func PostDelete(w http.ResponseWriter, app *http.Request) {

	var postDeleteData serializer.PostDeleteSerializer

	if err := json.NewDecoder(app.Body).Decode(&postDeleteData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	postID := postDeleteData.PostID
	userID := postDeleteData.UserID

	db := database.NewDb()

	post, err := models.PostGetByID(db, postID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if post.UserID != uint(userID) {
		http.Error(w, "You can't delete this post", http.StatusForbidden)
		return
	}

	if err := db.Delete(&post).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

func PostUpdate(w http.ResponseWriter, app *http.Request) {
	postContent := app.FormValue("content")
	postTitle := app.FormValue("title")
	postIDstr := app.FormValue("postID")
	userIDstr := app.FormValue("userID")

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
	conn := database.NewDb()

	file, fileData, err := app.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error to get file: %s", err.Error()), http.StatusInternalServerError)
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
		http.Error(w, fmt.Sprintf("Error to create file: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer filePath.Close()

	db := database.NewDb()

	upload := models.UserUploadPost{
		UserID:   uint(userIDUint),
		FilePath: filePath.Name(),
		FileSize: fileData.Size,
	}

	if err := db.Create(&upload).Error; err != nil {
		http.Error(w, fmt.Sprintf("Error while saving file to database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	updatedPost, err := models.PostGetByID(
		conn,
		uint(postIDUint),
	)

	if err != nil {
		http.Error(w, "Cannot retrive this post data", http.StatusInternalServerError)
	}

	if updatedPost.UserID != uint(userIDUint) {
		http.Error(w, "You cannot update this post", http.StatusUnauthorized)
		return
	}

	updatedPost.Content = postContent
	updatedPost.Title = postTitle
	updatedPost.ImageID = upload.ID
	updatedPost.Updated = true
	updatedPost.UpdatedAt = time.Now()

	_, err = updatedPost.Save(conn)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error to update info %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
