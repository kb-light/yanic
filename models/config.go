package models

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

//Config the config File of this daemon
type Config struct {
	Responedd struct {
		Enable          bool   `yaml:"enable"`
		Port            string `yaml:"port"`
		Address         string `yaml:"address"`
		CollectInterval int    `yaml:"collectinterval"`
	} `yaml:"responedd"`
	Webserver struct {
		Enable           bool   `yaml:"enable"`
		Port             string `yaml:"port"`
		Address          string `yaml:"address"`
		Webroot          string `yaml:"webroot"`
		WebsocketNode    bool   `yaml:"websocketnode"`
		WebsocketAliases bool   `yaml:"websocketaliases"`
	} `yaml:"webserver"`
	Nodes struct {
		Enable        bool     `yaml:"enable"`
		NodesPath     string   `yaml:"nodes_path"`
		GraphsPath    string   `yaml:"graphs_path"`
		AliasesEnable bool     `yaml:"aliases_enable"`
		AliasesPath   string   `yaml:"aliases_path"`
		SaveInterval  int      `yaml:"saveinterval"`
		VpnAddresses  []string `yaml:"vpn_addresses"`
	} `yaml:"nodes"`
}

//ConfigReadFile reads a Config models by path to a yml file
func ConfigReadFile(path string) *Config {
	config := &Config{}
	file, _ := ioutil.ReadFile(path)
	err := yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
