package main

import "testing"

var validYaml = []byte(`region: us-west-2
api_endpoint: http://127.0.0.1:8080/api
sync_interval_in_seconds: 5
upstreams:
  - name: backend1
    autoscaling_group: backend-group
    port: 80
    kind: http
  - name: backend2
    autoscaling_group: backend-group
    port: 80
    kind: http
`)

type testInput struct {
	cfg *config
	msg string
}

func getValidConfig() *config {
	upstreams := []upstream{
		upstream{
			Name:             "backend1",
			AutoscalingGroup: "backend-group",
			Port:             80,
			Kind:             "http",
		},
	}
	cfg := config{
		Region:                "us-west-2",
		APIEndpoint:           "http://127.0.0.1:8080/api",
		SyncIntervalInSeconds: 1,
		Upstreams:             upstreams,
	}

	return &cfg
}

func getInvalidConfigInput() []*testInput {
	var input []*testInput

	invalidRegionCfg := getValidConfig()
	invalidRegionCfg.Region = ""
	input = append(input, &testInput{invalidRegionCfg, "invalid region"})

	invalidAPIEndpointCfg := getValidConfig()
	invalidAPIEndpointCfg.APIEndpoint = ""
	input = append(input, &testInput{invalidAPIEndpointCfg, "invalid api_endpoint"})

	invalidSyncIntervalInSecondsCfg := getValidConfig()
	invalidSyncIntervalInSecondsCfg.SyncIntervalInSeconds = 0
	input = append(input, &testInput{invalidSyncIntervalInSecondsCfg, "invalid sync_interval_in_seconds"})

	invalidMissingUpstreamsCfg := getValidConfig()
	invalidMissingUpstreamsCfg.Upstreams = nil
	input = append(input, &testInput{invalidMissingUpstreamsCfg, "no upstreams"})

	invalidUpstreamNameCfg := getValidConfig()
	invalidUpstreamNameCfg.Upstreams[0].Name = ""
	input = append(input, &testInput{invalidUpstreamNameCfg, "invalid name of the upstream"})

	invalidUpstreamAutoscalingGroupCfg := getValidConfig()
	invalidUpstreamAutoscalingGroupCfg.Upstreams[0].AutoscalingGroup = ""
	input = append(input, &testInput{invalidUpstreamAutoscalingGroupCfg, "invalid autoscaling_group of the upstream"})

	invalidUpstreamPortCfg := getValidConfig()
	invalidUpstreamPortCfg.Upstreams[0].Port = 0
	input = append(input, &testInput{invalidUpstreamPortCfg, "invalid port of the upstream"})

	invalidUpstreamKindCfg := getValidConfig()
	invalidUpstreamKindCfg.Upstreams[0].Kind = ""
	input = append(input, &testInput{invalidUpstreamKindCfg, "invalid kind of the upstream"})

	return input
}

func TestUnmarshalConfig(t *testing.T) {
	_, err := unmarshalConfig(validYaml)
	if err != nil {
		t.Errorf("unmarshalConfig() failed for the valid yaml: %v", err)
	}
}

func TestValidateConfigNotValid(t *testing.T) {
	input := getInvalidConfigInput()

	for _, item := range input {
		err := validateConfig(item.cfg)
		if err == nil {
			t.Errorf("validateConfig() didn't fail for the invalid config file with %v", item.msg)
		}
	}
}

func TestValidateConfigValid(t *testing.T) {
	cfg := getValidConfig()

	err := validateConfig(cfg)
	if err != nil {
		t.Errorf("validateConfig() failed for the valid config: %v", err)
	}
}
