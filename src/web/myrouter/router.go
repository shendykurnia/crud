package myrouter

import (
    "regexp"
    "net/http"
    "strings"
)

type route struct {
    patternStr string
    pattern *regexp.Regexp
    methods map[string]interface{}
    handler http.Handler
}

type MyHandler struct {
    routes []*route
}

func (h *MyHandler) HandleFunc(patternStr string, methods map[string]interface{}, handler func(http.ResponseWriter, *http.Request)) {
    pattern := regexp.MustCompile(patternStr)
    h.routes = append(h.routes, &route{patternStr, pattern, methods, http.HandlerFunc(handler)})
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    for _, route := range h.routes {
        // check method first because it's cheaper compared to regex match
        validMethod := false
        for m, _ := range route.methods {
            if strings.ToLower(r.Method) == strings.ToLower(m) {
                validMethod = true
                break
            }
        }

        // change pattern from /a/{id}/b to /a/.*/b
        patternStr := route.patternStr
        re := regexp.MustCompile(`{.*}`)
        patternStrToMatch := re.ReplaceAllString(patternStr, `.*`)
        patternToMatch := regexp.MustCompile(patternStrToMatch)

        path := r.URL.Path
        if !validMethod || !patternToMatch.MatchString(path) {
            continue
        }

        route.handler.ServeHTTP(w, r)
        return
    }

    http.NotFound(w, r)
}

func GetUrlSlug(path string, ) map[string]string {

}