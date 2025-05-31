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
	return &BootstrapInterface{router: r}
}

func (b *BootstrapInterface) GetRouter() *router.Router {
	return b.router
}

func (b *BootstrapInterface) Static(relativePath string, root string) *BootstrapInterface {
	b.GetRouter().GetNativeRouter().Static(relativePath, root)
	return b
}

func (b *BootstrapInterface) Boot() {
	r := b.router
	n := r.GetNativeRouter()

	if !r.IsRelease() {
		d := diago.NewDiago()
		d.AddExtension(extensions.NewLatencyExtension())
		d.AddExtension(extensions2.NewDiagoRouteExtension(r))

		n.Use(diago.Middleware(r, d))
	}

	n.Static("/static", "./static")
	n.Static("/assets", "./static/assets")

	n.SetTrustedProxies([]string{"127.0.0.1"})
	renderer.RegisterToRouter(r, "./views/templates")

}
