package secretrotator

import (
	"github.com/sirupsen/logrus"
	gitlab "github.com/xanzy/go-gitlab"
)

func WriteGitlabSecret(gitlabClient *gitlab.Client, secretPath string, secretKey string, secretValue string) (bool, error) {

	_, resp, err := gitlabClient.ProjectVariables.GetVariable(secretPath, secretKey, &gitlab.GetProjectVariableOptions{})
	if err != nil {
		if resp.StatusCode == 404 {
			createVariableOptions := &gitlab.CreateProjectVariableOptions{
				Key:          gitlab.String(secretKey),
				Value:        gitlab.String(secretValue),
				Raw:          gitlab.Bool(true),
				VariableType: gitlab.VariableType("env_var"),
				Protected:    gitlab.Bool(true),
			}
			_, _, err = gitlabClient.ProjectVariables.CreateVariable(secretPath, createVariableOptions)
			if err != nil {
				logrus.Fatal(err)
			}
		} else {
			logrus.Fatal(err)
		}
	} else {
		opts := &gitlab.UpdateProjectVariableOptions{
			Value: gitlab.String(secretValue),
		}

		_, _, err = gitlabClient.ProjectVariables.UpdateVariable(secretPath, secretKey, opts)
		if err != nil {
			logrus.Errorf("Error when updating gitlab variable: %v\n", err)
			return false, err
		}
	}

	return true, nil
}
