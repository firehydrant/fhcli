// Code generated by go-swagger; DO NOT EDIT.

package changes

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/firehydrant/api-client-go/models"
)

// PatchV1ChangesChangeIDReader is a Reader for the PatchV1ChangesChangeID structure.
type PatchV1ChangesChangeIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PatchV1ChangesChangeIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewPatchV1ChangesChangeIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewPatchV1ChangesChangeIDOK creates a PatchV1ChangesChangeIDOK with default headers values
func NewPatchV1ChangesChangeIDOK() *PatchV1ChangesChangeIDOK {
	return &PatchV1ChangesChangeIDOK{}
}

/*PatchV1ChangesChangeIDOK handles this case with default header values.

Update a change entry
*/
type PatchV1ChangesChangeIDOK struct {
	Payload *models.ChangeEntity
}

func (o *PatchV1ChangesChangeIDOK) Error() string {
	return fmt.Sprintf("[PATCH /v1/changes/{change_id}][%d] patchV1ChangesChangeIdOK  %+v", 200, o.Payload)
}

func (o *PatchV1ChangesChangeIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ChangeEntity)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}