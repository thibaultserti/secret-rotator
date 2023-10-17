package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"secret-rotator/pkg/config"
	"secret-rotator/pkg/secretrotator"

	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

func main() {
	var configuration config.Configuration

	configPath := os.Getenv("SECRETROTATOR_CONFIG_PATH")
	if configPath == "" {
		configPath = "secretrotator.yaml"
	}

	configuration, err := config.LoadConfig(configPath)
	if err != nil {
		logrus.Fatal("Cannot load configuration")
	}
	level, err := logrus.ParseLevel(configuration.LogLevel)
	if err != nil {
		logrus.Fatal("Log level invalid")
	}

	if configuration.Env != "prod" && configuration.Env != "prd" {
		logrus.SetFormatter(&logrus.TextFormatter{})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(level)

	// validate configuration
	if !config.ValidateConfig(configuration) {
		logrus.Fatal("Invalid configuration")
	}

	// run secrets rotation
	var secret string
	for _, operation := range configuration.Secrets {
		logrus.Infof("Treating %s", operation.Name)

		// src
		switch secretType := operation.Source.Type; secretType {
		case "gcp_sa":
			logrus.Debugf("Generating new key for %s", operation.Source.SecretName)
			secret, err = secretrotator.GetGCPServiceAccountJSONKey(operation.Source.SecretName, operation.Source.Path, operation.Source.Options.NbMaxConcurrent)
			if err != nil {
				logrus.Errorf("%v", err)
			}
		}

		// for dest
		for _, dest := range operation.Dest {
			if dest.Options.Base64 {
				secret = base64.StdEncoding.EncodeToString([]byte(secret))
			}
			switch secretType := dest.Type; secretType {
			case "gitlab_project":
				logrus.Debugf("Writing key to %s", dest.SecretName)
				provider := configuration.Providers[dest.ID]
				gitlabClient, err := gitlab.NewClient(provider.ApiToken, gitlab.WithBaseURL(provider.RepoUrl+"/api/v4"))
				if err != nil {
					fmt.Printf("Error when creating GitLab client: %v\n", err)
				}

				_, err = secretrotator.WriteGitlabSecret(gitlabClient, dest.Path, dest.SecretName, secret)
				if err != nil {
					fmt.Printf("Error when writing secret in Gitlab variable: %v\n", err)
				}

			}
		}

	}
}
