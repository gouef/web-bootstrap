package web_bootstrap

import (
	"dario.cat/mergo"
	"github.com/gouef/diago"
	"github.com/gouef/diago/extensions"
	"github.com/gouef/renderer"
	"github.com/gouef/router"
	extensions2 "github.com/gouef/router/extensions"
	"log"
)

type BootstrapInterface struct {
	router  *router.Router
	configs []*Config
	config  *Config
}

func Bootstrap() *BootstrapInterface {
	return NewBootstrap()
}

func NewBootstrap() *BootstrapInterface {
	r := router.NewRouter()
	return &BootstrapInterface{router: r}
}

func (b *BootstrapInterface) AddConfig(path string) *BootstrapInterface {
	cfg, err := LoadConfig(path)

	if err != nil {
		log.Println("unable load config. ", err.Error())
		return b
	}

	b.configs = append(b.configs, cfg)
	return b
}

func (b *BootstrapInterface) LoadConfiguration() *Config {
	if len(b.configs) == 0 {
		b.config = DefaultConfig()
		return b.config
	}

	merged := *b.configs[0]

	for _, cfg := range b.configs[1:] {
		if err := mergo.Merge(&merged, cfg, mergo.WithOverride); err != nil {
			log.Println("error merging config:", err)
		}
	}

	b.config = &merged

	return b.config
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
