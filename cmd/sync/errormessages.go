package main

const (
	errorMsgFormat                 = "the mandatory field %v is either empty or missing in the config file"
	intervalErrorMsg               = "the mandatory field sync_interval_in_seconds is either 0 or missing in the config file"
	cloudProviderErrorMsg          = "the field cloud_provider has invalid value %v in the config file"
	defaultCloudProvider           = "AWS"
	upstreamNameErrorMsg           = "the mandatory field name is either empty or missing for an upstream in the config file"
	upstreamErrorMsgFormat         = "the mandatory field %v is either empty or missing for the upstream %v in the config file"
	upstreamPortErrorMsgFormat     = "the mandatory field port is either zero or missing for the upstream %v in the config file"
	upstreamKindErrorMsgFormat     = "the mandatory field kind is either not equal to http or tcp or missing for the upstream %v in the config file"
	upstreamMaxConnsErrorMsgFmt    = "the field max_conns has invalid value %v in the config file"
	upstreamMaxFailsErrorMsgFmt    = "the field max_fails has invalid value %v in the config file"
	upstreamFailTimeoutErrorMsgFmt = "the field fail_timeout has invalid value %v in the config file"
	upstreamSlowStartErrorMsgFmt   = "the field slow_start has invalid value %v in the config file"
)
