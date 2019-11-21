package main

const defaultFailTimeout = "10s"
const defaultSlowStart = "0s"

func getFailTimeoutOrDefault(failTimeout string) string {
	if failTimeout == "" {
		return defaultFailTimeout
	}

	return failTimeout
}

func getSlowStartOrDefault(slowStart string) string {
	if slowStart == "" {
		return defaultSlowStart
	}

	return slowStart
}