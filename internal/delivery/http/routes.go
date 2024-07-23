package routes

import (
	"net/http"

	"github.com/jedyEvgeny/time_tracker/internal/app/endpoint"
)

func SetupRoutes(e *endpoint.Endpoint) {
	http.HandleFunc("/addUser", e.StatusAdd)       // Обработчик с методом GET
	http.HandleFunc("/delUser", e.StatusDel)       // Обработчик с методом DELETE
	http.HandleFunc("/changeUser", e.StatusChange) // Обработчик с методом PATCH
}
