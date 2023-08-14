package subscriber

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	repos "github.com/sabrs0/L0_WB/dataAccess/repositories"
	pgs "github.com/sabrs0/L0_WB/db/postgres"
	ents "github.com/sabrs0/L0_WB/entities"
)

type Subscriber struct {
	cache        *ents.OrdersCache
	db           *sql.DB
	Subscription *nats.Subscription
	orderRepo    repos.OrdersRepository
	deliveryRepo repos.DeliverysRepository
	paymentRepo  repos.PaymentsRepository
	itemRepo     repos.ItemsRepository
}

func NewSubscriber(cache *ents.OrdersCache) (Subscriber, error) {
	db, err := pgs.NewDB()
	if err != nil {
		return Subscriber{}, err
	}

	return Subscriber{
		db:           db,
		orderRepo:    repos.NewOrdersRepository(db),
		deliveryRepo: repos.NewDeliverysRepository(db),
		paymentRepo:  repos.NewPaymentsRepository(db),
		itemRepo:     repos.NewItemsRepository(db),
		cache:        cache,
	}, nil
}

func (s *Subscriber) Subscribe(streamName string) error {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	nc, err := nats.Connect(url)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	fmt.Println("Connected")

	defer nc.Drain()
	js, err := nc.JetStream()
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	fmt.Println("Jet Stream Created")

	consumerName := "handler-1"
	js.DeleteConsumer(streamName, consumerName)
	js.AddConsumer(streamName, &nats.ConsumerConfig{
		Durable:   consumerName,
		AckPolicy: nats.AckExplicitPolicy,
	})

	sub, err := js.PullSubscribe("", consumerName, nats.BindStream(streamName))
	if err != nil {
		return fmt.Errorf("Cant subscribe : %s", err.Error())
	}
	s.Subscription = sub
	defer s.Subscription.Unsubscribe()
	for {

		msgs, err := s.Subscription.Fetch(1, nats.MaxWait(time.Second*6))
		if err != nil {
			if err == nats.ErrTimeout {
				fmt.Println("MSG TIMED OUT, CONTINUING")
				continue
			}
			panic("Cant fetch msg: " + err.Error())
		}
		ords := &ents.Orders{}
		json.Unmarshal(msgs[0].Data, ords)
		err = s.InsertOrder(*ords)
		if err != nil {
			panic("Cant insert order: " + err.Error())
		}
		fmt.Println(ords)
		msgs[0].Ack()
	}

}

/*
	func (s *Subscriber) RecMsgs() {
		defer s.Subscription.Unsubscribe()
		//fmt.Println(s.Subscription)
		for {

			msgs, err := s.Subscription.Fetch(1, nats.MaxWait(time.Second*6))
			if err != nil {

				panic("Cant fetch msg: " + err.Error())
			}
			ords := &ents.Orders{}
			json.Unmarshal(msgs[0].Data, ords)
			err = s.InsertOrder(*ords)
			if err != nil {
				panic("Cant insert order: " + err.Error())
			}
			msgs[0].Ack()
		}
	}
*/
func (s *Subscriber) InsertOrder(ord ents.Orders) error {
	//
	if _, isOK := (*s.cache)[ord.OrderUID]; !isOK {
		(*s.cache)[ord.OrderUID] = ord
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("Cant start transaction : %s", err.Error())
	}
	defer tx.Rollback()
	_, err = s.orderRepo.SelectById(ord.OrderUID)
	if err != nil {
		err = s.orderRepo.Insert(ord)
		if err != nil {

			return fmt.Errorf("Cant transaction : %s", err.Error())
		}
	}

	s.deliveryRepo.Insert(ord.Delivery, ord.OrderUID)
	if err != nil {
		return fmt.Errorf("Cant transaction : %s", err.Error())
	}
	s.paymentRepo.Insert(ord.Payment, ord.OrderUID)
	if err != nil {
		return fmt.Errorf("Cant transaction : %s", err.Error())
	}
	for _, item := range ord.Items {
		err = s.itemRepo.Insert(item, ord.OrderUID)
		if err != nil {
			return fmt.Errorf("Cant transaction : %s", err.Error())
		}

	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Cant finish transaction : %s", err.Error())
	}

	return nil
}
func RecoverCache() (ents.OrdersCache, error) {
	cache := make(ents.OrdersCache)
	db, err := pgs.NewDB()
	if err != nil {
		return nil, fmt.Errorf("Cant recover cache : %s", err.Error())
	}
	orderRepo := repos.NewOrdersRepository(db)
	deliveryRepo := repos.NewDeliverysRepository(db)
	paymentRepo := repos.NewPaymentsRepository(db)
	itemRepo := repos.NewItemsRepository(db)

	orders, err := orderRepo.Select()
	if err != nil {
		return nil, fmt.Errorf("Cant recover cache : %s", err.Error())
	}
	for _, ord := range orders {
		//не будет ли тут проблем с ord...
		delivery, err := deliveryRepo.SelectByOrderId(ord.OrderUID)
		if err != nil {
			return nil, fmt.Errorf("Cant recover cache : %s", err.Error())
		}
		payment, err := paymentRepo.SelectByOrderId(ord.OrderUID)
		if err != nil {
			return nil, fmt.Errorf("Cant recover cache : %s", err.Error())
		}
		items, err := itemRepo.SelectByOrderId(ord.OrderUID)
		if err != nil {
			return nil, fmt.Errorf("Cant recover cache : %s", err.Error())
		}
		ord.Payment = payment
		ord.Delivery = delivery
		ord.Items = make([]ents.Item, len(items))
		copy(ord.Items, items)
		cache[ord.OrderUID] = ord

	}
	return cache, nil
}
