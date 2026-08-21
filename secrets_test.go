package secrets

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	config "github.com/tommzn/go-config"
)

type SecretsManagerTestSuite struct {
	suite.Suite
}

func TestSecretsManagerTestSuite(t *testing.T) {
	suite.Run(t, new(SecretsManagerTestSuite))
}

func (suite *SecretsManagerTestSuite) TestNewManagerFromConfig() {

	config1 := suite.loadConfigForTest("fixtures/config/env_secrets.yml")
	manager1 := NewSecretsManagerByConfig(config1)
	suite.IsType(&EnvironmentSecretsManager{}, manager1)

	config2 := suite.loadConfigForTest("fixtures/config/docker_secrets.yml")
	manager2 := NewSecretsManagerByConfig(config2)
	suite.IsType(&DockerSecretsManager{}, manager2)

	config3 := suite.loadConfigForTest("fixtures/config/empty_secrets.yml")
	manager3 := NewSecretsManagerByConfig(config3)
	suite.IsType(&EnvironmentSecretsManager{}, manager3)
}

func (suite *SecretsManagerTestSuite) TestStaticSecretsManager() {

	secrets := make(map[string]string)
	testSecret := "TestSecret"
	secrets["TESTKEY"] = testSecret
	secretsmanager := NewStaticSecretsManager(secrets)

	secret1, err1 := secretsmanager.Obtain("TESTKEY")
	suite.Nil(err1)
	suite.NotNil(secret1)
	suite.Equal(testSecret, *secret1)

	secret2, err2 := secretsmanager.Obtain("TestKey")
	suite.Nil(err2)
	suite.NotNil(secret2)
	suite.Equal(testSecret, *secret2)

	secret3, err3 := secretsmanager.Obtain("XXX")
	suite.NotNil(err3)
	suite.Nil(secret3)
}

func (suite *SecretsManagerTestSuite) TestEnvironmentSecretsManager() {

	expectedSecret := "xxx123"
	key := "Test_Secret"
	suite.T().Setenv(key, expectedSecret)
	secretsmanager := NewSecretsManager()

	secret, err := secretsmanager.Obtain(key)
	suite.Nil(err)
	suite.NotNil(secret)
	suite.Equal(expectedSecret, *secret)

	notExistingSecret, err := secretsmanager.Obtain("xyz")
	suite.NotNil(err)
	suite.Nil(notExistingSecret)
}

func (suite *SecretsManagerTestSuite) TestDockerSecretsManager() {

	secretsmanager := NewDockerSecretsManager("./fixtures")
	secret, err := secretsmanager.Obtain("TestSecret")
	suite.Nil(err)
	suite.NotNil(secret)
	suite.Equal("xxxYYYzzz", *secret)

	notExistingSecret, err := secretsmanager.Obtain("xyz")
	suite.NotNil(err)
	suite.Nil(notExistingSecret)
}

func (suite *SecretsManagerTestSuite) TestAddToEnvironment() {

	secrets := make(map[string]string)
	secrets["TESTKEY"] = "TestSecret"
	secretsmanager := NewStaticSecretsManager(secrets)

	ExportToEnvironment([]string{"TESTKEY"}, secretsmanager)
	envValue, ok := os.LookupEnv("TESTKEY")
	suite.True(ok)
	suite.Equal(secrets["TESTKEY"], envValue)
}

func (suite *SecretsManagerTestSuite) TestStaticSecretsManagerClear() {

	secrets := map[string]string{"KEY1": "val1", "KEY2": "val2"}
	manager := NewStaticSecretsManager(secrets).(*StaticSecretsManager)

	_, err := manager.Obtain("KEY1")
	suite.Nil(err)

	manager.Clear()

	_, err = manager.Obtain("KEY1")
	suite.NotNil(err)
	_, err = manager.Obtain("KEY2")
	suite.NotNil(err)
}

func (suite *SecretsManagerTestSuite) TestExportToEnvironmentMissingKey() {

	manager := NewStaticSecretsManager(map[string]string{})
	failed := ExportToEnvironment([]string{"MISSING_KEY_XYZ_99"}, manager)
	suite.Equal([]string{"MISSING_KEY_XYZ_99"}, failed)
	_, ok := os.LookupEnv("MISSING_KEY_XYZ_99")
	suite.False(ok)
}

func (suite *SecretsManagerTestSuite) TestDockerSecretsManagerPathTraversal() {

	manager := NewDockerSecretsManager("./fixtures")

	for _, key := range []string{"../secrets_test", ".", "..", "path/traversal", "back\\slash"} {
		secret, err := manager.Obtain(key)
		suite.NotNil(err, "expected error for key %q", key)
		suite.Nil(secret, "expected nil secret for key %q", key)
	}
}

func (suite *SecretsManagerTestSuite) TestByteSliceAsStringPtr() {

	suite.Nil(byteSliceAsStringPtr([]byte{}))

	result := byteSliceAsStringPtr([]byte("hello"))
	suite.NotNil(result)
	suite.Equal("hello", *result)
}

func (suite *SecretsManagerTestSuite) TestGenerateSecretKeys() {

	// Mixed case produces original + lowercase + uppercase
	keys := generateSecretKeys("MyKey")
	suite.ElementsMatch([]string{"MyKey", "mykey", "MYKEY"}, keys)

	// Already lowercase produces original + uppercase only
	keys = generateSecretKeys("mykey")
	suite.ElementsMatch([]string{"mykey", "MYKEY"}, keys)

	// Already uppercase produces original + lowercase only
	keys = generateSecretKeys("MYKEY")
	suite.ElementsMatch([]string{"MYKEY", "mykey"}, keys)
}

func (suite *SecretsManagerTestSuite) TestIsValidSecretFileName() {

	for _, name := range []string{"validkey", "VALID_KEY", "valid-key", "secret.txt"} {
		suite.True(isValidSecretFileName(name), "expected valid: %q", name)
	}

	for _, name := range []string{"../secret", "path/secret", `path\secret`, ".", ".."} {
		suite.False(isValidSecretFileName(name), "expected invalid: %q", name)
	}
}

func (suite *SecretsManagerTestSuite) TestSecretNotFoundError() {

	err := asSecretNotFoundError("somekey")
	suite.Equal(`secret not found: "somekey"`, err.Error())
	suite.IsType(&SecretNotFoundError{}, err)
}

func (suite *SecretsManagerTestSuite) loadConfigForTest(configFile string) config.Config {
	configLoader := config.NewFileConfigSource(&configFile)
	conf, err := configLoader.Load()
	suite.Nil(err)
	return conf
}
