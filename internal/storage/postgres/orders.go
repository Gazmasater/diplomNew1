package postgres

import (
	"context"
	"errors"

	"diplom.com/go-musthave-diploma-tpl/internal/dto"
	"diplom.com/go-musthave-diploma-tpl/internal/handlers/addorder"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (d *Database) OrderSetter(ctx context.Context, userID, orderNumber string) error {
	var existingUserID string
	selectQuery := `SELECT userid FROM user_orders WHERE orderNumber = $1 LIMIT 1`
	err := d.db.QueryRow(ctx, selectQuery, orderNumber).Scan(&existingUserID)
	if err != nil {
		if err == pgx.ErrNoRows {
			insertQuery := `INSERT INTO user_orders (userid, orderNumber, status) VALUES ($1, $2, $3)`
			_, err = d.db.Exec(ctx, insertQuery, userID, orderNumber, "NEW")
			return err
		}
		return err
	}

	if existingUserID != userID {
		return addorder.ErrDifferentUser
	}

	return addorder.ErrOrderExists
}

func (d *Database) StatusSetter(ctx context.Context, orderNumber, orderStatus string, accrual float32) (string, error) {
	insertQuery := `
                UPDATE user_orders 
                SET status = $2, accrual = $3 
                WHERE orderNumber = $1
				RETURNING userid
            `
	var userID string
	err := d.db.QueryRow(ctx, insertQuery, orderNumber, orderStatus, accrual).Scan(&userID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return "", errors.New("this order already exists, please try another one")
			} else {
				return "", err
			}
		}
	}

	return userID, nil
}

func (d *Database) GetOrderNumbers(ctx context.Context) ([]string, error) {
	query := `
        SELECT ordernumber
        FROM user_orders
        WHERE status <> 'PROCESSED' AND status <> 'INVALID'
    `

	rows, err := d.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderNumbers []string
	for rows.Next() {
		var orderNumber string
		err = rows.Scan(&orderNumber)
		if err != nil {
			return nil, err
		}
		orderNumbers = append(orderNumbers, orderNumber)
	}
	return orderNumbers, nil
}

func (d *Database) GetOrders(ctx context.Context, userID string) ([]dto.FullOrder, error) {
	getOrdersQuery := `
		SELECT  ordernumber, status, COALESCE(accrual, 0) AS accrual, orderdate
		FROM user_orders 
		WHERE userid = $1
		ORDER BY orderdate DESC 
`

	rows, err := d.db.Query(ctx, getOrdersQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []dto.FullOrder{}
	for rows.Next() {
		var order dto.FullOrder
		err = rows.Scan(&order.OrderNumber, &order.OrderStatus, &order.Accrual, &order.OrderDate)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}
