package snmpsimclient

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"io/ioutil"
	"strconv"
	"strings"
)

/*
ManagementClient is a client for communicating with the management api.
*/
type ManagementClient struct {
	client
}

/*
NewManagementClient creates a new ManagementClient.
*/
func NewManagementClient(baseURL string) (*ManagementClient, error) {
	if baseURL == "" {
		return nil, errors.New("invalid base url")
	}

	//if baseURL does not end with an "/" it has to be added to the string
	if lastChar := baseURL[len(baseURL)-1:]; lastChar != "/" {
		baseURL += "/"
	}
	clientData := clientData{baseURL: baseURL, resty: resty.New(), useAuth: false}
	newClient := client{&clientData}
	return &ManagementClient{newClient}, nil
}

/*
LABS
*/

/*
GetLabs returns a list of labs, optionally filtered.
*/
func (c *ManagementClient) GetLabs(filter map[string]string) (Labs, error) {
	if !c.isValid() {
		return nil, &NotValidError{}
	}

	response, err := c.request("GET", mgmtEndpointPath+"labs", "", nil, filter)
	if err != nil {
		return nil, errors.Wrap(err, "error during search labs request")
	}
	if response.StatusCode() != 200 {
		return nil, getHTTPError(response)
	}

	var labs Labs
	err = json.Unmarshal(response.Body(), &labs)
	if err != nil {
		return nil, errors.Wrap(err, "error during unmarshalling http response")
	}
	return labs, nil
}

/*
GetLab returns the lab with the given id.
*/
func (c *ManagementClient) GetLab(id int) (Lab, error) {
	if !c.isValid() {
		return Lab{}, &NotValidError{}
	}

	response, err := c.request("GET", mgmtEndpointPath+"labs/"+strconv.Itoa(id), "", nil, nil)
	if err != nil {
		return Lab{}, errors.Wrap(err, "error during get labs request")
	}
	if response.StatusCode() != 200 {
		return Lab{}, getHTTPError(response)
	}

	var lab Lab
	err = json.Unmarshal(response.Body(), &lab)
	if err != nil {
		return Lab{}, errors.Wrap(err, "error during unmarshalling http response")
	}
	return lab, nil
}

/*
CreateLab creates a new lab.
*/
func (c *ManagementClient) CreateLab(name string) (Lab, error) {
	return c.createLab(&name, nil)
}

/*
CreateLabWithTag creates a new lab tagged with the given tag.
*/
func (c *ManagementClient) CreateLabWithTag(name string, tagID int) (Lab, error) {
	return c.createLab(&name, &tagID)
}

func (c *ManagementClient) createLab(name *string, tagID *int) (Lab, error) {
	if !c.isValid() {
		return Lab{}, &NotValidError{}
	}

	if *name == "" {
		return Lab{}, errors.New("invalid name")
	}

	type requestParams struct {
		Name string `json:"name"`
	}

	params := requestParams{*name}
	jsonString, err := json.Marshal(params)
	if err != nil {
		return Lab{}, errors.Wrap(err, "error during marshal")
	}

	path := mgmtEndpointPath + "labs"
	if tagID != nil {
		path = mgmtEndpointPath + "tags/" + strconv.Itoa(*tagID) + "/lab"
	}

	response, err := c.request("POST", path, string(jsonString), nil, nil)

	if err != nil {
		return Lab{}, errors.Wrap(err, "error during add lab request")
	}
	if response.StatusCode() != 201 { //TODO: right error code?
		return Lab{}, getHTTPError(response)
	}

	var lab Lab
	err = json.Unmarshal(response.Body(), &lab)
	if err != nil {
		return Lab{}, errors.Wrap(err, "error during unmarshalling http response")
	}
	return lab, nil
}

