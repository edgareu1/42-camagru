package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"backend/models"
	"backend/utils"
)

type CreateCommentInput struct {
	Content string `json:"content"`
}

func GetComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageIdStr := vars["id"]
	imageId, err := strconv.Atoi(imageIdStr)
	if err != nil {
		utils.SendError(w, "Invalid image ID")
		return
	}

	comments, err := models.GetComments(imageId)
	if err != nil {
		utils.SendError(w, err.Error())
		return
	}

	utils.SendMessage(w, comments)
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		utils.SendError(w, "Invalid user ID")
		return
	}

	vars := mux.Vars(r)
	imageIdStr := vars["id"]
	imageId, err := strconv.Atoi(imageIdStr)
	if err != nil {
		utils.SendError(w, "Invalid image ID")
		return
	}

	var input CreateCommentInput
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&input)
	if err != nil {
		utils.SendError(w, "Invalid input")
		return
	}

	err = models.CreateComment(userID, imageId, input.Content)
	if err != nil {
		utils.SendError(w, err.Error())
		return
	}

	imageDetails, err := models.GetImageDetails(imageId)
	if err != nil {
		utils.SendError(w, err.Error())
		return
	}

	err = models.SendCommentNotification(userID, imageDetails.UserID, input.Content)
	if err != nil {
		utils.SendError(w, err.Error())
		return
	}

	utils.SendMessage(w, "Comment created successfully")
}
