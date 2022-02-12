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

	// value for AWS_SECRET_ACCESS_KEY cintains a space suffix which leads to base64 decode error
	secret3, err3 := secretsmanager.Obtain("AWS_SECRET_ACCESS_KEY")
	suite.NotNil(err3)
	suite.IsType(&Base64DecodeError{}, err3)
	suite.Nil(secret3)
}

func (suite *FileSecretsManagerTestSuite) TestWithMissingFile() {

	secretsmanager := NewFileSecretsManager("xxx")
	secret, err := secretsmanager.Obtain("yxz")
	suite.NotNil(err)
	suite.Nil(secret)
}
