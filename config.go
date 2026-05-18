package web_bootstrap

import (
	"crypto/tls"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Custom map[string]any

type ConfigInterface interface {
	UnmarshalYAML(value *yaml.Node) error
}

type Config struct {
	Parameters map[string]any     `yaml:"parameters"`
	Renderer   RendererConfig     `yaml:"renderer"`
	Router     RouterConfig       `yaml:"router"`
	Cache      CacheConfig        `yaml:"cache"`
	Diago      DiagoConfig        `yaml:"diago"`
	Gorm       GormDatabaseConfig `yaml:"gorm"`
	Custom     Custom
}

type DiagoConfig struct {
	Enabled bool `yaml:"enabled"`
}

type RendererConfig struct {
	Dir    string   `yaml:"dir"`
	Layout []string `yaml:"layout"`
	Custom Custom
}

type RouterConfig struct {
	Statics []RouterStaticConfig `yaml:"statics"`
	Proxy   RouterProxyConfig    `yaml:"proxy"`
	Custom  Custom
}

type RouterStaticConfig struct {
	Path   string `yaml:"path"`
	Root   string `yaml:"root"`
	Custom Custom
}

type RouterProxyConfig struct {
	Trust  []string `yaml:"trust"`
	Custom Custom
}

type CacheConfig struct {
	Storages []CacheStorageItemConfig `yaml:"storages"`
	Custom   Custom
}

type CacheStorageItemConfig struct {
	Type     string            `yaml:"type"`
	Instance string            `yaml:"instance"`
	Name     string            `yaml:"name"`
	File     CacheFileConfig   `yaml:"file"`
	Memory   CacheMemoryConfig `yaml:"memory"`
	Redis    CacheRedisConfig  `yaml:"redis"`
	Custom   Custom
}

type CacheRedisConfig struct {
	Name    string            `yaml:"name"`
	Options CacheRedisOptions `yaml:"options"`
}

type CacheRedisOptions struct {
	// The network type, either tcp or unix.
	// Default is tcp.
	Network string
	// host:port address.
	Addr string

	// ClientName will execute the `CLIENT SETNAME ClientName` command for each conn.
	ClientName string

	// Protocol 2 or 3. Use the version to negotiate RESP version with redis-server.
	// Default is 3.
	Protocol int
	// Use the specified Username to authenticate the current connection
	// with one of the connections defined in the ACL list when connecting
	// to a Redis 6.0 instance, or greater, that is using the Redis ACL system.
	Username string
	// Optional password. Must match the password specified in the
	// requirepass server configuration option (if connecting to a Redis 5.0 instance, or lower),
	// or the User Password when connecting to a Redis 6.0 instance, or greater,
	// that is using the Redis ACL system.
	Password string

	// Database to be selected after connecting to the server.
	DB int

	// Maximum number of retries before giving up.
	// Default is 3 retries; -1 (not 0) disables retries.
	MaxRetries int
	// Minimum backoff between each retry.
	// Default is 8 milliseconds; -1 disables backoff.
	MinRetryBackoff time.Duration
	// Maximum backoff between each retry.
	// Default is 512 milliseconds; -1 disables backoff.
	MaxRetryBackoff time.Duration

	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout time.Duration
	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Supported values:
	//   - `0` - default timeout (3 seconds).
	//   - `-1` - no timeout (block indefinitely).
	//   - `-2` - disables SetReadDeadline calls completely.
	ReadTimeout time.Duration
	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.  Supported values:
	//   - `0` - default timeout (3 seconds).
	//   - `-1` - no timeout (block indefinitely).
	//   - `-2` - disables SetWriteDeadline calls completely.
	WriteTimeout time.Duration
	// ContextTimeoutEnabled controls whether the client respects context timeouts and deadlines.
	// See https://redis.uptrace.dev/guide/go-redis-debugging.html#timeouts
	ContextTimeoutEnabled bool

	// Type of connection pool.
	// true for FIFO pool, false for LIFO pool.
	// Note that FIFO has slightly higher overhead compared to LIFO,
	// but it helps closing idle connections faster reducing the pool size.
	PoolFIFO bool
	// Base number of socket connections.
	// Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
	// If there is not enough connections in the pool, new connections will be allocated in excess of PoolSize,
	// you can limit it through MaxActiveConns
	PoolSize int
	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	PoolTimeout time.Duration
	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	// Default is 0. the idle connections are not closed by default.
	MinIdleConns int
	// Maximum number of idle connections.
	// Default is 0. the idle connections are not closed by default.
	MaxIdleConns int
	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActiveConns int
	// ConnMaxIdleTime is the maximum amount of time a connection may be idle.
	// Should be less than server's timeout.
	//
	// Expired connections may be closed lazily before reuse.
	// If d <= 0, connections are not closed due to a connection's idle time.
	//
	// Default is 30 minutes. -1 disables idle timeout check.
	ConnMaxIdleTime time.Duration
	// ConnMaxLifetime is the maximum amount of time a connection may be reused.
	//
	// Expired connections may be closed lazily before reuse.
	// If <= 0, connections are not closed due to a connection's age.
	//
	// Default is to not close idle connections.
	ConnMaxLifetime time.Duration

	// TLS Config to use. When set, TLS will be negotiated.
	TLSConfig *tls.Config

	// Enables read only queries on slave/follower nodes.
	readOnly bool

	// Disable set-lib on connect. Default is false.
	DisableIndentity bool

	// Add suffix to client name. Default is empty.
	IdentitySuffix string

	// UnstableResp3 enables Unstable mode for Redis Search module with RESP3.
	UnstableResp3 bool
}

type CacheRedisOptionsTLS struct {
}

type CacheMemoryConfig struct {
	Name string `yaml:"name"`
}

type CacheFileConfig struct {
	Name string `yaml:"name"`
	Dir  string `yaml:"dir"`
}

func DefaultConfig() *Config {
	rootDir, _ := filepath.Abs(".")
	cfg := Config{}
	cfg.Parameters = map[string]any{
		"rootDir": rootDir,
	}
	cfg.Renderer = RendererConfig{Dir: "./views/templates", Layout: []string{"@layout", "base", "layout"}}
	cfg.Router.Statics = []RouterStaticConfig{
		{Path: "/static", Root: "./static"},
		{Path: "/assets", Root: "./static/assets"},
	}
	cfg.Diago.Enabled = true
	return &cfg
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := DefaultConfig()
	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) UnmarshalYAML(value *yaml.Node) error {
	type rawConfig Config
	var raw rawConfig
	custom, err := ParseKnownAndCustomAuto(value, &raw)
	if err != nil {
		return err
	}
	*c = Config(raw)
	c.Custom = custom
	return nil
}

func (r *CacheConfig) UnmarshalYAML(value *yaml.Node) error {
	type rawConfig CacheConfig
	var raw rawConfig
	custom, err := ParseKnownAndCustomAuto(value, &raw)
	if err != nil {
		return err
	}
	*r = CacheConfig(raw)
	r.Custom = custom
	return nil
}

func (r *CacheStorageItemConfig) UnmarshalYAML(value *yaml.Node) error {
	type rawConfig CacheStorageItemConfig
	var raw rawConfig
	custom, err := ParseKnownAndCustomAuto(value, &raw)
	if err != nil {
		return err
	}
	*r = CacheStorageItemConfig(raw)
	r.Custom = custom
	return nil
}

func (r *RouterConfig) UnmarshalYAML(value *yaml.Node) error {
	type rawConfig RouterConfig
	var raw rawConfig
	custom, err := ParseKnownAndCustomAuto(value, &raw)
	if err != nil {
		return err
	}
	*r = RouterConfig(raw)
	r.Custom = custom
	return nil
}

func (c *RendererConfig) UnmarshalYAML(value *yaml.Node) error {
	type rawConfig RendererConfig
	var raw rawConfig
	custom, err := ParseKnownAndCustomAuto(value, &raw)
	if err != nil {
		return err
	}
	*c = RendererConfig(raw)
	c.Custom = custom
	return nil
}

func ParseKnownAndCustom(node *yaml.Node, out any, knownFields []string) (map[string]any, error) {
	if node.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("expected mapping node, got: %d", node.Kind)
	}

	raw := make(map[string]*yaml.Node)
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i].Value
		val := node.Content[i+1]
		raw[key] = val
	}

	t := reflect.TypeOf(out)
	if t.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("out must be pointer to struct")
	}

	alias := reflect.New(t.Elem()).Interface()

	tmp := &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: []*yaml.Node{},
	}
	for _, k := range knownFields {
		if v, ok := raw[k]; ok {
			tmp.Content = append(tmp.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: k})
			tmp.Content = append(tmp.Content, v)
		}
	}

	if err := tmp.Decode(alias); err != nil {
		return nil, fmt.Errorf("decode known fields: %w", err)
	}

	reflect.ValueOf(out).Elem().Set(reflect.ValueOf(alias).Elem())

	custom := make(map[string]any)
	for k, v := range raw {
		found := false
		for _, known := range knownFields {
			if k == known {
				found = true
				break
			}
		}
		if !found {
			custom[k] = ValueParse(k, v)
		}
	}

	return custom, nil
}

