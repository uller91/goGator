package config

import (
	"os"
	"fmt"
	"encoding/json"
	"io"
)

const (
	configFileName = ".gatorconfig.json"
	PathSeparator  = "/"
)

type Config struct {
	DbUrl		string	`json:"db_url"`
	UserName	string	`json:"current_user_name"`
	
}

func (c *Config) GetUrl() string {
	return c.DbUrl
}

func (c *Config) SetUser(userName string) error {
	c.UserName = userName

	err := write(*c)
	if err != nil {
        return err
    }

	return nil
}

func write(cfg Config) error {
	byteData, err := json.Marshal(cfg)
    if err != nil {
        return err
    }

	homePath, err := getConfigFilePath()
	if err != nil {
		fmt.Println(err)
		return err
	}

	absolutePath := homePath + PathSeparator + configFileName

	err = os.WriteFile(absolutePath, byteData, 0600) //permission "owner can read/write, others have no access"
	if err != nil {
		return err
	}

	return nil
}


func Read() (Config, error) {
	var cfg Config
	homePath, err := getConfigFilePath()
	//fmt.Println(homePath)
	if err != nil {
		fmt.Println(err)
		return cfg, err
	}

	absolutePath := homePath + PathSeparator + configFileName

	jsonFile, err := os.Open(absolutePath)
	if err != nil {
		//fmt.Println(err)
		return cfg, err
	}

	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return cfg, err
	}
	//fmt.Println(byteValue)

	
	err = json.Unmarshal(byteValue, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}


func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	//fmt.Println(homePath)

	return homePath, nil
}