package resource

import (
	"fmt"

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

	response := handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "Create complete",
		ResourceModel:   currentModel,
	}

	createSAMLReq := &iam.CreateSAMLProviderInput{
		Name:                 currentModel.Name.Value(),
		SAMLMetadataDocument: currentModel.SAMLMetadataDocument.Value(),
	}

	client := iam.New(req.Session)

	_, err := client.CreateSAMLProvider(createSAMLReq)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case iam.ErrCodeEntityAlreadyExistsException:
				return response, nil
			default:
				return failedResponse, fmt.Errorf("error creating SAML provider: %v", err)
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
		}, fmt.Errorf("error retriving SAML provider ARN: %v", err)
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
		return failedResponse, fmt.Errorf("error retriving SAML provider ARN: %v", err)
	}

	updateSAMLReq := &iam.UpdateSAMLProviderInput{
		SAMLProviderArn:      arn,
		SAMLMetadataDocument: currentModel.SAMLMetadataDocument.Value(),
	}

	client := iam.New(req.Session)

	_, err = client.UpdateSAMLProvider(updateSAMLReq)
	if err != nil {
		return failedResponse, fmt.Errorf("error updating SAML provider: %v", err)
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
		return failedResponse, fmt.Errorf("error retriving SAML provider ARN: %v", err)
	}

	deleteSAMLReq := &iam.DeleteSAMLProviderInput{
		SAMLProviderArn: arn,
	}

	client := iam.New(req.Session)

	_, err = client.DeleteSAMLProvider(deleteSAMLReq)
	if err != nil {
		return failedResponse, fmt.Errorf("error deleting SAML provider: %v", err)
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
