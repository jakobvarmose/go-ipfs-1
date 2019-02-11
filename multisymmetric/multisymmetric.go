package multisymmetric

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"

	"golang.org/x/crypto/nacl/secretbox"
)

var (
	None            uint64 = 0
	NACL_SECRETBOX  uint64 = 1
	AES_CTR_ZERO_IV uint64 = 2
)

var Names = map[string]uint64{
	"none":      None,
	"secretbox": NACL_SECRETBOX,
	"aes":       AES_CTR_ZERO_IV,
}

func GenerateKey(mode uint64) ([]byte, error) {

	switch mode {
	case NACL_SECRETBOX:
		return random(32)
	case AES_CTR_ZERO_IV:
		return random(16)
	default:
		return nil, errors.New("invalid encryption mode")
	}
}

func Encrypt(mode uint64, key, pt []byte) ([]byte, error) {
	switch mode {
	case NACL_SECRETBOX:
		var key2 [32]byte
		copy(key2[:], key)
		var nonce [24]byte
		_, err := rand.Read(nonce[:])
		if err != nil {
			return nil, err
		}

		return secretbox.Seal(nonce[:], pt, &nonce, &key2), nil

	case AES_CTR_ZERO_IV:
		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}

		// use zero iv - only safe if we don't reuse the key
		iv := make([]byte, aes.BlockSize)

		stream := cipher.NewCTR(block, iv)
		ct := make([]byte, len(pt))
		stream.XORKeyStream(ct, pt)
		return ct, nil
	default:
		return nil, errors.New("invalid encryption mode")
	}
}

func Decrypt(mode uint64, key, ct []byte) ([]byte, error) {
	switch mode {
	case NACL_SECRETBOX:
		if len(key) != 32 {
			return nil, errors.New("wrong key length")
		}

		if len(ct) < 24+secretbox.Overhead {
			return nil, errors.New("ciphertext too short")
		}

		var key2 [32]byte
		copy(key2[:], key)
		var nonce [24]byte
		copy(nonce[:], ct[:24])
		pt, ok := secretbox.Open(nil, ct[24:], &nonce, &key2)
		if !ok {
			return nil, errors.New("decryption failed")
		}

		return pt, nil

	case AES_CTR_ZERO_IV:
		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}

		// use zero iv - only safe if we don't reuse the key
		iv := make([]byte, aes.BlockSize)

		stream := cipher.NewCTR(block, iv)
		pt := make([]byte, len(ct))
		stream.XORKeyStream(pt, ct)
		return pt, nil
	default:
		return nil, errors.New("invalid encryption mode")
	}
}

func random(length int) ([]byte, error) {
	key := make([]byte, length)
	_, err := rand.Read(key[:])
	if err != nil {
		return nil, err
	}

	return key, nil
}
