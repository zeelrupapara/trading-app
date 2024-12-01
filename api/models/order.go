package models

import (
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

const OrderTable = "orders"

// Order model
type Order struct {
	ID        string    `json:"id"`
	Symbol    string    `json:"symbol"`
	Volume    float64   `json:"volume"`
	OrderType string    `json:"order_type" db:"order_type"` // "buy" or "sell", with default value
	Price     string    `json:"price"`
	UserId    string    `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type OrderModel struct {
	db *goqu.Database
}

func InitOrderModel(db *goqu.Database) (OrderModel, error) {
	return OrderModel{
		db: db,
	}, nil
}

func (model *OrderModel) GetById(id string) (Order, error) {
	order := Order{}
	found, err := model.db.From(OrderTable).Where(goqu.Ex{
		"id": id,
	}).Select(
		"id",
		"symbol",
		"volume",
		"order_type",
		"price",
		"created_at",
		"updated_at",
	).ScanStruct(&order)

	if err != nil {
		return order, err
	}

	if !found {
		return order, sql.ErrNoRows
	}

	return order, err
}

func (model *OrderModel) InsertOrder(order Order) (Order, error) {
	// do transaction
	tx, err := model.db.Begin()
	if err != nil {
		return order, err
	}

	defer tx.Rollback()

	order.ID = uuid.New().String()
	_, err = tx.Insert(OrderTable).Rows(
		goqu.Record{
			"id":         order.ID,
			"symbol":     order.Symbol,
			"volume":     order.Volume,
			"order_type": order.OrderType,
			"price":      order.Price,
			"user_id":    order.UserId,
			"created_at": order.CreatedAt,
			"updated_at": order.UpdatedAt,
		},
	).Executor().Exec()
	if err != nil {
		return order, err
	}

	err = tx.Commit()
	if err != nil {
		return order, err
	}

	order, err = model.GetById(order.ID)
	if err != nil {
		return order, err
	}
	return order, err
}

func (model *OrderModel) GetOrders(limit uint, offset uint, userId string) ([]Order, error) {
	var orders []Order
	query := model.db.From(OrderTable).Limit(limit).Offset(offset)
	if userId != "" {
		query = query.Where(goqu.Ex{"user_id": userId})
	}
	if err := query.ScanStructs(&orders); err != nil {
		return nil, err
	}
	return orders, nil
}
