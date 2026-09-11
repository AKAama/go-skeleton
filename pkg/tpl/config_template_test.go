package tpl

import (
	"strings"
	"testing"

	"github.com/AKAama/go-skeleton/config"
	"github.com/AKAama/go-skeleton/pkg/util"
)

func TestConfigTemplatesUseMapstructureTags(t *testing.T) {
	cfg := config.NewProjectConfig(
		config.WithProjectName("demo"),
		config.WithModulePath("example.com/demo"),
	)

	global, err := util.TemplateParseFS(ConfigGo, cfg, "config.go.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(global.String(), `mapstructure:"db"`) {
		t.Fatal("global config must map the db key explicitly")
	}

	database, err := util.TemplateParseFS(DBConfigGo, cfg, "db.config.go.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	for _, tag := range []string{"host", "port", "username", "password", "database", "maxConnections"} {
		if !strings.Contains(database.String(), `mapstructure:"`+tag+`"`) {
			t.Fatalf("database config is missing mapstructure tag %q", tag)
		}
	}
}
