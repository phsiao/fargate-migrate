package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	config, err := ParseConfig("testconfig.yaml")
	assert.NoError(t, err)

	assert.NotNil(t, config.Spec.KubernetesConfig.Namespace)
	assert.EqualValues(t, "sandbox", *config.Spec.KubernetesConfig.Namespace)
}
