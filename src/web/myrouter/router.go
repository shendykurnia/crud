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
    handler func(http.ResponseWriter, *http.Request, *map[string]string)
}

type MyHandler struct {
    routes []*route
}

func (h *MyHandler) HandleFunc(patternStr string, methods map[string]interface{}, handler func(http.ResponseWriter, *http.Request, *map[string]string)) {
    pattern := regexp.MustCompile(patternStr)
    h.routes = append(h.routes, &route{patternStr, pattern, methods, handler})
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

        path := r.URL.Path
        if !validMethod || !route.pattern.MatchString(path) {
            continue
        }

        slugs := map[string]string{}
        matches := route.pattern.FindStringSubmatch(path)
        names := route.pattern.SubexpNames()
        for i, match := range matches {
                if i != 0 {
                    slugs[names[i]] = match
                }
        }

        route.handler(w, r, &slugs)
        return
    }

    http.NotFound(w, r)
}