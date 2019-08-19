package main

// CloudProvider is the interface to connect with any cloud provider.
type CloudProvider interface {
	GetPrivateIPsForScalingGroup(name string) ([]string, error)
	CheckIfScalingGroupExists(name string) (bool, error)
	GetUpstreams() []Upstream
}

func validateCloudProvider(provider string) bool {
	providers := map[string]bool{
		"AWS":   true,
		"Azure": true,
	}

	return providers[provider]
}
