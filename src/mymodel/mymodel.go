package mymodel

import (
    "time"
)

const (
    StatusCreated = "created"
    StatusProcessed = "processed"
    StatusCanceled = "canceled"
    StatusFinished = "finished"
)

type Order struct {
    orderId int
    shopId int
    customerId int
    status string
    products []Product
    created time.Time
}

type Product struct {
    productId int
    name string
    created time.Time
}

type ProductJson struct {
    productId int `json:"product_id"`
}

type CreateOrderJson struct {
    shopId int `json:"shop_id"`
    customerId int `json:"customer_id"`
    products []ProductJson `json:"products"`
}

type Datastore interface {
    init() error
    getOrderCollection(search string, page int) ([]*Order, error)
    getOrder(id int) (*Order, error)
    createOrder(productJson *ProductJson) (*Order, error)
    updateOrderStatus(id int, status string) error
}

type PostgresDatastore struct {
    
}

func (p *PostgresDatastore) init() error {
    return nil
}

func (p *PostgresDatastore) getOrderCollection(search string, page int) ([]*Order, error) {
    orders := []*Order{}

    return orders, nil
}

func (p *PostgresDatastore) getOrder(id int) (*Order, error) {
    order := Order{}

    return &order, nil
}

func (p *PostgresDatastore) createOrder(productJson *ProductJson) (*Order, error) {
    order := Order{}

    return &order, nil
}

func (p *PostgresDatastore) updateOrderStatus(id int, status string) error {
    return nil
}


func GetOrders(search string, page int) ([]*Order, error) {
    orders := []*Order{}

    return orders, nil
}

func CreateOrder(productJson *ProductJson) (*Order, error) {
    order := Order{}

    return &order, nil
}

func GetOrder(id int) (*Order, error) {
    order := Order{}

    return &order, nil
}

func UpdateOrderStatus(id int, status string) error {
    return nil
}