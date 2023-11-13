package web

import (
	"fmt"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Web server is up and running!!!")
}