/*
DeleteLab deletes the Lab with the given id.
*/
func (c *ManagementClient) DeleteLab(id int) error {
	if !c.isValid() {
		return &NotValidError{}
	}

	response, err := c.request("DELETE", mgmtEndpointPath+"labs/"+strconv.Itoa(id), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 204 {
		return getHTTPError(response)
	}
	return nil
}

/*
AddAgentToLab adds an Agent to a Lab.
*/
func (c *ManagementClient) AddAgentToLab(labID, agentID int) error {
	if !c.isValid() {
		return &NotValidError{}
	}

	response, err := c.request("PUT", mgmtEndpointPath+"labs/"+strconv.Itoa(labID)+"/agent/"+strconv.Itoa(agentID), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 200 {
		return getHTTPError(response)
	}
	return nil
}

/*
RemoveAgentFromLab removes an Agent from a Lab.
*/
func (c *ManagementClient) RemoveAgentFromLab(labID, agentID int) error {
	if !c.isValid() {
		return &NotValidError{}
	}

	response, err := c.request("DELETE", mgmtEndpointPath+"labs/"+strconv.Itoa(labID)+"/agent/"+strconv.Itoa(agentID), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 204 {
		return getHTTPError(response)
	}
	return nil
}

/*
SetLabPower activates or deactivates a lab.
*/
func (c *ManagementClient) SetLabPower(labID int, power bool) error {
	if !c.isValid() {
		return &NotValidError{}
	}

	var labPowerState string

	if power {
		labPowerState = "on"
	} else {
		labPowerState = "off"
	}

	response, err := c.request("PUT", mgmtEndpointPath+"labs/"+strconv.Itoa(labID)+"/power/"+labPowerState, "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 200 {
		return getHTTPError(response)
	}

	return nil
}

/*
AddTagToLab adds a tag to a lab.
*/
func (c *ManagementClient) AddTagToLab(labID, tagID int) error {
	if !c.isValid() {
		return &NotValidError{}
	}
	response, err := c.request("PUT", mgmtEndpointPath+"tags/"+strconv.Itoa(tagID)+"/lab/"+strconv.Itoa(labID), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 200 {
		return getHTTPError(response)
	}
	return nil
}

/*
RemoveTagFromLab removes a tag from a lab.
*/
func (c *ManagementClient) RemoveTagFromLab(labID, tagID int) error {
	if !c.isValid() {
		return &NotValidError{}
	}
	response, err := c.request("DELETE", mgmtEndpointPath+"tags/"+strconv.Itoa(tagID)+"/lab/"+strconv.Itoa(labID), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 200 {
		return getHTTPError(response)
	}
	return nil
}

/*
ENGINES
*/

/*
GetEngines returns a list of all engines.
*/
func (c *ManagementClient) GetEngines(filter map[string]string) (Engines, error) {
	if !c.isValid() {
		return nil, &NotValidError{}
	}

	response, err := c.request("GET", mgmtEndpointPath+"engines", "", nil, filter)
	if err != nil {
		return nil, errors.Wrap(err, "error during get engines request")
	}
	if response.StatusCode() != 200 {
		return nil, getHTTPError(response)
	}

	var engines Engines
	err = json.Unmarshal(response.Body(), &engines)
	if err != nil {
		return nil, errors.Wrap(err, "error during unmarshalling http response")
	}
	return engines, nil
}

/*
GetEngine returns the engine with the given id.
*/
func (c *ManagementClient) GetEngine(id int) (Engine, error) {
	if !c.isValid() {
		return Engine{}, &NotValidError{}
	}

	response, err := c.request("GET", mgmtEndpointPath+"engines/"+strconv.Itoa(id), "", nil, nil)
	if err != nil {
		return Engine{}, errors.Wrap(err, "error during get labs request")
	}
	if response.StatusCode() != 200 {
		return Engine{}, getHTTPError(response)
	}

	var engine Engine
	err = json.Unmarshal(response.Body(), &engine)
	if err != nil {
		return Engine{}, errors.Wrap(err, "error during unmarshalling http response")
	}
	return engine, nil
}

/*
CreateEngine creates a new engine.
*/
func (c *ManagementClient) CreateEngine(name, engineID string) (Engine, error) {
	return c.createEngine(&name, &engineID, nil)
}

/*
CreateEngineWithTag creates a new engine tagged with the given tag.
*/
func (c *ManagementClient) CreateEngineWithTag(name, engineID string, tagID int) (Engine, error) {
	return c.createEngine(&name, &engineID, &tagID)
}

func (c *ManagementClient) createEngine(name, engineID *string, tagID *int) (Engine, error) {
	if !c.isValid() {
		return Engine{}, &NotValidError{}
	}
	if *name == "" {
		return Engine{}, errors.New("invalid name")
	}

	//TODO: engine id should be removed! it should always be auto generated!
	if *engineID == "" {
		*engineID = "auto"
	}

	type requestParams struct {
		Name     string `json:"name"`
		EngineID string `json:"engine_id"`
	}

	params := requestParams{*name, *engineID}
	jsonString, err := json.Marshal(params)
	if err != nil {
		return Engine{}, errors.Wrap(err, "error during marshal")
	}

	path := mgmtEndpointPath + "engines"
	if tagID != nil {
		path = mgmtEndpointPath + "tags/" + strconv.Itoa(*tagID) + "/engine"
	}

	response, err := c.request("POST", path, string(jsonString), nil, nil)
	if err != nil {
		return Engine{}, errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 201 {
		return Engine{}, getHTTPError(response)
	}
	var engine Engine
	err = json.Unmarshal(response.Body(), &engine)
	if err != nil {
		return Engine{}, errors.Wrap(err, "error during unmarshalling http response")
	}
	return engine, nil
}

/*
DeleteEngine deletes the engine with the given id.
*/
func (c *ManagementClient) DeleteEngine(id int) error {
	if !c.isValid() {
		return &NotValidError{}
	}

	response, err := c.request("DELETE", mgmtEndpointPath+"engines/"+strconv.Itoa(id), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 204 {
		return getHTTPError(response)
	}
	return nil
}

/*
AddUserToEngine adds an User to an Engine
*/
func (c *ManagementClient) AddUserToEngine(engineID, userID int) error {
	if !c.isValid() {
		return &NotValidError{}
	}

	response, err := c.request("PUT", mgmtEndpointPath+"engines/"+strconv.Itoa(engineID)+"/user/"+strconv.Itoa(userID), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 200 {
		return getHTTPError(response)
	}

	return nil
}

/*
RemoveUserFromEngine removes an User from an Engine.
*/
func (c *ManagementClient) RemoveUserFromEngine(engineID, userID int) error {
	if !c.isValid() {
		return &NotValidError{}
	}

	response, err := c.request("DELETE", mgmtEndpointPath+"engines/"+strconv.Itoa(engineID)+"/user/"+strconv.Itoa(userID), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 204 {
		return getHTTPError(response)
	}

	return nil
}

/*
AddEndpointToEngine adds an Endpoint to an Engine.
*/
func (c *ManagementClient) AddEndpointToEngine(engineID, endpointID int) error {
	if !c.isValid() {
		return &NotValidError{}
	}

	response, err := c.request("PUT", mgmtEndpointPath+"engines/"+strconv.Itoa(engineID)+"/endpoint/"+strconv.Itoa(endpointID), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 200 {
		return getHTTPError(response)
	}

	return nil
}

/*
RemoveEndpointFromEngine removes an Endpoint from an Engine.
*/
func (c *ManagementClient) RemoveEndpointFromEngine(engineID, endpointID int) error {
	if !c.isValid() {
		return &NotValidError{}
	}

	response, err := c.request("DELETE", mgmtEndpointPath+"engines/"+strconv.Itoa(engineID)+"/endpoint/"+strconv.Itoa(endpointID), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 204 {
		return getHTTPError(response)
	}
	return nil
}

/*
AddTagToEngine adds a tag to a engine.
*/
func (c *ManagementClient) AddTagToEngine(engineID, tagID int) error {
	if !c.isValid() {
		return &NotValidError{}
	}
	response, err := c.request("PUT", mgmtEndpointPath+"tags/"+strconv.Itoa(tagID)+"/engine/"+strconv.Itoa(engineID), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 200 {
		return getHTTPError(response)
	}
	return nil
}

/*
RemoveTagFromEngine removes a tag from a engine.
*/
func (c *ManagementClient) RemoveTagFromEngine(engineID, tagID int) error {
	if !c.isValid() {
		return &NotValidError{}
	}
	response, err := c.request("DELETE", mgmtEndpointPath+"tags/"+strconv.Itoa(tagID)+"/engine/"+strconv.Itoa(engineID), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 200 {
		return getHTTPError(response)
	}
	return nil
}

/*
AGENTS
*/

/*
GetAgents returns a list of agents, optionally filtered.
*/
func (c *ManagementClient) GetAgents(filters map[string]string) (Agents, error) {
	if !c.isValid() {
		return nil, &NotValidError{}
	}

	response, err := c.request("GET", mgmtEndpointPath+"agents", "", nil, filters)
	if err != nil {
		return nil, errors.Wrap(err, "error during get agents request")
	}
	if response.StatusCode() != 200 {
		return nil, getHTTPError(response)
	}

	var agents Agents
	err = json.Unmarshal(response.Body(), &agents)
	if err != nil {
		return nil, errors.Wrap(err, "error during unmarshalling http response")
	}
	return agents, nil
}

/*
GetAgent returns the agent with the given id.
*/
func (c *ManagementClient) GetAgent(id int) (Agent, error) {
	if !c.isValid() {
		return Agent{}, &NotValidError{}
	}

	response, err := c.request("GET", mgmtEndpointPath+"agents/"+strconv.Itoa(id), "", nil, nil)
	if err != nil {
		return Agent{}, errors.Wrap(err, "error during get labs request")
	}
	if response.StatusCode() != 200 {
		return Agent{}, getHTTPError(response)
	}

	var agent Agent
	err = json.Unmarshal(response.Body(), &agent)
	if err != nil {
		return Agent{}, errors.Wrap(err, "error during unmarshalling http response")
	}
	return agent, nil
}

/*
CreateAgent creates a new agent.
*/
func (c *ManagementClient) CreateAgent(name, dataDir string) (Agent, error) {
	return c.createAgent(&name, &dataDir, nil)
}

/*
CreateAgentWithTag creates a new agent tagged with the given tag.
*/
func (c *ManagementClient) CreateAgentWithTag(name, dataDir string, tagID int) (Agent, error) {
	return c.createAgent(&name, &dataDir, &tagID)
}

func (c *ManagementClient) createAgent(name, dataDir *string, tagID *int) (Agent, error) {
	if !c.isValid() {
		return Agent{}, &NotValidError{}
	}

	if *name == "" {
		return Agent{}, errors.New("invalid name")
	}

	if *dataDir == "" {
		*dataDir = "."
	}

	type requestParams struct {
		Name    string `json:"name"`
		DataDir string `json:"data_dir"`
	}

	params := requestParams{*name, *dataDir}
	jsonString, err := json.Marshal(params)
	if err != nil {
		return Agent{}, errors.Wrap(err, "error during marshal")
	}

	path := mgmtEndpointPath + "agents"
	if tagID != nil {
		path = mgmtEndpointPath + "tags/" + strconv.Itoa(*tagID) + "/agent"
	}

	response, err := c.request("POST", path, string(jsonString), nil, nil)
	if err != nil {
		return Agent{}, errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 201 {
		return Agent{}, getHTTPError(response)
	}

	var agent Agent
	err = json.Unmarshal(response.Body(), &agent)
	if err != nil {
		return Agent{}, errors.Wrap(err, "error during unmarshalling http response")
	}
	return agent, nil
}

/*
DeleteAgent deletes the agent with the given id.
*/
func (c *ManagementClient) DeleteAgent(id int) error {
	if !c.isValid() {
		return &NotValidError{}
	}

	response, err := c.request("DELETE", mgmtEndpointPath+"agents/"+strconv.Itoa(id), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 204 {
		return getHTTPError(response)
	}
	return nil
}

/*
AddEngineToAgent adds an Engine to an Agent.
*/
func (c *ManagementClient) AddEngineToAgent(agentID, engineID int) error {
	if !c.isValid() {
		return &NotValidError{}
	}

	response, err := c.request("PUT", mgmtEndpointPath+"agents/"+strconv.Itoa(agentID)+"/engine/"+strconv.Itoa(engineID), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 200 {
		return getHTTPError(response)
	}

	return nil
}

/*
RemoveEngineFromAgent removes an Engine from an Agent.
*/
func (c *ManagementClient) RemoveEngineFromAgent(agentID, engineID int) error {
	if !c.isValid() {
		return &NotValidError{}
	}

	response, err := c.request("DELETE", mgmtEndpointPath+"agents/"+strconv.Itoa(agentID)+"/engine/"+strconv.Itoa(engineID), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 204 {
		return getHTTPError(response)
	}

	return nil
}

/*
AddSelectorToAgent adds a Selector to an Agent.
*/
func (c *ManagementClient) AddSelectorToAgent(agentID, selectorID int) (Agent, error) {
	//TODO: Not implemented yet!!! Selectors are not implemented yet, see SELECTORS section
	return Agent{}, errors.New("Not implemented yet")
}

/*
RemoveSelectorFromAgent removes a Selector from an Agent.
*/
func (c *ManagementClient) RemoveSelectorFromAgent(agentID, selectorID int) (Agent, error) {
	//TODO: Not implemented yet!!! Selectors are not implemented yet, see SELECTORS section
	return Agent{}, errors.New("Not implemented yet")
}

/*
AddTagToAgent adds a tag to a agent.
*/
func (c *ManagementClient) AddTagToAgent(agentID, tagID int) error {
	if !c.isValid() {
		return &NotValidError{}
	}
	response, err := c.request("PUT", mgmtEndpointPath+"tags/"+strconv.Itoa(tagID)+"/agent/"+strconv.Itoa(agentID), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 200 {
		return getHTTPError(response)
	}
	return nil
}

/*
RemoveTagFromAgent removes a tag from a agent.
*/
func (c *ManagementClient) RemoveTagFromAgent(agentID, tagID int) error {
	if !c.isValid() {
		return &NotValidError{}
	}
	response, err := c.request("DELETE", mgmtEndpointPath+"tags/"+strconv.Itoa(tagID)+"/agent/"+strconv.Itoa(agentID), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 200 {
		return getHTTPError(response)
	}
	return nil
}

/*
ENDPOINTS
*/

/*
GetEndpoints returns a list of endpoints, optionally filtered.
*/
func (c *ManagementClient) GetEndpoints(filters map[string]string) (Endpoints, error) {
	if !c.isValid() {
		return nil, &NotValidError{}
	}

	response, err := c.request("GET", mgmtEndpointPath+"endpoints", "", nil, filters)
	if err != nil {
		return nil, errors.Wrap(err, "error during get endpoints request")
	}
	if response.StatusCode() != 200 {
		return nil, getHTTPError(response)
	}

	var endpoints Endpoints
	err = json.Unmarshal(response.Body(), &endpoints)
	if err != nil {
		return nil, errors.Wrap(err, "error during unmarshalling http response")
	}
	return endpoints, nil
}

/*
GetEndpoint returns the endpoint with the given id.
*/
func (c *ManagementClient) GetEndpoint(id int) (Endpoint, error) {
	if !c.isValid() {
		return Endpoint{}, &NotValidError{}
	}

	response, err := c.request("GET", mgmtEndpointPath+"endpoints/"+strconv.Itoa(id), "", nil, nil)
	if err != nil {
		return Endpoint{}, errors.Wrap(err, "error during get labs request")
	}
	if response.StatusCode() != 200 {
		return Endpoint{}, getHTTPError(response)
	}

	var endpoint Endpoint
	err = json.Unmarshal(response.Body(), &endpoint)
	if err != nil {
		return Endpoint{}, errors.Wrap(err, "error during unmarshalling http response")
	}
	return endpoint, nil
}

/*
CreateEndpoint creates a new endpoint.
*/
func (c *ManagementClient) CreateEndpoint(name, address, protocol string) (Endpoint, error) {
	return c.createEndpoint(&name, &address, &protocol, nil)
}

/*
CreateEndpointWithTag creates a new endpoint tagged with the given tag.
*/
func (c *ManagementClient) CreateEndpointWithTag(name, address, protocol string, tagID int) (Endpoint, error) {
	return c.createEndpoint(&name, &address, &protocol, &tagID)
}

func (c *ManagementClient) createEndpoint(name, address, protocol *string, tagID *int) (Endpoint, error) {
	if !c.isValid() {
		return Endpoint{}, &NotValidError{}
	}

	if *name == "" {
		return Endpoint{}, errors.New("invalid name")
	}

	if *protocol == "" {
		*protocol = "udpv4"
	}

	if *address == "" {
		return Endpoint{}, errors.New("invalid address")
	}

	type requestParams struct {
		Name     string `json:"name"`
		Address  string `json:"address"`
		Protocol string `json:"protocol"`
	}

	params := requestParams{*name, *address, *protocol}
	jsonString, err := json.Marshal(params)
	if err != nil {
		return Endpoint{}, errors.Wrap(err, "error during marshal")
	}

	path := mgmtEndpointPath + "endpoints"
	if tagID != nil {
		path = mgmtEndpointPath + "tags/" + strconv.Itoa(*tagID) + "/endpoint"
	}

	response, err := c.request("POST", path, string(jsonString), nil, nil)

	if err != nil {
		return Endpoint{}, errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 201 {
		return Endpoint{}, getHTTPError(response)
	}

	var endpoint Endpoint
	err = json.Unmarshal(response.Body(), &endpoint)
	if err != nil {
		return Endpoint{}, errors.Wrap(err, "error during unmarshalling http response")
	}
	return endpoint, nil
}

/*
DeleteEndpoint the Lab with the given id.
*/
func (c *ManagementClient) DeleteEndpoint(id int) error {
	if !c.isValid() {
		return &NotValidError{}
	}

	response, err := c.request("DELETE", mgmtEndpointPath+"endpoints/"+strconv.Itoa(id), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 204 {
		return getHTTPError(response)
	}
	return nil
}

/*
AddTagToEndpoint adds a tag to a endpoint.
*/
func (c *ManagementClient) AddTagToEndpoint(endpointID, tagID int) error {
	if !c.isValid() {
		return &NotValidError{}
	}
	response, err := c.request("PUT", mgmtEndpointPath+"tags/"+strconv.Itoa(tagID)+"/endpoint/"+strconv.Itoa(endpointID), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 200 {
		return getHTTPError(response)
	}
	return nil
}

/*
RemoveTagFromEndpoint removes a tag from a endpoint.
*/
func (c *ManagementClient) RemoveTagFromEndpoint(endpointID, tagID int) error {
	if !c.isValid() {
		return &NotValidError{}
	}
	response, err := c.request("DELETE", mgmtEndpointPath+"tags/"+strconv.Itoa(tagID)+"/endpoint/"+strconv.Itoa(endpointID), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 200 {
		return getHTTPError(response)
	}
	return nil
}

/*
RECORD FILES
*/

/*
GetRecordFiles returns a list of all record files.
*/
func (c *ManagementClient) GetRecordFiles() (Recordings, error) {
	if !c.isValid() {
		return nil, &NotValidError{}
	}

	response, err := c.request("GET", mgmtEndpointPath+"recordings", "", nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error during get labs request")
	}
	if response.StatusCode() != 200 {
		return nil, getHTTPError(response)
	}

	var recordings Recordings
	err = json.Unmarshal(response.Body(), &recordings)
	if err != nil {
		return nil, errors.Wrap(err, "error during unmarshalling http response")
	}
	return recordings, nil
}

/*
UploadRecordFile uploads the given record file to the api and saves it at the given remote path inside of the data dir.
*/
func (c *ManagementClient) UploadRecordFile(localPath, remotePath string) error {
	localPath = strings.TrimSpace(localPath)
	if !strings.HasSuffix(localPath, ".snmprec") {
		return errors.New("file is not an snmprec file")
	}
	b, err := ioutil.ReadFile(localPath)
	if err != nil {
		return errors.Wrap(err, "error while reading file")
	}
	s := string(b)
	return c.UploadRecordFileString(&s, remotePath)
}

/*
UploadRecordFileString uploads the given record data to the api and saves it as a .snmprec file at the given remote path inside of the data dir.
*/
func (c *ManagementClient) UploadRecordFileString(recordContents *string, remotePath string) error {
	headerMap := make(map[string]string)
	headerMap["Content-Type"] = "text/plain"
	response, err := c.request("POST", mgmtEndpointPath+"recordings/"+remotePath, *recordContents, headerMap, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 204 {
		return getHTTPError(response)
	}
	return nil
}

/*
DeleteRecordFile deletes the record file at the given path.
*/
func (c *ManagementClient) DeleteRecordFile(remotePath string) error {
	remotePath = strings.TrimSpace(remotePath)
	if !strings.HasSuffix(remotePath, ".snmprec") {
		return errors.New("file is not an snmprec file")
	}
	headerMap := make(map[string]string)
	headerMap["Content-Type"] = "text/plain"
	response, err := c.request("DELETE", mgmtEndpointPath+"recordings/"+remotePath, "", headerMap, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 204 {
		return getHTTPError(response)
	}
	return nil
}

/*
GetRecordFile returns the record file at the given path.
*/
func (c *ManagementClient) GetRecordFile(remotePath string) (string, error) {
	remotePath = strings.TrimSpace(remotePath)
	if !strings.HasSuffix(remotePath, ".snmprec") {
		return "", errors.New("file is not an snmprec file")
	}
	headerMap := make(map[string]string)
	headerMap["Content-Type"] = "text/plain"
	response, err := c.request("GET", mgmtEndpointPath+"recordings/"+remotePath, "", headerMap, nil)
	if err != nil {
		return "", errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 200 {
		return "", getHTTPError(response)
	}
	return string(response.Body()), nil
}

/*
USERS
*/

/*
CreateUser creates a new user.
*/
func (c *ManagementClient) CreateUser(user, name, authKey, authProto, privKey, privProto string) (User, error) {
	return c.createUser(&user, &name, &authKey, &authProto, &privKey, &privProto, nil)
}

/*
CreateUserWithTag creates a new user tagged with the given tag.
*/
func (c *ManagementClient) CreateUserWithTag(user, name, authKey, authProto, privKey, privProto string, tagID int) (User, error) {
	return c.createUser(&user, &name, &authKey, &authProto, &privKey, &privProto, &tagID)
}

func (c *ManagementClient) createUser(user, name, authKey, authProto, privKey, privProto *string, tagID *int) (User, error) {
	if !c.isValid() {
		return User{}, &NotValidError{}
	}

	if *name == "" {
		return User{}, errors.New("invalid name")
	}

	if *user == "" {
		return User{}, errors.New("invalid user")
	}

	if *authProto == "" {
		*authProto = "none"
	}
	if *privProto == "" {
		*privProto = "none"
	}

	type requestParams struct {
		User      string  `json:"user"`
		Name      string  `json:"name"`
		AuthKey   *string `json:"auth_key"`
		AuthProto string  `json:"auth_proto"`
		PrivKey   *string `json:"priv_key"`
		PrivProto string  `json:"priv_proto"`
	}

	params := requestParams{
		User:      *user,
		Name:      *name,
		AuthProto: *authProto,
		PrivProto: *privProto,
	}

	if *authKey == "" {
		params.AuthKey = nil
	} else {
		params.AuthKey = authKey
	}

	if *privKey == "" {
		params.PrivKey = nil
	} else {
		params.PrivKey = privKey
	}

	jsonString, err := json.Marshal(params)
	if err != nil {
		return User{}, errors.Wrap(err, "error during marshal")
	}

	path := mgmtEndpointPath + "users"
	if tagID != nil {
		path = mgmtEndpointPath + "tags/" + strconv.Itoa(*tagID) + "/user"
	}

	response, err := c.request("POST", path, string(jsonString), nil, nil)
	if err != nil {
		return User{}, errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 201 {
		return User{}, getHTTPError(response)
	}

	var newUser User
	err = json.Unmarshal(response.Body(), &newUser)
	if err != nil {
		return User{}, errors.Wrap(err, "error during unmarshalling http response")
	}
	return newUser, nil
}

/*
GetUsers returns a list of users, optionally filtered.
*/
func (c *ManagementClient) GetUsers(filters map[string]string) (Users, error) {
	if !c.isValid() {
		return nil, &NotValidError{}
	}

	response, err := c.request("GET", mgmtEndpointPath+"users", "", nil, filters)
	if err != nil {
		return nil, errors.Wrap(err, "error during get users request")
	}
	if response.StatusCode() != 200 {
		return nil, getHTTPError(response)
	}

	var users Users
	err = json.Unmarshal(response.Body(), &users)
	if err != nil {
		return nil, errors.Wrap(err, "error during unmarshalling http response")
	}
	return users, nil
}

/*
GetUser returns the user with the given id.
*/
func (c *ManagementClient) GetUser(id int) (User, error) {
	if !c.isValid() {
		return User{}, &NotValidError{}
	}

	response, err := c.request("GET", mgmtEndpointPath+"users/"+strconv.Itoa(id), "", nil, nil)
	if err != nil {
		return User{}, errors.Wrap(err, "error during get labs request")
	}
	if response.StatusCode() != 200 {
		return User{}, getHTTPError(response)
	}

	var user User
	err = json.Unmarshal(response.Body(), &user)
	if err != nil {
		return User{}, errors.Wrap(err, "error during unmarshalling http response")
	}
	return user, nil
}

/*
DeleteUser deletes the user with the given id.
*/
func (c *ManagementClient) DeleteUser(id int) error {
	if !c.isValid() {
		return &NotValidError{}
	}

	response, err := c.request("DELETE", mgmtEndpointPath+"users/"+strconv.Itoa(id), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 204 {
		return getHTTPError(response)
	}
	return nil
}

/*
AddTagToUser adds a tag to a user.
*/
func (c *ManagementClient) AddTagToUser(userID, tagID int) error {
	if !c.isValid() {
		return &NotValidError{}
	}
	response, err := c.request("PUT", mgmtEndpointPath+"tags/"+strconv.Itoa(tagID)+"/user/"+strconv.Itoa(userID), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 200 {
		return getHTTPError(response)
	}
	return nil
}

/*
RemoveTagFromUser removes a tag from a user.
*/
func (c *ManagementClient) RemoveTagFromUser(userID, tagID int) error {
	if !c.isValid() {
		return &NotValidError{}
	}
	response, err := c.request("DELETE", mgmtEndpointPath+"tags/"+strconv.Itoa(tagID)+"/user/"+strconv.Itoa(userID), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 200 {
		return getHTTPError(response)
	}
	return nil
}

/*
SELECTORS
*/
//TODO: Not implemented yet in api, there is always a default selector for snmp v2c and v3. This has to be implemented when its possible to configure own selectors.

/*
CreateSelector creates a new selector.
*/
func (c *ManagementClient) CreateSelector(comment, template string) (Selector, error) {
	//TODO: Not implemented yet!!! This also isnt implemented in the api so we have to wait until its ready.
	return Selector{}, errors.New("Not implemented yet")
}

/*
GetSelectors returns a list of all selectors.
*/
func (c *ManagementClient) GetSelectors() (Selectors, error) {
	//TODO: Not implemented yet!!! This also isnt implemented in the api so we have to wait until its ready.
	return nil, errors.New("Not implemented yet")
}

/*
GetSelector returns the selector with the given id.
*/
func (c *ManagementClient) GetSelector(id int) (Selector, error) {
	//TODO: Not implemented yet!!! This also isnt implemented in the api so we have to wait until its ready.
	return Selector{}, errors.New("Not implemented yet")
}

/*
DeleteSelector deletes the selector with the given id.
*/
func (c *ManagementClient) DeleteSelector(id int) error {
	//TODO: Not implemented yet!!! This also isnt implemented in the api so we have to wait until its ready.
	return errors.New("Not implemented yet")
}

/*
 TAGS
*/

/*
CreateTag creates a new tag.
*/
func (c *ManagementClient) CreateTag(name, description string) (Tag, error) {
	if !c.isValid() {
		return Tag{}, &NotValidError{}
	}
	if name == "" {
		return Tag{}, errors.New("invalid name")
	}

	type requestParams struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	params := requestParams{name, description}
	jsonString, err := json.Marshal(params)
	if err != nil {
		return Tag{}, errors.Wrap(err, "error during marshal")
	}

	response, err := c.request("POST", mgmtEndpointPath+"tags", string(jsonString), nil, nil)
	if err != nil {
		return Tag{}, errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 201 {
		return Tag{}, getHTTPError(response)
	}

	var tag Tag
	err = json.Unmarshal(response.Body(), &tag)
	if err != nil {
		return Tag{}, errors.Wrap(err, "error during unmarshalling http response")
	}
	return tag, nil
}

/*
GetTag returns the lab with the given id.
*/
func (c *ManagementClient) GetTag(id int) (Tag, error) {
	if !c.isValid() {
		return Tag{}, &NotValidError{}
	}

	response, err := c.request("GET", mgmtEndpointPath+"tags/"+strconv.Itoa(id), "", nil, nil)
	if err != nil {
		return Tag{}, errors.Wrap(err, "error during get tags request")
	}
	if response.StatusCode() != 200 {
		return Tag{}, getHTTPError(response)
	}

	var tag Tag
	err = json.Unmarshal(response.Body(), &tag)
	if err != nil {
		return Tag{}, errors.Wrap(err, "error during unmarshalling http response")
	}
	return tag, nil
}

/*
GetTags returns a list of users, optionally filtered.
*/
func (c *ManagementClient) GetTags(filters map[string]string) (Tags, error) {
	if !c.isValid() {
		return nil, &NotValidError{}
	}

	response, err := c.request("GET", mgmtEndpointPath+"tags", "", nil, filters)
	if err != nil {
		return nil, errors.Wrap(err, "error during get users request")
	}
	if response.StatusCode() != 200 {
		return nil, getHTTPError(response)
	}

	var tags Tags
	err = json.Unmarshal(response.Body(), &tags)
	if err != nil {
		return nil, errors.Wrap(err, "error during unmarshalling http response")
	}
	return tags, nil
}

/*
DeleteTag deletes the tag with the given id.
*/
func (c *ManagementClient) DeleteTag(id int) error {
	if !c.isValid() {
		return &NotValidError{}
	}

	response, err := c.request("DELETE", mgmtEndpointPath+"tags/"+strconv.Itoa(id), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 204 {
		return getHTTPError(response)
	}
	return nil
}

/*
DeleteAllObjectsWithTag deletes all objects with the given tag.
*/
func (c *ManagementClient) DeleteAllObjectsWithTag(tagID int) (Tag, error) {
	if !c.isValid() {
		return Tag{}, &NotValidError{}
	}

	response, err := c.request("DELETE", mgmtEndpointPath+"tags/"+strconv.Itoa(tagID)+"/objects", "", nil, nil)
	if err != nil {
		return Tag{}, errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 200 {
		return Tag{}, getHTTPError(response)
	}

	var tag Tag
	err = json.Unmarshal(response.Body(), &tag)
	if err != nil {
		return Tag{}, errors.Wrap(err, "error during unmarshalling http response")
	}

	return tag, nil
}
