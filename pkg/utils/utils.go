package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"sync"
	"time"

	"github.com/goccy/go-json"
)

// Pointer return pointer for value
func Pointer[T any](v T) *T {
	return &v
}

func ApplyToPointer[T any, V any](v *T, applier func(T) V) *V {
	if v == nil {
		return nil
	}

	return Pointer(applier(*v))
}

func ExistInArray[T comparable](arr []T, value T) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

func GetRandomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

// JsonToStruct return populates the fields of the dst struct from the fields
// of the src struct using json tags
func JsonToStruct(src interface{}, dst interface{}) error {
	result, err := json.Marshal(src)
	if err != nil {
		return err
	}

	return json.Unmarshal(result, dst)
}

// StructToJson return populates the fields of the dst struct from the fields
// of the src struct using json tags
func StructToJson(src interface{}) (string, error) {
	result, err := json.Marshal(src)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// UpdateStruct return populates the fields of the dst struct from the fields
func UpdateStruct(src interface{}, dst interface{}, fieldMask []string) error {
	var srcMap map[string]interface{}
	if err := JsonToStruct(src, &srcMap); err != nil {
		return err
	}

	var resultMap map[string]interface{}
	if err := JsonToStruct(src, &resultMap); err != nil {
		return err
	}

	for _, field := range fieldMask {
		if _, ok := srcMap[field]; ok {
			resultMap[field] = srcMap[field]
		}
	}

	return JsonToStruct(resultMap, dst)
}

func Contains[T any](array []T, value T, predicate func(T, T) bool) bool {
	for _, item := range array {
		if predicate(value, item) {
			return true
		}
	}
	return false
}

func Find[T any](array []*T, predicate func(item *T) bool) *T {
	for _, item := range array {
		if predicate(item) {
			return item
		}
	}
	return nil
}

func Map[T, U any](data []T, f func(T) U) []U {
	res := make([]U, 0, len(data))

	for _, e := range data {
		res = append(res, f(e))
	}

	return res
}

func MapWithError[T, U any](data []T, f func(T) (U, error)) ([]U, error) {
	res := make([]U, 0, len(data))

	for _, e := range data {
		item, err := f(e)
		if err != nil {
			return nil, err
		}

		res = append(res, item)
	}

	return res, nil
}

func IsPointersValuesEqual[T string](pointerA, pointerB *T) bool {
	if pointerA == nil && pointerB == nil {
		return true
	}

	if pointerA == nil || pointerB == nil {
		return false
	}

	return *pointerA == *pointerB
}

func GetFirstOrNull[T any](array []*T) *T {
	if len(array) > 0 {
		return array[0]
	}
	return nil
}

func GenerateDigitCode(len int) string {
	rand.Seed(time.Now().UnixNano())
	var result string
	for i := 0; i < len; i++ {
		result += fmt.Sprint(rand.Intn(10))
	}
	return result
}

type Trottler struct {
	mu       sync.Mutex
	duration time.Duration
}

func NewTrottler(duration time.Duration) *Trottler {
	return &Trottler{
		mu:       sync.Mutex{},
		duration: duration,
	}
}

func (t *Trottler) Lock() {
	t.mu.Lock()

	go func() {
		time.Sleep(t.duration)
		t.mu.Unlock()
	}()
}

const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// Encode encodes the given version and data to a base58check encoded string
func Base58CheckEncode(version, data string) (string, error) {
	prefix, err := hex.DecodeString(version)
	if err != nil {
		return "", err
	}
	dataBytes, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}
	dataBytes = append(prefix, dataBytes...)

	// Performing SHA256 twice
	sha256hash := sha256.New()
	sha256hash.Write(dataBytes)
	middleHash := sha256hash.Sum(nil)
	sha256hash = sha256.New()
	sha256hash.Write(middleHash)
	hash := sha256hash.Sum(nil)

	checksum := hash[:4]
	dataBytes = append(dataBytes, checksum...)

	// For all the "00" versions or any prepended zeros as base58 removes them
	zeroCount := 0
	for i := 0; i < len(dataBytes); i++ {
		if dataBytes[i] == 0 {
			zeroCount++
		} else {
			break
		}
	}

	// Performing base58 encoding
	encoded := b58encode(dataBytes)

	for i := 0; i < zeroCount; i++ {
		encoded = "1" + encoded
	}

	return encoded, nil
}

// Decode decodes the given base58check encoded string and returns the version prepended decoded string
func Base58CheckDecode(encoded string) (string, error) {
	zeroCount := 0
	for i := 0; i < len(encoded); i++ {
		if encoded[i] == 49 {
			zeroCount++
		} else {
			break
		}
	}

	dataBytes, err := b58decode(encoded)
	if err != nil {
		return "", err
	}

	dataBytesLen := len(dataBytes)
	if dataBytesLen <= 4 {
		return "", errors.New("base58check data cannot be less than 4 bytes")
	}

	data, checksum := dataBytes[:dataBytesLen-4], dataBytes[dataBytesLen-4:]

	for i := 0; i < zeroCount; i++ {
		data = append([]byte{0}, data...)
	}

	// Performing SHA256 twice to validate checksum
	sha256hash := sha256.New()
	sha256hash.Write(data)
	middleHash := sha256hash.Sum(nil)
	sha256hash = sha256.New()
	sha256hash.Write(middleHash)
	hash := sha256hash.Sum(nil)

	if !reflect.DeepEqual(checksum, hash[:4]) {
		return "", errors.New("Data and checksum don't match")
	}

	return hex.EncodeToString(data), nil
}

func b58encode(data []byte) string {
	var encoded string
	decimalData := new(big.Int)
	decimalData.SetBytes(data)
	divisor, zero := big.NewInt(58), big.NewInt(0)

	for decimalData.Cmp(zero) > 0 {
		mod := new(big.Int)
		decimalData.DivMod(decimalData, divisor, mod)
		encoded = string(alphabet[mod.Int64()]) + encoded
	}

	return encoded
}

func b58decode(data string) ([]byte, error) {
	decimalData := new(big.Int)
	alphabetBytes := []byte(alphabet)
	multiplier := big.NewInt(58)

	for _, value := range data {
		pos := bytes.IndexByte(alphabetBytes, byte(value))
		if pos == -1 {
			return nil, errors.New("Character not found in alphabet")
		}
		decimalData.Mul(decimalData, multiplier)
		decimalData.Add(decimalData, big.NewInt(int64(pos)))
	}

	return decimalData.Bytes(), nil
}

func IsEqualValues[T comparable](str1, str2 *T) bool {
	if str1 == nil && str2 == nil {
		return true
	}

	if str1 != nil && str2 != nil {
		return *str1 == *str2
	}

	return false
}
