package analytics

import (
	"github.com/bitrise-io/go-utils/log"
)

const stepName = "bitrise-init"

func initData(data map[string]interface{}) map[string]interface{} {
	if data == nil {
		data = map[string]interface{}{}
	}
	data["source"] = "scanner"
	return data
}

// LogError sends analytics log using log.RErrorf by setting the stepID.
// Used for errors, returned to the consumer.
func LogError(tag string, data map[string]interface{}, format string, v ...interface{}) {
	log.RErrorf(stepName, tag, initData(data), format, v...)
}

// LogWarn sends analytics log using log.RInfof by setting the stepID.
// Used for warnings, returned to the consumer.
func LogWarn(tag string, data map[string]interface{}, format string, v ...interface{}) {
	log.RWarnf(stepName, tag, initData(data), format, v...)
}

// LogInfo sends analytics log using log.RInfof by setting the stepID.
// Used for internal errors (not returned to the consumer).
func LogInfo(tag string, data map[string]interface{}, format string, v ...interface{}) {
	log.RInfof(stepName, tag, initData(data), format, v...)
}

// DetectorErrorData creates analytics data that includes the platform and error
func DetectorErrorData(detector string, err error) map[string]interface{} {
	return map[string]interface{}{
		"detector": detector,
		"error":    err.Error(),
	}
}
