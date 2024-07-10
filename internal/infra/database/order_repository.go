package database

import (
	"database/sql"

	"github.com/dmarins/clean-arch-challenge-go/internal/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	cmd, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = cmd.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) GetAll() ([]*entity.Order, error) {
	rows, err := r.Db.Query("SELECT id, price, tax, final_price FROM orders ORDER BY final_price DESC")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orders := make([]*entity.Order, 0)
	for rows.Next() {
		var o entity.Order
		err = rows.Scan(
			&o.ID,
			&o.Price,
			&o.Tax,
			&o.FinalPrice,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, &o)
	}

	return orders, nil
}
