// Code generated by go-swagger; DO NOT EDIT.

package incidents

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetV1IncidentsIncidentIDEventsParams creates a new GetV1IncidentsIncidentIDEventsParams object
// with the default values initialized.
func NewGetV1IncidentsIncidentIDEventsParams() *GetV1IncidentsIncidentIDEventsParams {
	var ()
	return &GetV1IncidentsIncidentIDEventsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetV1IncidentsIncidentIDEventsParamsWithTimeout creates a new GetV1IncidentsIncidentIDEventsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetV1IncidentsIncidentIDEventsParamsWithTimeout(timeout time.Duration) *GetV1IncidentsIncidentIDEventsParams {
	var ()
	return &GetV1IncidentsIncidentIDEventsParams{

		timeout: timeout,
	}
}

// NewGetV1IncidentsIncidentIDEventsParamsWithContext creates a new GetV1IncidentsIncidentIDEventsParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetV1IncidentsIncidentIDEventsParamsWithContext(ctx context.Context) *GetV1IncidentsIncidentIDEventsParams {
	var ()
	return &GetV1IncidentsIncidentIDEventsParams{

		Context: ctx,
	}
}

// NewGetV1IncidentsIncidentIDEventsParamsWithHTTPClient creates a new GetV1IncidentsIncidentIDEventsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetV1IncidentsIncidentIDEventsParamsWithHTTPClient(client *http.Client) *GetV1IncidentsIncidentIDEventsParams {
	var ()
	return &GetV1IncidentsIncidentIDEventsParams{
		HTTPClient: client,
	}
}

/*GetV1IncidentsIncidentIDEventsParams contains all the parameters to send to the API endpoint
for the get v1 incidents incident Id events operation typically these are written to a http.Request
*/
type GetV1IncidentsIncidentIDEventsParams struct {

	/*IncidentID*/
	IncidentID string
	/*Page*/
	Page *int32
	/*PerPage*/
	PerPage *int32

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get v1 incidents incident Id events params
func (o *GetV1IncidentsIncidentIDEventsParams) WithTimeout(timeout time.Duration) *GetV1IncidentsIncidentIDEventsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get v1 incidents incident Id events params
func (o *GetV1IncidentsIncidentIDEventsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get v1 incidents incident Id events params
func (o *GetV1IncidentsIncidentIDEventsParams) WithContext(ctx context.Context) *GetV1IncidentsIncidentIDEventsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get v1 incidents incident Id events params
func (o *GetV1IncidentsIncidentIDEventsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get v1 incidents incident Id events params
func (o *GetV1IncidentsIncidentIDEventsParams) WithHTTPClient(client *http.Client) *GetV1IncidentsIncidentIDEventsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get v1 incidents incident Id events params
func (o *GetV1IncidentsIncidentIDEventsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithIncidentID adds the incidentID to the get v1 incidents incident Id events params
func (o *GetV1IncidentsIncidentIDEventsParams) WithIncidentID(incidentID string) *GetV1IncidentsIncidentIDEventsParams {
	o.SetIncidentID(incidentID)
	return o
}

// SetIncidentID adds the incidentId to the get v1 incidents incident Id events params
func (o *GetV1IncidentsIncidentIDEventsParams) SetIncidentID(incidentID string) {
	o.IncidentID = incidentID
}

// WithPage adds the page to the get v1 incidents incident Id events params
func (o *GetV1IncidentsIncidentIDEventsParams) WithPage(page *int32) *GetV1IncidentsIncidentIDEventsParams {
	o.SetPage(page)
	return o
}

// SetPage adds the page to the get v1 incidents incident Id events params
func (o *GetV1IncidentsIncidentIDEventsParams) SetPage(page *int32) {
	o.Page = page
}

// WithPerPage adds the perPage to the get v1 incidents incident Id events params
func (o *GetV1IncidentsIncidentIDEventsParams) WithPerPage(perPage *int32) *GetV1IncidentsIncidentIDEventsParams {
	o.SetPerPage(perPage)
	return o
}

// SetPerPage adds the perPage to the get v1 incidents incident Id events params
func (o *GetV1IncidentsIncidentIDEventsParams) SetPerPage(perPage *int32) {
	o.PerPage = perPage
}

// WriteToRequest writes these params to a swagger request
func (o *GetV1IncidentsIncidentIDEventsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param incident_id
	if err := r.SetPathParam("incident_id", o.IncidentID); err != nil {
		return err
	}

	if o.Page != nil {

		// query param page
		var qrPage int32
		if o.Page != nil {
			qrPage = *o.Page
		}
		qPage := swag.FormatInt32(qrPage)
		if qPage != "" {
			if err := r.SetQueryParam("page", qPage); err != nil {
				return err
			}
		}

	}

	if o.PerPage != nil {

		// query param per_page
		var qrPerPage int32
		if o.PerPage != nil {
			qrPerPage = *o.PerPage
		}
		qPerPage := swag.FormatInt32(qrPerPage)
		if qPerPage != "" {
			if err := r.SetQueryParam("per_page", qPerPage); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}