// Code generated by protoc-gen-twirp v5.2.0, DO NOT EDIT.
// source: lokahiadmin.proto

/*
Package lokahiadmin is a generated twirp stub package.
This code was generated with github.com/twitchtv/twirp/protoc-gen-twirp v5.2.0.

lokahiadmin is the administrative/backend API for lokahi power usage.


It is generated from these files:
	lokahiadmin.proto
*/
package lokahiadmin

import bytes "bytes"
import strings "strings"
import context "context"
import fmt "fmt"
import ioutil "io/ioutil"
import log "log"
import http "net/http"

import jsonpb "github.com/golang/protobuf/jsonpb"
import proto "github.com/golang/protobuf/proto"
import twirp "github.com/twitchtv/twirp"
import ctxsetters "github.com/twitchtv/twirp/ctxsetters"

// Imports only used by utility functions:
import io "io"
import strconv "strconv"
import json "encoding/json"
import url "net/url"

// ==================
// RunLocal Interface
// ==================

type RunLocal interface {
	Checks(context.Context, *CheckIDs) (*Run, error)

	Stats(context.Context, *Nil) (*HistogramData, error)
}

// ========================
// RunLocal Protobuf Client
// ========================

type runLocalProtobufClient struct {
	client HTTPClient
	urls   [2]string
}

// NewRunLocalProtobufClient creates a Protobuf client that implements the RunLocal interface.
// It communicates using Protobuf and can be configured with a custom HTTPClient.
func NewRunLocalProtobufClient(addr string, client HTTPClient) RunLocal {
	prefix := urlBase(addr) + RunLocalPathPrefix
	urls := [2]string{
		prefix + "Checks",
		prefix + "Stats",
	}
	if httpClient, ok := client.(*http.Client); ok {
		return &runLocalProtobufClient{
			client: withoutRedirects(httpClient),
			urls:   urls,
		}
	}
	return &runLocalProtobufClient{
		client: client,
		urls:   urls,
	}
}

func (c *runLocalProtobufClient) Checks(ctx context.Context, in *CheckIDs) (*Run, error) {
	ctx = ctxsetters.WithPackageName(ctx, "github.xe.lokahi.admin")
	ctx = ctxsetters.WithServiceName(ctx, "RunLocal")
	ctx = ctxsetters.WithMethodName(ctx, "Checks")
	out := new(Run)
	err := doProtobufRequest(ctx, c.client, c.urls[0], in, out)
	return out, err
}

func (c *runLocalProtobufClient) Stats(ctx context.Context, in *Nil) (*HistogramData, error) {
	ctx = ctxsetters.WithPackageName(ctx, "github.xe.lokahi.admin")
	ctx = ctxsetters.WithServiceName(ctx, "RunLocal")
	ctx = ctxsetters.WithMethodName(ctx, "Stats")
	out := new(HistogramData)
	err := doProtobufRequest(ctx, c.client, c.urls[1], in, out)
	return out, err
}

// ====================
// RunLocal JSON Client
// ====================

type runLocalJSONClient struct {
	client HTTPClient
	urls   [2]string
}

// NewRunLocalJSONClient creates a JSON client that implements the RunLocal interface.
// It communicates using JSON and can be configured with a custom HTTPClient.
func NewRunLocalJSONClient(addr string, client HTTPClient) RunLocal {
	prefix := urlBase(addr) + RunLocalPathPrefix
	urls := [2]string{
		prefix + "Checks",
		prefix + "Stats",
	}
	if httpClient, ok := client.(*http.Client); ok {
		return &runLocalJSONClient{
			client: withoutRedirects(httpClient),
			urls:   urls,
		}
	}
	return &runLocalJSONClient{
		client: client,
		urls:   urls,
	}
}

func (c *runLocalJSONClient) Checks(ctx context.Context, in *CheckIDs) (*Run, error) {
	ctx = ctxsetters.WithPackageName(ctx, "github.xe.lokahi.admin")
	ctx = ctxsetters.WithServiceName(ctx, "RunLocal")
	ctx = ctxsetters.WithMethodName(ctx, "Checks")
	out := new(Run)
	err := doJSONRequest(ctx, c.client, c.urls[0], in, out)
	return out, err
}

