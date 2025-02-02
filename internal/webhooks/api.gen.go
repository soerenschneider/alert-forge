//go:build go1.22

// Package webhooks provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package webhooks

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/runtime"
)

// Defines values for GetAllAlertsParamsContentType.
const (
	GetAllAlertsParamsContentTypeApplicationjson GetAllAlertsParamsContentType = "application/json"
	GetAllAlertsParamsContentTypeTexthtml        GetAllAlertsParamsContentType = "text/html"
)

// Defines values for GetAlertsByInstanceParamsContentType.
const (
	GetAlertsByInstanceParamsContentTypeApplicationjson GetAlertsByInstanceParamsContentType = "application/json"
	GetAlertsByInstanceParamsContentTypeTexthtml        GetAlertsByInstanceParamsContentType = "text/html"
)

// Defines values for GetAlertsBySeverityParamsContentType.
const (
	GetAlertsBySeverityParamsContentTypeApplicationjson GetAlertsBySeverityParamsContentType = "application/json"
	GetAlertsBySeverityParamsContentTypeTexthtml        GetAlertsBySeverityParamsContentType = "text/html"
)

// Defines values for GetAlertsTodayParamsContentType.
const (
	GetAlertsTodayParamsContentTypeApplicationjson GetAlertsTodayParamsContentType = "application/json"
	GetAlertsTodayParamsContentTypeTexthtml        GetAlertsTodayParamsContentType = "text/html"
)

// Defines values for GetAlertsYesterdayParamsContentType.
const (
	GetAlertsYesterdayParamsContentTypeApplicationjson GetAlertsYesterdayParamsContentType = "application/json"
	GetAlertsYesterdayParamsContentTypeTexthtml        GetAlertsYesterdayParamsContentType = "text/html"
)

// Defines values for StatisticsParamsAccept.
const (
	StatisticsParamsAcceptApplicationjson StatisticsParamsAccept = "application/json"
	StatisticsParamsAcceptTexthtml        StatisticsParamsAccept = "text/html"
)

