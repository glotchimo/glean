package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Meta struct {
	Title  string            `yaml:"title"`
	Author string            `yaml:"author"`
	Email  string            `yaml:"email"`
	Links  map[string]string `yaml:"links"`
}

type SMTP struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Sender   string `yaml:"sender"`
}

type Conf struct {
	PostsPath  string `yaml:"posts_path"`
	EmailsPath string `yaml:"emails_path"`
	Meta       Meta   `yaml:"meta"`
	SMTP       SMTP   `yaml:"smtp"`
}

func LoadConf() error {
	f, err := os.Open(PATH)
	if err != nil {
		return fmt.Errorf("error opening config: %w", err)
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&CONF); err != nil {
		return fmt.Errorf("error parsing config: %w", err)
	}

	return nil
}
