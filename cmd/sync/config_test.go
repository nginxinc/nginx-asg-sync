package main

import "testing"

var validYaml = []byte(`
cloud_provider: AWS
api_endpoint: http://127.0.0.1:8080/api
sync_interval_in_seconds: 5
`)

type testInputCommon struct {
	cfg *commonConfig
	msg string
}

func getValidCommonConfig() *commonConfig {
	return &commonConfig{
		APIEndpoint:           "http://127.0.0.1:8080/api",
		SyncIntervalInSeconds: 1,
	}
}

func getInvalidCommonConfigInput() []*testInputCommon {
	var input []*testInputCommon

	invalidAPIEndpointCfg := getValidCommonConfig()
	invalidAPIEndpointCfg.APIEndpoint = ""
	input = append(input, &testInputCommon{invalidAPIEndpointCfg, "invalid api_endpoint"})

	invalidSyncIntervalInSecondsCfg := getValidCommonConfig()
	invalidSyncIntervalInSecondsCfg.SyncIntervalInSeconds = 0
	input = append(input, &testInputCommon{invalidSyncIntervalInSecondsCfg, "invalid sync_interval_in_seconds"})

	return input
}

func TestValidateCommonConfigNotValid(t *testing.T) {
	input := getInvalidCommonConfigInput()

	for _, item := range input {
		err := validateCommonConfig(item.cfg)
		if err == nil {
			t.Errorf("validateCommonConfig() didn't fail for the invalid config file with %v", item.msg)
		}
	}
}

func TestValidateCommonConfigValid(t *testing.T) {
	cfg := getValidCommonConfig()

	err := validateCommonConfig(cfg)
	if err != nil {
		t.Errorf("validateCommonConfig() failed for the valid config: %v", err)
	}
}

func TestParseCommonConfig(t *testing.T) {
	_, err := parseCommonConfig(validYaml)
	if err != nil {
		t.Errorf("parseCommonConfig() failed for the valid config yaml: %v", string(validYaml))
	}
}
