package registry

import (
	"github.com/orlandorode97/go-disptach/infrastructure/api"
	"github.com/orlandorode97/go-disptach/interface/controller"
)

type registry struct {
	ImaggaClient *api.Imagga
}
type Register interface {
	NewAppController() controller.AppController
}

func NewRegistry(imagga *api.Imagga) Register {
	return &registry{imagga}
}

func (r *registry) NewAppController() controller.AppController {
	return controller.AppController{
		Words: r.NewWordsController(),
	}
}
