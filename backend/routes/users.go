package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"backend/models"
	"backend/utils"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetUsers()
	if err != nil {
		utils.SendError(w, "Invalid users query")
		return
	}

	utils.SendMessage(w, users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, "Invalid user ID")
		return
	}

	user, err := models.GetUser(id)
	if err != nil {
		utils.SendError(w, "Invalid user query")
		return
	}

	utils.SendMessage(w, user)
}

type EditUserInput struct {
	Username             string `json:"username"`
	Email                string `json:"email"`
	ReceiveNotifications bool   `json:"receive-notifications"`
	Password             string `json:"password"`
	ConfirmPassword      string `json:"confirm-password"`
}

func EditUser(w http.ResponseWriter, r *http.Request) {
	var inputUser EditUserInput

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, "Invalid user ID")
		return
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&inputUser)
	if err != nil {
		utils.SendError(w, "Invalid user input")
		return
	}

	user, err := models.GetUser(id)
	if err != nil {
		utils.SendError(w, "Invalid user query")
		return
	}

	if inputUser.Username+inputUser.Email+inputUser.Password == "" && inputUser.ReceiveNotifications == user.ReceiveNotifications {
		utils.SendError(w, "Nothing to edit")
		return
	}

	if inputUser.Password != "" && inputUser.Password != inputUser.ConfirmPassword {
		utils.SendError(w, "Passwords do not match")
		return
	}

	err = models.EditUser(id, inputUser.Username, inputUser.Email, inputUser.Password, inputUser.ReceiveNotifications)
	if err != nil {
		utils.SendError(w, err.Error())
		return
	}

	utils.SendMessage(w, "User edited successfully")
}

type VerifyUserInput struct {
	VerificationCode string `json:"verification-code"`
}

func VerifyUser(w http.ResponseWriter, r *http.Request) {
	var inputUser VerifyUserInput

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, "Invalid user ID")
		return
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&inputUser)
	if err != nil {
		utils.SendError(w, "Invalid user input")
		return
	}

	err = models.VerifyUser(id, inputUser.VerificationCode)
	if err != nil {
		utils.SendError(w, err.Error())
		return
	}

	utils.SendMessage(w, "User verified successfully")
}

type RenewPasswordInput struct {
	VerificationCode string `json:"verification-code"`
	Password         string `json:"password"`
	ConfirmPassword  string `json:"confirm-password"`
}

func RenewPassword(w http.ResponseWriter, r *http.Request) {
	var inputUser RenewPasswordInput

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, "Invalid user ID")
		return
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&inputUser)
	if err != nil {
		utils.SendError(w, "Invalid input")
		return
	}

	if inputUser.Password != "" && inputUser.Password != inputUser.ConfirmPassword {
		utils.SendError(w, "Passwords do not match")
		return
	}

	err = models.RenewPassword(id, inputUser.VerificationCode, inputUser.Password)
	if err != nil {
		utils.SendError(w, err.Error())
		return
	}

	utils.SendMessage(w, "Password renewed successfully")
}

type CreateUserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var inputUser CreateUserInput

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&inputUser)
	if err != nil {
		utils.SendError(w, "Invalid user input")
		return
	}

	err = models.CreateUser(inputUser.Username, inputUser.Email, inputUser.Password)
	if err != nil {
		utils.SendError(w, err.Error())
		return
	}

	err = models.SendVerificationEmail(inputUser.Username)
	if err != nil {
		utils.SendError(w, err.Error())
		return
	}

	utils.SendMessage(w, "User created successfully")
}

type SignInUserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInUserResponse struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
}

func SignInUser(w http.ResponseWriter, r *http.Request) {
	var inputUser SignInUserInput

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&inputUser)
	if err != nil {
		utils.SendError(w, "Invalid user input")
		return
	}

	user, err := models.GetUserByUsername(inputUser.Username)
	if err != nil {
		utils.SendError(w, "Invalid username")
		return
	}

	if !utils.IsPasswordValid(user.Password, inputUser.Password) {
		utils.SendError(w, "Invalid password")
		return
	}

	if !user.WasEmailVerified {
		utils.SendError(w, "Email not verified")
		return
	}

	authToken, err := utils.GenerateJWT(user.ID)
	if err != nil {
		utils.SendError(w, "Something went wrong")
		return
	}

	response := SignInUserResponse{
		ID:    user.ID,
		Token: authToken,
	}
	utils.SendMessage(w, response)
}

func AuthUser(w http.ResponseWriter, r *http.Request) {
	if utils.IsRequestAuthenticated(r) {
		utils.SendMessage(w, "Authorized")
	} else {
		utils.SendError(w, "Unauthorized")
	}
}

type RequestPasswordChangeRequest struct {
	Username string `json:"username"`
}

type RequestPasswordChangeResponse struct {
	RenewPasswordURL string `json:"renew-password-url"`
}

func RequestPasswordChange(w http.ResponseWriter, r *http.Request) {
	var inputUser RequestPasswordChangeRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&inputUser)
	if err != nil {
		utils.SendError(w, "Invalid input")
		return
	}

	user, err := models.GetUserByUsername(inputUser.Username)
	if err != nil {
		utils.SendError(w, "Invalid username")
		return
	}

	RenewPasswordURL, err := models.RequestPasswordChange(user.ID, user.Email)
	if err != nil {
		utils.SendError(w, err.Error())
		return
	}

	response := RequestPasswordChangeResponse{
		RenewPasswordURL: RenewPasswordURL,
	}
	utils.SendMessage(w, response)
}
