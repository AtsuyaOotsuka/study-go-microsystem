package test_funcs

import (
	"os"
	"testing"
)

type Envs map[string]string

func WithEnvMap(envs Envs, t *testing.T, fn func()) {
	originals := make(map[string]string)
	for key, value := range envs {
		originals[key] = os.Getenv(key)
		if err := os.Setenv(key, value); err != nil {
			t.Fatalf("Failed to set environment variable %s: %v", key, err)
		}
	}

	defer func() {
		for key, original := range originals {
			os.Setenv(key, original)
		}
	}()

	fn()
}
