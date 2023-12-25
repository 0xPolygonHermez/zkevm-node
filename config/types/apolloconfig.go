package types

// ApolloConfig is the config for apollo
type ApolloConfig struct {
	Enable        bool   `mapstructure:"Enable"`
	IP            string `mapstructure:"IP"`
	AppID         string `mapstructure:"AppID"`
	NamespaceName string `mapstructure:"NamespaceName"`
}
