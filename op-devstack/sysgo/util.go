package sysgo

import (
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
)

// getEnvVarOrDefault returns the value of the provided env var or the provided default value if unset.
func getEnvVarOrDefault(envVarName string, defaultValue string) string {
	val, found := os.LookupEnv(envVarName)
	if !found {
		val = defaultValue
	}
	return val
}

// propagateEnvVarOrDefault returns a string in the format "ENV_VAR_NAME=VALUE", with the ENV_VAR_NAME being
// the provided env var name and the value being the value of that env var, or the provided default
// value if that env var is unset.
func propagateEnvVarOrDefault(envVarName string, defaultValue string) string {
	if val := getEnvVarOrDefault(envVarName, defaultValue); val == "" {
		return ""
	} else {
		return fmt.Sprintf("%s=%s", envVarName, val)
	}
}

// setEnvFromEnvOrDefault sets the provided key in the provided env map to the value of the
// corresponding environment variable, or to the provided default value if that environment variable is unset.
func setEnvFromEnvOrDefault(env map[string]string, key, def string) {
	if v := os.Getenv(key); v != "" {
		env[key] = v
	} else if def != "" {
		env[key] = def
	}
}

// setEnvIfNotNil sets the provided key in the provided envVars map to the string representation
// of the provided val if val is not nil.
func setEnvIfNotNil[T any](envVars map[string]string, key string, val *T) {
	if val == nil {
		return
	}
	switch v := any(*val).(type) {
	case string:
		envVars[key] = v
	case bool:
		envVars[key] = fmt.Sprintf("%t", v)
	case float64:
		envVars[key] = fmt.Sprintf("%v", v)
	default:
		envVars[key] = fmt.Sprintf("%d", v)
	}
}

// NB: arbitrary start port with a low probability of conflict
var availableLocalPortStart = 20_000
var availableLocalPortMutex sync.Mutex

// getAvailableLocalPort searches for and returns a currently unused local port.
// Note: this function is threadsafe.
func getAvailableLocalPort() (string, error) {
	availableLocalPortMutex.Lock()
	defer availableLocalPortMutex.Unlock()

	for port := availableLocalPortStart; port < 65_535; port++ {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			continue
		}
		_ = ln.Close()
		availableLocalPortStart = port + 1
		return fmt.Sprintf("%d", port), nil
	}

	return "", errors.New("could not find open port")
}
