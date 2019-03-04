// Code generated by go-swagger; DO NOT EDIT.

package post_mortems

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/firehydrant/api-client-go/models"
)

// PatchV1PostMortemsReportsReportIDEventsNotesNoteIDReader is a Reader for the PatchV1PostMortemsReportsReportIDEventsNotesNoteID structure.
type PatchV1PostMortemsReportsReportIDEventsNotesNoteIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PatchV1PostMortemsReportsReportIDEventsNotesNoteIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewPatchV1PostMortemsReportsReportIDEventsNotesNoteIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewPatchV1PostMortemsReportsReportIDEventsNotesNoteIDOK creates a PatchV1PostMortemsReportsReportIDEventsNotesNoteIDOK with default headers values
func NewPatchV1PostMortemsReportsReportIDEventsNotesNoteIDOK() *PatchV1PostMortemsReportsReportIDEventsNotesNoteIDOK {
	return &PatchV1PostMortemsReportsReportIDEventsNotesNoteIDOK{}
}

/*PatchV1PostMortemsReportsReportIDEventsNotesNoteIDOK handles this case with default header values.

Update a post mortem note
*/
type PatchV1PostMortemsReportsReportIDEventsNotesNoteIDOK struct {
	Payload *models.EventEntity
}

func (o *PatchV1PostMortemsReportsReportIDEventsNotesNoteIDOK) Error() string {
	return fmt.Sprintf("[PATCH /v1/post_mortems/reports/{report_id}/events/notes/{note_id}][%d] patchV1PostMortemsReportsReportIdEventsNotesNoteIdOK  %+v", 200, o.Payload)
}

func (o *PatchV1PostMortemsReportsReportIDEventsNotesNoteIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.EventEntity)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}