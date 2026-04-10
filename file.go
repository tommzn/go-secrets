package secrets

import (
	"bufio"
	"encoding/base64"
	"fmt"
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
// Expects secrets as a key:value pair, separator is ":", where secrets value
// is base64 encoded.
func (s *FileSecretsManager) Obtain(key string) (*string, error) {

	info, err := os.Lstat(s.secretsFile)
	if err != nil {
		return nil, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("secret file must not be a symlink: %s", s.secretsFile)
	}
	if info.Mode().Perm()&0077 != 0 {
		return nil, fmt.Errorf("secret file has insecure permissions %04o, expected 0600 or stricter", info.Mode().Perm())
	}

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
				return byteSliceAsStringPtr(decoded), nil
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return nil, asSecretNotFoundError(key)
}

// splitCredentials splits a credentials file line into key and value on the
// first ":" separator. Lines without a separator return two empty strings.
func splitCredentials(line string) (string, string) {
	if splitted := strings.Split(line, ":"); len(splitted) >= 2 {
		return splitted[0], splitted[1]
	}
	return "", ""
}

// assertKeyIsEqual reports whether key and credentialsKey are identical.
func assertKeyIsEqual(key, credentialsKey string) bool {
	return key == credentialsKey
}
