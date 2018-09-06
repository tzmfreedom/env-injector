package envinjector

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// For lazy initialization
var services *awsServices

type awsServices struct {
	ssm            *ssm.SSM
	secretsManager *secretsmanager.SecretsManager
}

func newAWSServices() *awsServices {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		logger.Fatalf("failed to create a new session.\n %v", err)
	}
	if aws.StringValue(sess.Config.Region) == "" {
		trace("no explicit region configurations. So now retrieving ec2metadata...")
		region, err := ec2metadata.New(sess).Region()
		if err != nil {
			trace(err)
			logger.Fatalf("could not find region configurations")
		}
		sess.Config.Region = aws.String(region)
	}

	if arn := os.Getenv("ENV_INJECTOR_ASSUME_ROLE_ARN"); arn != "" {
		creds := stscreds.NewCredentials(sess, arn)
		return &awsServices{
			ssm:            ssm.New(sess, &aws.Config{Credentials: creds}),
			secretsManager: secretsmanager.New(sess, &aws.Config{Credentials: creds}),
		}
	}
	return &awsServices{
		ssm:            ssm.New(sess),
		secretsManager: secretsmanager.New(sess),
	}
}

// Initialize services lazily, and return it.
func getService() *awsServices {
	if services == nil {
		services = newAWSServices()
	}
	return services
}
