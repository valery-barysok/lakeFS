// Code generated by go-swagger; DO NOT EDIT.

package objects

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/validate"

	strfmt "github.com/go-openapi/strfmt"
)

// NewDeleteObjectParams creates a new DeleteObjectParams object
// no default values defined in spec.
func NewDeleteObjectParams() DeleteObjectParams {

	return DeleteObjectParams{}
}

// DeleteObjectParams contains all the bound params for the delete object operation
// typically these are obtained from a http.Request
//
// swagger:parameters deleteObject
type DeleteObjectParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: path
	*/
	BranchID string
	/*
	  Required: true
	  In: query
	*/
	Path string
	/*
	  Required: true
	  In: path
	*/
	RepositoryID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewDeleteObjectParams() beforehand.
func (o *DeleteObjectParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	rBranchID, rhkBranchID, _ := route.Params.GetOK("branchId")
	if err := o.bindBranchID(rBranchID, rhkBranchID, route.Formats); err != nil {
		res = append(res, err)
	}

	qPath, qhkPath, _ := qs.GetOK("path")
	if err := o.bindPath(qPath, qhkPath, route.Formats); err != nil {
		res = append(res, err)
	}

	rRepositoryID, rhkRepositoryID, _ := route.Params.GetOK("repositoryId")
	if err := o.bindRepositoryID(rRepositoryID, rhkRepositoryID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindBranchID binds and validates parameter BranchID from path.
func (o *DeleteObjectParams) bindBranchID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.BranchID = raw

	return nil
}

// bindPath binds and validates parameter Path from query.
func (o *DeleteObjectParams) bindPath(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("path", "query")
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false
	if err := validate.RequiredString("path", "query", raw); err != nil {
		return err
	}

	o.Path = raw

	return nil
}

// bindRepositoryID binds and validates parameter RepositoryID from path.
func (o *DeleteObjectParams) bindRepositoryID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.RepositoryID = raw

	return nil
}
