package endpoint

import (
	"net/http"
	"time"

	"github.com/jedyEvgeny/time_tracker/pkg/storage"
)

type Decoder interface {
	DecodeJSON(*http.Request) (string, string, error)
	EnrichUserData(*http.Response, string, string) (storage.EnrichedUser, error)
	ChangeUserData(*http.Request) (storage.EnrichedUser, error)
	GetPudding([]storage.EnrichedUser) ([]byte, error)
	DecodeJSONTask(*http.Request) (storage.TaskEnrichedUser, error)
	Now() (time.Time, error)
	DecodeJsonTaskDur(*http.Request) (storage.TaskEnrichedUser, error)
	GetSortTasks([]storage.UserTask) ([]byte, error)
}

type Adder interface {
	AddPerson(storage.EnrichedUser) error
	DelPerson(string, string) error
	ChangePerson(storage.EnrichedUser) error
	GetUsersByFilter(storage.EnrichedUser) ([]storage.EnrichedUser, error)
	AddStartTask(storage.TaskEnrichedUser, time.Time) error
	AddFinishTask(storage.TaskEnrichedUser, time.Time) error
	FindTimeTask(storage.TaskEnrichedUser) ([]storage.UserTask, error)
}

type EndpointCaller interface {
	CallEndpoint(string, string) (*http.Response, error)
}
