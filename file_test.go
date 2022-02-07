package secrets

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type FileSecretsManagerTestSuite struct {
	suite.Suite
}

func TestFileSecretsManagerTestSuite(t *testing.T) {
	suite.Run(t, new(FileSecretsManagerTestSuite))
}

func (suite *FileSecretsManagerTestSuite) TestObtainSecrets() {

	secretsmanager := NewFileSecretsManager("fixtures/config/credentials")

	secret1, err1 := secretsmanager.Obtain("secrets_key")
	suite.Nil(err1)
	suite.NotNil(secret1)
	suite.Equal("SecretValue", *secret1)

	secret2, err2 := secretsmanager.Obtain("yxz")
	suite.NotNil(err2)
	suite.Nil(secret2)
}

func (suite *FileSecretsManagerTestSuite) TestWithMissingFile() {

	secretsmanager := NewFileSecretsManager("xxx")
	secret, err := secretsmanager.Obtain("yxz")
	suite.NotNil(err)
	suite.Nil(secret)
}
