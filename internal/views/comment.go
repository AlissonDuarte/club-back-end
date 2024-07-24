package views

import (
	"clube/infraestructure/database"
	"clube/infraestructure/models"
	"clube/internal/serializer"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"
)

func CommentRead(w http.ResponseWriter, app *http.Request) {
	postIDstr := chi.URLParam(app, "id")
	postID, err := strconv.Atoi(postIDstr)

	if err != nil {
		http.Error(w, "Invalid post id", http.StatusBadRequest)
		return
	}

	db := database.NewDb()

	postData, err := models.GetPostComment(db, postID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error to get comments data: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(postData)
}

func CommentCreate(w http.ResponseWriter, app *http.Request) {

	var errorMessage string
	var commentData serializer.CommentSerializer

	err := json.NewDecoder(app.Body).Decode(&commentData)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error to read data %s", err.Error()), http.StatusBadRequest)
		return
	}
	contentLenght := utf8.RuneCountInString(commentData.Content)
	if contentLenght > 500 {
		errorMessage = fmt.Sprintf("Error: Max lenght of comments is 500, yours is: %s", strconv.Itoa(contentLenght))
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	if err := validate.Struct(commentData); err != nil {
		errorMessage = fmt.Sprintf("Error to validate data due to: %s", err.Error())
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	newComment := models.NeWComment(
		uint(commentData.UserID),
		uint(commentData.PostID),
		commentData.Content,
	)

	conn := database.NewDb()
	commentID, err := newComment.Save(conn)

	if err != nil {
		errorMessage = fmt.Sprintf("Error to save comment: %s", err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}
	err = models.AddCommentToPost(conn, uint(commentData.PostID), commentID)
	if err != nil {
		errorMessage = fmt.Sprintf("Error to save comment in post: %s", err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func CommentUpdate(w http.ResponseWriter, app *http.Request) {
	conn := database.NewDb()

	var errorMessage string
	var commentData serializer.CommentUpdateSerializer

	err := json.NewDecoder(app.Body).Decode(&commentData)

	if err != nil {
		http.Error(w, "Error to read data", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(commentData); err != nil {
		errorMessage = fmt.Sprintf("Error to validate data due to: %s", err.Error())
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	comment, err := models.GetCommentByID(
		conn,
		uint(commentData.CommentID),
		uint(commentData.UserID),
		uint(commentData.PostID),
	)

	if err != nil {
		errorMessage = fmt.Sprintf("Error to get comment: %s", err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	if comment == nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	comment.Content = commentData.Content
	comment.Updated = true

	err = comment.Update(conn)
	if err != nil {
		errorMessage = fmt.Sprintf("Error to update comment: %s", err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func CommentDelete(w http.ResponseWriter, app *http.Request) {
	conn := database.NewDb()

	var errorMessage string
	var commentData serializer.CommentUpdateSerializer

	err := json.NewDecoder(app.Body).Decode(&commentData)

	if err != nil {
		http.Error(w, "Error to read data", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(commentData); err != nil {
		errorMessage = fmt.Sprintf("Error to validate data due to: %s", err.Error())
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	comment, err := models.GetCommentByID(
		conn,
		uint(commentData.CommentID),
		uint(commentData.UserID),
		uint(commentData.PostID),
	)

	if err != nil {
		errorMessage = fmt.Sprintf("Error to get comment: %s", err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	if comment == nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	err = conn.Delete(comment).Error
	if err != nil {
		errorMessage = fmt.Sprintf("Error to delete comment: %s", err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
