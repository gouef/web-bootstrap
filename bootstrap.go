package web_bootstrap

import (
	"log"

	"dario.cat/mergo"
	"github.com/gouef/diago"
	"github.com/gouef/diago/extensions"
	"github.com/gouef/gorm"
	"github.com/gouef/renderer"
	"github.com/gouef/router"
	extensions2 "github.com/gouef/router/extensions"
	o_gorm "gorm.io/gorm"
)

type BootstrapInterface struct {
	router  *router.Router
	configs []*Config
	config  *Config
	gorm    *o_gorm.DB
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
	cfg := b.LoadConfiguration()

	b.LoadRouter(cfg)
	b.LoadGorm(cfg)

}

func (b *BootstrapInterface) LoadRouter(cfg *Config) {
	r := b.router
	n := r.GetNativeRouter()

	if !r.IsRelease() && cfg.Diago.Enabled {
		d := diago.NewDiago()
		d.AddExtension(extensions.NewLatencyExtension())
		d.AddExtension(extensions2.NewDiagoRouteExtension(r))

		n.Use(diago.Middleware(r, d))
	}

	for _, s := range cfg.Router.Statics {
		n.Static(s.Path, s.Root)
	}

	if len(cfg.Router.Proxy.Trust) > 0 {
		err := n.SetTrustedProxies(cfg.Router.Proxy.Trust)
		if err != nil {
			log.Println("unable set trusted proxies.", err.Error())
		}
	}

	rend := renderer.NewRenderer(cfg.Renderer.Dir, cfg.Renderer.Layout)
	rend.RegisterRouter(r)
}

func (b *BootstrapInterface) LoadGorm(cfg *Config) {

	if cfg.Gorm.Driver != "" {

		if cfg.Gorm.TimeZone == "" {
			cfg.Gorm.TimeZone = "UTC"
		}

		db, err := gorm.New(cfg.Gorm.ToGormConfig())
		if err != nil {
			log.Println("unable initialize gorm. ", err.Error())
		} else {
			b.gorm = db
			log.Println("connected to gorm database")
		}
	}
}
