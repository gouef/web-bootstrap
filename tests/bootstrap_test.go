package tests

import (
	web_bootstrap "github.com/gouef/web-bootstrap"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestConfig(t *testing.T) {
	cfg, err := web_bootstrap.LoadConfig("./config/config.yml")

	assert.NoError(t, err)
	assert.Equal(t, "./views/templates", cfg.Renderer.Dir)
	assert.Equal(t, "lalala", cfg.Renderer.Custom["test"])
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := web_bootstrap.LoadConfig("nonexistent.yaml")
	assert.Error(t, err)
}

func TestUnmarshalConfig_InvalidYAML(t *testing.T) {
	var c web_bootstrap.Config
	err := yaml.Unmarshal([]byte(`- invalid`), &c)
	assert.Error(t, err)
}

func TestParseKnownAndCustomAuto_InvalidKind(t *testing.T) {
	node := &yaml.Node{Kind: yaml.SequenceNode}
	_, err := web_bootstrap.ParseKnownAndCustomAuto(node, &web_bootstrap.RouterConfig{})
	assert.Error(t, err)
}

func TestParseScalarValue_StringFallback(t *testing.T) {
	assert.Equal(t, "hello", web_bootstrap.ParseScalarValue("hello"))
}
func TestIndexComma(t *testing.T) {
	assert.Equal(t, 3, web_bootstrap.IndexComma("abc,def"))
	assert.Equal(t, -1, web_bootstrap.IndexComma("abcdef"))
}
func TestParseKnownAndCustomAuto_NotPointer(t *testing.T) {
	node := &yaml.Node{Kind: yaml.MappingNode}
	_, err := web_bootstrap.ParseKnownAndCustomAuto(node, web_bootstrap.RouterConfig{}) // není pointer
	assert.Error(t, err)
}

func TestParseKnownAndCustomAuto_NotStruct(t *testing.T) {
	node := &yaml.Node{Kind: yaml.MappingNode}
	x := 123 // není struct
	_, err := web_bootstrap.ParseKnownAndCustomAuto(node, &x)
	assert.Error(t, err)
}
func TestParseKnownAndCustom_NotPointer(t *testing.T) {
	node := &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "dir"},
			{Kind: yaml.ScalarNode, Value: "./views"},
		},
	}
	_, err := web_bootstrap.ParseKnownAndCustom(node, web_bootstrap.RendererConfig{}, []string{"dir"}) // není pointer
	assert.Error(t, err)
}
func TestValueParse_UnknownKind(t *testing.T) {
	node := &yaml.Node{Kind: 999}
	val := web_bootstrap.ValueParse("test", node)
	assert.Equal(t, node, val)
}
func TestParseScalarValue_Types(t *testing.T) {
	assert.Equal(t, true, web_bootstrap.ParseScalarValue("true"))
	assert.Equal(t, int64(123), web_bootstrap.ParseScalarValue("123"))
	assert.Equal(t, 3.14, web_bootstrap.ParseScalarValue("3.14"))
	assert.Equal(t, "text", web_bootstrap.ParseScalarValue("text"))
}
func TestRendererConfig_CustomField(t *testing.T) {
	var cfg web_bootstrap.RendererConfig
	err := yaml.Unmarshal([]byte(`dir: "./views/templates"
unknown: "something"`), &cfg)

	assert.NoError(t, err)
	assert.Equal(t, "./views/templates", cfg.Dir)
	assert.Equal(t, "something", cfg.Custom["unknown"])
}