func (c *runLocalJSONClient) Stats(ctx context.Context, in *Nil) (*HistogramData, error) {
	ctx = ctxsetters.WithPackageName(ctx, "github.xe.lokahi.admin")
	ctx = ctxsetters.WithServiceName(ctx, "RunLocal")
	ctx = ctxsetters.WithMethodName(ctx, "Stats")
	out := new(HistogramData)
	err := doJSONRequest(ctx, c.client, c.urls[1], in, out)
	return out, err
}

// =======================
// RunLocal Server Handler
// =======================

type runLocalServer struct {
	RunLocal
	hooks *twirp.ServerHooks
}

func NewRunLocalServer(svc RunLocal, hooks *twirp.ServerHooks) TwirpServer {
	return &runLocalServer{
		RunLocal: svc,
		hooks:    hooks,
	}
}

// writeError writes an HTTP response with a valid Twirp error format, and triggers hooks.
// If err is not a twirp.Error, it will get wrapped with twirp.InternalErrorWith(err)
func (s *runLocalServer) writeError(ctx context.Context, resp http.ResponseWriter, err error) {
	writeError(ctx, resp, err, s.hooks)
}

// RunLocalPathPrefix is used for all URL paths on a twirp RunLocal server.
// Requests are always: POST RunLocalPathPrefix/method
// It can be used in an HTTP mux to route twirp requests along with non-twirp requests on other routes.
const RunLocalPathPrefix = "/twirp/github.xe.lokahi.admin.RunLocal/"

func (s *runLocalServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ctx = ctxsetters.WithPackageName(ctx, "github.xe.lokahi.admin")
	ctx = ctxsetters.WithServiceName(ctx, "RunLocal")
	ctx = ctxsetters.WithResponseWriter(ctx, resp)

	var err error
	ctx, err = callRequestReceived(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	if req.Method != "POST" {
		msg := fmt.Sprintf("unsupported method %q (only POST is allowed)", req.Method)
		err = badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, err)
		return
	}

	switch req.URL.Path {
	case "/twirp/github.xe.lokahi.admin.RunLocal/Checks":
		s.serveChecks(ctx, resp, req)
		return
	case "/twirp/github.xe.lokahi.admin.RunLocal/Stats":
		s.serveStats(ctx, resp, req)
		return
	default:
		msg := fmt.Sprintf("no handler for path %q", req.URL.Path)
		err = badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, err)
		return
	}
}

func (s *runLocalServer) serveChecks(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	header := req.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}
	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		s.serveChecksJSON(ctx, resp, req)
	case "application/protobuf":
		s.serveChecksProtobuf(ctx, resp, req)
	default:
		msg := fmt.Sprintf("unexpected Content-Type: %q", req.Header.Get("Content-Type"))
		twerr := badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, twerr)
	}
}

func (s *runLocalServer) serveChecksJSON(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "Checks")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	defer closebody(req.Body)
	reqContent := new(CheckIDs)
	unmarshaler := jsonpb.Unmarshaler{AllowUnknownFields: true}
	if err = unmarshaler.Unmarshal(req.Body, reqContent); err != nil {
		err = wrapErr(err, "failed to parse request json")
		s.writeError(ctx, resp, twirp.InternalErrorWith(err))
		return
	}

	// Call service method
	var respContent *Run
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.Checks(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *Run and nil error while calling Checks. nil responses are not supported"))
		return
	}

	ctx = callResponsePrepared(ctx, s.hooks)

	var buf bytes.Buffer
	marshaler := &jsonpb.Marshaler{OrigName: true}
	if err = marshaler.Marshal(&buf, respContent); err != nil {
		err = wrapErr(err, "failed to marshal json response")
		s.writeError(ctx, resp, twirp.InternalErrorWith(err))
		return
	}

	ctx = ctxsetters.WithStatusCode(ctx, http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	if _, err = resp.Write(buf.Bytes()); err != nil {
		log.Printf("errored while writing response to client, but already sent response status code to 200: %s", err)
	}
	callResponseSent(ctx, s.hooks)
}

