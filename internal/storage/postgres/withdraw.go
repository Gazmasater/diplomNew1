package postgres

import (
	"context"
	"errors"
	"fmt"

	"diplom.com/internal/dto"
)

func (d *Database) GetWithdrawals(ctx context.Context, userID string) ([]*dto.Billing, error) {
	getOrdersQuery := `
		SELECT  ordernumber, sum, processed_at
		FROM billing
		WHERE userid = $1
		ORDER BY processed_at ASC
`
	rows, err := d.db.Query(ctx, getOrdersQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []*dto.Billing{}

	for rows.Next() {
		billing := dto.Billing{}
		err = rows.Scan(&billing.Order, &billing.Sum, &billing.ProcessedAt)
		if err != nil {
			return nil, err
		}

		orders = append(orders, &billing)
	}

	if len(orders) == 0 {
		return nil, errors.New("no withdrawals found")
	}

	return orders, nil

}

func (d *Database) ProcessOrder(ctx context.Context, order, userID string, sum float32) error {
	// Проверяем баланс в таблице user_balance
	checkBalanceQuery := `
        SELECT current
        FROM user_balance
        WHERE current >= $1
    `
	var currentBalance float32
	err := d.db.QueryRow(ctx, checkBalanceQuery, sum).Scan(&currentBalance)
	if err != nil {
		return fmt.Errorf("failed to check balance: %s", err)
	}

	// Вычитаем сумму из баланса и обновляем значение в withdrawn
	updateBalanceQuery := `
        UPDATE user_balance
        SET current = current - $1, withdrawn = COALESCE(withdrawn, 0) + $1
    `
	_, err = d.db.Exec(ctx, updateBalanceQuery, sum)
	if err != nil {
		return fmt.Errorf("failed to update balance: %s", err)
	}

	// Вставляем информацию о заказе в таблицу billing
	insertBillingQuery := `
        INSERT INTO billing (userID, ordernumber, sum)
        VALUES ($1, $2, $3)
    `
	_, err = d.db.Exec(ctx, insertBillingQuery, userID, order, sum)
	if err != nil {
		return fmt.Errorf("failed to insert billing record: %s", err)
	}

	return nil
}
