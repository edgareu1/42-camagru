package routes

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/nfnt/resize"

	"backend/models"
	"backend/utils"
)

type CreateImageInput struct {
	Image   string `json:"image"`
	Overlay string `json:"overlay"`
}

func CreateImage(w http.ResponseWriter, r *http.Request) {
	var input CreateImageInput

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, "Invalid user ID")
		return
	}

	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr != idStr {
		utils.SendError(w, "Invalid user ID")
		return
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&input)
	if err != nil {
		utils.SendError(w, "Invalid input")
		return
	}

	backgroundData, err := base64.StdEncoding.DecodeString(input.Image)
	if err != nil {
		utils.SendError(w, "Invalid image")
		return
	}

	background, err := png.Decode(bytes.NewReader(backgroundData))
	if err != nil {
		utils.SendError(w, "Invalid image")
		return
	}

	overlayFilepath := fmt.Sprintf("/app/images/%s", input.Overlay)
	overlayFile, err := os.Open(overlayFilepath)
	if err != nil {
		utils.SendError(w, "Invalid overlay image")
		return
	}
	defer overlayFile.Close()

	overlay, err := png.Decode(overlayFile)
	if err != nil {
		utils.SendError(w, "Invalid overlay image")
		return
	}
	overlayResized := resize.Resize(
		uint(background.Bounds().Dx()),
		uint(background.Bounds().Dy()),
		overlay,
		resize.Lanczos3,
	)

	finalImage := image.NewRGBA(background.Bounds())

	draw.Draw(finalImage, background.Bounds(), background, image.Point{}, draw.Over)

	draw.Draw(finalImage, overlayResized.Bounds(), overlayResized, image.Point{}, draw.Over)

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, finalImage, nil)
	if err != nil {
		utils.SendError(w, "Failed to generate final image")
		return
	}

	finalImageData := buf.Bytes()
	err = models.CreateImage(id, finalImageData)
	if err != nil {
		utils.SendError(w, err.Error())
		return
	}

	utils.SendMessage(w, "Image created successfully")
}

type GetImagesResponse struct {
	Images      []models.Image `json:"images"`
	CurrentPage int            `json:"current_page"`
	PageSize    int            `json:"page_size"`
	TotalItems  int            `json:"total_items"`
	TotalPages  int            `json:"total_pages"`
}

func GetImages(w http.ResponseWriter, r *http.Request) {
	imagesPerPage := 12
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	images, err := models.GetImages(page, imagesPerPage)
	if err != nil {
		utils.SendError(w, err.Error())
		return
	}

	numImages, err := models.GetNumImages()
	if err != nil {
		utils.SendError(w, err.Error())
		return
	}

	totalPages := (numImages + imagesPerPage - 1) / imagesPerPage

	response := GetImagesResponse{
		Images:      images,
		CurrentPage: page,
		PageSize:    imagesPerPage,
		TotalItems:  numImages,
		TotalPages:  totalPages,
	}

	utils.SendMessage(w, response)
}

func GetUserImages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, "Invalid user ID")
		return
	}

	images, err := models.GetUserImages(id)
	if err != nil {
		utils.SendError(w, err.Error())
		return
	}

	utils.SendMessage(w, images)
}

func GetImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, "Invalid image ID")
		return
	}

	imageData, err := models.GetImage(id)
	if err != nil {
		utils.SendError(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.WriteHeader(http.StatusOK)
	w.Write(imageData)
}

type GetImageDetailsResponse struct {
	Id          int       `json:"id"`
	UserId      int       `json:"user_id"`
	Username    string    `json:"username"`
	CreatedAt   time.Time `json:"created_at"`
	NumComments int       `json:"num_comments"`
	NumLikes    int       `json:"num_likes"`
	WasLiked    bool      `json:"was_liked"`
}

func GetImageDetails(w http.ResponseWriter, r *http.Request) {
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

	user, err := models.GetUser(imageDetails.UserID)
	if err != nil {
		utils.SendError(w, "Invalid image user ID")
		return
	}

	numComments, err := models.GetCommentsCount(imageId)
	if err != nil {
		utils.SendError(w, err.Error())
		return
	}

	numLikes, err := models.GetLikesCount(imageId)
	if err != nil {
		utils.SendError(w, err.Error())
		return
	}

	wasLiked := false
	userIDStr := r.Header.Get("X-User-ID")
	userID, err := strconv.Atoi(userIDStr)
	if err == nil {
		wasLiked, _ = models.WasImageLiked(userID, imageId)
	}

	response := GetImageDetailsResponse{
		Id:          imageDetails.ID,
		UserId:      user.ID,
		Username:    user.Username,
		CreatedAt:   imageDetails.CreatedAt,
		NumComments: numComments,
		NumLikes:    numLikes,
		WasLiked:    wasLiked,
	}

	utils.SendMessage(w, response)
}

func DeleteImage(w http.ResponseWriter, r *http.Request) {
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
	if imageDetails.UserID != userID {
		utils.SendError(w, "Invalid user ID")
		return
	}

	err = models.DeleteImage(imageId)
	if err != nil {
		utils.SendError(w, err.Error())
	}

	utils.SendMessage(w, "Image deleted successfully")
}
