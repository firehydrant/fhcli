// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// PostV1PostMortemsReportsReportIDAffectedComponents Add an component to a post mortem report
// swagger:model postV1PostMortemsReportsReportIdAffectedComponents
type PostV1PostMortemsReportsReportIDAffectedComponents struct {

	// component id
	// Required: true
	ComponentID *string `json:"component_id"`
}

// Validate validates this post v1 post mortems reports report Id affected components
func (m *PostV1PostMortemsReportsReportIDAffectedComponents) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateComponentID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PostV1PostMortemsReportsReportIDAffectedComponents) validateComponentID(formats strfmt.Registry) error {

	if err := validate.Required("component_id", "body", m.ComponentID); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PostV1PostMortemsReportsReportIDAffectedComponents) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PostV1PostMortemsReportsReportIDAffectedComponents) UnmarshalBinary(b []byte) error {
	var res PostV1PostMortemsReportsReportIDAffectedComponents
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
