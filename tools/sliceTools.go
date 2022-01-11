package tools

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"sort"
)

func UniqueSlice(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

func UniqueStructSlice(input []interface{}) []interface{} {
	u := make([]interface{}, 0, len(input))
	m := make(map[uint64]bool)
	for _, val := range input {
		code := HashCode(val)
		if _, ok := m[code]; !ok {
			m[code] = true
			u = append(u, val)
		}
	}

	return u
}

func Contains(array []string, s string) bool {
	copy := append([]string{}, array...)
	sort.Strings(copy)
	i := sort.SearchStrings(copy, s)
	return i < len(copy) && copy[i] == s
}

func HashCode(h interface{}) uint64 {
	bytes, _ := json.Marshal(h)
	hash := sha256.Sum256(bytes)
	return binary.BigEndian.Uint64(hash[:])
}
