package main

import "github.com/kelseyhightower/envconfig"

type config struct {
	Database  string `required:"true"`
	LogLevel  string `default:"info" split_words:"true"`
	LogPretty bool   `default:"false" split_words:"true"`
	Port      string `default:"8080" envconfig:"PORT"`
}

func initConfig() (*config, error) {
	c := &config{}

	if err := envconfig.Process("", c); err != nil {
		return nil, err
	}

	return c, nil
}
