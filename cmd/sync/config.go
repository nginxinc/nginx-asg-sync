package main

import (
	"errors"
	"fmt"
	"time"

	yaml "gopkg.in/yaml.v2"
)

// commonConfig stores the configuration parameters common to all providers
type commonConfig struct {
	APIEndpoint           string        `yaml:"api_endpoint"`
	CloudProvider         string        `yaml:"cloud_provider"`
	SyncIntervalInSeconds time.Duration `yaml:"sync_interval_in_seconds"`
}

func parseCommonConfig(data []byte) (*commonConfig, error) {
	cfg := &commonConfig{}
	err := yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	err = validateCommonConfig(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func validateCommonConfig(cfg *commonConfig) error {
	if cfg.APIEndpoint == "" {
		return fmt.Errorf(errorMsgFormat, "api_endpoint")
	}

	if cfg.SyncIntervalInSeconds == 0 {
		return errors.New(intervalErrorMsg)
	}

	if cfg.CloudProvider == "" {
		cfg.CloudProvider = defaultCloudProvider
	}

	if !validateCloudProvider(cfg.CloudProvider) {
		return fmt.Errorf(cloudProviderErrorMsg, cfg.CloudProvider)
	}

	return nil
}

// Upstream is the cloud agnostic representation of an Upstream (eg, common fields for every cloud provider)
type Upstream struct {
	MaxConns     *int
	MaxFails     *int
	Name         string
	ScalingGroup string
	Kind         string
	FailTimeout  string
	SlowStart    string
	Port         int
	InService    bool
}
