package web_bootstrap

import (
	"github.com/gouef/diago"
	"github.com/gouef/diago/extensions"
	"github.com/gouef/renderer"
	"github.com/gouef/router"
	extensions2 "github.com/gouef/router/extensions"
)

type BootstrapInterface struct {
	router *router.Router
}

func Bootstrap() *BootstrapInterface {
	return NewBootstrap()
}

func NewBootstrap() *BootstrapInterface {
	r := router.NewRouter()
	n := r.GetNativeRouter()

	if !r.IsRelease() {
		d := diago.NewDiago()
		d.AddExtension(extensions.NewLatencyExtension())
		d.AddExtension(extensions2.NewDiagoRouteExtension(r))

		n.Use(diago.Middleware(r, d))
	}

	n.SetTrustedProxies([]string{"127.0.0.1"})
	renderer.RegisterToRouter(r, "./views/templates")

	return &BootstrapInterface{router: r}
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
