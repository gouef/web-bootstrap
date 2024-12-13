package web_bootstrap

import (
	"github.com/gouef/diago"
	"github.com/gouef/diago/extensions"
	"github.com/gouef/router"
	extensions2 "github.com/gouef/router/extensions"
)

type Bootstrap struct {
	router          *router.Router
	diago           *diago.Diago
	diagoExtensions []diago.DiagoExtension
}

func NewBootstrap() *Bootstrap {
	r := router.NewRouter()
	d := diago.NewDiago()

	return &Bootstrap{router: r, diago: d}
}

func (b *Bootstrap) GetRouter() *router.Router {
	return b.router
}

func (b *Bootstrap) GetDiago() *diago.Diago {
	return b.diago
}

func (b *Bootstrap) AddDiagoExtension(extension diago.DiagoExtension) {
	b.diagoExtensions = append(b.diagoExtensions, extension)
}

func (b Bootstrap) Boot() {
	r := b.router

	if !r.IsRelease() {
		d := b.diago
		d.AddExtension(extensions.NewDiagoLatencyExtension())
		d.AddExtension(extensions2.NewDiagoRouteExtension(r))

		n := r.GetNativeRouter()
		n.Use(diago.DiagoMiddleware(r, d))
	}
}
