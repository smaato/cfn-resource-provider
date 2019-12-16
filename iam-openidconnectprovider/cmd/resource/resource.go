package resource

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/encoding"
	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
)

// Create handles the Create event from the Cloudformation service.
func Create(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {

	failedResponse := handler.ProgressEvent{
		OperationStatus: handler.Failed,
		Message:         "Create failed",
		ResourceModel:   currentModel,
	}

	issuerUrl := currentModel.Url.Value()

	createOIDCReq := &iam.CreateOpenIDConnectProviderInput{
		Url:          issuerUrl,
		ClientIDList: listEncodingToList(currentModel.ClientIDList),
	}

	// get Thumbprint if not provided
	if currentModel.ThumbprintList == nil {
		thumbprint, err := getIssuerCAThumbprint(*issuerUrl)
		if err != nil {
			return failedResponse, fmt.Errorf("could not retrive thumbprint: %v", err)
		}
		createOIDCReq.ThumbprintList = []*string{thumbprint}
	} else {
		createOIDCReq.ThumbprintList = listEncodingToList(currentModel.ThumbprintList)
	}

	// set primary identifier (physical ID)
	currentModel.Name = encoding.NewString(strings.TrimPrefix(*issuerUrl, "https://"))

	response := handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "Create complete",
		ResourceModel:   currentModel,
	}

	client := iam.New(req.Session)

	_, err := client.CreateOpenIDConnectProvider(createOIDCReq)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case iam.ErrCodeEntityAlreadyExistsException:
				return response, nil
			default:
				return failedResponse, fmt.Errorf("error creating OpenID Connect provider: %v", err)
			}
		}
	}

	return response, nil
}

// Read handles the Read event from the Cloudformation service.
func Read(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {

	// retrive provider ARN
	arn, err := getProviderArn(req.Session, currentModel.Name.Value())
	if err != nil {
		return handler.ProgressEvent{
			OperationStatus: handler.Failed,
			Message:         "Read failed",
			ResourceModel:   currentModel,
		}, fmt.Errorf("error retriving OpenID connect provider ARN: %v", err)
	}

	// set read-only property
	currentModel.Arn = encoding.NewString(*arn)

	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "Read complete",
		ResourceModel:   currentModel,
	}, nil
}

// Update handles the Update event from the Cloudformation service.
func Update(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {

	failedResponse := handler.ProgressEvent{
		OperationStatus: handler.Failed,
		Message:         "Update failed",
		ResourceModel:   currentModel,
	}

	// retrive provider ARN
	arn, err := getProviderArn(req.Session, currentModel.Name.Value())
	if err != nil {
		return failedResponse, fmt.Errorf("error retriving OpenID connect provider ARN: %v", err)
	}

	client := iam.New(req.Session)

	// updating OIDC is not straightforward, can't updated in batch
	// update Thumbprints if changed
	if !reflect.DeepEqual(prevModel.ThumbprintList, currentModel.ThumbprintList) {
		updateThumbprintReq := &iam.UpdateOpenIDConnectProviderThumbprintInput{
			OpenIDConnectProviderArn: arn,
			ThumbprintList:           listEncodingToList(currentModel.ThumbprintList),
		}
		_, err := client.UpdateOpenIDConnectProviderThumbprint(updateThumbprintReq)
		if err != nil {
			return failedResponse, fmt.Errorf("error updating thumbprint of OpenID connect provider: %v", err)
		}
	}

	// update client IDs if changed
	// AWS allows only either append or remove of client ID to current IDs
	// so we remove the old IDs and then add the new IDs
	if !reflect.DeepEqual(prevModel.ClientIDList, currentModel.ClientIDList) {
		// remove old IDs
		for _, c := range prevModel.ClientIDList {
			removeClientReq := &iam.RemoveClientIDFromOpenIDConnectProviderInput{
				OpenIDConnectProviderArn: arn,
				ClientID:                 c.Value(),
			}
			_, err := client.RemoveClientIDFromOpenIDConnectProvider(removeClientReq)
			if err != nil {
				return failedResponse, fmt.Errorf("error removing client ID from OpenID connect provider: %v", err)
			}
		}
		// add new IDs
		for _, c := range currentModel.ClientIDList {
			addClientReq := &iam.AddClientIDToOpenIDConnectProviderInput{
				OpenIDConnectProviderArn: arn,
				ClientID:                 c.Value(),
			}
			_, err := client.AddClientIDToOpenIDConnectProvider(addClientReq)
			if err != nil {
				return failedResponse, fmt.Errorf("error adding client ID to OpenID connect provider: %v", err)
			}
		}
	}

	response := handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "Update complete",
		ResourceModel:   currentModel,
	}

	return response, nil
}

// Delete handles the Delete event from the Cloudformation service.
func Delete(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {

	failedResponse := handler.ProgressEvent{
		OperationStatus: handler.Failed,
		Message:         "Delete failed",
		ResourceModel:   currentModel,
	}

	// retrive provider ARN
	arn, err := getProviderArn(req.Session, currentModel.Name.Value())
	if err != nil {
		return failedResponse, fmt.Errorf("error retriving OpenID connect provider ARN: %v", err)
	}

	deleteOIDCReq := &iam.DeleteOpenIDConnectProviderInput{
		OpenIDConnectProviderArn: arn,
	}

	client := iam.New(req.Session)

	_, err = client.DeleteOpenIDConnectProvider(deleteOIDCReq)
	if err != nil {
		return failedResponse, fmt.Errorf("error deleting OpenID connect provider: %v", err)
	}

	response := handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "Delete complete",
		ResourceModel:   currentModel,
	}

	return response, nil
}

// List handles the List event from the Cloudformation service.
func List(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "List complete",
		ResourceModel:   currentModel,
	}, nil
}
