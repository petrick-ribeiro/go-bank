package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/petrick-ribeiro/go-bank/types"
)

type Storage interface {
	CreateAccount(*types.Account) error
	GetAccounts() ([]*types.Account, error)
	GetAccountByID(int) (*types.Account, error)
	GetAccountByNumber(int) (*types.Account, error)
	UpdateAccount(*types.Account) error
	DeleteAccount(int) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostGresStore() (*PostgresStore, error) {
	// TODO: Get DB info from a dotenv
	connStr := "user=postgres dbname=postgres password=foo sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS accounts 
        (id SERIAL PRIMARY KEY,
        first_name VARCHAR(50),
        last_name VARCHAR(50),
        password VARCHAR(60),
        number INT,
        balance INT,
        created_at TIMESTAMP)`

	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) CreateAccount(acc *types.Account) error {
	query := `
        INSERT INTO accounts 
        (first_name, last_name, password, number, balance, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        `
	_, err := s.db.Query(
		query,
		acc.FirstName,
		acc.LastName,
		acc.EncryptedPassword,
		acc.Number,
		acc.Balance, acc.CreatedAt)
	if err != nil {
		return err
	}

	return err
}

func (s *PostgresStore) GetAccounts() ([]*types.Account, error) {
	rows, err := s.db.Query("SELECT * from accounts")
	if err != nil {
		return nil, err
	}

	accounts := []*types.Account{}
	for rows.Next() {
		account, err := ScanIntoAccount(rows)
		accounts = append(accounts, account)
		if err != nil {
			return nil, err
		}
	}

	return accounts, nil
}

func ScanIntoAccount(rows *sql.Rows) (*types.Account, error) {
	account := new(types.Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.EncryptedPassword,
		&account.Number,
		&account.Balance,
		&account.CreatedAt)

	return account, err
}

func (s *PostgresStore) GetAccountByNumber(number int) (*types.Account, error) {
	rows, err := s.db.Query("SELECT * FROM accounts WHERE number = $1", number)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return ScanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account with number [%d] not found", number)
}

func (s *PostgresStore) GetAccountByID(id int) (*types.Account, error) {
	rows, err := s.db.Query("SELECT * FROM accounts WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return ScanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account with id [%d] not found", id)
}

func (s *PostgresStore) UpdateAccount(*types.Account) error {
	// TODO: Update Account info using SQL Query
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Query("DELETE FROM accounts WHERE id = $1", id)

	return err
}
