package secretrotator

import (
	"context"
	"encoding/base64"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
)

func GetGCPServiceAccountJSONKey(serviceAccountName string, projectID string, maxNbConcurrent int) (string, error) {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, iam.CloudPlatformScope)
	if err != nil {
		logrus.Errorf("Error when creating GCP client: %v\n", err)
		return "", err
	}

	iamService, err := iam.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		logrus.Errorf("Error when creating IAM service: %v\n", err)
		return "", err
	}

	keyRequest := &iam.CreateServiceAccountKeyRequest{}
	key, err := iamService.Projects.ServiceAccounts.Keys.Create("projects/"+projectID+"/serviceAccounts/"+serviceAccountName, keyRequest).Do()
	if err != nil {
		logrus.Errorf("Error when creating key: %v\n", err)
		return "", err
	}

	keys, err := iamService.Projects.ServiceAccounts.Keys.List("projects/" + projectID + "/serviceAccounts/" + serviceAccountName).Do()
	if err != nil {
		logrus.Errorf("Error when getting existing keys: %v\n", err)
		return "", err
	}
	nbKeys := len(keys.Keys) - 1 // there is always one system managed key not visible in the UI

	sort.Slice(keys.Keys, func(i, j int) bool {
		timeI, _ := time.Parse(time.RFC3339Nano, keys.Keys[i].ValidAfterTime)
		timeJ, _ := time.Parse(time.RFC3339Nano, keys.Keys[j].ValidAfterTime)
		return timeI.Before(timeJ)
	})

	if nbKeys > maxNbConcurrent {
		// i=1 to ignore system managed key
		for i := 1; i < nbKeys-maxNbConcurrent+1; i++ {
			keyToDelete := keys.Keys[i]
			_, err := iamService.Projects.ServiceAccounts.Keys.Delete(keyToDelete.Name).Do()
			if err != nil {
				logrus.Errorf("Error when deleting key: %v\n", err)
				return "", err
			}
		}
	}

	keyJSON, err := base64.StdEncoding.DecodeString(key.PrivateKeyData)
	if err != nil {
		logrus.Errorf("Error when deconding key base64 private key: %v\n", err)
		return "", err
	}
	return string(keyJSON), nil

}
