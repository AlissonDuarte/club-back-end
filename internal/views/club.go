package views

import (
	"clube/infraestructure/database"
	"clube/infraestructure/models"
	"clube/internal/serializer"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func ClubCreate(w http.ResponseWriter, r *http.Request) {
	conn := database.NewDb()

	var clubData serializer.GroupSerializer
	err := json.NewDecoder(r.Body).Decode(&clubData)
	if err != nil {
		fmt.Println("Error aqui")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := validate.Struct(clubData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newClub := models.NewClub(
		clubData.Name,
		clubData.Description,
		clubData.Users,
		clubData.Owner,
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
