// Code generated by go-swagger; DO NOT EDIT.

package environments

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// DeleteV1EnvironmentsIDReader is a Reader for the DeleteV1EnvironmentsID structure.
type DeleteV1EnvironmentsIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteV1EnvironmentsIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 204:
		result := NewDeleteV1EnvironmentsIDNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewDeleteV1EnvironmentsIDNoContent creates a DeleteV1EnvironmentsIDNoContent with default headers values
func NewDeleteV1EnvironmentsIDNoContent() *DeleteV1EnvironmentsIDNoContent {
	return &DeleteV1EnvironmentsIDNoContent{}
}

/*DeleteV1EnvironmentsIDNoContent handles this case with default header values.

Archive an environment
*/
type DeleteV1EnvironmentsIDNoContent struct {
}

func (o *DeleteV1EnvironmentsIDNoContent) Error() string {
	return fmt.Sprintf("[DELETE /v1/environments/{id}][%d] deleteV1EnvironmentsIdNoContent ", 204)
}

func (o *DeleteV1EnvironmentsIDNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
