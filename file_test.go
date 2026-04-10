package secrets

import (
	"os"
	"strings"
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

	// value for AWS_SECRET_ACCESS_KEY contains a space suffix which causes base64 decode to fail;
	// the entry is skipped and secret not found is returned.
	secret3, err3 := secretsmanager.Obtain("AWS_SECRET_ACCESS_KEY")
	suite.NotNil(err3)
	suite.IsType(&SecretNotFoundError{}, err3)
	suite.Nil(secret3)
}

func (suite *FileSecretsManagerTestSuite) TestWithMissingFile() {

	secretsmanager := NewFileSecretsManager("xxx")
	secret, err := secretsmanager.Obtain("yxz")
	suite.NotNil(err)
	suite.Nil(secret)
}

func (suite *FileSecretsManagerTestSuite) TestSymlinkRejection() {

	// Create a valid credentials file
	tmpFile, err := os.CreateTemp("", "credentials-*")
	suite.Require().NoError(err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()
	os.Chmod(tmpFile.Name(), 0600)

	// Point a symlink at it
	symlinkPath := tmpFile.Name() + ".link"
	suite.Require().NoError(os.Symlink(tmpFile.Name(), symlinkPath))
	defer os.Remove(symlinkPath)

	secretsmanager := NewFileSecretsManager(symlinkPath)
	secret, err := secretsmanager.Obtain("anykey")
	suite.NotNil(err)
	suite.Nil(secret)
	suite.Contains(err.Error(), "symlink")
}

func (suite *FileSecretsManagerTestSuite) TestInsecurePermissions() {

	tmpFile, err := os.CreateTemp("", "credentials-*")
	suite.Require().NoError(err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()
	os.Chmod(tmpFile.Name(), 0644)

	secretsmanager := NewFileSecretsManager(tmpFile.Name())
	secret, err := secretsmanager.Obtain("anykey")
	suite.NotNil(err)
	suite.Nil(secret)
	suite.Contains(err.Error(), "insecure permissions")
}

func (suite *FileSecretsManagerTestSuite) TestExpandHome() {

	// Path without tilde is returned unchanged
	suite.Equal("/absolute/path", expandHome("/absolute/path"))
	suite.Equal("relative/path", expandHome("relative/path"))

	// Path with ~/ prefix is expanded to home directory
	expanded := expandHome("~/somefile")
	suite.False(strings.HasPrefix(expanded, "~"), "tilde should be expanded")
	suite.True(strings.HasSuffix(expanded, "somefile"))
}
