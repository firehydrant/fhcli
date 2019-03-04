// Code generated by go-swagger; DO NOT EDIT.

package incidents

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/firehydrant/api-client-go/models"
)

// PostV1IncidentsIncidentIDImpactTypeReader is a Reader for the PostV1IncidentsIncidentIDImpactType structure.
type PostV1IncidentsIncidentIDImpactTypeReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostV1IncidentsIncidentIDImpactTypeReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 201:
		result := NewPostV1IncidentsIncidentIDImpactTypeCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 400:
		result := NewPostV1IncidentsIncidentIDImpactTypeBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewPostV1IncidentsIncidentIDImpactTypeCreated creates a PostV1IncidentsIncidentIDImpactTypeCreated with default headers values
func NewPostV1IncidentsIncidentIDImpactTypeCreated() *PostV1IncidentsIncidentIDImpactTypeCreated {
	return &PostV1IncidentsIncidentIDImpactTypeCreated{}
}

/*PostV1IncidentsIncidentIDImpactTypeCreated handles this case with default header values.

Add a piece of infrastructure to an incident as impact
*/
type PostV1IncidentsIncidentIDImpactTypeCreated struct {
	Payload *models.IncidentImpactEntity
}

func (o *PostV1IncidentsIncidentIDImpactTypeCreated) Error() string {
	return fmt.Sprintf("[POST /v1/incidents/{incident_id}/impact/{type}][%d] postV1IncidentsIncidentIdImpactTypeCreated  %+v", 201, o.Payload)
}

func (o *PostV1IncidentsIncidentIDImpactTypeCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.IncidentImpactEntity)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostV1IncidentsIncidentIDImpactTypeBadRequest creates a PostV1IncidentsIncidentIDImpactTypeBadRequest with default headers values
func NewPostV1IncidentsIncidentIDImpactTypeBadRequest() *PostV1IncidentsIncidentIDImpactTypeBadRequest {
	return &PostV1IncidentsIncidentIDImpactTypeBadRequest{}
}

/*PostV1IncidentsIncidentIDImpactTypeBadRequest handles this case with default header values.

Bad Request
*/
type PostV1IncidentsIncidentIDImpactTypeBadRequest struct {
	Payload *models.ErrorEntity
}

func (o *PostV1IncidentsIncidentIDImpactTypeBadRequest) Error() string {
	return fmt.Sprintf("[POST /v1/incidents/{incident_id}/impact/{type}][%d] postV1IncidentsIncidentIdImpactTypeBadRequest  %+v", 400, o.Payload)
}

func (o *PostV1IncidentsIncidentIDImpactTypeBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorEntity)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}