func (s *runLocalServer) serveChecksProtobuf(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "Checks")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	defer closebody(req.Body)
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		err = wrapErr(err, "failed to read request body")
		s.writeError(ctx, resp, twirp.InternalErrorWith(err))
		return
	}
	reqContent := new(CheckIDs)
	if err = proto.Unmarshal(buf, reqContent); err != nil {
		err = wrapErr(err, "failed to parse request proto")
		s.writeError(ctx, resp, twirp.InternalErrorWith(err))
		return
	}

	// Call service method
	var respContent *Run
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.Checks(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *Run and nil error while calling Checks. nil responses are not supported"))
		return
	}

	ctx = callResponsePrepared(ctx, s.hooks)

	respBytes, err := proto.Marshal(respContent)
	if err != nil {
		err = wrapErr(err, "failed to marshal proto response")
		s.writeError(ctx, resp, twirp.InternalErrorWith(err))
		return
	}

	ctx = ctxsetters.WithStatusCode(ctx, http.StatusOK)
	resp.Header().Set("Content-Type", "application/protobuf")
	resp.WriteHeader(http.StatusOK)
	if _, err = resp.Write(respBytes); err != nil {
		log.Printf("errored while writing response to client, but already sent response status code to 200: %s", err)
	}
	callResponseSent(ctx, s.hooks)
}

func (s *runLocalServer) serveStats(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	header := req.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}
	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		s.serveStatsJSON(ctx, resp, req)
	case "application/protobuf":
		s.serveStatsProtobuf(ctx, resp, req)
	default:
		msg := fmt.Sprintf("unexpected Content-Type: %q", req.Header.Get("Content-Type"))
		twerr := badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, twerr)
	}
}

func (s *runLocalServer) serveStatsJSON(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "Stats")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	defer closebody(req.Body)
	reqContent := new(Nil)
	unmarshaler := jsonpb.Unmarshaler{AllowUnknownFields: true}
	if err = unmarshaler.Unmarshal(req.Body, reqContent); err != nil {
		err = wrapErr(err, "failed to parse request json")
		s.writeError(ctx, resp, twirp.InternalErrorWith(err))
		return
	}

	// Call service method
	var respContent *HistogramData
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.Stats(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *HistogramData and nil error while calling Stats. nil responses are not supported"))
		return
	}

	ctx = callResponsePrepared(ctx, s.hooks)

	var buf bytes.Buffer
	marshaler := &jsonpb.Marshaler{OrigName: true}
	if err = marshaler.Marshal(&buf, respContent); err != nil {
		err = wrapErr(err, "failed to marshal json response")
		s.writeError(ctx, resp, twirp.InternalErrorWith(err))
		return
	}

	ctx = ctxsetters.WithStatusCode(ctx, http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	if _, err = resp.Write(buf.Bytes()); err != nil {
		log.Printf("errored while writing response to client, but already sent response status code to 200: %s", err)
	}
	callResponseSent(ctx, s.hooks)
}

func (s *runLocalServer) serveStatsProtobuf(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "Stats")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	defer closebody(req.Body)
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		err = wrapErr(err, "failed to read request body")
		s.writeError(ctx, resp, twirp.InternalErrorWith(err))
		return
	}
	reqContent := new(Nil)
	if err = proto.Unmarshal(buf, reqContent); err != nil {
		err = wrapErr(err, "failed to parse request proto")
		s.writeError(ctx, resp, twirp.InternalErrorWith(err))
		return
	}

	// Call service method
	var respContent *HistogramData
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.Stats(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *HistogramData and nil error while calling Stats. nil responses are not supported"))
		return
	}

	ctx = callResponsePrepared(ctx, s.hooks)

	respBytes, err := proto.Marshal(respContent)
	if err != nil {
		err = wrapErr(err, "failed to marshal proto response")
		s.writeError(ctx, resp, twirp.InternalErrorWith(err))
		return
	}

	ctx = ctxsetters.WithStatusCode(ctx, http.StatusOK)
	resp.Header().Set("Content-Type", "application/protobuf")
	resp.WriteHeader(http.StatusOK)
	if _, err = resp.Write(respBytes); err != nil {
		log.Printf("errored while writing response to client, but already sent response status code to 200: %s", err)
	}
	callResponseSent(ctx, s.hooks)
}

func (s *runLocalServer) ServiceDescriptor() ([]byte, int) {
	return twirpFileDescriptor0, 0
}

func (s *runLocalServer) ProtocGenTwirpVersion() string {
	return "v5.2.0"
}

// =====
// Utils
// =====

