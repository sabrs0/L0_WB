package services

import (
	"fmt"

	repos "github.com/sabrs0/L0_WB/dataAccess/repositories"
	ents "github.com/sabrs0/L0_WB/entities"
)

type OrdersService struct {
	repo  repos.OrdersRepository
	cache *ents.OrdersCache
}

func NewOrdersService(repo repos.OrdersRepository, cache *ents.OrdersCache) OrdersService {
	return OrdersService{
		repo:  repo,
		cache: cache,
	}
}

func (OS *OrdersService) Insert(ords ents.Orders) error {
	err := OS.repo.Insert(ords)
	if err != nil {
		return err
	}
	(*OS.cache)[ords.OrderUID] = ords
	return nil
}

func (OS *OrdersService) Delete(id string) error {
	return OS.repo.Delete(id)
	//return nil
}

func (OR *OrdersService) SelectById(id string) (ents.Orders, error) {
	ord, isOK := (*OR.cache)[id]
	if !isOK {
		return ents.Orders{}, fmt.Errorf("No such elem in cache")
	}
	return ord, nil
}
