package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"golang.org/x/crypto/ssh"
)

// KeyPair is a public and private key pair that can be used for SSH access.
type KeyPair struct {
	PublicKey  string
	PrivateKey string
}

// GenerateRSAKeyPair generates an RSA Keypair and return the public and private keys.
func GenerateRSAKeyPair(t *testing.T, keySize int) *KeyPair {
	keyPair, err := GenerateRSAKeyPairE(t, keySize)
	if err != nil {
		t.Fatal(err)
	}
	return keyPair
}

// GenerateRSAKeyPairE generates an RSA Keypair and return the public and private keys.
func GenerateRSAKeyPairE(t *testing.T, keySize int) (*KeyPair, error) {
	logger.Logf(t, "Generating new public/private key of size %d", keySize)

	rsaKeyPair, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, err
	}

	// Extract the private key
	keyPemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(rsaKeyPair),
	}

	keyPem := string(pem.EncodeToMemory(keyPemBlock))

	// Extract the public key
	sshPubKey, err := ssh.NewPublicKey(rsaKeyPair.Public())
	if err != nil {
		return nil, err
	}

	sshPubKeyBytes := ssh.MarshalAuthorizedKey(sshPubKey)
	sshPubKeyStr := string(sshPubKeyBytes)

	// Return
	return &KeyPair{PublicKey: sshPubKeyStr, PrivateKey: keyPem}, nil
}

// TODO Vit: document this
func LoadRSAKeyPair(t *testing.T, privateKeyPath string, publicKeyPath string) *KeyPair {
	keyPair, err := LoadRSAKeyPairE(t, privateKeyPath, publicKeyPath)
	if err != nil {
		t.Fatal(err)
	}
	return keyPair
}

// TODO Vit: document this
func LoadRSAKeyPairE(t *testing.T, publicKeyPath string, privateKeyPath string) (*KeyPair, error) {
	publicKey, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	privateKey, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	return &KeyPair{
		PublicKey:  string(publicKey),
		PrivateKey: string(privateKey),
	}, nil
}
