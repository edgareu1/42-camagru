package models

import (
	"backend/utils"
	"time"
)

type Comment struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ImageID   int       `json:"image_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func GetCommentsCount(imageId int) (int, error) {
	db := utils.GetDB()
	row := db.QueryRow(`
			SELECT COUNT(*)
			FROM comments
			WHERE image_id = $1
		`, imageId)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

type CommentDetail struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	ImageID   int       `json:"image_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func GetComments(imageId int) ([]CommentDetail, error) {
	db := utils.GetDB()
	rows, err := db.Query(`
			SELECT id, user_id, image_id, content, created_at
			FROM comments
			WHERE image_id = $1
			ORDER BY created_at DESC
		`, imageId)
	if err != nil {
		return nil, err
	}

	comments := []CommentDetail{}
	for rows.Next() {
		var comment Comment
		err = rows.Scan(&comment.ID, &comment.UserID, &comment.ImageID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}

		user, err := GetUser(comment.UserID)
		if err != nil {
			return nil, err
		}

		commentDetail := CommentDetail{
			ID:        comment.ID,
			UserID:    comment.UserID,
			Username:  user.Username,
			ImageID:   comment.ImageID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
		}

		comments = append(comments, commentDetail)
	}

	return comments, nil
}

func CreateComment(userID, imageID int, content string) error {
	db := utils.GetDB()
	_, err := db.Exec(`
			INSERT INTO comments (user_id, image_id, content)
			VALUES ($1, $2, $3)
		`, userID, imageID, content)

	return err
}

func SendCommentNotification(fromUserID, toUserID int, content string) error {
	fromUser, err := GetUser(fromUserID)
	if err != nil {
		return err
	}

	toUser, err := GetUser(toUserID)
	if err != nil {
		return err
	}

	if !toUser.ReceiveNotifications {
		return nil
	}

	emailSubject := "New Comment on Image"
	emailContent := "The user " + fromUser.Username + " has commented on your image:\n\n" + content

	utils.SendEmail(toUser.Email, emailSubject, emailContent)

	return nil
}
