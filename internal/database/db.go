package database

import (
	"WBtestL0/internal/models"
	"context"
	"database/sql"
	"fmt"
)

type DbConn struct {
	db *sql.DB
}

// создаём подключение
func NewConn(cfg Config) (*DbConn, error) {
	db, err := sql.Open("pgx", cfg.connStr())
	if err != nil {
		return nil, err
	}

	conn := &DbConn{db}

	return conn, nil
}

// закрываем БД
func (d DbConn) CloseDB() error {
	err := d.db.Close()
	if err != nil {
		return fmt.Errorf("problem with close DB: %w", err)
	}
	return nil
}

func (d *DbConn) CreateOrder(ctx context.Context, order models.Order) error {
	tx, err := d.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin is %w", err)
	}

	defer tx.Rollback()

	//проводим транзакции в БД
	_, err = tx.ExecContext(ctx, "INSERT INTO orders(order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)",
		order.OrderUid, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerId, order.DeliveryService, order.ShardKey, order.SmId, order.DateCreated, order.OofShard)
	if err != nil {
		return fmt.Errorf("wrong orders - error is: %w", err)
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO deliveries(name, phone, zip, city, address, region, email, order_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)",
		order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email, order.OrderUid)
	if err != nil {
		return fmt.Errorf("wrong deliveries - error is: %w", err)
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO payments(transaction, request_id, currency , provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee, order_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)",
		order.Payment.Transaction, order.Payment.RequestId, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee, order.OrderUid)
	if err != nil {
		return fmt.Errorf("wrong payments - error is: %w", err)
	}

	for i, item := range order.Items {
		_, err = tx.ExecContext(ctx, "INSERT INTO items(chrt_id, track_number, price , rid, name, sale, size, total_price, nm_id, brand, status, order_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)",
			item.ChrtId, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmId, item.Brand, item.Status, order.OrderUid)
		if err != nil {
			return fmt.Errorf("insert item №%v:%w", i, err)
		}
	}

	//если все транзакции прошли успешно, комитим
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("tx commit:%w", err)
	}

	return nil
}

func (d DbConn) SelectOrders(ctx context.Context) ([]models.Order, error) {
	ids, err := d.db.QueryContext(ctx, "SELECT order_uid FROM orders")
	if err != nil {
		return nil, fmt.Errorf("something wrong with ids. Error is %w", err)
	}

	//сохраняем в памяти
	orders := make([]models.Order, 0, 10)

	//проходимся по всем айдишникам и сохраняем их в памяти
	for ids.Next() {
		var id string

		//сканим в id
		err = ids.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("wrong scan.Error is :%w", err)
		}
		order, err := d.SelectOrderById(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("select order by id:%w", err)
		}

		orders = append(orders, order)

	}

	return nil, err
}

func (d DbConn) SelectOrderById(ctx context.Context, orderId string) (models.Order, error) {
	tx, err := d.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return models.Order{}, fmt.Errorf("error in begin: %w", err)
	}

	defer tx.Rollback()

	var order models.Order
	if err = tx.QueryRowContext(ctx, "SELECT * FROM orders WHERE order_uid = $1", orderId).
		Scan(&order.OrderUid, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerId, &order.DeliveryService, &order.ShardKey, &order.SmId, &order.DateCreated, &order.OofShard); err != nil {
		return models.Order{}, fmt.Errorf("select order:%w", err)
	}

	var delivery models.Delivery
	if err = tx.QueryRowContext(ctx, "SELECT * FROM delivery WHERE order_uid = $1", orderId).
		Scan(&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address, &delivery.Region, &delivery.Email, &orderId); err != nil {
		return models.Order{}, fmt.Errorf("select delivery:%w", err)
	}

	order.Delivery = delivery

	var payment models.Payment
	if err = tx.QueryRowContext(ctx, "SELECT * FROM payments WHERE order_id = $1", orderId).
		Scan(&payment.Transaction, &payment.RequestId, &payment.Currency, &payment.Provider, &payment.Amount, &payment.PaymentDt, &payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee, &orderId); err != nil {
		return models.Order{}, fmt.Errorf("select delivery:%w", err)
	}

	order.Payment = payment

	rows, err := tx.QueryContext(ctx, "SELECT * FROM items WHERE order_id = $1", orderId)
	if err != nil {
		return models.Order{}, fmt.Errorf("wrong items: %w", err)
	}

	items := make([]models.Item, 0, 10)

	for rows.Next() {
		var item models.Item

		if err := rows.Scan(&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status); err != nil {
			return models.Order{}, fmt.Errorf("wrong scan item: %w", err)
		}

		items = append(items, item)
	}

	order.Items = items

	if err = tx.Commit(); err != nil {
		return models.Order{}, fmt.Errorf("tx commit: %w", err)
	}

	return order, nil
}
