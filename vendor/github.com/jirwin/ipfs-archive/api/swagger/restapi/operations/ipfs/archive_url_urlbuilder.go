// Code generated by go-swagger; DO NOT EDIT.

package ipfs

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"errors"
	"net/url"
	golangswaggerpaths "path"
)

// ArchiveURLURL generates an URL for the archive Url operation
type ArchiveURLURL struct {
	_basePath string
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *ArchiveURLURL) WithBasePath(bp string) *ArchiveURLURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *ArchiveURLURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *ArchiveURLURL) Build() (*url.URL, error) {
	var result url.URL

	var _path = "/archive"

	_basePath := o._basePath
	if _basePath == "" {
		_basePath = "/api"
	}
	result.Path = golangswaggerpaths.Join(_basePath, _path)

	return &result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *ArchiveURLURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *ArchiveURLURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *ArchiveURLURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on ArchiveURLURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on ArchiveURLURL")
	}

	base, err := o.Build()
	if err != nil {
		return nil, err
	}

	base.Scheme = scheme
	base.Host = host
	return base, nil
}

// StringFull returns the string representation of a complete url
func (o *ArchiveURLURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}