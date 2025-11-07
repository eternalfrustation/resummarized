package middlewares

import (
	"net/http"

	"github.com/gorilla/schema"
)

func FormHandler[T any](next func(http.ResponseWriter, *http.Request, T)) http.Handler {
	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("Invalid Submission"))
		}
		var formData T

		err = decoder.Decode(&formData, r.Form)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		next(w, r, formData)
	})
}
