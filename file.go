package secrets

import (
	"bufio"
	"encoding/base64"
	"os"
	"strings"
)

// DEFAULT_SECRETS_FILE defines default path to a credentials file.
const DEFAULT_SECRETS_FILE = "~/.credentials"

// FileSecretsManager reads secrets from defined file. Secrets have to be
// added as key:value pair in this credentials file.
type FileSecretsManager struct {
	secretsFile string
}

// Obtain will try to read secret from defined credentials file.
// Expects secrets as a key:value pair, separatir is ":", where secrets value
// is base64 encoded.
func (s *FileSecretsManager) Obtain(key string) (*string, error) {

	file, err := os.Open(s.secretsFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		secretsKey, secretsValue := splitCredentials(line)
		if assertKeyIsEqual(key, secretsKey) {
			if decoded, err := base64.StdEncoding.DecodeString(secretsValue); err == nil {
				decodedSecretsValue := string(decoded)
				return &decodedSecretsValue, nil
			}
		}
	}
	return nil, newSecretNotFoundError(key)
}

func splitCredentials(line string) (string, string) {
	if splitted := strings.Split(line, ":"); len(splitted) >= 2 {
		return splitted[0], splitted[1]
	}
	return "", ""
}

func assertKeyIsEqual(key, credentialsKey string) bool {
	return key == credentialsKey ||
		strings.ToUpper(key) == strings.ToUpper(credentialsKey)
}
