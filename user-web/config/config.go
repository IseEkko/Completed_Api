package config

type UserSrvConfig struct {
	Name string `mapstructure:"name"`
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type CosulConfig struct {
	Host string `json:"host" mapstructure:"host"`
	Port int    `json:"port" mapstructure:"port"`
}

type ServerConfig struct {
	Name            string        `mapstructure:"name"`
	UserSrv         UserSrvConfig `mapstructure:"user_srv"`
	JWTInfo         JWTConfig     `mapstructure:"jwt"`
	CosulConfiginfo CosulConfig   `mapstructure:"consul"`
}
