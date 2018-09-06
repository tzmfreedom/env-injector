package envinjector

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func injectEnvironSecretManager(name string, decorator envKeyDecorator) {
	tracef("secret name: %s", name)

	svc := getService().secretsManager
	ret, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(name),
	})
	if err != nil {
		logger.Fatalf("secretsmanager:GetSecretValue failed. (name: %s)\n %v", name, err)
	}
	secrets := make(map[string]string)
	err = json.Unmarshal([]byte(aws.StringValue(ret.SecretString)), &secrets)
	if err != nil {
		logger.Fatalf("secretsmanager:GetSecretValue returns invalid json. (name: %s)\n %v", name, err)
	}
	for key, val := range secrets {
		key = decorator.decorate(key)
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
			tracef("env injected: %s", key)
		}
	}
}
