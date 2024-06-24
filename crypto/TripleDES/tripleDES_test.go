package tripleDES

import (
	"fmt"
	"testing"
)

func TestUsernameCodec(t *testing.T) {
	key := "0123456789abcdef01234567"
	iv := "01234567"
	username := "testuser"

	codec := NewTripleDES(key, iv)

	// 测试编码器
	encoded, err := codec.Encode(username)
	fmt.Println("encoded:", encoded)
	if err != nil {
		t.Errorf("Encode failed: %v", err)
	}

	if encoded == "" {
		t.Error("Encode result is empty")
	}

	// 测试解码器
	decoded, err := codec.Decode(encoded)
	if err != nil {
		t.Errorf("Decode failed: %v", err)
	}

	if decoded != username {
		t.Errorf("Decode result is not equal to original: %s != %s", decoded, username)
	}
}

func TestDecodec(t *testing.T) {
	key := "P6ieR6N8iTRXzZHkMoZBlKfD"
	iv := "feXU1TAS"
	encoded := "tB3UVXrTslQ="

	decoded, err := NewTripleDES(key, iv).Decode(encoded)
	if err != nil {
		t.Errorf("Decode failed: %v", err)
	}

	fmt.Println(decoded)
}
