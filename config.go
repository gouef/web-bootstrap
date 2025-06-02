package web_bootstrap

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type Custom map[string]any

type ConfigInterface interface {
	UnmarshalYAML(value *yaml.Node) error
}

type Config struct {
	Parameters map[string]any `yaml:"parameters"`
	Renderer   RendererConfig `yaml:"renderer"`
	Router     RouterConfig   `yaml:"router"`
	Cache      CacheConfig    `yaml:"cache"`
	Custom     Custom
}

type RendererConfig struct {
	Dir    string `yaml:"dir"`
	Custom Custom
}

type RouterConfig struct {
	Statics []RouterStaticConfig `yaml:"statics"`
	Custom  Custom
}

type RouterStaticConfig struct {
	Path   string `yaml:"path"`
	Root   string `yaml:"root"`
	Custom Custom
}

type CacheConfig struct {
	Storages []CacheStorageItemConfig `yaml:"storages"`
	Custom   Custom
}

type CacheStorageItemConfig struct {
	Type     string `yaml:"type"`
	Instance string `yaml:"instance"`
	Name     string `yaml:"name"`
	Custom   Custom
}

func DefaultConfig() *Config {
	rootDir, _ := filepath.Abs(".")
	cfg := Config{}
	cfg.Parameters = map[string]any{
		"rootDir": rootDir,
	}
	cfg.Renderer = RendererConfig{Dir: "./views/templates"}
	cfg.Router.Statics = []RouterStaticConfig{
		{Path: "/static", Root: "./static"},
		{Path: "/assets", Root: "./static/assets"},
	}
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
