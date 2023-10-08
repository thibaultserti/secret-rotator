package test

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"secret-rotator/pkg/secretrotator"
	"testing"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

func TestWriteGitlabSecret(t *testing.T) {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	GITLAB_TOKEN := os.Getenv("GITLAB_TOKEN")
	URL := "https://gitlab.com"
	SECRET_PATH := "thibaultserti/secretrotator-test"

	gitlabClient, err := gitlab.NewClient(GITLAB_TOKEN, gitlab.WithBaseURL(URL+"/api/v4"))
	if err != nil {
		fmt.Printf("Error when creating GitLab client: %v\n", err)
	}

	PREFIX_KEY := "SECRETROTATOR_TEST_"
	PREFIX_VALUE := "TEST_"

	t.Run("Write new secret", func(t *testing.T) {
		SECRET_KEY := PREFIX_KEY + fmt.Sprint(rand.Intn(1000))
		SECRET_VALUE := PREFIX_VALUE
		_, err := secretrotator.WriteGitlabSecret(gitlabClient, SECRET_PATH, SECRET_KEY, SECRET_VALUE)
		if err != nil {
			logrus.Fatalf("%v", err)
		}

		variable, _, err := gitlabClient.ProjectVariables.GetVariable(SECRET_PATH, SECRET_KEY, &gitlab.GetProjectVariableOptions{})
		if err != nil {
			logrus.Fatalf("%v", err)
		}

		_, err = gitlabClient.ProjectVariables.RemoveVariable(SECRET_PATH, SECRET_KEY, &gitlab.RemoveProjectVariableOptions{})
		if err != nil {
			logrus.Fatalf("%v", err)
		}

		got := variable.Value
		want := SECRET_VALUE

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("Update secret", func(t *testing.T) {
		SECRET_KEY := PREFIX_KEY
		SECRET_VALUE := PREFIX_VALUE + fmt.Sprint(rand.Intn(1000))
		_, err := secretrotator.WriteGitlabSecret(gitlabClient, SECRET_PATH, SECRET_KEY, SECRET_VALUE)
		if err != nil {
			logrus.Fatalf("%v", err)
		}

		variable, _, err := gitlabClient.ProjectVariables.GetVariable(SECRET_PATH, SECRET_KEY, &gitlab.GetProjectVariableOptions{})
		if err != nil {
			logrus.Fatalf("%v", err)
		}

		got := variable.Value
		want := SECRET_VALUE

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

}
