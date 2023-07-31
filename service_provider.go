package cloudinary

import (
	"github.com/goravel/framework/contracts/foundation"
	"golang.org/x/net/context"
)

const Binding = "goravel.cloudinary"

var App foundation.Application

type ServiceProvider struct {
}

func (receiver *ServiceProvider) Register(app foundation.Application) {
	App = app

	app.BindWith(Binding, func(app foundation.Application, parameters map[string]any) (any, error) {
		return NewCloudinary(context.Background(), app.MakeConfig(), parameters["disk"].(string))
	})
}
func (receiver *ServiceProvider) Boot(app foundation.Application) {

}