// HTTPClient is the interface used by generated clients to send HTTP requests.
// It is fulfilled by *(net/http).Client, which is sufficient for most users.
// Users can provide their own implementation for special retry policies.
//
// HTTPClient implementations should not follow redirects. Redirects are
// automatically disabled if *(net/http).Client is passed to client
// constructors. See the withoutRedirects function in this file for more
// details.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// TwirpServer is the interface generated server structs will support: they're
// HTTP handlers with additional methods for accessing metadata about the
// service. Those accessors are a low-level API for building reflection tools.
// Most people can think of TwirpServers as just http.Handlers.
type TwirpServer interface {
	http.Handler
	// ServiceDescriptor returns gzipped bytes describing the .proto file that
	// this service was generated from. Once unzipped, the bytes can be
	// unmarshalled as a
	// github.com/golang/protobuf/protoc-gen-go/descriptor.FileDescriptorProto.
	//
	// The returned integer is the index of this particular service within that
	// FileDescriptorProto's 'Service' slice of ServiceDescriptorProtos. This is a
	// low-level field, expected to be used for reflection.
	ServiceDescriptor() ([]byte, int)
	// ProtocGenTwirpVersion is the semantic version string of the version of
	// twirp used to generate this file.
	ProtocGenTwirpVersion() string
}

// WriteError writes an HTTP response with a valid Twirp error format.
// If err is not a twirp.Error, it will get wrapped with twirp.InternalErrorWith(err)
func WriteError(resp http.ResponseWriter, err error) {
	writeError(context.Background(), resp, err, nil)
}

// writeError writes Twirp errors in the response and triggers hooks.
func writeError(ctx context.Context, resp http.ResponseWriter, err error, hooks *twirp.ServerHooks) {
	// Non-twirp errors are wrapped as Internal (default)
	twerr, ok := err.(twirp.Error)
	if !ok {
		twerr = twirp.InternalErrorWith(err)
	}

	statusCode := twirp.ServerHTTPStatusFromErrorCode(twerr.Code())
	ctx = ctxsetters.WithStatusCode(ctx, statusCode)
	ctx = callError(ctx, hooks, twerr)

	resp.Header().Set("Content-Type", "application/json") // Error responses are always JSON (instead of protobuf)
	resp.WriteHeader(statusCode)                          // HTTP response status code

	respBody := marshalErrorToJSON(twerr)
	_, err2 := resp.Write(respBody)
	if err2 != nil {
		log.Printf("unable to send error message %q: %s", twerr, err2)
	}

	callResponseSent(ctx, hooks)
}

// urlBase helps ensure that addr specifies a scheme. If it is unparsable
// as a URL, it returns addr unchanged.
func urlBase(addr string) string {
	// If the addr specifies a scheme, use it. If not, default to
	// http. If url.Parse fails on it, return it unchanged.
	url, err := url.Parse(addr)
	if err != nil {
		return addr
	}
	if url.Scheme == "" {
		url.Scheme = "http"
	}
	return url.String()
}

// getCustomHTTPReqHeaders retrieves a copy of any headers that are set in
// a context through the twirp.WithHTTPRequestHeaders function.
// If there are no headers set, or if they have the wrong type, nil is returned.
func getCustomHTTPReqHeaders(ctx context.Context) http.Header {
	header, ok := twirp.HTTPRequestHeaders(ctx)
	if !ok || header == nil {
		return nil
	}
	copied := make(http.Header)
	for k, vv := range header {
		if vv == nil {
			copied[k] = nil
			continue
		}
		copied[k] = make([]string, len(vv))
		copy(copied[k], vv)
	}
	return copied
}

// closebody closes a response or request body and just logs
// any error encountered while closing, since errors are
// considered very unusual.
func closebody(body io.Closer) {
	if err := body.Close(); err != nil {
		log.Printf("error closing body: %q", err)
	}
}

// newRequest makes an http.Request from a client, adding common headers.
func newRequest(ctx context.Context, url string, reqBody io.Reader, contentType string) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if customHeader := getCustomHTTPReqHeaders(ctx); customHeader != nil {
		req.Header = customHeader
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Twirp-Version", "v5.2.0")
	return req, nil
}

// JSON serialization for errors
type twerrJSON struct {
	Code string            `json:"code"`
	Msg  string            `json:"msg"`
	Meta map[string]string `json:"meta,omitempty"`
}

