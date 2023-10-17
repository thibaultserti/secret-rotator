package config

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mcuadros/go-defaults"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ProviderType string

const (
	Gitlab ProviderType = "gitlab"
	GCP    ProviderType = "gcp"
)

type SecretType string

const (
	GitlabProjectVar     SecretType = "gitlab_project"
	GCPServiceAccountKey SecretType = "gcp_sa"
)

type Configuration struct {
	Env       string              `mapstructure:"env" default:"prod"`
	LogLevel  string              `mapstructure:"logLevel" default:"info" validate:"oneof=trace debug info warn error fatal panic"`
	Providers map[string]Provider `mapstructure:"providers" validate:"required,dive"`
	Secrets   []Secret            `mapstructure:"secrets" validate:"required,dive"`
}

type Provider struct {
	Type     ProviderType `mapstructure:"type" validate:"required,oneof=gitlab gcp"`
	RepoUrl  string       `mapstructure:"repoUrl" default:""`
	ApiToken string       `mapstructure:"apiToken" default:""`
}

type Secret struct {
	Name   string       `mapstructure:"name" validate:"required"`
	Source SecretSource `mapstructure:"source" validate:"required"`
	Dest   []SecretDest `mapstructure:"dest" validate:"required,dive"`
}

type SecretSource struct {
	ID         string        `mapstructure:"id" validate:"required"`
	Type       SecretType    `mapstructure:"type" validate:"required,oneof=gitlab_project gcp_sa"`
	Path       string        `mapstructure:"path"`
	SecretName string        `mapstructure:"secretName"`
	Options    SecretOptions `mapstructure:"options"`
}

type SecretDest struct {
	ID         string        `mapstructure:"id" validate:"required"`
	Type       SecretType    `mapstructure:"type" validate:"required,oneof=gitlab_project gcp_sa"`
	Path       string        `mapstructure:"path"`
	SecretName string        `mapstructure:"secretName"`
	Options    SecretOptions `mapstructure:"options"`
}

type SecretOptions struct {
	NbMaxConcurrent int  `mapstructure:"nbMaxConcurrent" default:"2"`
	Base64          bool `mapstructure:"base64" default:"false"`
}

func LoadConfig(config_file string) (config Configuration, err error) {
	viper := viper.New()

	viper.SetConfigFile(config_file)
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("SECRETROTATOR")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	config = Configuration{}
	err = viper.Unmarshal(&config)
	if err != nil {
		logrus.Fatalf("%v", err)
	}
	defaults.SetDefaults(&config)

	return
}

func ValidateConfig(config Configuration) bool {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(config)
	if err != nil {
		logrus.Fatalf("%v", err)
	}

	providerIds := make(map[string]bool)

	for providerId := range config.Providers {
		providerIds[providerId] = false
	}

	// validate source id and dest id related to provide ids
	for _, operation := range config.Secrets {
		_, ok := providerIds[operation.Source.ID]
		if ok {
			providerIds[operation.Source.ID] = true
		} else {
			logrus.Errorf("Invalid configuration, provider ID %s not defined", operation.Source.ID)
			return false
		}
		for _, dest := range operation.Dest {
			_, ok := providerIds[dest.ID]
			if ok {
				providerIds[dest.ID] = true
			} else {
				logrus.Errorf("Invalid configuration, provider ID %s not defined", operation.Source.ID)
				return false
			}
		}
	}

	// check unused provider
	for providerId, used := range providerIds {
		if !used {
			logrus.Warnf("Provider Id %s defined but not used", providerId)
		}
	}

	return true
}
