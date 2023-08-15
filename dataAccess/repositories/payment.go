package repositories

import (
	"database/sql"
	"fmt"

	pgs "github.com/sabrs0/L0_WB/db/postgres"
	ents "github.com/sabrs0/L0_WB/entities"
)

type PaymentsRepository struct {
	db *sql.DB
}

func NewPaymentsRepository(db *sql.DB) PaymentsRepository {
	return PaymentsRepository{
		db: db,
	}
}

func (OR *PaymentsRepository) Insert(pmnt ents.Payment, orderID string) error {
	_, err := OR.db.Exec(pgs.InsertPayment, orderID, pmnt.Transaction,
		pmnt.RequestID, pmnt.Currency, pmnt.Provider, pmnt.Amount, pmnt.PaymentDt, pmnt.Bank, pmnt.DeliveryCost,
		pmnt.GoodsTotal, pmnt.CustomFee)
	if err != nil {
		return fmt.Errorf("Cant insert Payments : %s", err.Error())
	}
	return nil
}

func (OR *PaymentsRepository) Delete(id int) error {
	_, err := OR.db.Exec(pgs.DeletePayment, id)
	if err != nil {
		return fmt.Errorf("Cant insert Payments : %s", err.Error())
	}
	return nil
}

func (OR *PaymentsRepository) SelectById(id int) (ents.Payment, error) {
	var pmnt ents.Payment
	row := OR.db.QueryRow(pgs.SelectByIDPayment, id)
	var pmtID int
	err := row.Scan(&pmtID, &pmnt.OrderUID, &pmnt.Transaction, &pmnt.RequestID, &pmnt.Currency,
		&pmnt.Provider, &pmnt.Amount, &pmnt.PaymentDt, &pmnt.Bank, &pmnt.DeliveryCost,
		&pmnt.GoodsTotal, &pmnt.CustomFee)
	if err != nil {
		return ents.Payment{}, fmt.Errorf("Cant get Payment by id %d : %s", id, err.Error())
	}
	return pmnt, nil
}
func (OR *PaymentsRepository) SelectByOrderId(id string) (ents.Payment, error) {
	var pmnt ents.Payment
	row := OR.db.QueryRow(pgs.SelectByOrderIDPayment, id)
	var pmtID int
	err := row.Scan(&pmtID, &pmnt.OrderUID, &pmnt.Transaction, &pmnt.RequestID, &pmnt.Currency,
		&pmnt.Provider, &pmnt.Amount, &pmnt.PaymentDt, &pmnt.Bank, &pmnt.DeliveryCost,
		&pmnt.GoodsTotal, &pmnt.CustomFee)
	if err != nil {
		return ents.Payment{}, fmt.Errorf("Cant get Payment by orderID %s : %s", id, err.Error())
	}
	return pmnt, nil
}