// marshalErrorToJSON returns JSON from a twirp.Error, that can be used as HTTP error response body.
// If serialization fails, it will use a descriptive Internal error instead.
func marshalErrorToJSON(twerr twirp.Error) []byte {
	// make sure that msg is not too large
	msg := twerr.Msg()
	if len(msg) > 1e6 {
		msg = msg[:1e6]
	}

	tj := twerrJSON{
		Code: string(twerr.Code()),
		Msg:  msg,
		Meta: twerr.MetaMap(),
	}

	buf, err := json.Marshal(&tj)
	if err != nil {
		buf = []byte("{\"type\": \"" + twirp.Internal + "\", \"msg\": \"There was an error but it could not be serialized into JSON\"}") // fallback
	}

	return buf
}

// errorFromResponse builds a twirp.Error from a non-200 HTTP response.
// If the response has a valid serialized Twirp error, then it's returned.
// If not, the response status code is used to generate a similar twirp
// error. See twirpErrorFromIntermediary for more info on intermediary errors.
func errorFromResponse(resp *http.Response) twirp.Error {
	statusCode := resp.StatusCode
	statusText := http.StatusText(statusCode)

	if isHTTPRedirect(statusCode) {
		// Unexpected redirect: it must be an error from an intermediary.
		// Twirp clients don't follow redirects automatically, Twirp only handles
		// POST requests, redirects should only happen on GET and HEAD requests.
		location := resp.Header.Get("Location")
		msg := fmt.Sprintf("unexpected HTTP status code %d %q received, Location=%q", statusCode, statusText, location)
		return twirpErrorFromIntermediary(statusCode, msg, location)
	}

	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return clientError("failed to read server error response body", err)
	}
	var tj twerrJSON
	if err := json.Unmarshal(respBodyBytes, &tj); err != nil {
		// Invalid JSON response; it must be an error from an intermediary.
		msg := fmt.Sprintf("Error from intermediary with HTTP status code %d %q", statusCode, statusText)
		return twirpErrorFromIntermediary(statusCode, msg, string(respBodyBytes))
	}

	errorCode := twirp.ErrorCode(tj.Code)
	if !twirp.IsValidErrorCode(errorCode) {
		msg := "invalid type returned from server error response: " + tj.Code
		return twirp.InternalError(msg)
	}

	twerr := twirp.NewError(errorCode, tj.Msg)
	for k, v := range tj.Meta {
		twerr = twerr.WithMeta(k, v)
	}
	return twerr
}

// twirpErrorFromIntermediary maps HTTP errors from non-twirp sources to twirp errors.
// The mapping is similar to gRPC: https://github.com/grpc/grpc/blob/master/doc/http-grpc-status-mapping.md.
// Returned twirp Errors have some additional metadata for inspection.
func twirpErrorFromIntermediary(status int, msg string, bodyOrLocation string) twirp.Error {
	var code twirp.ErrorCode
	if isHTTPRedirect(status) { // 3xx
		code = twirp.Internal
	} else {
		switch status {
		case 400: // Bad Request
			code = twirp.Internal
		case 401: // Unauthorized
			code = twirp.Unauthenticated
		case 403: // Forbidden
			code = twirp.PermissionDenied
		case 404: // Not Found
			code = twirp.BadRoute
		case 429, 502, 503, 504: // Too Many Requests, Bad Gateway, Service Unavailable, Gateway Timeout
			code = twirp.Unavailable
		default: // All other codes
			code = twirp.Unknown
		}
	}

	twerr := twirp.NewError(code, msg)
	twerr = twerr.WithMeta("http_error_from_intermediary", "true") // to easily know if this error was from intermediary
	twerr = twerr.WithMeta("status_code", strconv.Itoa(status))
	if isHTTPRedirect(status) {
		twerr = twerr.WithMeta("location", bodyOrLocation)
	} else {
		twerr = twerr.WithMeta("body", bodyOrLocation)
	}
	return twerr
}
func isHTTPRedirect(status int) bool {
	return status >= 300 && status <= 399
}

// wrappedError implements the github.com/pkg/errors.Causer interface, allowing errors to be
// examined for their root cause.
type wrappedError struct {
	msg   string
	cause error
}

func wrapErr(err error, msg string) error { return &wrappedError{msg: msg, cause: err} }
func (e *wrappedError) Cause() error      { return e.cause }
func (e *wrappedError) Error() string     { return e.msg + ": " + e.cause.Error() }

// clientError adds consistency to errors generated in the client
func clientError(desc string, err error) twirp.Error {
	return twirp.InternalErrorWith(wrapErr(err, desc))
}

