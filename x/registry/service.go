package registry

import (
	"github.com/wenewzhang/core/service"
)

type serviceRegistry struct {
	registry[service.Service]
}
