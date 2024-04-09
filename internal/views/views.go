package views

import (
	"clube/infraestructure/database"
	"clube/infraestructure/models"
	"clube/internal/functions"
	"clube/internal/serializer"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var validate = validator.New()

var response struct {
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

func UserCreate(w http.ResponseWriter, app *http.Request) {
	conn := database.NewDb()
	var userData serializer.UserSerializer

	err := json.NewDecoder(app.Body).Decode(&userData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		userData.Passwd, userData.Email,
		userData.Phone,
	)

	err = newUser.Save(conn)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Message = "User created successfully!!"
	response.Status = "success"
	response.Code = http.StatusCreated

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
	response.Status = "success"
	response.Code = http.StatusOK

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
	response.Status = "success"
	response.Code = http.StatusOK

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/home/home.html")
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

func UserLogin(w http.ResponseWriter, r *http.Request) {

	var userLoginData serializer.UserLoginSerializer

	err := json.NewDecoder(r.Body).Decode(&userLoginData)

	if err != nil {
		response.Message = "Invalid JSON data"
		response.Status = "error"
		response.Code = http.StatusBadRequest

		jsonData, err := json.Marshal(response)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonData)
		return
	}

	if err := validate.Struct(userLoginData); err != nil {
		response.Message = "Invalid data format"
		response.Status = "error"
		response.Code = http.StatusBadRequest

		jsonData, err := json.Marshal(response)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonData)
		return
	}

	// buscar usuário pelo email filtrando

	conn := database.NewDb()
	user, err := functions.FindUserByEmail(conn, userLoginData.Email)

	if err != nil {

		response.Message = "User with email not found: " + userLoginData.Email
		response.Status = "error"
		response.Code = http.StatusBadRequest

		jsonData, err := json.Marshal(response)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonData)
		return
	}

	// verificar se a senha está correta

	err = functions.VerifyPassword(userLoginData.Passwd, user.PasswdHash)

	if err != nil {
		response.Message = "User logged in successfully!!"
		response.Status = "success"
		response.Code = http.StatusOK

		jsonData, err := json.Marshal(response)

		if err != nil {
			http.Error(w, "Error to format response data, try again later", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		_, err = w.Write(jsonData)

		if err != nil {
			http.Error(w, "Error to response data, try again later", http.StatusInternalServerError)
			return
		}
	}
	userJWT, err := functions.GenerateJWT(int(user.ID))

	if err != nil {
		response.Message = "Cannot generate jwt token"
		response.Status = "error"
		response.Code = http.StatusInternalServerError

		jsonData, err := json.Marshal(response)

		if err != nil {
			http.Error(w, "Error to format data to response, try again later", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(jsonData)

		if err != nil {
			http.Error(w, "Error to respnse data, try again later", http.StatusInternalServerError)
		}
	}

	response.Message = "User logged in successfully!!"
	response.Status = "success"
	response.Code = http.StatusOK

	jsonData, err := json.Marshal(response)

	if err != nil {
		http.Error(w, "Error to format response data, try again later", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer "+userJWT)
	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, "Error to response data, try again later", http.StatusInternalServerError)
		return
	}

}
