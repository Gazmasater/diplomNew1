package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	db *pgxpool.Pool
}

const createUsersTableQuery = `
      CREATE TABLE IF NOT EXISTS users (
    UserId VARCHAR(255) PRIMARY KEY,
    Username VARCHAR(255) UNIQUE,
    Password VARCHAR(255),
    CONSTRAINT unique_person UNIQUE (Username)
);
    `

const createUserOrdersTableQuery = `
	CREATE TABLE IF NOT EXISTS user_orders (
	OrderId SERIAL PRIMARY KEY,
	UserId VARCHAR(255),
	OrderNumber VARCHAR(255),
	Status VARCHAR(255),
	Accrual DOUBLE PRECISION DEFAULT 0 NOT NULL,
	OrderDate VARCHAR(255) DEFAULT TO_CHAR(CURRENT_TIMESTAMP, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
	FOREIGN KEY (UserId) REFERENCES users(UserId),
	CONSTRAINT unique_order UNIQUE (OrderNumber)
	    );
`

const createUserBalanceTableQuery = `
	CREATE TABLE IF NOT EXISTS user_balance (
	    UserID VARCHAR(255),
	    Current DOUBLE PRECISION DEFAULT 0 NOT NULL,
	    CONSTRAINT user_balance_userId_unique UNIQUE (UserId),
	    Withdrawn DOUBLE PRECISION DEFAULT 0 NOT NULL 
	);
	`

const createBillingTableQuery = `
CREATE TABLE IF NOT EXISTS billing (
    UserID VARCHAR(255),
    OrderNumber VARCHAR(255),
    Sum DOUBLE PRECISION, 
    Processed_at TIMESTAMPTZ DEFAULT TO_TIMESTAMP(TO_CHAR(CURRENT_TIMESTAMP, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
    FOREIGN KEY (UserID) REFERENCES users(UserID),
    CONSTRAINT unique_order_number UNIQUE (OrderNumber)
);
`

func New(connString string) (*Database, error) {
	conn, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to db: %v", err)
	}

	_, err = conn.Exec(context.Background(), createUsersTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating user table: %v", err)
	}

	_, err = conn.Exec(context.Background(), createUserOrdersTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating user orders table: %v", err)
	}

	_, err = conn.Exec(context.Background(), createUserBalanceTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating user balance table: %v", err)
	}

	_, err = conn.Exec(context.Background(), createBillingTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating billing table: %v", err)
	}

	return &Database{db: conn}, nil
}

func (d *Database) UserSetter(ctx context.Context, user, password, id string) error {

	insertQuery := `
                INSERT INTO users (Username, Password, UserID)
				VALUES ($1, $2, $3)
			
            `
	_, err := d.db.Exec(ctx, insertQuery, user, password, id)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return errors.New("this user already exists")
			} else {
				return fmt.Errorf("error when trying to add user to database: %v", err)
			}
		}
	}

	return nil
}

func (d *Database) CredentialsGetter(ctx context.Context, user string) (string, string, error) {
	var password string
	var id string
	err := d.db.QueryRow(ctx, "SELECT password, userid FROM users WHERE username = $1", user).Scan(&password, &id)

	if err != nil {
		if err == pgx.ErrNoRows {
			return "", "", errors.New("user not found")
		}
		return "", "", fmt.Errorf("error when getting hash password from database: %v", err)
	}

	return password, id, nil
}
