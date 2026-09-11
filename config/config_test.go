package config

import "testing"

func TestNewProjectConfigDoesNotMutateDefaultModules(t *testing.T) {
	withoutGin := NewProjectConfig(WithOutModules("gin"))
	if _, ok := withoutGin.Modules["gin"]; ok {
		t.Fatal("gin should be removed from this config")
	}

	fresh := NewProjectConfig()
	if _, ok := fresh.Modules["gin"]; !ok {
		t.Fatal("creating one config must not mutate DefaultModules")
	}
}
