package routes

import (
	"net/http"

	"github.com/jedyEvgeny/time_tracker/internal/app/endpoint"
)

func SetupRoutes(e *endpoint.Endpoint) {
	http.HandleFunc("/addUser", e.StatusAdd)       // Обработчик с методом POST
	http.HandleFunc("/delUser", e.StatusDel)       // Обработчик с методом DELETE
	http.HandleFunc("/changeUser", e.StatusChange) // Обработчик с методом PATCH
	http.HandleFunc("/getUsers", e.StatusFilter)   //Обработчик с методом GET
	http.HandleFunc("/startTask", e.StatusStart)   //Обработчик с методом PUT
	http.HandleFunc("/finishTask", e.StatusFinish) //Обработчик с методом PUT
	http.HandleFunc("/duration", e.StatusDur)      //Обработчик с методом GET
}
