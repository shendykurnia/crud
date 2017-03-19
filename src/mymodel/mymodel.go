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
    Id int `json:"id"`
    ShopId int `json:"shop_id"`
    CustomerId int `json:"customer_id"`
    Status string `json:"status"`
    Products []Product `json:"products"` // tempting, but too scary to use []*Product here
    Created time.Time `json:"created"`
}

type Product struct {
    Id int `json:"id"`
    Name string `json:"name"`
    Created time.Time `json:"created"`
}

type ProductJson struct {
    Id int `json:"id"`
}

type CreateOrderJson struct {
    ShopId int `json:"shop_id"`
    CustomerId int `json:"customer_id"`
    Products []ProductJson `json:"products"`
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
        &Product{1, "Table", dummyTime()},
        &Product{2, "Chair", dummyTime()},
        &Product{3, "Pencil", dummyTime()},
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
                if strings.Index(strings.ToLower(product.Name), strings.ToLower(search)) != -1 {
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
    if len(collection) < limit {
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
    if len(createOrderJson.Products) == 0 {
        return nil, &ModelError{"Order needs at least 1 product"}
    }

    for _, product := range createOrderJson.Products {
        var foundProduct *Product
        foundProduct = nil
        for _, _product := range p.products {
            if _product.Id == product.Id {
                foundProduct = _product
            }
        }
        if foundProduct == nil {
            return nil, &ModelError{fmt.Sprintf("Invalid product (id: %d )", product.Id)}
        }
        products = append(products, foundProduct)
    }

    order := EfficientOrder{
        id: potentialId,
        shopId: createOrderJson.ShopId,
        customerId: createOrderJson.CustomerId,
        status: StatusCreated,
        products: products,
        created: dummyTime(),
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

    if _, isValidStatusChange := statusMap[order.Status][status]; !isValidStatusChange {
        return &ModelError{"Invalid status change"}
    }

    for _, _order := range p.orders {
        if order.Id == _order.id {
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
        Id: order.id,
        ShopId: order.shopId,
        CustomerId: order.customerId,
        Status: order.status,
        Products: products,
        Created: order.created,
    }
}

func cloneOrder(order *Order) *Order {
    products := []Product{}
    for _, product := range order.Products {
        products = append(products, (*cloneProduct(&product)))
    }
    return &Order{
        Id: order.Id,
        ShopId: order.ShopId,
        CustomerId: order.CustomerId,
        Status: order.Status,
        Products: products,
        Created: order.Created,
    }
}

func cloneProduct(product *Product) *Product {
    return &Product{
        Id: product.Id,
        Name: product.Name,
        Created: product.Created,
    }
}

// Cache
type Cache interface {
    init() error
    set(key string, obj *interface{}) error
    get(key string) (*interface{}, error)
}

type RedisCache struct {

}

func (c *RedisCache) init() error {
    return nil
}

func (c *RedisCache) set(key string, obj *interface{}) error {
    return nil
}

func (c *RedisCache) get(key string) (*interface{}, error) {
    return nil, nil
}

type MockCache struct {
    
}

func (c *MockCache) init() error {
    return nil
}

func (c *MockCache) set(key string, obj *interface{}) error {
    return nil
}

func (c *MockCache) get(key string) (*interface{}, error) {
    return nil, nil
}

// Queue
type Queue interface {
    init() error
    enqueue(obj *interface{}) error
    dequeue() (*interface{}, error)
}

type NsqQueue struct {

}

func (q *NsqQueue) init() error {
    return nil
}

func (q *NsqQueue) enqueue(obj *interface{}) error {
    return nil
}

func (q *NsqQueue) dequeue() (*interface{}, error) {
    return nil, nil
}

type MockQueue struct {

}

func (q *MockQueue) init() error {
    return nil
}

func (q *MockQueue) enqueue(obj *interface{}) error {
    return nil
}

func (q *MockQueue) dequeue() (*interface{}, error) {
    return nil, nil
}

// Stack
type Stack struct {
    datastore *Datastore
    cache *Cache
    queue *Queue
}

// functions
func GetOrders(stack *Stack, search string, page int) ([]*Order, error) {
    return []*Order{}, nil
}

func CreateOrder(stack *Stack, productJson *ProductJson) (*Order, error) {
    return &Order{}, nil
}

func GetOrder(stack *Stack, id int) (*Order, error) {
    return &Order{}, nil
}

func UpdateOrderStatus(stack *Stack, id int, status string) error {
    return nil
}

// helper functions
func dummyTime() time.Time {
    location, err := time.LoadLocation("UTC")
    if err != nil {
        panic(err)
    }
    return time.Date(2017, 3, 19, 0, 0, 0, 0, location)
}