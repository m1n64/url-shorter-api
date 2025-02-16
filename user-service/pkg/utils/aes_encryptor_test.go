package utils

import (
	"encoding/base64"
	"os"
	"testing"
)

func TestAESEncryptor_Encrypt(t *testing.T) {
	os.Setenv("APP_SECRET_KEY", "base64:Y3qM+TZfAf7Be41uqOeGwhwTRrLAkSYW7nsVlBYa3XA=")

	encryptor := NewAESEncryptor()

	plainText := "9c4db945-c6d1-4ff4-bbc8-e79e8950762e:1725379112"

	key := encryptor.GetKey()

	iv, _ := base64.StdEncoding.DecodeString("5ua2s0Pr6q50wsWJFB1DQw==")

	result := "5ua2s0Pr6q50wsWJFB1DQxUrkIt8Hb11Cg2Mp7PVTgj8CdRjZWe2JeafYt4wDWhq+v70hWIaEVZiL5bRpbpWgg=="

	encrypted, err := encryptor.Encrypt(plainText, key, iv)

	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	if result != encrypted {
		t.Errorf("Failed to encrypt: %v", encrypted)
	}

	t.Log(encrypted)
}

func TestAESEncryptor_Decrypt(t *testing.T) {
	os.Setenv("APP_SECRET_KEY", "base64:Y3qM+TZfAf7Be41uqOeGwhwTRrLAkSYW7nsVlBYa3XA=")

	encryptor := NewAESEncryptor()

	plainText := "Y5yZJsbwzHIxAYfpv8CdPVdR1YEmgtUUOAW69eXX8xrLWhbr37BmQdnpq4FCGH1OWsMRm5CFZ1UGRvF5Cp0YgQ==" // Замените на ваши данные

	result := "9c4db945-c6d1-4ff4-bbc8-e79e8950762e:1725379112"

	key := encryptor.GetKey()

	iv, encBaseStr := encryptor.GetIVAndCipher(plainText)

	decrypted, err := encryptor.Decrypt(encBaseStr, key, iv)

	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	if result != decrypted {
		t.Errorf("Decrypted data is incorrect: got %s, want %s", decrypted, result)
	}

	t.Log("Decrypted data:", decrypted)
}
