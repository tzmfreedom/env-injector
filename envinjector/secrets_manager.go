package envinjector

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"gopkg.in/yaml.v2"
	"os"
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
	converts := make(map[string]string)
	if path := os.Getenv("ENV_INJECTOR_CONVERT_YAML"); path != "" {
		fmt.Printf("hogehoge")
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			logger.Fatalf("secretsmanager:GetSecretValue returns invalid json. (name: %s)\n %v", name, err)
		}
		err = yaml.Unmarshal(buf, &converts)
		if err != nil {
			logger.Fatalf("secretsmanager:GetSecretValue returns invalid json. (name: %s)\n %v", name, err)
		}
		fmt.Printf("%v", converts)
	}
	for key, val := range secrets {
		key = decorator.decorate(key)
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
			tracef("env injected: %s", key)
		}
	}
}
