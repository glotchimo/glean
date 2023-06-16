package main

import (
	"encoding/json"
	"os"
	"strconv"
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
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	PostsPath  string `yaml:"posts_path"`
	EmailsPath string `yaml:"emails_path"`
	Meta       Meta   `yaml:"meta"`
	SMTP       SMTP   `yaml:"smtp"`
}

func LoadConf() error {
	CONF.Meta.Title = os.Getenv("GLEAN_META_TITLE")
	CONF.Meta.Author = os.Getenv("GLEAN_META_AUTHOR")
	CONF.Meta.Email = os.Getenv("GLEAN_META_EMAIL")
	json.Unmarshal([]byte(os.Getenv("GLEAN_META_LINKS")), &CONF.Meta.Links)

	CONF.SMTP.Host = os.Getenv("RAILWAY_STATIC_URL")
	CONF.SMTP.Port, _ = strconv.Atoi(os.Getenv("GLEAN_SMTP_PORT"))
	CONF.SMTP.Username = os.Getenv("GLEAN_SMTP_USERNAME")
	CONF.SMTP.Password = os.Getenv("GLEAN_SMTP_PASSWORD")
	CONF.SMTP.Sender = os.Getenv("GLEAN_SMTP_SENDER")

	CONF.Host = os.Getenv("RAILWAY_STATIC_URL")
	CONF.Port = os.Getenv("GLEAN_PORT")

	return nil
}
