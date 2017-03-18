package mymodel

import (
    "time"
)

const (
    STATUS_CREATED = "created"
    STATUS_PROCESSED = "processed"
    STATUS_CANCELED = "canceled"
    STATUS_FINISHED = "finished"
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