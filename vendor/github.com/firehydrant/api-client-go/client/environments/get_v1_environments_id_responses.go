// Code generated by go-swagger; DO NOT EDIT.

package environments

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/firehydrant/api-client-go/models"
)

// GetV1EnvironmentsIDReader is a Reader for the GetV1EnvironmentsID structure.
type GetV1EnvironmentsIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetV1EnvironmentsIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetV1EnvironmentsIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetV1EnvironmentsIDOK creates a GetV1EnvironmentsIDOK with default headers values
func NewGetV1EnvironmentsIDOK() *GetV1EnvironmentsIDOK {
	return &GetV1EnvironmentsIDOK{}
}

/*GetV1EnvironmentsIDOK handles this case with default header values.

Retrieve a single environment
*/
type GetV1EnvironmentsIDOK struct {
	Payload *models.EnvironmentEntity
}

func (o *GetV1EnvironmentsIDOK) Error() string {
	return fmt.Sprintf("[GET /v1/environments/{id}][%d] getV1EnvironmentsIdOK  %+v", 200, o.Payload)
}

func (o *GetV1EnvironmentsIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.EnvironmentEntity)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}