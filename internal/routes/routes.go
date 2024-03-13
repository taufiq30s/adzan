package routes

import (
	"net/http"
)

func NewRoute() http.Handler {
	mux := http.NewServeMux()

	return mux
}
