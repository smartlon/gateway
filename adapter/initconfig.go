package main

import (
	"github.com/smartlon/gateway/config"
	"github.com/smartlon/gateway/log"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	FlagHome   = "home"
	FlagConfig = "config"
	FlagLog    = "log"
	FlagQueue  = "queue"
)



func initConfig() error {
	// init config
	defaultHome := os.ExpandEnv("$HOME/.gateway")
	defaultConfig := filepath.Join(defaultHome, "config/gateway.yml")
	//defaultLog := filepath.Join(defaultHome, "config/log.conf")

	log.Debug("home: ", defaultHome)
	viper.Set(FlagHome, defaultHome)

	// // Sets name for the config file.
	// // Does not include extension.
	// viper.SetConfigName("config")
	// // Adds a path for Viper to search for the config file in.
	// viper.AddConfigPath(filepath.Join(homeDir, "config"))
	// // Can be called multiple times to define multiple search paths.
	// // viper.AddConfigPath(homeDir)

	log.Debug("Init config: ", defaultConfig)
	viper.Set(FlagConfig, defaultConfig)
	viper.SetConfigFile(viper.GetString(FlagConfig))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
	} else if _, ok := err.(viper.ConfigFileNotFoundError); ok ||
		os.IsNotExist(err) {
		log.Warn(err.Error())
		if !strings.HasPrefix(
			viper.GetString(FlagConfig), viper.GetString(FlagHome)) {
			return err
		}
		if err := cmn.EnsureDir(
			viper.GetString(FlagHome), 0700); err != nil {
			panic(err.Error())
		}
		if err := cmn.EnsureDir(
			filepath.Join(viper.GetString(FlagHome), "config"),
			0700); err != nil {
			panic(err.Error())
		}
		// create & write default config
		var bytes []byte
		bytes, err = yaml.Marshal(config.DefaultConfig())
		if err != nil {
			log.Error("Marshal config error: ", err.Error())
			return err
		}
		err = ioutil.WriteFile(viper.GetString(FlagConfig), bytes, 0644)
		if err != nil {
			log.Error("write config file error: ", err.Error())
			return err
		}
	} else {
		log.Error("Load config error: ", err.Error())
		return err
	}

	return config.GetConfig().Load()
}