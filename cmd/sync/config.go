package main

import (
	"fmt"
	"time"

	yaml "gopkg.in/yaml.v2"
)

// commonConfig stores the configuration parameters common to all providers
type commonConfig struct {
	APIEndpoint           string        `yaml:"api_endpoint"`
	SyncIntervalInSeconds time.Duration `yaml:"sync_interval_in_seconds"`
	CloudProvider         string        `yaml:"cloud_provider"`
}

const errorMsgFormat = "The mandatory field %v is either empty or missing in the config file"
const intervalErrorMsg = "The mandatory field sync_interval_in_seconds is either 0 or missing in the config file"
const cloudProviderErrorMsg = "The field cloud_provider has invalid value %v in the config file"
const defaultCloudProvider = "AWS"

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
		return fmt.Errorf(intervalErrorMsg)
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
	Name         string
	Port         int
	ScalingGroup string
	Kind         string
}
