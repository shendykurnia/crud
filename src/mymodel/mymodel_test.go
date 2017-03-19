package mymodel

import (
    "testing"
    "encoding/json"
    "strings"
)

type TestCaseGetOrderCollection struct {
    search string
    page int
    prepare func(testCase *TestCaseGetOrderCollection)
    checker func(testCase *TestCaseGetOrderCollection, orders []*Order)
}

func Test(t *testing.T) {
    // TODO test GetOrders, CreateOrder, GetOrder, UpdateOrderStatus, but since those functions are not utilizing cache or queue, datastore tests suffice for now
}

func TestMockDatastore(t *testing.T) {
    datastore := MockDatastore{}
    err := datastore.init()
    if err != nil {
        t.Fatal(err)
    }

    // test getOrderCollection
    for _, testCase := range []TestCaseGetOrderCollection{
        // empty database
        {"", 1,
            func(testCase *TestCaseGetOrderCollection) {},
            func(testCase *TestCaseGetOrderCollection, orders []*Order) {
                if len(orders) != 0 {
                    t.Errorf("Getting order collection %v %v returned unexpected result: orders is not empty", testCase.search, testCase.page)
                }
            },
        },

        {"", 2,
            func(testCase *TestCaseGetOrderCollection) {},
            func(testCase *TestCaseGetOrderCollection, orders []*Order) {
                if len(orders) != 0 {
                    t.Errorf("Getting order collection %v %v returned unexpected result: orders is not empty", testCase.search, testCase.page)
                }
            },
        },

        {"table", 1,
            func(testCase *TestCaseGetOrderCollection) {},
            func(testCase *TestCaseGetOrderCollection, orders []*Order) {
                if len(orders) != 0 {
                    t.Errorf("Getting order collection %v %v returned unexpected result: orders is not empty", testCase.search, testCase.page)
                }
            },
        },

        {"table", 2,
            func(testCase *TestCaseGetOrderCollection) {},
            func(testCase *TestCaseGetOrderCollection, orders []*Order) {
                if len(orders) != 0 {
                    t.Errorf("Getting order collection %v %v returned unexpected result: orders is not empty", testCase.search, testCase.page)
                }
            },
        },

        // non empty database
        {"", 1,
            func(testCase *TestCaseGetOrderCollection) {
                datastore.createOrder(&CreateOrderJson{1, 1, []ProductJson{ProductJson{1}}})
            },
            func(testCase *TestCaseGetOrderCollection, orders []*Order) {
                expected := []*Order{&Order{1, 1, 1, StatusCreated, []Product{Product{1, "Table", dummyTime()}}, dummyTime()}}
                if ok, jsonStr1, jsonStr2 := compareByJson(&orders, &expected); !ok {
                    t.Errorf("Getting order collection %v %v returned unexpected result: got %v expected %v", testCase.search, testCase.page, jsonStr1, jsonStr2)
                }
            },
        },

        {"", 2,
            func(testCase *TestCaseGetOrderCollection) {},
            func(testCase *TestCaseGetOrderCollection, orders []*Order) {
                expected := []*Order{}
                if ok, jsonStr1, jsonStr2 := compareByJson(&orders, &expected); !ok {
                    t.Errorf("Getting order collection %v %v returned unexpected result: got %v expected %v", testCase.search, testCase.page, jsonStr1, jsonStr2)
                }
            },
        },

        {"table", 1,
            func(testCase *TestCaseGetOrderCollection) {},
            func(testCase *TestCaseGetOrderCollection, orders []*Order) {
                expected := []*Order{&Order{1, 1, 1, StatusCreated, []Product{Product{1, "Table", dummyTime()}}, dummyTime()}}
                if ok, jsonStr1, jsonStr2 := compareByJson(&orders, &expected); !ok {
                    t.Errorf("Getting order collection %v %v returned unexpected result: got %v expected %v", testCase.search, testCase.page, jsonStr1, jsonStr2)
                }
            },
        },

        {"table", 2,
            func(testCase *TestCaseGetOrderCollection) {},
            func(testCase *TestCaseGetOrderCollection, orders []*Order) {
                expected := []*Order{}
                if ok, jsonStr1, jsonStr2 := compareByJson(&orders, &expected); !ok {
                    t.Errorf("Getting order collection %v %v returned unexpected result: got %v expected %v", testCase.search, testCase.page, jsonStr1, jsonStr2)
                }
            },
        },

        {"z", 1,
            func(testCase *TestCaseGetOrderCollection) {},
            func(testCase *TestCaseGetOrderCollection, orders []*Order) {
                expected := []*Order{}
                if ok, jsonStr1, jsonStr2 := compareByJson(&orders, &expected); !ok {
                    t.Errorf("Getting order collection %v %v returned unexpected result: got %v expected %v", testCase.search, testCase.page, jsonStr1, jsonStr2)
                }
            },
        },

        {"z", 2,
            func(testCase *TestCaseGetOrderCollection) {},
            func(testCase *TestCaseGetOrderCollection, orders []*Order) {
                expected := []*Order{}
                if ok, jsonStr1, jsonStr2 := compareByJson(&orders, &expected); !ok {
                    t.Errorf("Getting order collection %v %v returned unexpected result: got %v expected %v", testCase.search, testCase.page, jsonStr1, jsonStr2)
                }
            },
        },

        // testing reference
        {"", 1,
            func(testCase *TestCaseGetOrderCollection) {},
            func(testCase *TestCaseGetOrderCollection, orders []*Order) {
                if len(orders) < 1 {
                    t.Errorf("Getting order collection %v %v returned unexpected result: size is unexpected", testCase.search, testCase.page)
                    return
                }
                // change value of the returned orders
                orders[0].ShopId = 99

                // refetch
                expected, err := datastore.getOrderCollection(testCase.search, testCase.page)
                if err != nil {
                    t.Errorf("Getting order %v %v collection returned error: %v", testCase.search, testCase.page, err)
                }

                // expecting results to be different
                if ok, jsonStr1, jsonStr2 := compareByJson(&orders, &expected); ok {
                    t.Errorf("Getting order collection %v %v returned unexpected result: got %v expected %v", testCase.search, testCase.page, jsonStr1, jsonStr2)
                }

                if len(orders) != len(expected) {
                    t.Errorf("Getting order collection %v %v returned unexpected result: got %v expected %v", testCase.search, testCase.page, len(orders), len(expected))
                }

                if len(expected) < 1 {
                    t.Errorf("Getting order collection %v %v returned unexpected result: size expected is unexpected", testCase.search, testCase.page)
                    return
                }
                expected[0].ShopId = 99
                // now that I set the same value, expecting results to be the same
                if ok, jsonStr1, jsonStr2 := compareByJson(&orders, &expected); !ok {
                    t.Errorf("Getting order collection %v %v returned unexpected result: got %v expected %v", testCase.search, testCase.page, jsonStr1, jsonStr2)
                }
            },
        },

    } {
        testCase.prepare(&testCase)

        orders, err := datastore.getOrderCollection(testCase.search, testCase.page)
        if err != nil {
            t.Errorf("Getting order %v %v collection returned error: %v", testCase.search, testCase.page, err)
        }

        testCase.checker(&testCase, orders)
    }


    // test createOrder and getOrder
    createOrderJson := CreateOrderJson{1, 1, []ProductJson{ProductJson{99}}}
    order, err := datastore.createOrder(&createOrderJson)
    if err == nil {
        t.Errorf("Creating order %v is unexpected: an error should be returned", createOrderJson)
    }

    // test reference
    createOrderJson = CreateOrderJson{1, 1, []ProductJson{ProductJson{1}}}
    order, err = datastore.createOrder(&createOrderJson)
    if err != nil {
        t.Errorf("Creating order %v returned an error: %v", createOrderJson, err)
    }
    if order == nil {
        t.Errorf("Creating order %v returned unexpected result: order is empty", createOrderJson)
    }
    expected := Order{order.Id, 1, 1, StatusCreated, []Product{Product{1, "Table", dummyTime()}}, dummyTime()}
    if ok, jsonStr1, jsonStr2 := compareByJson(order, &expected); !ok {
        t.Errorf("Creating order %v returned unexpected result: got %v expected %v", createOrderJson, jsonStr1, jsonStr2)
    }

    order.ShopId = 99

    _order, err := datastore.getOrder(order.Id)
    if err != nil {
        t.Errorf("Getting an order %v returned an error: %v", order.Id, err)
    }
    if ok, jsonStr1, jsonStr2 := compareByJson(_order, &expected); !ok {
        t.Errorf("Getting an order %v returned unexpected result: got %v expected %v", order.Id, jsonStr1, jsonStr2)
    }
    if ok, jsonStr1, jsonStr2 := compareByJson(_order, order); ok {
        t.Errorf("Getting an order %v returned unexpected result: object is changed by reference", order.Id, jsonStr1, jsonStr2)
    }

    _order, err = datastore.getOrder(99)
    if err == nil {
        t.Errorf("Getting an order %v returned unexpected result: should return an error", order.Id)
    }

    // test updateOrderStatus
    _order, err = datastore.getOrder(1)
    if err != nil {
        t.Errorf("Getting an order %v returned an error: %v", 1, err)
    }

    // created cannot change to created
    if err := datastore.updateOrderStatus(_order.Id, StatusCreated); err == nil {
        t.Errorf("Updating order %v status returned an unexpected result: should return an error", order.Id)
    }

    if err := datastore.updateOrderStatus(_order.Id, StatusProcessed); err != nil {
        t.Errorf("Updating order %v status returned an unexpected result: %e", order.Id, err)
    }

    if err := datastore.updateOrderStatus(_order.Id, "z"); err == nil {
        t.Errorf("Updating order %v status returned an unexpected result: should return an error", order.Id)
    }
}

func TestMockCache(t *testing.T) {
    // TODO
}

func TestMockQueue(t *testing.T) {
    // TODO
}

func compareByJson(obj1 interface{}, obj2 interface{}) (bool, string, string) {
    json1, err := json.Marshal(obj1)
    if err != nil {
        panic(err)
    }
    json2, err := json.Marshal(obj2)
    if err != nil {
        panic(err)
    }
    jsonStr1 := string(json1)
    jsonStr2 := string(json2)
    return strings.Compare(jsonStr1, jsonStr2) == 0, jsonStr1, jsonStr2
}