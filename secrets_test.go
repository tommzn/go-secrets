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
	os.Setenv(key, expectedSecret)
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

	secretsmanager := NewDockerecretsManager("./fixtures")
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

func (suite *SecretsManagerTestSuite) loadConfigForTest(configFile string) config.Config {
	configLoader := config.NewFileConfigSource(&configFile)
	conf, err := configLoader.Load()
	suite.Nil(err)
	return conf
}
