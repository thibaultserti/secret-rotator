package test

import (
	"encoding/json"
	"log"
	"reflect"
	"secret-rotator/pkg/secretrotator"
	"testing"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func TestGetGCPServiceAccountJSONKey(t *testing.T) {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	t.Run("Create new key and keep only the 2 latest", func(t *testing.T) {

		SERVICE_ACCOUNT := "secretrotator-test@infra-390806.iam.gserviceaccount.com"
		PROJECT_ID := "infra-390806"
		key, _ := secretrotator.GetGCPServiceAccountJSONKey(SERVICE_ACCOUNT, PROJECT_ID, 2)
		expectedKey := `{
			"type": "service_account",
			"project_id": "infra-390806",
			"private_key_id": "",
			"private_key": "",
			"client_email": "secretrotator-test@infra-390806.iam.gserviceaccount.com",
			"client_id": "116668024987070255659",
			"auth_uri": "https://accounts.google.com/o/oauth2/auth",
			"token_uri": "https://oauth2.googleapis.com/token",
			"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
			"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/secretrotator-test%40infra-390806.iam.gserviceaccount.com",
			"universe_domain": "googleapis.com"
		  }`

		var expectedKeyJSON, keyJSON map[string]interface{}
		err = json.Unmarshal([]byte(expectedKey), &expectedKeyJSON)
		if err != nil {
			logrus.Fatalf("%v", err)
		}
		err = json.Unmarshal([]byte(key), &keyJSON)
		if err != nil {
			logrus.Fatalf("%v", err)
		}

		isEqual := reflect.DeepEqual(IgnoreFields(expectedKeyJSON, "private_key", "private_key_id"), IgnoreFields(keyJSON, "private_key", "private_key_id"))

		got := keyJSON
		want := expectedKeyJSON

		if !isEqual {
			t.Errorf("got %v want %v", got, want)
		}
	})

}
