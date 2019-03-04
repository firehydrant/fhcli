// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// PostV1Teams Create a team
// swagger:model postV1Teams
type PostV1Teams struct {

	// description
	Description string `json:"description,omitempty"`

	// memberships
	Memberships []*PostV1TeamsMembershipsItems0 `json:"memberships"`

	// name
	// Required: true
	Name *string `json:"name"`
}

// Validate validates this post v1 teams
func (m *PostV1Teams) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateMemberships(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PostV1Teams) validateMemberships(formats strfmt.Registry) error {

	if swag.IsZero(m.Memberships) { // not required
		return nil
	}

	for i := 0; i < len(m.Memberships); i++ {
		if swag.IsZero(m.Memberships[i]) { // not required
			continue
		}

		if m.Memberships[i] != nil {
			if err := m.Memberships[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("memberships" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *PostV1Teams) validateName(formats strfmt.Registry) error {

	if err := validate.Required("name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PostV1Teams) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PostV1Teams) UnmarshalBinary(b []byte) error {
	var res PostV1Teams
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// PostV1TeamsMembershipsItems0 post v1 teams memberships items0
// swagger:model PostV1TeamsMembershipsItems0
type PostV1TeamsMembershipsItems0 struct {

	// An incident role ID that this user will automatically assigned if this team is assigned to an incident
	IncidentRoleID string `json:"incident_role_id,omitempty"`

	// user id
	// Required: true
	UserID *string `json:"user_id"`
}

// Validate validates this post v1 teams memberships items0
func (m *PostV1TeamsMembershipsItems0) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateUserID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PostV1TeamsMembershipsItems0) validateUserID(formats strfmt.Registry) error {

	if err := validate.Required("user_id", "body", m.UserID); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PostV1TeamsMembershipsItems0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PostV1TeamsMembershipsItems0) UnmarshalBinary(b []byte) error {
	var res PostV1TeamsMembershipsItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}