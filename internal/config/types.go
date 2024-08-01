package config

type RpcConfig struct {
	ChainId   uint   `yaml:"chain_id"`
	Url       string `yaml:"url"` // todo: consider making this a different type if needed
	StreamUrl string `yaml:"stream_url"`
	// todo: consider backup urls
	Blocktime uint `yaml:"blocktime"` // todo this can be replaced by a dynamic block-time determination
}

type LogConfig struct {
	Level    uint   `yaml:"level"` // todo: this should match log level
	Format   string `yaml:"format"`
	FilePath string `yaml:"file_path"`
}

// TODO: add fields for outbound handler struct
type OutboundHandlerConfig struct {
}

type AppConfig struct {
	Rpcs    []RpcConfig `yaml:"rpcs"`
	Logging LogConfig   `yaml:"logging"`
	// OutboundHandlers []OutboundHandlerConfig `yaml:"outbound_handlers"`
}
