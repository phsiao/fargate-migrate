// fargate only supports specific configuration of cpu and memory
// so this module provides a function to find minimal matching configuration
// for the given cpu and memory
// https://docs.aws.amazon.com/AmazonECS/latest/developerguide/AWS_Fargate.html
package fargate

type cpuMemoryConfig struct {
	cpuConfig     int
	memoryConfigs []int
}

var (
	supportedCPUMemory = []cpuMemoryConfig{
		{
			cpuConfig:     256,
			memoryConfigs: []int{512, 1024, 2048}},
	}
)

func init() {
	var config cpuMemoryConfig

	config = cpuMemoryConfig{
		cpuConfig:     512,
		memoryConfigs: []int{},
	}
	for i := 1024; i <= 4096; i += 1024 {
		config.memoryConfigs = append(config.memoryConfigs, i)
	}
	supportedCPUMemory = append(supportedCPUMemory, config)

	config = cpuMemoryConfig{
		cpuConfig:     1024,
		memoryConfigs: []int{},
	}
	for i := 2048; i <= 8192; i += 1024 {
		config.memoryConfigs = append(config.memoryConfigs, i)
	}
	supportedCPUMemory = append(supportedCPUMemory, config)

	config = cpuMemoryConfig{
		cpuConfig:     2048,
		memoryConfigs: []int{},
	}
	for i := 4096; i <= 16384; i += 1024 {
		config.memoryConfigs = append(config.memoryConfigs, i)
	}
	supportedCPUMemory = append(supportedCPUMemory, config)

	config = cpuMemoryConfig{
		cpuConfig:     4096,
		memoryConfigs: []int{},
	}
	for i := 8192; i <= 30720; i += 1024 {
		config.memoryConfigs = append(config.memoryConfigs, i)
	}
	supportedCPUMemory = append(supportedCPUMemory, config)
}

func MinCPUMemroyConfiguration(cpu, memory int) (matchCPU, matchMemory int) {
	if matchCPU == 0 {
		maxConfig := supportedCPUMemory[len(supportedCPUMemory)-1]
		matchCPU = maxConfig.cpuConfig
		matchMemory = maxConfig.memoryConfigs[len(maxConfig.memoryConfigs)-1]
	}

	for _, config := range supportedCPUMemory {
		if cpu > config.cpuConfig {
			continue
		}
		matchCPU = config.cpuConfig
		matchMemory = config.memoryConfigs[len(config.memoryConfigs)-1]
		for _, memoryConfig := range config.memoryConfigs {
			if memory > memoryConfig {
				continue
			} else {
				matchMemory = memoryConfig
				return
			}
		}
	}

	return
}
