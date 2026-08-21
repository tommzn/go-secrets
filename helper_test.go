package secrets

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HelperTestSuite struct {
	suite.Suite
}

func TestHelperTestSuite(t *testing.T) {
	suite.Run(t, new(HelperTestSuite))
}

// --- generateSecretKeys ---

func (suite *HelperTestSuite) TestGenerateSecretKeys_MixedCase() {
	keys := generateSecretKeys("MyKey")
	suite.Contains(keys, "MyKey")
	suite.Contains(keys, "mykey")
	suite.Contains(keys, "MYKEY")
	suite.Len(keys, 3)
}

func (suite *HelperTestSuite) TestGenerateSecretKeys_AlreadyLower() {
	keys := generateSecretKeys("mykey")
	suite.Contains(keys, "mykey")
	suite.Contains(keys, "MYKEY")
	suite.Len(keys, 2) // no duplicate lowercase added
}

func (suite *HelperTestSuite) TestGenerateSecretKeys_AlreadyUpper() {
	keys := generateSecretKeys("MYKEY")
	suite.Contains(keys, "MYKEY")
	suite.Contains(keys, "mykey")
	suite.Len(keys, 2)
}

// --- byteSliceAsStringPtr ---

func (suite *HelperTestSuite) TestByteSliceAsStringPtr_NonEmpty() {
	val := byteSliceAsStringPtr([]byte("hello"))
	suite.NotNil(val)
	suite.Equal("hello", *val)
}

func (suite *HelperTestSuite) TestByteSliceAsStringPtr_Empty() {
	val := byteSliceAsStringPtr([]byte{})
	suite.Nil(val)
}

// --- generateSecretFilePath ---

func (suite *HelperTestSuite) TestGenerateSecretFilePath() {
	path := generateSecretFilePath("/run/secrets", "mykey")
	suite.Equal("/run/secrets/mykey", path)
}

// --- error types ---

func (suite *HelperTestSuite) TestSecretNotFoundError() {
	err := asSecretNotFoundError("mykey")
	suite.EqualError(err, `secret not found: "mykey"`)

	var target *SecretNotFoundError
	suite.True(errors.As(err, &target))
	suite.Equal("mykey", target.Key())
}

// --- splitCredentials ---

func (suite *HelperTestSuite) TestSplitCredentials_Valid() {
	k, v := splitCredentials("mykey:myvalue")
	suite.Equal("mykey", k)
	suite.Equal("myvalue", v)
}

func (suite *HelperTestSuite) TestSplitCredentials_ValueWithColon() {
	// The value part should be kept intact even if it contains colons.
	k, v := splitCredentials("mykey:val:ue")
	suite.Equal("mykey", k)
	suite.Equal("val:ue", v)
}

func (suite *HelperTestSuite) TestSplitCredentials_NoColon() {
	k, v := splitCredentials("nokeyvalue")
	suite.Equal("", k)
	suite.Equal("", v)
}

// --- assertKeyIsEqual ---

func (suite *HelperTestSuite) TestAssertKeyIsEqual_ExactMatch() {
	suite.True(assertKeyIsEqual("MYKEY", "MYKEY"))
}

func (suite *HelperTestSuite) TestAssertKeyIsEqual_CaseSensitive() {
	suite.False(assertKeyIsEqual("mykey", "MYKEY"))
	suite.False(assertKeyIsEqual("MYKEY", "mykey"))
	suite.False(assertKeyIsEqual("MyKey", "mykey"))
}

func (suite *HelperTestSuite) TestAssertKeyIsEqual_NoMatch() {
	suite.False(assertKeyIsEqual("key1", "key2"))
}

// --- ExportToEnvironment error path ---

func (suite *HelperTestSuite) TestExportToEnvironmentMissingKey() {
	secrets := map[string]string{"PRESENT": "value"}
	manager := NewStaticSecretsManager(secrets)

	failed := ExportToEnvironment([]string{"MISSING"}, manager)
	suite.Equal([]string{"MISSING"}, failed)
}

func (suite *HelperTestSuite) TestExportToEnvironmentAllPresent() {
	secrets := map[string]string{"KEY1": "val1", "KEY2": "val2"}
	manager := NewStaticSecretsManager(secrets)

	failed := ExportToEnvironment([]string{"KEY1", "KEY2"}, manager)
	suite.Empty(failed)
}

func (suite *HelperTestSuite) TestExportToEnvironmentSetenvFailure() {
	// "=" is not a valid character in an environment variable name; os.Setenv
	// rejects it, which should surface as a failed key even though the lookup succeeded.
	invalidKey := "INVALID=KEY"
	secrets := map[string]string{invalidKey: "value"}
	manager := NewStaticSecretsManager(secrets)

	failed := ExportToEnvironment([]string{invalidKey}, manager)
	suite.Equal([]string{invalidKey}, failed)
}
