//go:generate wire
// +build wireinject

package router

import (
	"github.com/google/wire"
	"nine-village-road/internal/service"
)

func newUserService() (*service.UserService, func(), error) {
	panic(wire.Build(service.ProviderUserSet))
}

