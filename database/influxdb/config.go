package influxdb

type config struct {
	Address  string
	Database string
	Username string
	Password string
}

func toConfig(configMap map[string]interface{}) *config {
	return &config{
		Address:  configMap["address"].(string),
		Database: configMap["database"].(string),
		Username: configMap["username"].(string),
		Password: configMap["password"].(string),
	}
}
