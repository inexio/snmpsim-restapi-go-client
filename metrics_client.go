package snmpsimclient

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
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
GetProcessesMetrics returns process metrics.
*/
func (c *MetricsClient) GetProcessesMetrics(filters map[string]string) (ProcessesMetrics, error) {
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
GetPacketMetrics returns packet metrics.
*/
func (c *MetricsClient) GetPacketMetrics(filters map[string]string) (PacketMetrics, error) {
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
func (c *MetricsClient) GetPacketFilters() (map[string]string, error) {
	response, err := c.request("GET", metricsEndpointPath+"activity/packets/filters", "", nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 200 {
		return nil, getHttpError(response)
	}

	var packetFilters map[string]string
	err = json.Unmarshal(response.Body(), &packetFilters)
	if err != nil {
		return nil, errors.Wrap(err, "error during unmarshalling http response")
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
GetMessageMetrics returns message metrics.
*/
func (c *MetricsClient) GetMessageMetrics(filters map[string]string) (MessageMetrics, error) {
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
func (c *MetricsClient) GetMessageFilters() (map[string]string, error) {
	response, err := c.request("GET", metricsEndpointPath+"activity/messages/filters", "", nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 200 {
		return nil, getHttpError(response)
	}

	var messageFilters map[string]string
	err = json.Unmarshal(response.Body(), &messageFilters)
	if err != nil {
		return nil, errors.Wrap(err, "error during unmarshalling http response")
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
