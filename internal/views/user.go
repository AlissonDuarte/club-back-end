package views

import (
	"clube/infraestructure/database"
	"clube/infraestructure/models"
	"clube/internal/functions"
	"clube/internal/serializer"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	//	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Teste(w http.ResponseWriter, app *http.Request) {
	// endpoint criado para testes e debug de funcoes
	type teste struct {
		Password string `json:"password"`
	}
	var body teste

	err := json.NewDecoder(app.Body).Decode(&body)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	hash, err := models.GeneratePasswordHash(body.Password)

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("hash gerado", hash)

	err = functions.VerifyPassword(body.Password, hash)

	if err != nil {
		fmt.Println(err.Error())
	}

}

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
		response.Message = "Cannot validated your data"
		response.Status = "Failed"
		response.Code = http.StatusBadRequest
		response.Data = err.Error()
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

	if err := validate.Struct(userData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !functions.ValidGender(userData.Gender) {
		http.Error(w, "Invalid Gender", http.StatusBadRequest)
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

	userID, err := newUser.Save(conn)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Message = "User created successfully!!"
	response.Status = "success"
	response.Code = http.StatusCreated
	response.Data = map[string]interface{}{
		"userID": userID,
	}

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

func UserProfilePicture(w http.ResponseWriter, app *http.Request) {
	userIDStr := chi.URLParam(app, "id")
	userID, err := strconv.Atoi(userIDStr)

	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	db := database.NewDb()

	user_picture, err := models.GetUserUploadByUserID(db, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user_picture.FilePath == "" {
		http.Error(w, "No profile picture found", http.StatusNotFound)
		return
	}

	file, err := os.Open(user_picture.FilePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	w.Header().Set("Content-Type", "image/jpeg")
	io.Copy(w, file)
}

func UserPostsPictures(w http.ResponseWriter, app *http.Request) {
	userIDStr := chi.URLParam(app, "id")
	postIDStr := chi.URLParam(app, "imageID")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	db := database.NewDb()

	post_picture, err := models.GetPostUploadByPostID(db, uint(postID), uint(userID))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if post_picture.FilePath == "" {
		http.Error(w, "No post picture found", http.StatusNotFound)
		return
	}

	file, err := os.Open(post_picture.FilePath)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "image/png")
	io.Copy(w, file)
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

	var userData serializer.UserUpdateSerializer

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
	user.Bio = userData.Bio
	user.UpdatedAt = time.Now()

	if !functions.PhoneCheck(user.Phone) {
		http.Error(w, "Invalid phone number", http.StatusBadRequest)
		return
	}

	err = user.Update(db, "")
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

func UserChangePassword(w http.ResponseWriter, app *http.Request) {
	var body serializer.UserChangePasswordSerliazer

	err := json.NewDecoder(app.Body).Decode(&body)

	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot decode data due to: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	userID, err := strconv.Atoi(body.Id)
	if err != nil {
		http.Error(w, "Error to read user id", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.NewPasswd != body.NewPasswdCheck {
		http.Error(w, "New password not match", http.StatusBadRequest)
		return
	}

	db := database.NewDb()

	user, err := models.UserGetById(db, userID)

	if err != nil {
		http.Error(w, "Cannot find user", http.StatusInternalServerError)
		return
	}
	old_passwd_hash, err := models.UserGetPassword(db, userID)
	fmt.Println("old password: ", old_passwd_hash)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error to confirm old password: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	err = functions.VerifyPassword(body.OldPasswd, old_passwd_hash)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error %s", err.Error()), http.StatusBadRequest)
		return
	}

	err = user.ChangePassword(db, body.NewPasswd)

	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot change password due to: %s", err.Error()), http.StatusInternalServerError)
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

	if err := user.Update(db, ""); err != nil {
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
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(userLoginData); err != nil {
		http.Error(w, "Invalid data format", http.StatusBadRequest)
		return
	}

	conn := database.NewDb()

	user, err := functions.FindUserByEmail(conn, userLoginData.Email)
	if err != nil {
		http.Error(w, "User with email not found: "+userLoginData.Email, http.StatusBadRequest)
		return
	}

	fmt.Println(userLoginData.Passwd, user.PasswdHash)
	err = functions.VerifyPassword(userLoginData.Passwd, user.PasswdHash)
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	userJWT, err := functions.GenerateJWT(int(user.ID))
	if err != nil {
		http.Error(w, "Cannot generate jwt token", http.StatusInternalServerError)
		return
	}

	response := struct {
		Message string      `json:"message"`
		Status  string      `json:"status"`
		Code    int         `json:"code"`
		Data    interface{} `json:"data,omitempty"`
	}{
		Message: "User logged in successfully!!",
		Status:  "success",
		Code:    http.StatusOK,
		Data: map[string]interface{}{
			"userID":  user.ID,
			"userJWT": userJWT,
		},
	}
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

func UserUploadProfilePicture(w http.ResponseWriter, app *http.Request) {
	userIDStr := chi.URLParam(app, "id")
	userID, err := strconv.Atoi(userIDStr)

	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	var errorMessage string

	if app.Method != "POST" {
		errorMessage = "Method Not Allowed"
		errorMessage = fmt.Sprintf("Method Not Allowed: %s", app.Method)
		http.Error(w, errorMessage, http.StatusMethodNotAllowed)
		return
	}

	err = app.ParseMultipartForm(10 << 20) // Limite de 10 MB
	if err != nil {
		errorMessage = "Error to read file"
		errorMessage = fmt.Sprintf("Error to read file: %s", err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	// Obtém o arquivo enviado
	file, fileData, err := app.FormFile("file")
	if err != nil {
		errorMessage = "Error to get file"
		errorMessage = fmt.Sprintf("Error to get file: %s", err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Diretório onde os arquivos serão salvos
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, fmt.Sprintf("Error to create upload directory: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Cria um novo arquivo no servidor
	filePath, err := os.OpenFile(filepath.Join(uploadDir, functions.GenerateKeys(16)+".png"), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		errorMessage = "Error to create file"
		errorMessage = fmt.Sprintf("Error to create file: %s", err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}
	defer filePath.Close()

	// Copia o conteúdo do arquivo recebido para o novo arquivo no servidor
	_, err = io.Copy(filePath, file)
	if err != nil {
		errorMessage = "Error to copy file content"
		errorMessage = fmt.Sprintf("Error to copy file content: %s", err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	// Salva as informações do arquivo no banco de dados
	db := database.NewDb()

	upload := models.UserUpload{
		UserID:      userID,
		FilePath:    filePath.Name(),
		FileSize:    fileData.Size,
		ContentType: fileData.Header.Get("Content-Type"),
	}

	if err := db.Create(&upload).Error; err != nil {
		http.Error(w, fmt.Sprintf("Error saving file to database: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	profilePictureID := upload.ID

	// Atualiza o ID da imagem de perfil do usuário
	user := models.User{}
	if err := db.First(&user, userID).Error; err != nil {
		http.Error(w, fmt.Sprintf("Error to find user: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	user.ProfilePictureID = profilePictureID
	if err := db.Save(&user).Error; err != nil {
		http.Error(w, fmt.Sprintf("Error to update user: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Responde com uma mensagem de sucesso
	fmt.Fprintf(w, "Upload %s with success", fileData.Filename)
}

func UserFollow(w http.ResponseWriter, app *http.Request) {

	var followData serializer.FollowAndUnfollowSerializer

	conn := database.NewDb()

	userIDStr := chi.URLParam(app, "id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	err = json.NewDecoder(app.Body).Decode(&followData)

	if err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	user, err := models.UserGetById(conn, userID)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	follow, err := models.UserGetById(conn, int(followData.FollowedID))

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	err = models.Follow(conn, user.ID, follow.ID)

	if err != nil {
		http.Error(w, "Error to follow user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

func UserUnfollow(w http.ResponseWriter, app *http.Request) {

	var followData serializer.FollowAndUnfollowSerializer

	conn := database.NewDb()

	userIDStr := chi.URLParam(app, "id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	err = json.NewDecoder(app.Body).Decode(&followData)

	if err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	user, err := models.UserGetById(conn, userID)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	follow, err := models.UserGetById(conn, int(followData.FollowedID))

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	err = models.Unfollow(conn, user.ID, follow.ID)

	if err != nil {
		http.Error(w, "Error to unfollow user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

func UserGetFollowers(w http.ResponseWriter, app *http.Request) {

	conn := database.NewDb()

	userIDStr := chi.URLParam(app, "id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := models.UserGetById(conn, userID)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	followers, err := models.GetFollowers(conn, user.ID)

	if err != nil {
		http.Error(w, "Error to get followers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)

}
func UserGetFollowing(w http.ResponseWriter, app *http.Request) {

	conn := database.NewDb()

	userIDStr := chi.URLParam(app, "id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := models.UserGetById(conn, userID)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	following, err := models.GetFollowing(conn, user.ID)

	if err != nil {
		http.Error(w, "Error to get following", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(following)

}

func UserFeed(w http.ResponseWriter, app *http.Request) {
	userIDStr := chi.URLParam(app, "id")
	userID, err := strconv.Atoi(userIDStr)

	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	conn := database.NewDb()

	_, err = models.UserGetById(conn, userID)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Obter parâmetros de consulta para paginação
	pageStr := app.URL.Query().Get("page")
	pageSizeStr := app.URL.Query().Get("pageSize")

	var page, pageSize int
	// Definir valores padrão para página e tamanho da página
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

	// Calcular o índice inicial para a consulta
	offset := (page - 1) * pageSize

	posts, err := models.GetFeed(
		conn,
		uint(userID),
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
