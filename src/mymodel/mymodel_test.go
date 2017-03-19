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

type TestCaseCreateOrder struct {
    createOrderJson *CreateOrderJson
    prepare func(testCase *TestCaseCreateOrder)
    checker func(testCase *TestCaseCreateOrder, order *Order)
}

func Test(t *testing.T) {
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
                var obj1 interface{}
                var obj2 interface{}
                obj1 = orders
                obj2 = expected
                if ok, jsonStr1, jsonStr2 := compareByJson(&obj1, &obj2); !ok {
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


    // test createOrder
    // for _, testCase := range []struct{
    //     createOrderJson *CreateOrderJson,
    //     prepare func(),
    //     checker func(order *Order),
    // }{
    //     {&CreateOrderJson{1, 1, []ProductJson{2}},
    //         func() {},
    //         func(order *Order) {
    //             if len(orders) != 0 {
    //                 t.Errorf("Getting order collection returned unexpected result: orders is not empty")
    //             }
    //         }
    //     },
        
    // } {
    //     testCase.prepare()

    //     order, err := datastore.createOrder(testCase.createOrderJson)
    //     if err != nil {
    //         t.Errorf("Creating order %v collection returned error: %v", testCase.createOrderJson, err)
    //     }

    //     testCase.checker(order)
    // }

    // datastore.getOrderCollection("", 2)
    // datastore.getOrderCollection("table", 1)
    // datastore.getOrderCollection("table", 2)
    // datastore.getOrderCollection("z", 1)
    // datastore.getOrderCollection("z", 2)

    // getOrder
    // updateOrderStatus
    // datastore.getOrderCollection("", 1)
    // test if changing returned value change underlying data
    // t.Errorf("Handler for %v %v returned wrong HTTP status code: got %v expected %v", testCase.method, testCase.path, status, testCase.expectedStatus)
}

func TestMockCache(t *testing.T) {
}

func TestMockQueue(t *testing.T) {
}

func compareByJson(obj1 *interface{}, obj2 *interface{}) (bool, string, string) {
    json1, err := json.Marshal(*obj1)
    if err != nil {
        panic(err)
    }
    json2, err := json.Marshal(*obj2)
    if err != nil {
        panic(err)
    }
    jsonStr1 := string(json1)
    jsonStr2 := string(json2)
    return strings.Compare(jsonStr1, jsonStr2) == 0, jsonStr1, jsonStr2
}