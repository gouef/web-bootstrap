package web_bootstrap

import (
	"github.com/gouef/diago"
	"github.com/gouef/diago/extensions"
	"github.com/gouef/router"
	extensions2 "github.com/gouef/router/extensions"
)

type BootstrapInterface struct {
	router          *router.Router
	diago           *diago.Diago
	diagoExtensions []diago.Extension
}

func Bootstrap() *BootstrapInterface {
	return NewBootstrap()
}

func NewBootstrap() *BootstrapInterface {
	r := router.NewRouter()
	d := diago.NewDiago()

	return &BootstrapInterface{router: r, diago: d}
}

func (b *BootstrapInterface) GetRouter() *router.Router {
	return b.router
}

func (b *BootstrapInterface) GetDiago() *diago.Diago {
	return b.diago
}

func (b *BootstrapInterface) AddDiagoExtension(extension diago.Extension) {
	b.diagoExtensions = append(b.diagoExtensions, extension)
}

func (b *BootstrapInterface) Boot() {
	r := b.router

	if !r.IsRelease() {
		d := b.diago
		d.AddExtension(extensions.NewDiagoLatencyExtension())
		d.AddExtension(extensions2.NewDiagoRouteExtension(r))

		n := r.GetNativeRouter()
		n.Use(diago.DiagoMiddleware(r, d))
	}
}
