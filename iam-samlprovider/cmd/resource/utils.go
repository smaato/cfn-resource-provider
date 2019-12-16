package resource

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

// retrive SAML provider ARN using sts
func getProviderArn(sess *session.Session, providerName *string) (*string, error) {

	client := sts.New(sess)

	input := &sts.GetCallerIdentityInput{}

	output, err := client.GetCallerIdentity(input)
	if err != nil {
		return nil, fmt.Errorf("could not get caller identity: %v", err)
	}

	arn := fmt.Sprintf("arn:aws:iam::%s:saml-provider/%s", *output.Account, *providerName)

	return &arn, nil
}