// badRouteError is used when the twirp server cannot route a request
func badRouteError(msg string, method, url string) twirp.Error {
	err := twirp.NewError(twirp.BadRoute, msg)
	err = err.WithMeta("twirp_invalid_route", method+" "+url)
	return err
}

// The standard library will, by default, redirect requests (including POSTs) if it gets a 302 or
// 303 response, and also 301s in go1.8. It redirects by making a second request, changing the
// method to GET and removing the body. This produces very confusing error messages, so instead we
// set a redirect policy that always errors. This stops Go from executing the redirect.
//
// We have to be a little careful in case the user-provided http.Client has its own CheckRedirect
// policy - if so, we'll run through that policy first.
//
// Because this requires modifying the http.Client, we make a new copy of the client and return it.
func withoutRedirects(in *http.Client) *http.Client {
	copy := *in
	copy.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if in.CheckRedirect != nil {
			// Run the input's redirect if it exists, in case it has side effects, but ignore any error it
			// returns, since we want to use ErrUseLastResponse.
			err := in.CheckRedirect(req, via)
			_ = err // Silly, but this makes sure generated code passes errcheck -blank, which some people use.
		}
		return http.ErrUseLastResponse
	}
	return &copy
}

// doProtobufRequest is common code to make a request to the remote twirp service.
func doProtobufRequest(ctx context.Context, client HTTPClient, url string, in, out proto.Message) (err error) {
	reqBodyBytes, err := proto.Marshal(in)
	if err != nil {
		return clientError("failed to marshal proto request", err)
	}
	reqBody := bytes.NewBuffer(reqBodyBytes)
	if err = ctx.Err(); err != nil {
		return clientError("aborted because context was done", err)
	}

	req, err := newRequest(ctx, url, reqBody, "application/protobuf")
	if err != nil {
		return clientError("could not build request", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return clientError("failed to do request", err)
	}

	defer func() {
		cerr := resp.Body.Close()
		if err == nil && cerr != nil {
			err = clientError("failed to close response body", cerr)
		}
	}()

	if err = ctx.Err(); err != nil {
		return clientError("aborted because context was done", err)
	}

	if resp.StatusCode != 200 {
		return errorFromResponse(resp)
	}

	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return clientError("failed to read response body", err)
	}
	if err = ctx.Err(); err != nil {
		return clientError("aborted because context was done", err)
	}

	if err = proto.Unmarshal(respBodyBytes, out); err != nil {
		return clientError("failed to unmarshal proto response", err)
	}
	return nil
}

