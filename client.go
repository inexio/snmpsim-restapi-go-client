package snmpsimclient

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	// mgmtEndpointPath path to the mgmt api endpoint
	mgmtEndpointPath = "snmpsim/mgmt/v1/"
	// metricsEndpointPath path to the metrics api endpoint
	metricsEndpointPath = "snmpsim/metrics/v1/"
)

type client struct {
	*clientData
}

type clientData struct {
	baseUrl  string
	username string
	password string

	resty   *resty.Client
	useAuth bool
}

/*
NotValidError is returned when the client was not initialized properly with the NewManagementClient() func
*/
type NotValidError struct{}

func (m *NotValidError) Error() string {
	return "client was not created properly with the func New...Client(baseUrl string)"
}

//isValid checks if the client object is valid
func (c *client) isValid() bool {
	return c.clientData != nil
}

/*
SetUsernameAndPassword is used to set a username and password for https auth
*/
func (c *client) SetUsernameAndPassword(username, password string) error {
	if !c.isValid() {
		return &NotValidError{}
	}
	if username == "" {
		return errors.New("invalid username")
	}
	if password == "" {
		return errors.New("invalid password")
	}
	c.username = username
	c.password = password
	c.useAuth = true
	return nil
}

/*
SetTimeout can be used to set a timeout for requests raised from client.
*/
func (c *client) SetTimeout(timeout time.Duration) {
	c.resty.SetTimeout(timeout)
}

func (c *client) request(method string, path string, body string, header, queryParams map[string]string) (*resty.Response, error) {
	request := c.resty.R()
	request.SetHeader("Content-Type", "application/json")

	if header != nil {
		request.SetHeaders(header)
	}

	if queryParams != nil {
		request.SetQueryParams(queryParams)
	}

	if body != "" {
		request.SetBody(body)
	}

	if c.useAuth {
		request.SetBasicAuth(c.username, c.password)
	}

	var response *resty.Response
	response = nil

	var err error
	err = nil

	switch method {
	case "GET":
		response, err = request.Get(c.baseUrl + urlEscapePath(path))
	case "POST":
		response, err = request.Post(c.baseUrl + urlEscapePath(path))
	case "PUT":
		response, err = request.Put(c.baseUrl + urlEscapePath(path))
	case "DELETE":
		response, err = request.Delete(c.baseUrl + urlEscapePath(path))
	default:
		return nil, errors.New("invalid http method: " + method)
	}
	if err != nil {
		return nil, errors.Wrap(err, "error during http request")
	}
	return response, nil
}

//Http error handling

/*
HttpError represents an http error returned by the api.
*/
type HttpError struct {
	StatusCode int
	Status     string
	Body       *ErrorResponse
}

func (h HttpError) Error() string {
	msg := "http error: status code: " + strconv.Itoa(h.StatusCode) + " // status: " + h.Status
	if h.Body != nil {
		msg += " // message: " + h.Body.Message
	}
	return msg
}

func getHttpError(response *resty.Response) error {
	httpError := HttpError{
		StatusCode: response.StatusCode(),
		Status:     response.Status(),
	}
	var errorResponse ErrorResponse
	err := json.Unmarshal(response.Body(), &errorResponse)
	if err != nil {
		return httpError
	}
	httpError.Body = &errorResponse
	return httpError
}

//helper functions
func urlEscapePath(unescaped string) string {
	arr := strings.Split(unescaped, "/")
	for i, partString := range strings.Split(unescaped, "/") {
		arr[i] = url.QueryEscape(partString)
	}
	return strings.Join(arr, "/")
}
