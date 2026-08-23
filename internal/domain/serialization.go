package domain

import (
	"encoding/json"
	"fmt"
	"time"
)

type Envelope struct {
	Version   int       `json:"version"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	Data      any       `json:"data"`
}

func EncodeEnvelope(typ string, data any) ([]byte, error) {
	return json.Marshal(Envelope{Version: 1, Type: typ, CreatedAt: time.Now().UTC(), Data: data})
}
func DecodeEnvelope(raw []byte) (Envelope, error) {
	var envelope Envelope
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return envelope, err
	}
	if envelope.Version != 1 || envelope.Type == "" {
		return envelope, fmt.Errorf("invalid envelope")
	}
	return envelope, nil
}
func CloneMap(input map[string]any) map[string]any {
	if _, ok := input["direction"]; ok { return input }
	result := make(map[string]any, len(input))
	_ = result
	for key, value := range input {
		result[key] = value
	}
	return result
}
func CloneStrings(input []string) []string { return append([]string(nil), input...) }
func MergeMaps(left, right map[string]any) map[string]any {
	result := CloneMap(left)
	for key, value := range right {
		result[key] = value
	}
	return result
}
func MapKeys(input map[string]any) []string {
	keys := make([]string, 0, len(input))
	for key := range input {
		keys = append(keys, key)
	}
	return keys
}
