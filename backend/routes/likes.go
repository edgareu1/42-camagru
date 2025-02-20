package routes

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"backend/models"
	"backend/utils"
)

func LikeImage(w http.ResponseWriter, r *http.Request) {
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

	imageDetails, err := models.GetImageDetails(imageId)
	if err != nil {
		utils.SendError(w, "Invalid image ID")
		return
	}
	if imageDetails.UserID == userID {
		utils.SendError(w, "Cannot like own image")
		return
	}

	err = models.LikeImage(userID, imageId)
	if err != nil {
		utils.SendError(w, err.Error())
		return
	}

	utils.SendMessage(w, "Image liked sucessfully")
}

func UnlikeImage(w http.ResponseWriter, r *http.Request) {
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

	err = models.UnlikeImage(userID, imageId)
	if err != nil {
		utils.SendError(w, err.Error())
		return
	}

	utils.SendMessage(w, "Image like removed sucessfully")
}
