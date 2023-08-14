package db

// INSERTS
const (
	InsertOrder = `
	INSERT INTO orders (order_uid, track_number, entry, locale,
						internal_signature, customer_id, delivery_service,
						shardkey, sm_id, date_created, off_shard)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
`
	InsertDelivery = `
	INSERT INTO delivery (order_uid, name, phone, zip,
						city, address, region, email)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8)
`
	InsertPayment = `
	INSERT INTO payment (order_uid, transaction, request_id, currency,
						provider, bank, delivery_cost,
						goods_total, custom_fee)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
`
	InsertItem = `
	INSERT INTO item (order_uid, chrt_id, track_number, price,
						rid, name, sale, size, total_price,
						nm_id, brand, status)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
`
)

// DELETES
const (
	DeleteOrder = `DELETE FROM orders WHERE id = $1`

	DeleteDelivery = `DELETE FROM delivery WHERE id = $1`

	DeletePayment = `DELETE FROM payment WHERE id = $1`

	DeleteItem = `DELETE FROM item WHERE id = $1`
)

// SELECTS BY ID
const (
	SelectByIDOrder = `SELECT * FROM orders WHERE order_uid = $1`

	SelectByIDDelivery = `SELECT * FROM delivery WHERE id = $1`

	SelectByIDPayment = `SELECT * FROM payment WHERE id = $1`

	SelectByIDItem = `SELECT * FROM item WHERE id = $1`
)

// SELECTS BY ORDERID
const (
	SelectByOrderIDDelivery = `SELECT * FROM delivery WHERE order_uid = $1`

	SelectByOrderIDPayment = `SELECT * FROM payment WHERE ORDER_UID = $1`

	SelectByOrderIDItem = `SELECT * FROM item WHERE ORDER_UID = $1`
)

// SELECTS
const (
	SelectOrders = `SELECT * FROM orders`
)
