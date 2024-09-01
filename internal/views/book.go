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

func BookRead(w http.ResponseWriter, app *http.Request) {
	bookIDString := chi.URLParam(app, "id")
	bookID, err := strconv.Atoi(bookIDString)
	if err != nil {
		http.Error(w, "invalid book 'id' format", http.StatusBadRequest)
		return
	}

	db := database.NewDb()

	book, err := models.BookGetByID(db, uint(bookID))

	if err != nil {
		http.Error(w, fmt.Sprintf("error to get object in database: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// serialized, err := serializer.BookGetSerialize(book)
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("error to get serialized data %s "), err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func BookCreate(w http.ResponseWriter, app *http.Request) {
	name := app.FormValue("name")
	resume := app.FormValue("resume")
	release := app.FormValue("release")
	authorIDstr := app.FormValue("authorID")
	userIDStr := app.FormValue("userID")
	authorID, err := strconv.Atoi(authorIDstr)

	if err != nil {
		http.Error(w, "invalid author id format", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(userIDStr)

	if err != nil {
		http.Error(w, "invalid user id format", http.StatusBadRequest)
		return
	}

	if err := app.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	file, fileHeader, err := app.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting file: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	uploadDir := "./uploads/book"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, fmt.Sprintf("error creating upload directory: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	filePath := filepath.Join(uploadDir, functions.GenerateKeys(16)+".png")
	outFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("error creating file: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("error copying file: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	db := database.NewDb()
	upload := models.UserUpload{
		UserID:      userID,
		FilePath:    filePath,
		FileSize:    fileHeader.Size,
		ContentType: fileHeader.Header.Get("Content-Type"),
	}
	if err := db.Create(&upload).Error; err != nil {
		http.Error(w, fmt.Sprintf("Error while saving file to database: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	// CRIAR O LIVRO
	layout := "2006-01-02 15:04:05"
	parsedTime, err := time.Parse(layout, release)
	if err != nil {
		http.Error(w, fmt.Sprintf("error while convert release time object: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	book := models.NewBook(name, resume, parsedTime, upload.ID, uint(authorID), false, db)

	if err := db.Create(&book).Error; err != nil {
		http.Error(w, fmt.Sprintf("error creating book: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Book created!",
		"status":  "ok",
	})
}