func ValueParse(k any, node *yaml.Node) any {
	switch node.Kind {
	case yaml.ScalarNode:
		return ParseScalarValue(node.Value)
	case yaml.SequenceNode:
		var values []any
		for kk, item := range node.Content {
			if item.Kind == yaml.ScalarNode {
				values = append(values, item.Value)
			} else if node.Kind == yaml.SequenceNode {
				values = append(values, ValueParse(kk, item))
			}
		}
		return values
	case yaml.MappingNode:
		m := make(map[string]any)
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valNode := node.Content[i+1]
			key := keyNode.Value
			m[key] = ValueParse(key, valNode)
		}
		return m
	default:
		return node
	}
}

func ParseScalarValue(s string) any {
	s = strings.TrimSpace(s)

	if b, err := strconv.ParseBool(s); err == nil {
		return b
	}

	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}

	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}

	return s
}

func ParseKnownAndCustomAuto(node *yaml.Node, out any) (map[string]any, error) {
	if node.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("expected mapping node")
	}

	t := reflect.TypeOf(out)
	if t.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("output must be pointer to struct")
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("output must point to struct")
	}

	var knownFields []string
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		yamlTag := f.Tag.Get("yaml")
		if yamlTag == "" || yamlTag == "-" {
			continue
		}

		yamlKey := yamlTag
		if idx := len(yamlTag); idx > 0 {
			if comma := IndexComma(yamlTag); comma > -1 {
				yamlKey = yamlTag[:comma]
			}
		}
		knownFields = append(knownFields, yamlKey)
	}

	return ParseKnownAndCustom(node, out, knownFields)
}

func IndexComma(tag string) int {
	for i := 0; i < len(tag); i++ {
		if tag[i] == ',' {
			return i
		}
	}
	return -1
}
