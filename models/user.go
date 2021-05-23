package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// User is used by pop to map your users database table to your go code.
type User struct {
	ID     uuid.UUID `json:"id,omitempty" db:"id"`
	ChatID string    `json:"chat_id,omitempty" db:"chat_id"`
	FEN    string    `json:"fen,omitempty" db:"fen"`
	// Side        int       `json:"side,omitempty" db:"side"`
	History string `json:"history,omitempty" db:"history"`
	// Orientation int       `json:"orientation,omitempty" db:"orientation"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
	// Board       string    `json:"board" db:"board"`
}

func AddUser(user *User) (*User, error) {
	_, err = Tx.ValidateAndCreate(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateUserFEN(chatID, fen, history string) (*User, error) {
	// fmt.Println("Update user")

	user := User{}
	query := Tx.Where("chat_id = ?", chatID)
	err = query.First(&user)
	if err != nil {
		return nil, err
	}
	user.FEN = fen
	user.History = history
	// Tx.ValidateAndCreate()
	// Tx.ValidateAndSave()
	// Tx.ValidateAndUpdate()
	_, err = Tx.ValidateAndSave(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUser(userID string) (*User, error) {

	user := User{}
	query := Tx.Where("chat_id = ?", userID)
	err = query.First(&user)
	if err != nil {
		return nil, err
	}
	Tx.Load(&user)
	return &user, nil
}

// func UpdateUserBoard()
