package main

import "testing"

type testInputAWS struct {
	cfg *awsConfig
	msg string
}

func getValidAWSConfig() *awsConfig {
	upstreams := []awsUpstream{
		{
			Name:             "backend1",
			AutoscalingGroup: "backend-group",
			Port:             80,
			Kind:             "http",
			MaxConns:         1000,
			MaxFails:         10,
			FailTimeout:      "5s",
			SlowStart:        "30s",
		},
	}
	cfg := awsConfig{
		Region:    "us-west-2",
		Upstreams: upstreams,
	}

	return &cfg
}

func getInvalidAWSConfigInput() []*testInputAWS {
	var input []*testInputAWS

	invalidRegionCfg := getValidAWSConfig()
	invalidRegionCfg.Region = ""
	input = append(input, &testInputAWS{invalidRegionCfg, "invalid region"})

	invalidMissingUpstreamsCfg := getValidAWSConfig()
	invalidMissingUpstreamsCfg.Upstreams = nil
	input = append(input, &testInputAWS{invalidMissingUpstreamsCfg, "no upstreams"})

	invalidUpstreamNameCfg := getValidAWSConfig()
	invalidUpstreamNameCfg.Upstreams[0].Name = ""
	input = append(input, &testInputAWS{invalidUpstreamNameCfg, "invalid name of the upstream"})

	invalidUpstreamAutoscalingGroupCfg := getValidAWSConfig()
	invalidUpstreamAutoscalingGroupCfg.Upstreams[0].AutoscalingGroup = ""
	input = append(input, &testInputAWS{invalidUpstreamAutoscalingGroupCfg, "invalid autoscaling_group of the upstream"})

	invalidUpstreamPortCfg := getValidAWSConfig()
	invalidUpstreamPortCfg.Upstreams[0].Port = 0
	input = append(input, &testInputAWS{invalidUpstreamPortCfg, "invalid port of the upstream"})

	invalidUpstreamKindCfg := getValidAWSConfig()
	invalidUpstreamKindCfg.Upstreams[0].Kind = ""
	input = append(input, &testInputAWS{invalidUpstreamKindCfg, "invalid kind of the upstream"})

	invalidUpstreamMaxConnsCfg := getValidAWSConfig()
	invalidUpstreamMaxConnsCfg.Upstreams[0].MaxConns = -10
	input = append(input, &testInputAWS{invalidUpstreamMaxConnsCfg, "invalid max_conns of the upstream"})

	invalidUpstreamMaxFailsCfg := getValidAWSConfig()
	invalidUpstreamMaxFailsCfg.Upstreams[0].MaxFails = -10
	input = append(input, &testInputAWS{invalidUpstreamMaxFailsCfg, "invalid max_fails of the upstream"})

	invalidUpstreamFailTimeoutCfg := getValidAWSConfig()
	invalidUpstreamFailTimeoutCfg.Upstreams[0].FailTimeout = "-10s"
	input = append(input, &testInputAWS{invalidUpstreamFailTimeoutCfg, "invalid fail_timeout of the upstream"})

	invalidUpstreamSlowStartCfg := getValidAWSConfig()
	invalidUpstreamSlowStartCfg.Upstreams[0].SlowStart = "-10s"
	input = append(input, &testInputAWS{invalidUpstreamSlowStartCfg, "invalid slow_start of the upstream"})

	return input
}

func TestValidateAWSConfigNotValid(t *testing.T) {
	input := getInvalidAWSConfigInput()

	for _, item := range input {
		err := validateAWSConfig(item.cfg)
		if err == nil {
			t.Errorf("validateAWSConfig() didn't fail for the invalid config file with %v", item.msg)
		}
	}
}

func TestValidateAWSConfigValid(t *testing.T) {
	cfg := getValidAWSConfig()

	err := validateAWSConfig(cfg)
	if err != nil {
		t.Errorf("validateAWSConfig() failed for the valid config: %v", err)
	}
}
