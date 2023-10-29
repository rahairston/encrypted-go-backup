package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"

	"golang.org/x/crypto/ssh"
)

type KeyHandler struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func BuildKeyHandler(fileName string) (*KeyHandler, error) {
	pubKey, err := os.ReadFile(fileName + ".pub")
	if err != nil {
		return nil, err
	}

	parsed, _, _, _, err := ssh.ParseAuthorizedKey(pubKey)
	if err != nil {
		return nil, err
	}
	// To get back to an *rsa.PublicKey, we need to first upgrade to the
	// ssh.CryptoPublicKey interface
	parsedCryptoKey := parsed.(ssh.CryptoPublicKey)

	// Then, we can call CryptoPublicKey() to get the actual crypto.PublicKey
	pubCrypto := parsedCryptoKey.CryptoPublicKey()

	// Finally, we can convert back to an *rsa.PublicKey
	pub := pubCrypto.(*rsa.PublicKey)

	privKey, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privKey)
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &KeyHandler{
		publicKey:  pub,
		privateKey: priv,
	}, nil
}

func (keys KeyHandler) Encrypt(data []byte) ([]byte, error) {
	hash := sha512.New()
	dataLength := len(data)
	step := keys.publicKey.Size() - 2*hash.Size() - 2
	var encryptedBytes []byte

	for start := 0; start < dataLength; start += step {
		end := start + step
		if end > dataLength {
			end = dataLength
		}

		encryptedBlockBytes, err := rsa.EncryptOAEP(
			hash,
			rand.Reader,
			keys.publicKey,
			data[start:end],
			nil)

		if err != nil {
			return nil, err
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}

	return []byte(base64.StdEncoding.EncodeToString(encryptedBytes)), nil
}

func (keys KeyHandler) Decrypt(data []byte) ([]byte, error) {

	data2, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}

	hash := sha512.New()
	dataLength := len(data2)
	step := keys.publicKey.Size()
	var decryptedBytes []byte

	for start := 0; start < dataLength; start += step {
		end := start + step
		if end > dataLength {
			end = dataLength
		}

		decryptedBlockBytes, err := rsa.DecryptOAEP(
			hash,
			rand.Reader,
			keys.privateKey,
			data2[start:end],
			nil)

		if err != nil {
			return nil, err
		}

		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}
