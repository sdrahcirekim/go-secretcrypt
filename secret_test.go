package secretcrypt

import (
	"flag"
	"os"
	"testing"

	"github.com/Zemanta/go-secretcrypt/internal"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	internal.CryptersMap["plain"] = internal.PlainCrypter{}
	flag.Parse()
	os.Exit(m.Run())
}

func assertStrictSecretValid(t *testing.T, secret StrictSecret) {
	assert.Equal(t, "plain", secret.crypter.Name())
	assert.Equal(t, "my-abc", string(secret.ciphertext))
	assert.Equal(t, internal.DecryptParams{
		"k1": "v1",
		"k2": "v2",
	}, secret.decryptParams)
}

func TestUnmarshalText(t *testing.T) {
	var secret StrictSecret
	err := secret.UnmarshalText([]byte("plain:k1=v1&k2=v2:my-abc"))
	assert.Nil(t, err)
	assertStrictSecretValid(t, secret)
}

func TestNewStrictSecret(t *testing.T) {
	secret, err := LoadStrictSecret("plain:k1=v1&k2=v2:my-abc")
	assert.Nil(t, err)
	assertStrictSecretValid(t, secret)
}

func TestDecrypt(t *testing.T) {
	mockCrypter := &internal.MockCrypter{}
	internal.CryptersMap["mock"] = mockCrypter
	mockCrypter.On(
		"Decrypt",
		internal.Ciphertext("my-abc"),
		internal.DecryptParams{
			"k1": "v1",
			"k2": "v2",
		}).Return("myplaintext", nil)

	secret, err := LoadStrictSecret("mock:k1=v1&k2=v2:my-abc")
	assert.NoError(t, err)

	plaintext, err := secret.Decrypt()
	assert.NoError(t, err)
	assert.Equal(t, plaintext, "myplaintext")
	mockCrypter.AssertExpectations(t)
}

func TestNoCaching(t *testing.T) {
	mockCrypter := &internal.MockCrypter{}
	internal.CryptersMap["mock"] = mockCrypter
	mockCrypter.On(
		"Decrypt",
		internal.Ciphertext("my-abc"),
		internal.DecryptParams{
			"k1": "v1",
			"k2": "v2",
		}).Return("myplaintext", nil)

	secret, err := LoadStrictSecret("mock:k1=v1&k2=v2:my-abc")
	assert.NoError(t, err)

	plaintext, err := secret.Decrypt()
	assert.NoError(t, err)
	assert.Equal(t, plaintext, "myplaintext")
	plaintext2, err := secret.Decrypt()
	assert.NoError(t, err)
	assert.Equal(t, plaintext2, "myplaintext")

	mockCrypter.AssertExpectations(t)
	mockCrypter.AssertNumberOfCalls(t, "Decrypt", 2)
}

func TestEmptyStrictSecret(t *testing.T) {
	zero := StrictSecret{}
	emptyStr, err := LoadStrictSecret("")
	assert.NoError(t, err)
	for _, secret := range []StrictSecret{zero, emptyStr} {
		plaintext, err := secret.Decrypt()
		assert.NoError(t, err)
		assert.Equal(t, plaintext, "")
	}
}

func TestSecret(t *testing.T) {
	var secret Secret
	err := secret.UnmarshalText([]byte("plain:k1=v1&k2=v2:my-abc"))
	assert.Nil(t, err)
	assert.Equal(t, "my-abc", secret.Get())
}
