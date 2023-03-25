package types

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginResponse struct {
	Number int64  `json:"number"`
	Token  string `json:"token"`
}

type LoginRequest struct {
	Number   int64  `json:"number"`
	Password string `json:"password"`
}

type TransferRequest struct {
	ToAccount int `json:"to_account"`
	Amount    int `json:"amount"`
}

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

type Account struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	EncryptedPassword string    `json:"-"`
	Number            int64     `json:"number"`
	Balance           int64     `json:"balance"`
	CreatedAt         time.Time `json:"created_at"`
}

func (a *Account) ValidadePassword(pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(pw))
}

func NewAccount(firstName, lastName, password string) (*Account, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Initialize with a seed based on current time
	rand.Seed(time.Now().UnixNano())

	return &Account{
		FirstName:         firstName,
		LastName:          lastName,
		EncryptedPassword: string(hashedPassword),
		Number:            int64(rand.Intn(100000)),
		CreatedAt:         time.Now().UTC(),
	}, nil
}
