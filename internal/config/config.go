package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DB_url string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func (c Config) SetUser(user string) {
	c.Current_user_name = user
	err := write(c)
	if err != nil {
		fmt.Printf("%v", err)
	}
}

func write(c Config) error {
	path, err := getConfigPath()
	if err != nil{
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData, err := json.Marshal(&c)
	if err != nil {
		return err
	}


	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}

func Read() Config {
	path, err := getConfigPath()
	if err != nil{
		fmt.Printf("%v", err)
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("%v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Printf("%v", err)
	}

	return config
}

func getConfigPath() (string, error) {
	dir, err := os.UserHomeDir() 
	if err != nil {
		return "", err
	}

	return dir + "/" + configFileName, nil
}