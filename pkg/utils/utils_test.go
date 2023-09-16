package utils

import (
	"testing"

	"github.com/goccy/go-json"
)

// Test JsonToStruct
func TestJsonToStruct(t *testing.T) {
	type Value struct {
		Value string `json:"value"`
	}

	a := `{"value":"add_integration_exchange"}`
	var b Value

	err := json.Unmarshal([]byte(a), &b)
	if err != nil {
		t.Error(err)
	}
}

// Test StructToJson
func TestStructToJson(t *testing.T) {
	a := "test"
	b, err := StructToJson(a)
	if err != nil {
		t.Error(err)
	}

	if b != `"test"` {
		t.Error("wrong struct to json")
	}
}
