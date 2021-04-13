package fargate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinCPUMemroyConfiguration(t *testing.T) {
	tests := []struct {
		cpu    int
		memory int
		want   [2]int
	}{
		{
			cpu:    100,
			memory: 100,
			want:   [2]int{256, 512},
		},
		{
			cpu:    100,
			memory: 1034,
			want:   [2]int{256, 2048},
		},
		{
			cpu:    1047,
			memory: 2020,
			want:   [2]int{2048, 4096},
		},
		{
			cpu:    10000,
			memory: 103400,
			want:   [2]int{4096, 30720},
		},
	}

	for _, tc := range tests {
		gotCPU, gotMemory := MinCPUMemroyConfiguration(tc.cpu, tc.memory)
		assert.Equal(t, tc.want[0], gotCPU)
		assert.Equal(t, tc.want[1], gotMemory)
	}
}
