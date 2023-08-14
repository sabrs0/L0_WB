package repositories

import (
	"database/sql"
	"fmt"

	pgs "github.com/sabrs0/L0_WB/db/postgres"
	ents "github.com/sabrs0/L0_WB/entities"
)

type OrdersRepository struct {
	db *sql.DB
}

func NewOrdersRepository(db *sql.DB) OrdersRepository {
	return OrdersRepository{
		db: db,
	}
}

func (OR *OrdersRepository) Insert(ords ents.Orders) error {
	_, err := OR.db.Exec(pgs.InsertOrder, ords.OrderUID, ords.TrackNumber, ords.Entry, ords.Locale,
		ords.InternalSignature, ords.CustomerID, ords.DeliveryService, ords.ShardKey, ords.SmID,
		ords.DateCreated, ords.OofShard)
	if err != nil {
		return fmt.Errorf("Cant insert orders : %s", err.Error())
	}
	return nil
}

func (OR *OrdersRepository) Delete(id string) error {
	_, err := OR.db.Exec(pgs.DeleteOrder, id)
	if err != nil {
		return fmt.Errorf("Cant insert orders : %s", err.Error())
	}
	return nil
}

func (OR *OrdersRepository) SelectById(id string) (ents.Orders, error) {
	row := OR.db.QueryRow(pgs.SelectByIDOrder, id)
	post := &ents.Orders{}
	err := row.Scan(&post.OrderUID, &post.TrackNumber, &post.Entry, &post.Locale, &post.InternalSignature,
		&post.CustomerID, &post.DeliveryService, &post.ShardKey, &post.SmID, &post.DateCreated,
		&post.OofShard)
	if err != nil {
		return ents.Orders{}, fmt.Errorf("Cant get order by id %s : %s", id, err.Error())
	}
	return *post, nil
}
func (OR *OrdersRepository) Select() ([]ents.Orders, error) {
	orders := []ents.Orders{}
	rows, err := OR.db.Query(pgs.SelectOrders)
	if err != nil {
		return nil, fmt.Errorf("Cant get orders : %s", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		post := &ents.Orders{}
		err = rows.Scan(&post.OrderUID, &post.TrackNumber, &post.Entry, &post.Locale, &post.InternalSignature,
			&post.CustomerID, &post.DeliveryService, &post.ShardKey, &post.SmID, &post.DateCreated,
			&post.OofShard)
		if err != nil {
			return nil, fmt.Errorf("Cant get orders: Cant Scan : %s", err.Error())
		}
		orders = append(orders, *post)
	}

	return orders, nil
}
