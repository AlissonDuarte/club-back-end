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
	"unicode/utf8"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func ClubCreate(w http.ResponseWriter, r *http.Request) {
	clubDescription := r.FormValue("description")
	clubName := r.FormValue("name")
	ownerIDstr := r.FormValue("owner")

	ownerID, err := strconv.Atoi(ownerIDstr)
	if err != nil {
		http.Error(w, "Invalid owner ID format", http.StatusBadRequest)
		return
	}
	nameLenght := utf8.RuneCountInString(clubName)
	if nameLenght > 70 {
		errorMessage := fmt.Sprintf("Error: Max lenght of comments is 70, yours is: %s", strconv.Itoa(nameLenght))
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting file: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	uploadDir := "./uploads/club"
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

	conn := database.NewDb()

	upload := models.UserUploadClub{
		UserID:   uint(ownerID),
		FilePath: filePath,
		FileSize: fileHeader.Size,
	}

	if err := conn.Create(&upload).Error; err != nil {
		http.Error(w, fmt.Sprintf("Error to save file into database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	newClub := models.NewClub(
		clubName,
		clubDescription,
		[]int{ownerID},
		uint(ownerID),
		upload.ID,
		conn,
	)

	if err := conn.Create(&newClub).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Message = "Club creted successfully!!"
	response.Status = "success"
	response.Code = http.StatusCreated
	response.Data = map[string]interface{}{
		"clubID": newClub.ID,
	}

	jsonData, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(jsonData)
	w.WriteHeader(http.StatusCreated)
}

func ClubRead(w http.ResponseWriter, app *http.Request) {
	clubIDStr := chi.URLParam(app, "id")
	clubID, err := strconv.Atoi(clubIDStr)

	if err != nil {
		http.Error(w, "Invalid club ID", http.StatusBadRequest)
		return
	}

	db := database.NewDb()

	club, err := models.ClubGetById(db, clubID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(club)
}

func ClubReadAll(w http.ResponseWriter, app *http.Request) {
	userIDStr := chi.URLParam(app, "userId")
	userID, err := strconv.Atoi(userIDStr)

	var clubs []*models.Club

	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot read id value due to: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	conn := database.NewDb()

	err = conn.Where("owner_id = ?", userID).
		Preload("OwnerRefer", func(db *gorm.DB) *gorm.DB {
			return db.Omit("passwd_hash")
		}).
		Preload("Users", func(db *gorm.DB) *gorm.DB {
			return db.Omit("passwd_hash")
		}).
		Find(&clubs).Error

	if err != nil {
		http.Error(w, fmt.Sprintf("Error to get data: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clubs)

}

func ClubUpdate(w http.ResponseWriter, r *http.Request) {
	clubIDStr := chi.URLParam(r, "id")
	clubID, err := strconv.Atoi(clubIDStr)

	if err != nil {
		http.Error(w, "Invalid club ID", http.StatusBadRequest)
		return
	}

	conn := database.NewDb()

	var clubData serializer.GroupSerializer
	err = json.NewDecoder(r.Body).Decode(&clubData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := validate.Struct(clubData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	club, err := models.ClubGetById(conn, clubID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	club.Name = clubData.Name
	club.Description = clubData.Description

	if err := conn.Save(club).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Message = "Club updated successfully!!"
	response.Status = "success"
	response.Code = http.StatusOK
	response.Data = map[string]interface{}{
		"clubID": club.ID,
	}

	jsonData, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(jsonData)
	w.WriteHeader(http.StatusOK)
}

func ClubSoftDelete(w http.ResponseWriter, r *http.Request) {
	clubIDStr := chi.URLParam(r, "id")
	clubID, err := strconv.Atoi(clubIDStr)

	if err != nil {
		http.Error(w, "Invalid club ID", http.StatusBadRequest)
		return
	}

	conn := database.NewDb()

	club, err := models.ClubGetById(conn, clubID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := conn.Delete(club).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Message = "Club deleted successfully!!"
	response.Status = "success"
	response.Code = http.StatusOK
	response.Data = map[string]interface{}{
		"clubID": club.ID,
	}

	jsonData, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(jsonData)
	w.WriteHeader(http.StatusOK)
}

func ClubFeed(w http.ResponseWriter, app *http.Request) {
	clubIDStr := chi.URLParam(app, "id")
	clubID, err := strconv.Atoi(clubIDStr)

	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot format club id: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	userID, err := functions.UserIdFromToken(app)

	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot decode token: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	conn := database.NewDb()

	allowed, err := models.IsUserIDInClub(conn, uint(userID), uint(clubID))

	if err != nil {
		http.Error(w, "Cannot check if you're member of this group, try later!", http.StatusInternalServerError)
		return
	}

	if !allowed {
		http.Error(w, "You're not a member of this group.", http.StatusUnauthorized)
		return
	}

	pageStr := app.URL.Query().Get("page")
	pageSizeStr := app.URL.Query().Get("pageSize")

	var page, pageSize int

	if pageStr == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(pageStr)

		if err != nil || page < 1 {
			http.Error(w, "Invalid page number", http.StatusBadRequest)
			return
		}
	}

	if pageSizeStr == "" {
		pageSize = 2
	} else {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil || pageSize < 1 {
			http.Error(w, "Invalid page size", http.StatusBadRequest)
			return
		}
	}

	offset := (page - 1) * pageSize

	posts, err := models.GetClubFeed(
		conn,
		uint(clubID),
		offset,
		pageSize,
	)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error to get feed due to: %s", err.Error()), http.StatusInternalServerError)
		return

	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func ClubPictures(w http.ResponseWriter, app *http.Request) {
	clubIDstr := chi.URLParam(app, "id")

	clubID, err := strconv.Atoi(clubIDstr)

	if err != nil {
		http.Error(w, "Invalid club id format", http.StatusInternalServerError)
		return
	}

	conn := database.NewDb()

	club_picture, err := models.GetClubUploadByID(conn, uint(clubID))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if club_picture.FilePath == "" {
		http.Error(w, "No post picture found", http.StatusNotFound)
		return
	}

	file, err := os.Open(club_picture.FilePath)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "image/png")
	io.Copy(w, file)
}
