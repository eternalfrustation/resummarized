package middlewares

import (
	"fmt"
	"net/http"
	"strings"
)

type MethodMux struct {
	handlers map[string]http.Handler
}

func NewMethodMux() MethodMux {
	return MethodMux{
		handlers: map[string]http.Handler{},
	}
}

func (mux *MethodMux) With(method string, handler http.Handler) *MethodMux {
	mux.handlers[method] = handler
	return mux
}

func (mux *MethodMux) ServerHttp(w http.ResponseWriter, r *http.Request) {
	handler, ok := mux.handlers[r.Method]
	if !ok {
		available_methods := []string{}
		for k := range mux.handlers {
			available_methods = append(available_methods, k)
		}
		http.Error(w, fmt.Sprintf("%s Method not supported, supported methods are: %s", r.Method, strings.Join(available_methods, ",")), http.StatusNotFound)
		return
	}
	handler.ServeHTTP(w, r)
}
