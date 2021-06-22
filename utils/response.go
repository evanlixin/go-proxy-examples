package utils

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// APIVersion xxx
type APIVersion string

// const values
const (
	V2APIVersion = APIVersion("1.0")
)

type commResp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NotFoundHander is the default NotFoundHandler
type NotFoundHander struct{}

func (NotFoundHander) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ResourceNotFound(w, r, r.URL.String())
}

// CommReply can be used for replying some common data
func CommReply(w http.ResponseWriter, r *http.Request, status int, message string) {
	resp := commResp{
		Code:    http.StatusText(status),
		Message: message,
	}
	Reply(w, r, status, resp)
}

// ResponseReply xxx
func ResponseReply(w http.ResponseWriter, r *http.Request, status int, version APIVersion, kind string, v interface{}) {
	resp := map[string]interface{}{}
	resp[`apiVersion`] = version
	resp[kind] = v
	Reply(w, r, status, resp)
}

// Reply can be used for replying response
// Rename this to "reply" in the future because should call ResponseReply instead
func Reply(w http.ResponseWriter, r *http.Request, status int, v interface{}) {
	data, _ := json.Marshal(v)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

// OK reply
func OK(w http.ResponseWriter, r *http.Request, message string) {
	CommReply(w, r, http.StatusOK, message)
}

// ResourceNotFound will return an error message indicating that the resource is not exist
func ResourceNotFound(w http.ResponseWriter, r *http.Request, message string) {
	CommReply(w, r, http.StatusNotFound, message)
}

// BadRequest will return an error message indicating that the request is invalid
func BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	CommReply(w, r, http.StatusBadRequest, err.Error())
}

// Forbidden will block user access the resource, not authorized
func Forbidden(w http.ResponseWriter, r *http.Request, err error) {
	CommReply(w, r, http.StatusForbidden, err.Error())
}

// Unauthorized will block user access the api, not login
func Unauthorized(w http.ResponseWriter, r *http.Request, err error) {
	CommReply(w, r, http.StatusUnauthorized, err.Error())
}

// InternalError will return an error message indicating that the something is error inside the controller
func InternalError(w http.ResponseWriter, r *http.Request, err error) {
	CommReply(w, r, http.StatusInternalServerError, err.Error())
}

// ServiceUnavailable will return an error message indicating that the service is not available now
func ServiceUnavailable(w http.ResponseWriter, r *http.Request, err error) {
	CommReply(w, r, http.StatusServiceUnavailable, err.Error())
}

// Conflict xxx
func Conflict(w http.ResponseWriter, r *http.Request, err error) {
	CommReply(w, r, http.StatusConflict, err.Error())
}

// NotAcceptable xxx
func NotAcceptable(w http.ResponseWriter, r *http.Request, err error) {
	CommReply(w, r, http.StatusNotAcceptable, err.Error())
}

// SetRequestID will set the response header of the requestID
func SetRequestID(w http.ResponseWriter, requestID string) {
	w.Header().Set("x-request-id", requestID)
}

// QueryString will return the query string value
func QueryString(r *http.Request, query string) string {
	return mux.Vars(r)[query]
}
