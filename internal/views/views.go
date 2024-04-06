package views

import (
	"fmt"
	"time"
	"strconv"
    "net/http"
	"html/template"
    "encoding/json"
	"gorm.io/gorm"
	"clube/internal/serializer"
	"clube/internal/functions"
	"clube/infraestructure/models"
	"clube/infraestructure/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

var response struct {
	Message string `json:"message"`
}

func UserCreate(w http.ResponseWriter, app *http.Request) {
	conn := database.NewDb()
	var userData serializer.UserSerializer

	fmt.Println(app.Body)
	
	err := json.NewDecoder(app.Body).Decode(&userData)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error aqui1")
		return
	}

    if err := validate.Struct(userData); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }


	if !functions.PhoneCheck(userData.Phone) {
		http.Error(w, "Invalid phone", http.StatusBadRequest)
		return
	}

	newUser := models.NewUser(
		userData.Name, userData.Username, 
		userData.Gender, userData.BirthDate, 
		userData.Passwd,  userData.Email,
		userData.Phone,
	)

	err = newUser.Save(conn)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error aqui2")
		return
	}

	response.Message = "User created successfully!!"
	jsonData, err := json.Marshal(response)

	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error aqui3")
        return
    }
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jsonData)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error aqui4")
        return
    }
}

func UserRead(w http.ResponseWriter, app *http.Request) {
    userIDStr := chi.URLParam(app, "id")
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    db := database.NewDb()

    user, err := models.UserGetById(db, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func UserUpdate(w http.ResponseWriter, app *http.Request) {
	userIDStr := chi.URLParam(app, "id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	db := database.NewDb()

	user, err := models.UserGetById(db, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userData serializer.UserSerializer

	err = json.NewDecoder(app.Body).Decode(&userData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(userData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.Name = userData.Name
	user.Username = userData.Username
	user.Gender = userData.Gender
	user.BirthDate = userData.BirthDate
	user.Email = userData.Email
	user.Phone = userData.Phone

	if !functions.PhoneCheck(user.Phone) {
		http.Error(w, "Invalid phone number", http.StatusBadRequest)
		return
	}

	err = user.Update(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Message = "User updated successfully!!"
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func UserSoftDelete(w http.ResponseWriter, app *http.Request) {
    userIDStr := chi.URLParam(app, "id")
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    db := database.NewDb()

    user, err := models.UserGetById(db, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    user.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}

    if err := user.Update(db); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    response.Message = "User deleted successfully!!"

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func SignUp(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("templates/signup/signup.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    err = tmpl.Execute(w, nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}
