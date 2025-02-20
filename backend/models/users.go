package models

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"backend/utils"
	"backend/validity"
)

type User struct {
	ID                   int       `json:"id"`
	Username             string    `json:"username"`
	Email                string    `json:"email"`
	Password             string    `json:"-"`
	WasEmailVerified     bool      `json:"was_email_verified"`
	ReceiveNotifications bool      `json:"receive_notifications"`
	VerificationCode     string    `json:"verification_code"`
	CreatedAt            time.Time `json:"created_at"`
}

func GetUsers() ([]User, error) {
	db := utils.GetDB()
	rows, err := db.Query(`
			SELECT id, username, email, password, was_email_verified, receive_notifications, verification_code, created_at
			FROM users
		`)
	if err != nil {
		return nil, err
	}

	users := []User{}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.WasEmailVerified, &user.ReceiveNotifications, &user.VerificationCode, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func GetUser(id int) (User, error) {
	db := utils.GetDB()
	row := db.QueryRow(`
			SELECT id, username, email, password, was_email_verified, receive_notifications, verification_code, created_at
			FROM users
			WHERE id = $1
		`, id)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.WasEmailVerified, &user.ReceiveNotifications, &user.VerificationCode, &user.CreatedAt)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func GetUserByUsername(username string) (User, error) {
	db := utils.GetDB()
	row := db.QueryRow(`
			SELECT id, username, email, password, was_email_verified, receive_notifications, verification_code, created_at
			FROM users
			WHERE username = $1
		`, username)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.WasEmailVerified, &user.ReceiveNotifications, &user.VerificationCode, &user.CreatedAt)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func EditUser(id int, username, email, password string, receiveNotifications bool) error {
	err := validity.ValidateUser(username, email, password, false)
	if err != nil {
		return err
	}

	fieldsToUpdate := map[string]interface{}{}
	if username != "" {
		fieldsToUpdate["username"] = username
	}
	if email != "" {
		fieldsToUpdate["email"] = email
	}
	if password != "" {
		hashPassword, err := utils.HashPassword(password)
		if err != nil {
			return err
		}
		fieldsToUpdate["password"] = hashPassword
	}

	fieldsToUpdate["receive_notifications"] = receiveNotifications

	var setClauses []string
	var params []interface{}
	paramCount := 1

	for field, value := range fieldsToUpdate {
		setClauses = append(setClauses, field+" = $"+fmt.Sprint(paramCount))
		params = append(params, value)
		paramCount++
	}
	params = append(params, id)

	db := utils.GetDB()
	query := fmt.Sprintf(
		"UPDATE users SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "),
		paramCount,
	)
	_, err = db.Exec(query, params...)

	return err
}

func VerifyUser(id int, verificationCode string) error {
	user, err := GetUser(id)
	if err != nil {
		return errors.New("invalid user query")
	}

	if user.VerificationCode != verificationCode {
		return errors.New("invalid verification code")
	}

	db := utils.GetDB()
	_, err = db.Exec(`
		UPDATE users
		SET was_email_verified = $1
		WHERE id = $2
	`, true, id)

	return err
}

func RenewPassword(id int, verificationCode, password string) error {
	user, err := GetUser(id)
	if err != nil {
		return errors.New("invalid user query")
	}

	if user.VerificationCode != verificationCode {
		return errors.New("invalid verification code")
	}

	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	db := utils.GetDB()
	_, err = db.Exec(`
		UPDATE users
		SET password = $1
		WHERE id = $2
	`, hashPassword, id)

	return err
}

var LETTER_RUNES = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateVerificationCode() string {
	length := 16
	res := make([]rune, length)
	for i := range res {
		res[i] = LETTER_RUNES[rand.Intn(len(LETTER_RUNES))]
	}
	return string(res)
}

func CreateUser(username, email, password string) error {
	err := validity.ValidateUser(username, email, password, true)
	if err != nil {
		return err
	}

	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	verificationCode := GenerateVerificationCode()

	db := utils.GetDB()
	_, err = db.Exec(`
			INSERT INTO users (username, email, password, verification_code)
			VALUES ($1, $2, $3, $4)
		`, username, email, hashPassword, verificationCode)

	return err
}

func SendVerificationEmail(username string) error {
	user, err := GetUserByUsername(username)
	if err != nil {
		return errors.New("something went wrong")
	}

	emailSubject := "Email Verification"
	emailContent := "Please follow the link below:\n" +
		"http://localhost:3000/users/verify?userId=" + strconv.Itoa(user.ID) + "\n\n" +
		"And use the following code to verify your email:\n" +
		user.VerificationCode

	utils.SendEmail(user.Email, emailSubject, emailContent)

	return nil
}

func RequestPasswordChange(id int, email string) (string, error) {
	newVerificationCode := GenerateVerificationCode()
	renewPasswordURL := "http://localhost:3000/users/renew-password?userId=" + strconv.Itoa(id)

	db := utils.GetDB()
	_, err := db.Exec(`
		UPDATE users
		SET verification_code = $1
		WHERE id = $2
	`, newVerificationCode, id)
	if err != nil {
		return "", err
	}

	emailSubject := "Password Change"
	emailContent := "Please follow the link below:\n" +
		renewPasswordURL + "\n\n" +
		"And use the following code to renew your email:\n" +
		newVerificationCode

	utils.SendEmail(email, emailSubject, emailContent)

	return renewPasswordURL, nil
}
