package sdk

import "os"

// GGEnvironment is a set of well known GreenGrass
// environment parameters.
type GGEnvironment struct {
	AuthToken                  string `env:"AWS_CONTAINER_AUTHORIZATION_TOKEN"`
	FunctionArn                string `env:"MY_FUNCTION_ARN"`
	EncodingType               string `env:"ENCODING_TYPE"`
	ShadowFunctionArn          string `env:"SHADOW_FUNCTION_ARN"`
	RouterFunctionArn          string `env:"ROUTER_FUNCTION_ARN"`
	GGCMaxInterfaceVersion     string `env:"GGC_MAX_INTERFACE_VERSION"`
	SecretesManagerFunctionArn string `env:"SECRETS_MANAGER_FUNCTION_ARN"`
	GGDeamonPort               string `env:"GG_DAEMON_PORT"`
}

// GetGGEnvironment fetches a new instance of `GetGGEnvironment` where it has
// be populated from the environment variables.
func GetGGEnvironment() GGEnvironment {

	return GGEnvironment{
		AuthToken:                  os.Getenv("AWS_CONTAINER_AUTHORIZATION_TOKEN"),
		FunctionArn:                os.Getenv("MY_FUNCTION_ARN"),
		EncodingType:               os.Getenv("ENCODING_TYPE"),
		ShadowFunctionArn:          os.Getenv("SHADOW_FUNCTION_ARN"),
		RouterFunctionArn:          os.Getenv("ROUTER_FUNCTION_ARN"),
		GGCMaxInterfaceVersion:     os.Getenv("GGC_MAX_INTERFACE_VERSION"),
		SecretesManagerFunctionArn: os.Getenv("SECRETS_MANAGER_FUNCTION_ARN"),
		GGDeamonPort:               os.Getenv("GG_DAEMON_PORT"),
	}

}
