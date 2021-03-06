// Code generated by go-swagger; DO NOT EDIT.

package incidents

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// DeleteV1IncidentsIncidentIDEventsEventIDStarsReader is a Reader for the DeleteV1IncidentsIncidentIDEventsEventIDStars structure.
type DeleteV1IncidentsIncidentIDEventsEventIDStarsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteV1IncidentsIncidentIDEventsEventIDStarsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 204:
		result := NewDeleteV1IncidentsIncidentIDEventsEventIDStarsNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewDeleteV1IncidentsIncidentIDEventsEventIDStarsNoContent creates a DeleteV1IncidentsIncidentIDEventsEventIDStarsNoContent with default headers values
func NewDeleteV1IncidentsIncidentIDEventsEventIDStarsNoContent() *DeleteV1IncidentsIncidentIDEventsEventIDStarsNoContent {
	return &DeleteV1IncidentsIncidentIDEventsEventIDStarsNoContent{}
}

/*DeleteV1IncidentsIncidentIDEventsEventIDStarsNoContent handles this case with default header values.

Remove a star from an incident event
*/
type DeleteV1IncidentsIncidentIDEventsEventIDStarsNoContent struct {
}

func (o *DeleteV1IncidentsIncidentIDEventsEventIDStarsNoContent) Error() string {
	return fmt.Sprintf("[DELETE /v1/incidents/{incident_id}/events/{event_id}/stars][%d] deleteV1IncidentsIncidentIdEventsEventIdStarsNoContent ", 204)
}

func (o *DeleteV1IncidentsIncidentIDEventsEventIDStarsNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
