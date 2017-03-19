package mymodel

import (
    "time"
    "sort"
    "fmt"
    "strings"
)

const (
    StatusCreated = "created"
    StatusProcessed = "processed"
    StatusCanceled = "canceled"
    StatusFinished = "finished"

    Limit = 10
)

type Order struct {
    id int
    shopId int
    customerId int
    status string
    products []Product // tempting, but too scary to use []*Product here
    created time.Time
}

type Product struct {
    id int
    name string
    created time.Time
}

type ProductJson struct {
    id int `json:"id"`
}

type CreateOrderJson struct {
    shopId int `json:"shop_id"`
    customerId int `json:"customer_id"`
    products []ProductJson `json:"products"`
}

type ModelError struct {
    message string
}

func (e *ModelError) Error() string {
    return fmt.Sprintf("Object error: %s", e.message)
}

// Datastores
type Datastore interface {
    init() error
    getOrderCollection(search string, page int) ([]*Order, error)
    getOrder(id int) (*Order, error)
    createOrder(createOrderJson *CreateOrderJson) (*Order, error)
    updateOrderStatus(id int, status string) error
}

// Postgres datastore

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

func (p *PostgresDatastore) createOrder(createOrderJson *CreateOrderJson) (*Order, error) {
    order := Order{}

    return &order, nil
}

func (p *PostgresDatastore) updateOrderStatus(id int, status string) error {
    return nil
}

// Mock datastore

// like Order but more space-efficient because of []*Product
type EfficientOrder struct {
    id int
    shopId int
    customerId int
    status string
    products []*Product
    created time.Time
}

type MockDatastore struct {
    orders []*EfficientOrder
    products []*Product
}

func (p *MockDatastore) init() error {
    p.orders = []*EfficientOrder{}
    p.products = []*Product{
        &Product{1, "Table", time.Now()},
        &Product{2, "Chair", time.Now()},
        &Product{3, "Pencil", time.Now()},
    }
    return nil
}

func (p *MockDatastore) getOrderCollection(search string, page int) ([]*Order, error) {
    var collection []*EfficientOrder

    if len(search) == 0 {
        collection = p.orders
    } else {
        for _, order := range p.orders {
            isMatch := false
            for _, product := range order.products {
                // simple matching logic, please don't judge my development skill based on this. this is just for simplicity sake
                if strings.Index(strings.ToLower(product.name), strings.ToLower(search)) != -1 {
                    isMatch = true
                    break
                }
            }
            if isMatch {
                collection = append(collection, order)
            }
        }
    }

    // pagination
    offset := ((page - 1) * Limit)
    if offset > len(collection) || offset < 0 {
        return []*Order{}, nil
    }
    limit := offset + Limit
    if len(collection) > limit {
        limit = len(collection)
    }
    collection = collection[offset:limit]

    // clone order to be returned, so caller cannot change stored object directly
    orders := []*Order{}
    for _, order := range collection {
        orders = append(orders, createOrder(order))
    }

    return orders, nil
}

func (p *MockDatastore) getOrder(id int) (*Order, error) {
    for _, order := range p.orders {
        if order.id == id {
            // clone order to be returned, so caller cannot change stored object directly
            return createOrder(order), nil
        }
    }

    return nil, &ModelError{"Invalid id"}
}

func (p *MockDatastore) createOrder(createOrderJson *CreateOrderJson) (*Order, error) {
    // generate id
    potentialId := 1
    orderIds := []int{}
    for _, order := range p.orders {
        orderIds = append(orderIds, order.id)
    }
    sort.Ints(orderIds)
    for _, orderId := range orderIds {
        if potentialId == orderId {
            potentialId++
            continue
        } else if potentialId < orderId {
            break
        }
    }

    var products []*Product
    if len(createOrderJson.products) == 0 {
        return nil, &ModelError{"Order needs at least 1 product"}
    }

    for _, product := range createOrderJson.products {
        var foundProduct *Product
        foundProduct = nil
        for _, _product := range p.products {
            if _product.id == product.id {
                foundProduct = _product
            }
        }
        if foundProduct == nil {
            return nil, &ModelError{fmt.Sprintf("Invalid product (id: %d )", product.id)}
        }
        products = append(products, foundProduct)
    }

    order := EfficientOrder{
        id: potentialId,
        shopId: createOrderJson.shopId,
        customerId: createOrderJson.customerId,
        status: StatusCreated,
        products: products,
        created: time.Now(),
    }

    p.orders = append(p.orders, &order)

    // clone order to be returned, so caller cannot change stored object directly
    return createOrder(&order), nil
}

func (p *MockDatastore) updateOrderStatus(id int, status string) error {
    order, err := p.getOrder(id)
    if err != nil {
        return err
    }

    // map holding the business logic governing status change, eg: finished order can be changed to created order
    statusMap := map[string]map[string]bool{
        StatusCreated: map[string]bool{StatusProcessed: true},
        StatusProcessed: map[string]bool{StatusCanceled: true, StatusFinished: true},
        StatusCanceled: map[string]bool{},
        StatusFinished: map[string]bool{},
    }

    if _, isValidStatus := statusMap[status]; !isValidStatus {
        return &ModelError{"Invalid status"}
    }

    if _, isValidStatusChange := statusMap[order.status][status]; !isValidStatusChange {
        return &ModelError{"Invalid status change"}
    }

    for _, _order := range p.orders {
        if order.id == _order.id {
            _order.status = status
            return nil
        }
    }

    return &ModelError{"Failed to update order status somehow"}
}

func createOrder(order *EfficientOrder) *Order {
    products := []Product{}
    for _, product := range order.products {
        products = append(products, (*cloneProduct(product)))
    }
    return &Order{
        id: order.id,
        shopId: order.shopId,
        customerId: order.customerId,
        status: order.status,
        products: products,
        created: order.created,
    }
}

func cloneOrder(order *Order) *Order {
    products := []Product{}
    for _, product := range order.products {
        products = append(products, (*cloneProduct(&product)))
    }
    return &Order{
        id: order.id,
        shopId: order.shopId,
        customerId: order.customerId,
        status: order.status,
        products: products,
        created: order.created,
    }
}

func cloneProduct(product *Product) *Product {
    return &Product{
        id: product.id,
        name: product.name,
    }
}

// Cache

// Queue