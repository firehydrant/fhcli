// Code generated by go-swagger; DO NOT EDIT.

package changes

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/firehydrant/api-client-go/models"
)

// NewPatchV1ChangesChangeIDParams creates a new PatchV1ChangesChangeIDParams object
// with the default values initialized.
func NewPatchV1ChangesChangeIDParams() *PatchV1ChangesChangeIDParams {
	var ()
	return &PatchV1ChangesChangeIDParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewPatchV1ChangesChangeIDParamsWithTimeout creates a new PatchV1ChangesChangeIDParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewPatchV1ChangesChangeIDParamsWithTimeout(timeout time.Duration) *PatchV1ChangesChangeIDParams {
	var ()
	return &PatchV1ChangesChangeIDParams{

		timeout: timeout,
	}
}

// NewPatchV1ChangesChangeIDParamsWithContext creates a new PatchV1ChangesChangeIDParams object
// with the default values initialized, and the ability to set a context for a request
func NewPatchV1ChangesChangeIDParamsWithContext(ctx context.Context) *PatchV1ChangesChangeIDParams {
	var ()
	return &PatchV1ChangesChangeIDParams{

		Context: ctx,
	}
}

// NewPatchV1ChangesChangeIDParamsWithHTTPClient creates a new PatchV1ChangesChangeIDParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewPatchV1ChangesChangeIDParamsWithHTTPClient(client *http.Client) *PatchV1ChangesChangeIDParams {
	var ()
	return &PatchV1ChangesChangeIDParams{
		HTTPClient: client,
	}
}

/*PatchV1ChangesChangeIDParams contains all the parameters to send to the API endpoint
for the patch v1 changes change Id operation typically these are written to a http.Request
*/
type PatchV1ChangesChangeIDParams struct {

	/*V1Changes*/
	V1Changes *models.PatchV1Changes
	/*ChangeID*/
	ChangeID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the patch v1 changes change Id params
func (o *PatchV1ChangesChangeIDParams) WithTimeout(timeout time.Duration) *PatchV1ChangesChangeIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the patch v1 changes change Id params
func (o *PatchV1ChangesChangeIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the patch v1 changes change Id params
func (o *PatchV1ChangesChangeIDParams) WithContext(ctx context.Context) *PatchV1ChangesChangeIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the patch v1 changes change Id params
func (o *PatchV1ChangesChangeIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the patch v1 changes change Id params
func (o *PatchV1ChangesChangeIDParams) WithHTTPClient(client *http.Client) *PatchV1ChangesChangeIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the patch v1 changes change Id params
func (o *PatchV1ChangesChangeIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithV1Changes adds the v1Changes to the patch v1 changes change Id params
func (o *PatchV1ChangesChangeIDParams) WithV1Changes(v1Changes *models.PatchV1Changes) *PatchV1ChangesChangeIDParams {
	o.SetV1Changes(v1Changes)
	return o
}

// SetV1Changes adds the v1Changes to the patch v1 changes change Id params
func (o *PatchV1ChangesChangeIDParams) SetV1Changes(v1Changes *models.PatchV1Changes) {
	o.V1Changes = v1Changes
}

// WithChangeID adds the changeID to the patch v1 changes change Id params
func (o *PatchV1ChangesChangeIDParams) WithChangeID(changeID string) *PatchV1ChangesChangeIDParams {
	o.SetChangeID(changeID)
	return o
}

// SetChangeID adds the changeId to the patch v1 changes change Id params
func (o *PatchV1ChangesChangeIDParams) SetChangeID(changeID string) {
	o.ChangeID = changeID
}

// WriteToRequest writes these params to a swagger request
func (o *PatchV1ChangesChangeIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.V1Changes != nil {
		if err := r.SetBodyParam(o.V1Changes); err != nil {
			return err
		}
	}

	// path param change_id
	if err := r.SetPathParam("change_id", o.ChangeID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}