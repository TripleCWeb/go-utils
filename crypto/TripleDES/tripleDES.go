package tripleDES

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
)

// TripleDESCodec 包含基于 3DES 的用户名编解码器
type TripleDESCodec struct {
	key []byte
	iv  []byte
}

// NewTripleDES 创建一个基于 3DES 的用户名编解码器实例
func NewTripleDES(key, iv string) *TripleDESCodec {
	return &TripleDESCodec{
		key: []byte(key),
		iv:  []byte(iv),
	}
}

// Encode 使用 3DES 算法将用户名字符串编码为 Base64 格式
func (codec *TripleDESCodec) Encode(username string) (string, error) {
	// 创建 3DES 加密算法实例
	block, err := des.NewTripleDESCipher(codec.key)
	if err != nil {
		return "", err
	}

	// 填充数据，使其长度为 8 的倍数
	plaintext := PKCS5Padding([]byte(username), block.BlockSize())

	// 创建 CBC 模式的加密器
	mode := cipher.NewCBCEncrypter(block, codec.iv)

	// 加密数据
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)

	// 将加密后的数据转换为 Base64 格式
	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	return encoded, nil
}

// Decode 使用 3DES 算法将 Base64 格式的用户名解码为字符串格式
func (codec *TripleDESCodec) Decode(encodedUsername string) (string, error) {
	// 将 Base64 格式的数据解码为二进制格式
	ciphertext, err := base64.StdEncoding.DecodeString(encodedUsername)
	if err != nil {
		return "", err
	}

	// 创建 3DES 加密算法实例
	block, err := des.NewTripleDESCipher(codec.key)
	if err != nil {
		return "", err
	}

	// 创建 CBC 模式的解密器
	mode := cipher.NewCBCDecrypter(block, codec.iv)

	// 解密数据
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// 去除填充的数据
	plaintext = PKCS5UnPadding(plaintext)

	return string(plaintext), nil
}

// PKCS5Padding 对数据进行 PKCS5 填充
func PKCS5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// PKCS5UnPadding 对数据进行 PKCS5 反填充
func PKCS5UnPadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}
