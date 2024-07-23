package config_test

import (
	"testing"

	"github.com/aapclark/go-indexer/m/v2/internal/config"
)

func TestValidLoadConfig(t *testing.T) {
	p := "../../config/config_example.yaml"
	expected := config.AppConfig{
		Rpcs: []config.RpcConfig{
			{ChainId: 1, Url: "http://example.com", Blocktime: 10},
			{ChainId: 2, Url: "http://example2.com", Blocktime: 20},
		},
		Logging: config.LogConfig{
			Level:    1,
			Format:   "json",
			FilePath: "/var/log/app.log",
		},
	}

	result, err := config.LoadConfig(p)
	if err != nil {
		t.Fatalf("failed with error: %v", err)
	}
	for i, v := range result.Rpcs {
		if v != expected.Rpcs[i] {
			t.Fatalf("rpc fields do not match, got %v, expected %v", v, expected.Rpcs[i])
		}
	}
	if result.Logging != expected.Logging {
		t.Fatalf("logging fields do not match, got %v, expected %v", result.Logging, expected.Logging)
	}

}

func TestInvalidLoadConfig(t *testing.T) {
	p := "../../config/not.yaml"

	result, err := config.LoadConfig(p)
	if err == nil {
		t.Fatal("expected error but did not recieve one")
	}
	if result != nil {
		t.Fatal("expected result to be nil")
	}

}
