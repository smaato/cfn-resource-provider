package resource

import (
	"crypto/sha1"
	"crypto/tls"
	"fmt"
	"net/url"

	"github.com/pkg/errors"

	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/encoding"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

// convert list of encoding.String to list of string
func listEncodingToList(l []*encoding.String) []*string {
	var r []*string
	for _, s := range l {
		r = append(r, s.Value())
	}
	return r
}

// retrive OIDC provider ARN using sts
func getProviderArn(sess *session.Session, providerName *string) (*string, error) {

	client := sts.New(sess)

	input := &sts.GetCallerIdentityInput{}

	output, err := client.GetCallerIdentity(input)
	if err != nil {
		return nil, fmt.Errorf("could not get caller identity: %v", err)
	}

	arn := fmt.Sprintf("arn:aws:iam::%s:oidc-provider/%s", *output.Account, *providerName)

	return &arn, nil
}

// get thumbprint
func getIssuerCAThumbprint(issuer string) (*string, error) {

	config := &tls.Config{InsecureSkipVerify: true}

	issuerURL, err := url.Parse(issuer)
	if err != nil {
		return nil, fmt.Errorf("unable to parse OIDC issuer's url")
	}

	if issuerURL.Scheme != "https" {
		return nil, fmt.Errorf("unsupported URL scheme %q", issuerURL.Scheme)
	}

	if issuerURL.Port() == "" {
		issuerURL.Host += ":443"
	}

	conn, err := tls.Dial("tcp", issuerURL.Host, config)
	if err != nil {
		return nil, errors.Wrapf(err, "connecting to issuer OIDC (%s)", issuerURL)
	}
	defer conn.Close()

	cs := conn.ConnectionState()
	if numCerts := len(cs.PeerCertificates); numCerts >= 1 {
		root := cs.PeerCertificates[numCerts-1]
		issuerCAThumbprint := fmt.Sprintf("%x", sha1.Sum(root.Raw))
		return &issuerCAThumbprint, nil
	}
	return nil, fmt.Errorf("unable to get OIDC issuer's certificate")
}
