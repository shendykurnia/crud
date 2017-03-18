package myutil

import (
    "testing"
)

func Test(t *testing.T) {
    for _, testCase := range []struct {
        url string
        params [][]string
        expectedUrl string
    } {
        {"http://a", [][]string{{"a", "b"}}, "http://a?a=b"},
        {"http://b", [][]string{{"c", "d"}, {"e", "f"}}, "http://b?c=d&e=f"},
        {"http://c", [][]string{{"a[]", "b"}, {"a[]", "c"}}, "http://c?a%5B%5D=b&a%5B%5D=c"},
    } {
        url, err := ConstructUrl(testCase.url, testCase.params)
        if err != nil {
            t.Errorf("Test case %v %v returned error: %v", testCase.url, testCase.params, err)
            continue
        }

        if url != testCase.expectedUrl {
            t.Errorf("Test case %v %v returned unexpected url: got %v expected %v", testCase.url, testCase.params, url, testCase.expectedUrl)
        }
    }
}