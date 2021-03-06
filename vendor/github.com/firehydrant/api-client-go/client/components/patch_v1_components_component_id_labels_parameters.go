// Code generated by go-swagger; DO NOT EDIT.

package components

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
)

// NewPatchV1ComponentsComponentIDLabelsParams creates a new PatchV1ComponentsComponentIDLabelsParams object
// with the default values initialized.
func NewPatchV1ComponentsComponentIDLabelsParams() *PatchV1ComponentsComponentIDLabelsParams {
	var ()
	return &PatchV1ComponentsComponentIDLabelsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewPatchV1ComponentsComponentIDLabelsParamsWithTimeout creates a new PatchV1ComponentsComponentIDLabelsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewPatchV1ComponentsComponentIDLabelsParamsWithTimeout(timeout time.Duration) *PatchV1ComponentsComponentIDLabelsParams {
	var ()
	return &PatchV1ComponentsComponentIDLabelsParams{

		timeout: timeout,
	}
}

// NewPatchV1ComponentsComponentIDLabelsParamsWithContext creates a new PatchV1ComponentsComponentIDLabelsParams object
// with the default values initialized, and the ability to set a context for a request
func NewPatchV1ComponentsComponentIDLabelsParamsWithContext(ctx context.Context) *PatchV1ComponentsComponentIDLabelsParams {
	var ()
	return &PatchV1ComponentsComponentIDLabelsParams{

		Context: ctx,
	}
}

// NewPatchV1ComponentsComponentIDLabelsParamsWithHTTPClient creates a new PatchV1ComponentsComponentIDLabelsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewPatchV1ComponentsComponentIDLabelsParamsWithHTTPClient(client *http.Client) *PatchV1ComponentsComponentIDLabelsParams {
	var ()
	return &PatchV1ComponentsComponentIDLabelsParams{
		HTTPClient: client,
	}
}

/*PatchV1ComponentsComponentIDLabelsParams contains all the parameters to send to the API endpoint
for the patch v1 components component Id labels operation typically these are written to a http.Request
*/
type PatchV1ComponentsComponentIDLabelsParams struct {

	/*ComponentID
	  Component UUID

	*/
	ComponentID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the patch v1 components component Id labels params
func (o *PatchV1ComponentsComponentIDLabelsParams) WithTimeout(timeout time.Duration) *PatchV1ComponentsComponentIDLabelsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the patch v1 components component Id labels params
func (o *PatchV1ComponentsComponentIDLabelsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the patch v1 components component Id labels params
func (o *PatchV1ComponentsComponentIDLabelsParams) WithContext(ctx context.Context) *PatchV1ComponentsComponentIDLabelsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the patch v1 components component Id labels params
func (o *PatchV1ComponentsComponentIDLabelsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the patch v1 components component Id labels params
func (o *PatchV1ComponentsComponentIDLabelsParams) WithHTTPClient(client *http.Client) *PatchV1ComponentsComponentIDLabelsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the patch v1 components component Id labels params
func (o *PatchV1ComponentsComponentIDLabelsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithComponentID adds the componentID to the patch v1 components component Id labels params
func (o *PatchV1ComponentsComponentIDLabelsParams) WithComponentID(componentID string) *PatchV1ComponentsComponentIDLabelsParams {
	o.SetComponentID(componentID)
	return o
}

// SetComponentID adds the componentId to the patch v1 components component Id labels params
func (o *PatchV1ComponentsComponentIDLabelsParams) SetComponentID(componentID string) {
	o.ComponentID = componentID
}

// WriteToRequest writes these params to a swagger request
func (o *PatchV1ComponentsComponentIDLabelsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param component_id
	if err := r.SetPathParam("component_id", o.ComponentID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
