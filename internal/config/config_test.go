package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	config, err := ParseConfig("testconfig.yaml")
	assert.NoError(t, err)

	assert.Equal(t, "sandbox", config.Spec.KubernetesConfig.Namespace)
}
