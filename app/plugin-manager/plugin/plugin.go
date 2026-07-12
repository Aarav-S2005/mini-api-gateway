package plugin

import (
	"net/http"
)

type Plugin interface {
	Name() string
	Validate(config map[string]interface{}) error
	CreateMiddleware(config map[string]interface{}) (func(next http.Handler) http.Handler, error)
}
