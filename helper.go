package secrets

import (
	"log"
	"os"
	"strings"
)

// ExportToEnvironment will export secrets identified by given keys to environment variables.
func ExportToEnvironment(keys []string, manager SecretsManager) {
	for _, key := range keys {
		if val, err := manager.Obtain(key); err == nil {
			os.Setenv(key, *val)
		} else {
			log.Println(err)
		}
	}
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

// byteSliceAsStringPtr returns passed byte slice as string point.
// If byte slice is empty nil will be retunred.
func byteSliceAsStringPtr(byteSlice []byte) *string {
	if len(byteSlice) == 0 {
		return nil
	}
	stringValue := string(byteSlice)
	return &stringValue
}
