package repositories

import (
	"database/sql"
	"fmt"

	pgs "github.com/sabrs0/L0_WB/db/postgres"
	ents "github.com/sabrs0/L0_WB/entities"
)

type ItemsRepository struct {
	db *sql.DB
}

func NewItemsRepository(db *sql.DB) ItemsRepository {
	return ItemsRepository{
		db: db,
	}
}

func (OR *ItemsRepository) Insert(itm ents.Item, orderID string) error {
	_, err := OR.db.Exec(pgs.InsertItem, orderID, itm.ChrtID, itm.TrackNumber,
		itm.Price, itm.Rid, itm.Name, itm.Sale, itm.Size, itm.TotalPrice,
		itm.NmID, itm.Brand, itm.Status)
	if err != nil {
		return fmt.Errorf("Cant insert Items : %s", err.Error())
	}
	return nil
}

func (OR *ItemsRepository) Delete(id int) error {
	_, err := OR.db.Exec(pgs.DeleteItem, id)
	if err != nil {
		return fmt.Errorf("Cant insert Items : %s", err.Error())
	}
	return nil
}

func (OR *ItemsRepository) SelectById(id int) (ents.Item, error) {
	post := &ents.Item{}
	row := OR.db.QueryRow(pgs.SelectByIDItem, id)
	var itmId int
	err := row.Scan(&itmId, &post.OrderUID, &post.ChrtID, &post.TrackNumber, &post.Price, &post.Rid, &post.Name, &post.Sale,
		&post.Size, &post.TotalPrice, &post.NmID, &post.Brand, &post.Status)
	if err != nil {
		return ents.Item{}, fmt.Errorf("Cant get Item by id %d : %s", id, err.Error())
	}
	return *post, nil
}
func (OR *ItemsRepository) SelectByOrderId(id string) ([]ents.Item, error) {
	items := []ents.Item{}
	rows, err := OR.db.Query(pgs.SelectByOrderIDItem, id)
	if err != nil {
		return nil, fmt.Errorf("Cant get items : %s", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		post := &ents.Item{}
		var itmId int
		err = rows.Scan(&itmId, &post.OrderUID, &post.ChrtID, &post.TrackNumber, &post.Price, &post.Rid, &post.Name, &post.Sale,
			&post.Size, &post.TotalPrice, &post.NmID, &post.Brand, &post.Status)
		if err != nil {
			return nil, fmt.Errorf("Cant get items: Cant Scan : %s", err.Error())
		}
		items = append(items, *post)
	}
	return items, nil
}
