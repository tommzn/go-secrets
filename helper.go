package secrets

import (
	"os"
	"strings"
)

// ExportToEnvironment will export secrets identified by given keys to environment variables.
// It returns a slice of keys for which the lookup or the environment assignment failed.
func ExportToEnvironment(keys []string, manager SecretsManager) []string {
	var failed []string
	for _, key := range keys {
		val, err := manager.Obtain(key)
		if err != nil {
			failed = append(failed, key)
			continue
		}
		if err := os.Setenv(key, *val); err != nil {
			failed = append(failed, key)
		}
	}
	return failed
}

// generateSecretKeys will create a slice of keys. This includes the passed key and a lower and upper case version of it.
func generateSecretKeys(key string) []string {

	keys := []string{key}
	lowerKey := strings.ToLower(key)
	if key != lowerKey {
		keys = append(keys, lowerKey)
	}
	upperKey := strings.ToUpper(key)
	if key != upperKey {
		keys = append(keys, upperKey)
	}
	return keys
}

// byteSliceAsStringPtr returns passed byte slice as string pointer.
// If byte slice is empty nil will be returned.
func byteSliceAsStringPtr(byteSlice []byte) *string {
	if len(byteSlice) == 0 {
		return nil
	}
	stringValue := string(byteSlice)
	return &stringValue
}
