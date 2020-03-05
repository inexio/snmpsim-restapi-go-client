package snmpsimclient

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"strconv"
)

/*
MetricsClient is a client for communicating with the metrics api.
*/
type MetricsClient struct {
	client
}

/*
NewMetricsClient creates a new NewMetricsClient.
*/
func NewMetricsClient(baseUrl string) (*MetricsClient, error) {
	if baseUrl == "" {
		return nil, errors.New("invalid base url")
	}
	//if baseUrl does not end with an "/" it has to be added to the string
	if lastChar := baseUrl[len(baseUrl)-1:]; lastChar != "/" {
		baseUrl += "/"
	}
	clientData := clientData{baseUrl: baseUrl, resty: resty.New(), useAuth: false}
	newClient := client{&clientData}
	return &MetricsClient{newClient}, nil
}

/*
GetProcesses returns process metrics.
*/
func (c *MetricsClient) GetProcesses(filters map[string]string) (ProcessesMetrics, error) {
	response, err := c.request("GET", metricsEndpointPath+"processes", "", nil, filters)
	if err != nil {
		return nil, errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 200 {
		return nil, getHttpError(response)
	}
	var processes ProcessesMetrics
	err = json.Unmarshal(response.Body(), &processes)
	if err != nil {
		return nil, errors.Wrap(err, "error during unmarshalling http response")
	}
	return processes, nil
}

/*
GetProcess returns the process with the given id.
*/
func (c *MetricsClient) GetProcess(id int) (ProcessMetrics, error) {
	response, err := c.request("GET", metricsEndpointPath+"processes/"+strconv.Itoa(id), "", nil, nil)
	if err != nil {
		return ProcessMetrics{}, errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 200 {
		return ProcessMetrics{}, getHttpError(response)
	}
	var process ProcessMetrics
	err = json.Unmarshal(response.Body(), &process)
	if err != nil {
		return ProcessMetrics{}, errors.Wrap(err, "error during unmarshalling http response")
	}
	return process, nil
}

/*
GetProcessEndpoints returns an array of endpoints for the given process-id.
*/
func (c *MetricsClient) GetProcessEndpoints(id int) (Endpoints, error) {
	response, err := c.request("GET", metricsEndpointPath+"processes/"+strconv.Itoa(id)+"/endpoints", "", nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 200 {
		return nil, getHttpError(response)
	}
	var endpoints Endpoints
	err = json.Unmarshal(response.Body(), &endpoints)
	if err != nil {
		return nil, errors.Wrap(err, "error during unmarshalling http response")
	}
	return endpoints, nil
}

/*
GetProcessEndpoint returns the endpoint for the given process- and endpoint-id.
*/
func (c *MetricsClient) GetProcessEndpoint(processId int, endpointId int) (Endpoint, error) {
	response, err := c.request("GET", metricsEndpointPath+"processes/"+strconv.Itoa(processId)+"/endpoints/"+strconv.Itoa(endpointId), "", nil, nil)
	if err != nil {
		return Endpoint{}, errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 200 {
		return Endpoint{}, getHttpError(response)
	}
	var endpoint Endpoint
	err = json.Unmarshal(response.Body(), &endpoint)
	if err != nil {
		return Endpoint{}, errors.Wrap(err, "error during unmarshalling http response")
	}
	return endpoint, nil
}

/*
GetProcessConsolePages returns an array of console-pages for the given process-id.
*/
func (c *MetricsClient) GetProcessConsolePages(processId int) (Consoles, error) {
	response, err := c.request("GET", metricsEndpointPath+"processes/"+strconv.Itoa(processId)+"/console", "", nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 200 {
		return nil, getHttpError(response)
	}
	var consolePages Consoles
	err = json.Unmarshal(response.Body(), &consolePages)
	if err != nil {
		return nil, errors.Wrap(err, "error during unmarshalling http response")
	}
	return consolePages, nil
}

/*
GetProcessConsolePage returns the console-pages for the given process- and console-page-id.
*/
func (c *MetricsClient) GetProcessConsolePage(processId int, pageId int) (Console, error) {
	response, err := c.request("GET", metricsEndpointPath+"processes/"+strconv.Itoa(processId)+"/console/"+strconv.Itoa(pageId), "", nil, nil)
	if err != nil {
		return Console{}, errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 200 {
		return Console{}, getHttpError(response)
	}
	var consolePages Console
	err = json.Unmarshal(response.Body(), &consolePages)
	if err != nil {
		return Console{}, errors.Wrap(err, "error during unmarshalling http response")
	}
	return consolePages, nil
}

/*
GetPackets returns packet metrics.
*/
func (c *MetricsClient) GetPackets(filters map[string]string) (PacketMetrics, error) {
	response, err := c.request("GET", metricsEndpointPath+"activity/packets", "", nil, filters)
	if err != nil {
		return PacketMetrics{}, errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 200 {
		return PacketMetrics{}, getHttpError(response)
	}
	var packetMetrics PacketMetrics
	err = json.Unmarshal(response.Body(), &packetMetrics)
	if err != nil {
		return PacketMetrics{}, errors.Wrap(err, "error during unmarshalling http response")
	}
	return packetMetrics, nil
}

/*
GetPacketFilters returns all packet filters.
*/
func (c *MetricsClient) GetPacketFilters() (PacketFilters, error) {
	response, err := c.request("GET", metricsEndpointPath+"activity/packets/filters", "", nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 200 {
		return nil, getHttpError(response)
	}

	var filters map[string]interface{}
	err = json.Unmarshal(response.Body(), &filters)
	if err != nil {
		return nil, errors.Wrap(err, "error during unmarshalling http response")
	}

	var packetFilters PacketFilters
	for key := range filters {
		packetFilters = append(packetFilters, key)
	}

	return packetFilters, nil
}

/*
GetPossibleValuesForPacketFilter returns a list of all values that can be used for the given filter.
*/
func (c *MetricsClient) GetPossibleValuesForPacketFilter(filter string) ([]string, error) {
	response, err := c.request("GET", metricsEndpointPath+"activity/packets/filters/"+filter, "", nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 200 {
		return nil, getHttpError(response)
	}

	var messageFilters []string
	err = json.Unmarshal(response.Body(), &messageFilters)
	if err != nil {
		return nil, errors.Wrap(err, "error during unmarshalling http response")
	}
	return messageFilters, nil
}

/*
GetMessages returns message metrics.
*/
func (c *MetricsClient) GetMessages(filters map[string]string) (MessageMetrics, error) {
	response, err := c.request("GET", metricsEndpointPath+"activity/messages", "", nil, filters)
	if err != nil {
		return MessageMetrics{}, errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 200 {
		return MessageMetrics{}, getHttpError(response)
	}
	var messageMetrics MessageMetrics
	err = json.Unmarshal(response.Body(), &messageMetrics)
	if err != nil {
		return MessageMetrics{}, errors.Wrap(err, "error during unmarshalling http response")
	}
	return messageMetrics, nil
}

/*
GetMessageFilters returns all message filters.
*/
func (c *MetricsClient) GetMessageFilters() (MessageFilters, error) {
	response, err := c.request("GET", metricsEndpointPath+"activity/messages/filters", "", nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 200 {
		return nil, getHttpError(response)
	}

	var filters map[string]interface{}
	err = json.Unmarshal(response.Body(), &filters)
	if err != nil {
		return nil, errors.Wrap(err, "error during unmarshalling http response")
	}


	var messageFilters MessageFilters

	for key := range filters {
		messageFilters = append(messageFilters, key)
	}

	return messageFilters, nil
}

/*
GetPossibleValuesForMessageFilter returns a list of all values that can be used for the given filter.
*/
func (c *MetricsClient) GetPossibleValuesForMessageFilter(filter string) ([]string, error) {
	response, err := c.request("GET", metricsEndpointPath+"activity/messages/filters/"+filter, "", nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 200 {
		return nil, getHttpError(response)
	}

	var messageFilters []string
	err = json.Unmarshal(response.Body(), &messageFilters)
	if err != nil {
		return nil, errors.Wrap(err, "error during unmarshalling http response")
	}
	return messageFilters, nil
}
