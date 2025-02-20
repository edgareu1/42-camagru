package models

import (
	"backend/utils"
)

func GetLikesCount(imageId int) (int, error) {
	db := utils.GetDB()
	row := db.QueryRow(`
			SELECT COUNT(*)
			FROM user_image_likes
			WHERE image_id = $1
		`, imageId)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func LikeImage(userID int, imageId int) error {
	db := utils.GetDB()
	_, err := db.Exec(`
			INSERT INTO user_image_likes (user_id, image_id)
			VALUES ($1, $2)
		`, userID, imageId)

	return err
}

func UnlikeImage(userID int, imageId int) error {
	db := utils.GetDB()
	_, err := db.Exec(`
			DELETE FROM user_image_likes 
			WHERE user_id = $1 AND image_id = $2
		`, userID, imageId)

	return err
}

func WasImageLiked(userID int, imageId int) (bool, error) {
	db := utils.GetDB()
	row := db.QueryRow(`
			SELECT EXISTS (
				SELECT 1
				FROM user_image_likes
				WHERE user_id = $1 AND image_id = $2
			)
		`, userID, imageId)
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
