package models

import (
	"backend/utils"
	"time"
)

type Image struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Data      string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateImage(id int, image []byte) error {
	db := utils.GetDB()
	_, err := db.Exec(`
			INSERT INTO images (data, user_id)
			VALUES ($1, $2)
		`, image, id)

	return err
}

func GetImages(page int, imagesPerPage int) ([]Image, error) {
	offset := (page - 1) * imagesPerPage
	db := utils.GetDB()
	rows, err := db.Query(`
			SELECT id, user_id, data, created_at
			FROM images
			ORDER BY created_at DESC
			LIMIT $1
			OFFSET $2
		`, imagesPerPage, offset)
	if err != nil {
		return nil, err
	}

	images := []Image{}
	for rows.Next() {
		var image Image
		err := rows.Scan(&image.ID, &image.UserID, &image.Data, &image.CreatedAt)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	return images, nil
}

func GetUserImages(id int) ([]Image, error) {
	db := utils.GetDB()
	rows, err := db.Query(`
			SELECT id, user_id, data, created_at
			FROM images
			WHERE user_id = $1
		`, id)
	if err != nil {
		return nil, err
	}

	images := []Image{}
	for rows.Next() {
		var image Image
		err := rows.Scan(&image.ID, &image.UserID, &image.Data, &image.CreatedAt)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	return images, nil
}

func GetImage(id int) ([]byte, error) {
	db := utils.GetDB()
	row := db.QueryRow(`
			SELECT data
			FROM images
			WHERE id = $1
		`, id)

	var imageData []byte
	err := row.Scan(&imageData)
	if err != nil {
		return nil, err
	}

	return imageData, nil
}

func GetImageDetails(id int) (Image, error) {
	db := utils.GetDB()
	row := db.QueryRow(`
			SELECT id, user_id, data, created_at
			from images
			WHERE id = $1
		`, id)

	var image Image
	err := row.Scan(&image.ID, &image.UserID, &image.Data, &image.CreatedAt)
	if err != nil {
		return Image{}, err
	}

	return image, nil
}

func GetNumImages() (int, error) {
	db := utils.GetDB()
	row := db.QueryRow(`
			SELECT COUNT(*)
			FROM images
		`)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func DeleteImage(id int) error {
	db := utils.GetDB()
	_, err := db.Exec(`
			DELETE FROM images
			WHERE id = $1
		`, id)

	return err
}
