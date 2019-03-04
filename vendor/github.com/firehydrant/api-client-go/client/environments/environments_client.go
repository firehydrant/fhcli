// Code generated by go-swagger; DO NOT EDIT.

package environments

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new environments API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for environments API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
DeleteV1EnvironmentsID Archive an environment
*/
func (a *Client) DeleteV1EnvironmentsID(params *DeleteV1EnvironmentsIDParams, authInfo runtime.ClientAuthInfoWriter) (*DeleteV1EnvironmentsIDNoContent, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteV1EnvironmentsIDParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "deleteV1EnvironmentsId",
		Method:             "DELETE",
		PathPattern:        "/v1/environments/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{""},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &DeleteV1EnvironmentsIDReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*DeleteV1EnvironmentsIDNoContent), nil

}

/*
GetV1Environments Retrieve all environments
*/
func (a *Client) GetV1Environments(params *GetV1EnvironmentsParams, authInfo runtime.ClientAuthInfoWriter) (*GetV1EnvironmentsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetV1EnvironmentsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getV1Environments",
		Method:             "GET",
		PathPattern:        "/v1/environments",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{""},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetV1EnvironmentsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetV1EnvironmentsOK), nil

}

/*
GetV1EnvironmentsID Retrieve a single environment
*/
func (a *Client) GetV1EnvironmentsID(params *GetV1EnvironmentsIDParams, authInfo runtime.ClientAuthInfoWriter) (*GetV1EnvironmentsIDOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetV1EnvironmentsIDParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getV1EnvironmentsId",
		Method:             "GET",
		PathPattern:        "/v1/environments/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{""},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetV1EnvironmentsIDReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetV1EnvironmentsIDOK), nil

}

/*
PatchV1EnvironmentsID Update an environment
*/
func (a *Client) PatchV1EnvironmentsID(params *PatchV1EnvironmentsIDParams, authInfo runtime.ClientAuthInfoWriter) (*PatchV1EnvironmentsIDOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPatchV1EnvironmentsIDParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "patchV1EnvironmentsId",
		Method:             "PATCH",
		PathPattern:        "/v1/environments/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &PatchV1EnvironmentsIDReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*PatchV1EnvironmentsIDOK), nil

}

/*
PostV1Environments Creates an environment
*/
func (a *Client) PostV1Environments(params *PostV1EnvironmentsParams, authInfo runtime.ClientAuthInfoWriter) (*PostV1EnvironmentsCreated, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPostV1EnvironmentsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "postV1Environments",
		Method:             "POST",
		PathPattern:        "/v1/environments",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &PostV1EnvironmentsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*PostV1EnvironmentsCreated), nil

}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}