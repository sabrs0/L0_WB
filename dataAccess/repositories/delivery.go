package repositories

import (
	"database/sql"
	"fmt"

	pgs "github.com/sabrs0/L0_WB/db/postgres"
	ents "github.com/sabrs0/L0_WB/entities"
)

type DeliverysRepository struct {
	db *sql.DB
}

func NewDeliverysRepository(db *sql.DB) DeliverysRepository {
	return DeliverysRepository{
		db: db,
	}
}

func (OR *DeliverysRepository) Insert(dlv ents.Delivery, orderID string) error {
	_, err := OR.db.Exec(pgs.InsertDelivery, orderID, dlv.Name, dlv.Phone,
		dlv.Zip, dlv.City, dlv.Address, dlv.Region, dlv.Email)
	if err != nil {
		return fmt.Errorf("Cant insert Deliverys : %s", err.Error())
	}
	return nil
}

func (OR *DeliverysRepository) Delete(id int) error {
	_, err := OR.db.Exec(pgs.DeleteDelivery, id)
	if err != nil {
		return fmt.Errorf("Cant insert Deliverys : %s", err.Error())
	}
	return nil
}

func (OR *DeliverysRepository) SelectById(id int) (ents.Delivery, error) {
	dlv := &ents.Delivery{}
	row := OR.db.QueryRow(pgs.SelectByIDDelivery, id)
	var dlvid int
	err := row.Scan(&dlvid, &dlv.OrderUID, &dlv.Name, &dlv.Phone, &dlv.Zip, &dlv.City, &dlv.Address,
		&dlv.Region, &dlv.Email)
	if err != nil {
		return ents.Delivery{}, fmt.Errorf("Cant get Delivery by id %d : %s", id, err.Error())
	}
	return *dlv, nil
}
func (OR *DeliverysRepository) SelectByOrderId(id string) (ents.Delivery, error) {
	dlv := &ents.Delivery{}
	row := OR.db.QueryRow(pgs.SelectByOrderIDDelivery, id)
	var dlvid int
	err := row.Scan(&dlvid, &dlv.OrderUID, &dlv.Name, &dlv.Phone, &dlv.Zip, &dlv.City, &dlv.Address,
		&dlv.Region, &dlv.Email)
	if err != nil {
		return ents.Delivery{}, fmt.Errorf("Cant get Delivery by orderID %s : %s", id, err.Error())
	}
	return *dlv, nil
}
