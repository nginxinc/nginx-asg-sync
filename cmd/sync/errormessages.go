package main

const (
	errorMsgFormat                 = "The mandatory field %v is either empty or missing in the config file"
	intervalErrorMsg               = "The mandatory field sync_interval_in_seconds is either 0 or missing in the config file"
	cloudProviderErrorMsg          = "The field cloud_provider has invalid value %v in the config file"
	defaultCloudProvider           = "AWS"
	upstreamNameErrorMsg           = "The mandatory field name is either empty or missing for an upstream in the config file"
	upstreamErrorMsgFormat         = "The mandatory field %v is either empty or missing for the upstream %v in the config file"
	upstreamPortErrorMsgFormat     = "The mandatory field port is either zero or missing for the upstream %v in the config file"
	upstreamKindErrorMsgFormat     = "The mandatory field kind is either not equal to http or tcp or missing for the upstream %v in the config file"
	upstreamMaxConnsErrorMsgFmt    = "The field max_conns has invalid value %v in the config file"
	upstreamMaxFailsErrorMsgFmt    = "The field max_fails has invalid value %v in the config file"
	upstreamFailTimeoutErrorMsgFmt = "The field fail_timeout has invalid value %v in the config file"
	upstreamSlowStartErrorMsgFmt   = "The field slow_start has invalid value %v in the config file"
)