// doJSONRequest is common code to make a request to the remote twirp service.
func doJSONRequest(ctx context.Context, client HTTPClient, url string, in, out proto.Message) (err error) {
	reqBody := bytes.NewBuffer(nil)
	marshaler := &jsonpb.Marshaler{OrigName: true}
	if err = marshaler.Marshal(reqBody, in); err != nil {
		return clientError("failed to marshal json request", err)
	}
	if err = ctx.Err(); err != nil {
		return clientError("aborted because context was done", err)
	}

	req, err := newRequest(ctx, url, reqBody, "application/json")
	if err != nil {
		return clientError("could not build request", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return clientError("failed to do request", err)
	}

	defer func() {
		cerr := resp.Body.Close()
		if err == nil && cerr != nil {
			err = clientError("failed to close response body", cerr)
		}
	}()

	if err = ctx.Err(); err != nil {
		return clientError("aborted because context was done", err)
	}

	if resp.StatusCode != 200 {
		return errorFromResponse(resp)
	}

	unmarshaler := jsonpb.Unmarshaler{AllowUnknownFields: true}
	if err = unmarshaler.Unmarshal(resp.Body, out); err != nil {
		return clientError("failed to unmarshal json response", err)
	}
	if err = ctx.Err(); err != nil {
		return clientError("aborted because context was done", err)
	}
	return nil
}

// Call twirp.ServerHooks.RequestReceived if the hook is available
func callRequestReceived(ctx context.Context, h *twirp.ServerHooks) (context.Context, error) {
	if h == nil || h.RequestReceived == nil {
		return ctx, nil
	}
	return h.RequestReceived(ctx)
}

// Call twirp.ServerHooks.RequestRouted if the hook is available
func callRequestRouted(ctx context.Context, h *twirp.ServerHooks) (context.Context, error) {
	if h == nil || h.RequestRouted == nil {
		return ctx, nil
	}
	return h.RequestRouted(ctx)
}

// Call twirp.ServerHooks.ResponsePrepared if the hook is available
func callResponsePrepared(ctx context.Context, h *twirp.ServerHooks) context.Context {
	if h == nil || h.ResponsePrepared == nil {
		return ctx
	}
	return h.ResponsePrepared(ctx)
}

// Call twirp.ServerHooks.ResponseSent if the hook is available
func callResponseSent(ctx context.Context, h *twirp.ServerHooks) {
	if h == nil || h.ResponseSent == nil {
		return
	}
	h.ResponseSent(ctx)
}

// Call twirp.ServerHooks.Error if the hook is available
func callError(ctx context.Context, h *twirp.ServerHooks, err twirp.Error) context.Context {
	if h == nil || h.Error == nil {
		return ctx
	}
	return h.Error(ctx, err)
}

var twirpFileDescriptor0 = []byte{
	// 562 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x94, 0xd1, 0x6e, 0xd3, 0x30,
	0x14, 0x86, 0x95, 0xa4, 0xe9, 0xb2, 0x53, 0xd6, 0x0d, 0x83, 0x46, 0x28, 0x48, 0x44, 0x95, 0x80,
	0x70, 0x13, 0xa6, 0x41, 0xb5, 0x75, 0x97, 0xdb, 0x10, 0x43, 0xa0, 0x5d, 0x18, 0xb8, 0xe1, 0x82,
	0xca, 0xad, 0xcd, 0x62, 0x35, 0xb1, 0xa3, 0xd8, 0x99, 0xda, 0xb7, 0x41, 0xbc, 0x06, 0x2f, 0xc0,
	0x63, 0xa1, 0x38, 0x0d, 0x24, 0xa3, 0x1d, 0x77, 0x3e, 0x3e, 0x9f, 0x7f, 0xfb, 0xfc, 0xc7, 0x36,
	0xdc, 0x4d, 0xe4, 0x9c, 0xc4, 0x9c, 0xd0, 0x94, 0x8b, 0x28, 0xcb, 0xa5, 0x96, 0x68, 0xff, 0x8a,
	0xeb, 0xb8, 0x98, 0x46, 0x0b, 0x16, 0x55, 0xc9, 0xc8, 0x64, 0x87, 0x8f, 0xc1, 0x3b, 0x8b, 0xd9,
	0x6c, 0xfe, 0xee, 0x5c, 0xa1, 0x3d, 0x70, 0x38, 0x55, 0xbe, 0x15, 0x38, 0xe1, 0x36, 0x2e, 0x87,
	0xc3, 0x1f, 0x1d, 0x70, 0x70, 0x21, 0x50, 0x1f, 0x6c, 0x4e, 0x7d, 0x2b, 0xb0, 0xc2, 0x6d, 0x6c,
	0x73, 0x8a, 0x5e, 0x43, 0x67, 0x56, 0xa2, 0x76, 0x60, 0x85, 0xbd, 0xc3, 0x20, 0x5a, 0x2f, 0x1e,
	0xd5, 0xca, 0xd8, 0xd0, 0x68, 0x00, 0xde, 0x37, 0x2e, 0xb8, 0x8a, 0x19, 0xf5, 0x9d, 0xc0, 0x0a,
	0x3d, 0xfc, 0x27, 0x46, 0xa7, 0xb0, 0x95, 0x33, 0x55, 0x24, 0x5a, 0xf9, 0x9d, 0xc0, 0x09, 0x7b,
	0x87, 0xe1, 0x26, 0x51, 0x5c, 0x88, 0x08, 0x57, 0xe8, 0x1b, 0xa1, 0xf3, 0x25, 0xae, 0x17, 0xa2,
	0x67, 0xb0, 0xab, 0x34, 0xc9, 0xf5, 0x44, 0xf3, 0x94, 0x4d, 0x0a, 0xc1, 0x17, 0xbe, 0x1b, 0x58,
	0xa1, 0x83, 0x77, 0xcc, 0xf4, 0x27, 0x9e, 0xb2, 0xcf, 0x82, 0x2f, 0xd0, 0x4b, 0xb8, 0xc7, 0x12,
	0x92, 0x29, 0x46, 0x27, 0x82, 0x08, 0xa9, 0xd8, 0x4c, 0x0a, 0xaa, 0xfc, 0xae, 0x61, 0xd1, 0x2a,
	0x75, 0xf9, 0x37, 0x33, 0xf8, 0x69, 0x41, 0xf7, 0x82, 0x91, 0x44, 0xc7, 0xa5, 0x47, 0x45, 0x9e,
	0xac, 0xac, 0x28, 0x87, 0xe8, 0x04, 0x1e, 0xe6, 0x4c, 0x65, 0x52, 0x28, 0x56, 0x6d, 0xdc, 0xd4,
	0xb4, 0x8d, 0xe6, 0x83, 0x1a, 0x28, 0x8f, 0xd0, 0x10, 0x46, 0x4f, 0xa0, 0xa7, 0x34, 0xd1, 0x85,
	0x9a, 0xcc, 0x24, 0x65, 0xc6, 0x14, 0x17, 0x43, 0x35, 0x75, 0x26, 0x29, 0x43, 0x08, 0x3a, 0x53,
	0x49, 0x97, 0x7e, 0xc7, 0xec, 0x67, 0xc6, 0xe8, 0x3e, 0xb8, 0x2c, 0xcf, 0x65, 0x6e, 0x8a, 0xdb,
	0xc6, 0x55, 0x80, 0x7c, 0xd8, 0x8a, 0xcd, 0x11, 0x97, 0xa6, 0x10, 0x0f, 0xd7, 0xe1, 0xe0, 0x2b,
	0xdc, 0x69, 0xfa, 0x55, 0x96, 0x30, 0x67, 0xcb, 0xba, 0x84, 0x39, 0x5b, 0xa2, 0x63, 0x70, 0xaf,
	0x49, 0x52, 0xb0, 0x55, 0x3f, 0x87, 0xb7, 0x59, 0x5f, 0xf9, 0x80, 0xab, 0x05, 0x27, 0xf6, 0xb1,
	0x35, 0x74, 0xc1, 0xb9, 0xe4, 0xc9, 0xf0, 0x97, 0x0d, 0x3b, 0x17, 0x5c, 0x69, 0x79, 0x95, 0x93,
	0xf4, 0x9c, 0x68, 0x82, 0x9e, 0xc3, 0x6e, 0x4a, 0x16, 0x2d, 0x3f, 0x2c, 0xe3, 0x47, 0x3f, 0x25,
	0x8b, 0xa6, 0x0d, 0x25, 0xc8, 0xc5, 0x1a, 0xe3, 0xfa, 0x29, 0x17, 0x4d, 0xf0, 0x05, 0xec, 0xa5,
	0x8c, 0xb4, 0x49, 0xc7, 0x90, 0xbb, 0xe5, 0x7c, 0x13, 0xdd, 0x87, 0xae, 0xd2, 0x94, 0xb2, 0x6b,
	0xe3, 0x9d, 0x83, 0x57, 0x51, 0xb9, 0x57, 0x36, 0x3a, 0x68, 0x29, 0x54, 0x97, 0xa4, 0x9f, 0x8d,
	0x0e, 0x6e, 0x1c, 0x2a, 0x3b, 0x1a, 0xad, 0xb9, 0x21, 0xfd, 0xec, 0x68, 0x74, 0x13, 0x1c, 0xb7,
	0xc1, 0xad, 0x15, 0x38, 0xfe, 0x17, 0x1c, 0xb7, 0x40, 0xaf, 0x06, 0xc7, 0x0d, 0xf0, 0xf0, 0xbb,
	0x05, 0x1e, 0x2e, 0xc4, 0x07, 0x39, 0x23, 0x09, 0x7a, 0x0b, 0x5d, 0xf3, 0x8e, 0x14, 0xfa, 0xef,
	0x3b, 0x1b, 0x3c, 0xba, 0xa5, 0x73, 0xe8, 0x3d, 0xb8, 0x1f, 0x35, 0xd1, 0x0a, 0x6d, 0xa4, 0x2e,
	0x79, 0x32, 0x78, 0xba, 0x29, 0xd9, 0xea, 0xed, 0xe9, 0xce, 0x97, 0x5e, 0xe3, 0x93, 0x99, 0x76,
	0xcd, 0x2f, 0xf3, 0xea, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0xe9, 0x4d, 0x8e, 0x55, 0x7a, 0x04,
	0x00, 0x00,
}
