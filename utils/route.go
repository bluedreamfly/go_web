package utils

import (
	"net/http"
	//"strings"
	"strings"
	"regexp"
	"fmt"
	//"strconv"
)

type Route struct {
	Handles map[string]map[string]HandleFn
}

type Req struct {
	*http.Request
	Params map[string]interface{}
}

type HandleFn func(w http.ResponseWriter, r *Req)

func (route *Route) addHandle(path string, method string, handle HandleFn)  {

	if route.Handles == nil {
		route.Handles = map[string]map[string]HandleFn{}
	}

	if route.Handles[method] == nil {
		route.Handles[method] = map[string]HandleFn{}
	}
	route.Handles[method][path] = handle
}



func (route *Route) Post(path string, handle HandleFn)  {
	route.addHandle(path, "POST", handle)
}

func (route *Route) Get(path string, handle HandleFn)  {
	route.addHandle(path, "GET", handle)
}

func (route *Route) PUT(path string, handle HandleFn) {
	route.addHandle(path, "PUT", handle)
}

func (route *Route) DELETE(path string, handle HandleFn) {
	route.addHandle(path, "DELETE", handle)
}

func (route *Route) Start(w http.ResponseWriter, req *http.Request) {

	method := req.Method
	path := req.URL.Path

	var handle HandleFn
	paramsMap := map[string]interface{}{}

	if route.Handles != nil && route.Handles[method] != nil {
		handle = route.Handles[method][path]

		if handle == nil {
			for pathStr := range route.Handles[method] {
				if strings.Index(pathStr, ":") > -1 {

					reg := regexp.MustCompile(`:[^/]+`)

					params := reg.FindAllString(pathStr, -1)

					newPathStr := reg.ReplaceAllString(pathStr, `([^/]+)`)

					if result, _ := regexp.MatchString(newPathStr, path); result {

						handle = route.Handles[method][pathStr]
						reg = regexp.MustCompile(newPathStr)

						result := reg.FindStringSubmatch(path)

						paramsValue := result[1:]

						fmt.Printf("%v", result)
						for index, fieldValue := range params {
							paramsMap[fieldValue[1:]] = paramsValue[index]
						}

						return
					}

					println(newPathStr)
				}
			}
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 no fond"))
		}

		fmt.Printf("%v", paramsMap)


		newReq := Req{req, paramsMap}

		if handle != nil {
			handle(w, &newReq)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 no fond"))
	}

}