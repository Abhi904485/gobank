package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage interface {
	createAccount(*Account) error
	isAccountExists(id int) (*Account, error)
	updateAccount(id int, account *UpdateAccountRequest) error
	deleteAccount(id int) error
	getAccountByID(id int) (*Account, error)
	getAccounts() ([]*Account, error)
}

type PostgresStore struct {
	Db *sql.DB
}

func newPostgresStore() (*PostgresStore, error) {
	connectionString := "user=postgres dbname=gobank password=postgres sslmode=disable"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{
		Db: db,
	}, nil
}

func (postgresStore PostgresStore) InitDb() error {
	tableQuery := `CREATE TABLE IF NOT EXISTS Account (
             id serial primary key NOT NULL,
             firstName varchar(100) NOT NULL,
             lastName varchar(100) NOT NULL,
             number serial NOT NULL,
             balance serial NOT NULL,
             created_at timestamp)`

	_, err := postgresStore.Db.Exec(tableQuery)
	if err != nil {
		return fmt.Errorf("table Creation Failed %s ", err)
	}
	return nil
}

func (postgresStore PostgresStore) createAccount(account *Account) error {
	insertQuery := "Insert into Account (firstname, lastname, number, balance, created_at) values ($1, $2, $3, $4, $5)"
	_, err := postgresStore.Db.Exec(insertQuery, account.FirstName, account.LastName, account.Number, account.Balance, account.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (postgresStore PostgresStore) updateAccount(id int, account *UpdateAccountRequest) error {
	_, err1 := postgresStore.isAccountExists(id)
	if err1 != nil {
		return err1
	}
	updateQuery := `Update Account set firstname = $1, lastname = $2, number = $3, balance =$4 where id = $5`
	_, err := postgresStore.Db.Exec(updateQuery, account.FirstName, account.LastName, account.Number, account.Balance, id)
	if err != nil {
		return err
	}
	return nil
}

func (postgresStore PostgresStore) deleteAccount(id int) error {
	_, err1 := postgresStore.isAccountExists(id)
	if err1 != nil {
		return err1
	}
	deleteQuery := `delete from account where id = $1`
	_, err := postgresStore.Db.Exec(deleteQuery, id)
	if err != nil {
		return fmt.Errorf("deletion Failed for id %d", id)
	}
	return nil

}

func (postgresStore PostgresStore) isAccountExists(id int) (*Account, error) {
	accountExistQuery := `select * from Account where id = $1`
	rows, err := postgresStore.Db.Query(accountExistQuery, id)
	if err != nil {
		return nil, fmt.Errorf("account does not Exist With id %d ", id)
	}

	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		return account, nil
	}
	return nil, fmt.Errorf("account does not Exist With id %d ", id)

}

func (postgresStore PostgresStore) getAccounts() ([]*Account, error) {
	query := `Select * from Account`
	rows, err := postgresStore.Db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("get Account Query Failed %s", err)
	}
	var accounts []*Account
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (postgresStore PostgresStore) getAccountByID(id int) (*Account, error) {
	account, err2 := postgresStore.isAccountExists(id)
	if err2 != nil {
		return nil, err2
	}
	return account, nil

}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := &Account{}
	err := rows.Scan(&account.Id, &account.FirstName, &account.LastName, &account.Number, &account.Balance, &account.CreatedAt)
	if err != nil {
		return nil, err
	}
	return account, nil
}
