package postgres

import (
	"context"
	"fmt"

	"diplom.com/go-musthave-diploma-tpl/internal/dto"
	"github.com/jackc/pgx/v5"
)

func (d *Database) AddUserBalance(ctx context.Context, userID string, value float32) error {
	insertQuery := `
        INSERT INTO user_balance (userID, current)
        VALUES ($1, $2)
        ON CONFLICT (userId)
        DO UPDATE SET current = user_balance.current + $2
    `

	_, err := d.db.Exec(ctx, insertQuery, userID, value)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetUserBalance(ctx context.Context, userID string) (*dto.UserBalance, error) {
	query := `
        SELECT COALESCE(current, 0), COALESCE(withdrawn, 0)
        FROM user_balance
        WHERE userid = $1
    `

	rows := d.db.QueryRow(ctx, query, userID)

	var userBalance dto.UserBalance
	err := rows.Scan(&userBalance.Current, &userBalance.Withdrawn)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user balance not found for UserID: %s", userID)
		}
		return nil, err
	}

	return &userBalance, nil
}
