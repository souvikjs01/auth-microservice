package utils

func GetConfigPath(configPath string) string {
	if configPath == "docker" {
		return "./config/config-docker.yml"
	}
	return "./config/config-local.yml"
}