// Alert defines model for Alert.
type Alert struct {
	Annotations  map[string]string `json:"annotations,omitempty"`
	EndsAt       time.Time         `json:"endsAt,omitempty"`
	Fingerprint  string            `json:"fingerprint,omitempty"`
	GeneratorURL string            `json:"generatorURL,omitempty"`
	Labels       map[string]string `json:"labels,omitempty"`
	Receivers    []struct {
		Name string `json:"name,omitempty"`
	} `json:"receivers,omitempty"`
	StartsAt  time.Time `json:"startsAt,omitempty"`
	Status    Status    `json:"status,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

// Status defines model for Status.
type Status struct {
	InhibitedBy []string `json:"inhibitedBy,omitempty"`
	SilencedBy  []string `json:"silencedBy,omitempty"`
	State       string   `json:"state,omitempty"`
}

// GetAllAlertsParams defines parameters for GetAllAlerts.
type GetAllAlertsParams struct {
	// ContentType Content type of the request (e.g., application/json or text/html)
	ContentType *GetAllAlertsParamsContentType `json:"Content-Type,omitempty"`
}

// GetAllAlertsParamsContentType defines parameters for GetAllAlerts.
type GetAllAlertsParamsContentType string

// GetAlertsByInstanceParams defines parameters for GetAlertsByInstance.
type GetAlertsByInstanceParams struct {
	// ContentType Content type of the request (e.g., application/json or text/html)
	ContentType *GetAlertsByInstanceParamsContentType `json:"Content-Type,omitempty"`
}

// GetAlertsByInstanceParamsContentType defines parameters for GetAlertsByInstance.
type GetAlertsByInstanceParamsContentType string

// GetAlertsBySeverityParams defines parameters for GetAlertsBySeverity.
type GetAlertsBySeverityParams struct {
	// ContentType Content type of the request (e.g., application/json or text/html)
	ContentType *GetAlertsBySeverityParamsContentType `json:"Content-Type,omitempty"`
}

// GetAlertsBySeverityParamsContentType defines parameters for GetAlertsBySeverity.
type GetAlertsBySeverityParamsContentType string

// GetAlertsTodayParams defines parameters for GetAlertsToday.
type GetAlertsTodayParams struct {
	// ContentType Content type of the request (e.g., application/json or text/html)
	ContentType *GetAlertsTodayParamsContentType `json:"Content-Type,omitempty"`
}

// GetAlertsTodayParamsContentType defines parameters for GetAlertsToday.
type GetAlertsTodayParamsContentType string

// GetAlertsYesterdayParams defines parameters for GetAlertsYesterday.
type GetAlertsYesterdayParams struct {
	// ContentType Content type of the request (e.g., application/json or text/html)
	ContentType *GetAlertsYesterdayParamsContentType `json:"Content-Type,omitempty"`
}

// GetAlertsYesterdayParamsContentType defines parameters for GetAlertsYesterday.
type GetAlertsYesterdayParamsContentType string

// StatisticsParams defines parameters for Statistics.
type StatisticsParams struct {
	// Accept Specify `application/json` or `text/html` to request JSON or HTML format.
	Accept *StatisticsParamsAccept `json:"Accept,omitempty"`
}

// StatisticsParamsAccept defines parameters for Statistics.
type StatisticsParamsAccept string

// CreateAlertJSONBody defines parameters for CreateAlert.
type CreateAlertJSONBody struct {
	Alerts []struct {
		Annotations  map[string]string `json:"annotations,omitempty"`
		EndsAt       time.Time         `json:"endsAt,omitempty"`
		Fingerprint  string            `json:"fingerprint,omitempty"`
		GeneratorURL string            `json:"generatorURL,omitempty"`
		Labels       map[string]string `json:"labels,omitempty"`
		StartsAt     time.Time         `json:"startsAt,omitempty"`
		Status       string            `json:"status,omitempty"`
	} `json:"alerts,omitempty"`
	CommonAnnotations map[string]string `json:"commonAnnotations,omitempty"`
	CommonLabels      map[string]string `json:"commonLabels,omitempty"`
	ExternalURL       string            `json:"externalURL,omitempty"`
	GroupKey          string            `json:"groupKey,omitempty"`
	GroupLabels       map[string]string `json:"groupLabels,omitempty"`
	Receiver          string            `json:"receiver,omitempty"`
	Status            string            `json:"status,omitempty"`
	TruncatedAlerts   int               `json:"truncatedAlerts,omitempty"`
	Version           string            `json:"version,omitempty"`
}

// CreateAlertJSONRequestBody defines body for CreateAlert for application/json ContentType.
type CreateAlertJSONRequestBody CreateAlertJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get all alerts
	// (GET /alerts)
	GetAllAlerts(w http.ResponseWriter, r *http.Request, params GetAllAlertsParams)
	// Get alerts for a specific instance
	// (GET /alerts/instances/{instance})
	GetAlertsByInstance(w http.ResponseWriter, r *http.Request, instance string, params GetAlertsByInstanceParams)
	// Get alerts filtered by severity
	// (GET /alerts/severity/{severity})
	GetAlertsBySeverity(w http.ResponseWriter, r *http.Request, severity string, params GetAlertsBySeverityParams)
	// Get alerts for today
	// (GET /alerts/today)
	GetAlertsToday(w http.ResponseWriter, r *http.Request, params GetAlertsTodayParams)
	// Get alerts for yesterday
	// (GET /alerts/yesterday)
	GetAlertsYesterday(w http.ResponseWriter, r *http.Request, params GetAlertsYesterdayParams)
	// Retrieve alert statistics
	// (GET /statistics)
	Statistics(w http.ResponseWriter, r *http.Request, params StatisticsParams)
	// Create a new alert
	// (POST /webhook)
	CreateAlert(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// GetAllAlerts operation middleware
func (siw *ServerInterfaceWrapper) GetAllAlerts(w http.ResponseWriter, r *http.Request) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetAllAlertsParams

	headers := r.Header

	// ------------- Optional header parameter "Content-Type" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("Content-Type")]; found {
		var ContentType GetAllAlertsParamsContentType
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "Content-Type", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithOptions("simple", "Content-Type", valueList[0], &ContentType, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false})
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "Content-Type", Err: err})
			return
		}

		params.ContentType = &ContentType

	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetAllAlerts(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetAlertsByInstance operation middleware
func (siw *ServerInterfaceWrapper) GetAlertsByInstance(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "instance" -------------
	var instance string

	err = runtime.BindStyledParameterWithOptions("simple", "instance", r.PathValue("instance"), &instance, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "instance", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetAlertsByInstanceParams

	headers := r.Header

	// ------------- Optional header parameter "Content-Type" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("Content-Type")]; found {
		var ContentType GetAlertsByInstanceParamsContentType
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "Content-Type", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithOptions("simple", "Content-Type", valueList[0], &ContentType, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false})
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "Content-Type", Err: err})
			return
		}

		params.ContentType = &ContentType

	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetAlertsByInstance(w, r, instance, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetAlertsBySeverity operation middleware
func (siw *ServerInterfaceWrapper) GetAlertsBySeverity(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "severity" -------------
	var severity string

	err = runtime.BindStyledParameterWithOptions("simple", "severity", r.PathValue("severity"), &severity, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "severity", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetAlertsBySeverityParams

	headers := r.Header

	// ------------- Optional header parameter "Content-Type" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("Content-Type")]; found {
		var ContentType GetAlertsBySeverityParamsContentType
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "Content-Type", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithOptions("simple", "Content-Type", valueList[0], &ContentType, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false})
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "Content-Type", Err: err})
			return
		}

		params.ContentType = &ContentType

	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetAlertsBySeverity(w, r, severity, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetAlertsToday operation middleware
func (siw *ServerInterfaceWrapper) GetAlertsToday(w http.ResponseWriter, r *http.Request) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetAlertsTodayParams

	headers := r.Header

	// ------------- Optional header parameter "Content-Type" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("Content-Type")]; found {
		var ContentType GetAlertsTodayParamsContentType
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "Content-Type", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithOptions("simple", "Content-Type", valueList[0], &ContentType, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false})
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "Content-Type", Err: err})
			return
		}

		params.ContentType = &ContentType

	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetAlertsToday(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetAlertsYesterday operation middleware
func (siw *ServerInterfaceWrapper) GetAlertsYesterday(w http.ResponseWriter, r *http.Request) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetAlertsYesterdayParams

	headers := r.Header

	// ------------- Optional header parameter "Content-Type" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("Content-Type")]; found {
		var ContentType GetAlertsYesterdayParamsContentType
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "Content-Type", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithOptions("simple", "Content-Type", valueList[0], &ContentType, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false})
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "Content-Type", Err: err})
			return
		}

		params.ContentType = &ContentType

	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetAlertsYesterday(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// Statistics operation middleware
func (siw *ServerInterfaceWrapper) Statistics(w http.ResponseWriter, r *http.Request) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params StatisticsParams

	headers := r.Header

	// ------------- Optional header parameter "Accept" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("Accept")]; found {
		var Accept StatisticsParamsAccept
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "Accept", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithOptions("simple", "Accept", valueList[0], &Accept, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false})
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "Accept", Err: err})
			return
		}

		params.Accept = &Accept

	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Statistics(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// CreateAlert operation middleware
func (siw *ServerInterfaceWrapper) CreateAlert(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateAlert(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{})
}

// ServeMux is an abstraction of http.ServeMux.
type ServeMux interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type StdHTTPServerOptions struct {
	BaseURL          string
	BaseRouter       ServeMux
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, m ServeMux) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseRouter: m,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, m ServeMux, baseURL string) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseURL:    baseURL,
		BaseRouter: m,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options StdHTTPServerOptions) http.Handler {
	m := options.BaseRouter

	if m == nil {
		m = http.NewServeMux()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	m.HandleFunc("GET "+options.BaseURL+"/alerts", wrapper.GetAllAlerts)
	m.HandleFunc("GET "+options.BaseURL+"/alerts/instances/{instance}", wrapper.GetAlertsByInstance)
	m.HandleFunc("GET "+options.BaseURL+"/alerts/severity/{severity}", wrapper.GetAlertsBySeverity)
	m.HandleFunc("GET "+options.BaseURL+"/alerts/today", wrapper.GetAlertsToday)
	m.HandleFunc("GET "+options.BaseURL+"/alerts/yesterday", wrapper.GetAlertsYesterday)
	m.HandleFunc("GET "+options.BaseURL+"/statistics", wrapper.Statistics)
	m.HandleFunc("POST "+options.BaseURL+"/webhook", wrapper.CreateAlert)

	return m
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"+Ju+yUXtubeJOaur1f8CAAD//zuFezy2OAAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
