package snmpsim_restapi_client

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
func NewManagementClient(baseUrl string) (*ManagementClient, error) {
	if baseUrl == "" {
		return nil, errors.New("invalid base url")
	}

	//if baseUrl does not end with an "/" it has to be added to the string
	if lastChar := baseUrl[len(baseUrl)-1:]; lastChar != "/" {
		baseUrl += "/"
	}
	clientData := clientData{baseUrl: baseUrl, resty: resty.New(), useAuth: false}
	newClient := client{&clientData}
	return &ManagementClient{newClient}, nil
}

/*
LABS
*/

/*
GetLabs returns a list of all labs.
*/
func (c *ManagementClient) GetLabs() (Labs, error) {
	if !c.isValid() {
		return nil, &NotValidError{}
	}

	response, err := c.request("GET", MGMT_ENDPOINT_PATH+"labs", "", nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error during get labs request")
	}
	if response.StatusCode() != 200 {
		return nil, getHttpError(response)
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

	response, err := c.request("GET", MGMT_ENDPOINT_PATH+"labs/"+strconv.Itoa(id), "", nil, nil)
	if err != nil {
		return Lab{}, errors.Wrap(err, "error during get labs request")
	}
	if response.StatusCode() != 200 {
		return Lab{}, getHttpError(response)
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
	if !c.isValid() {
		return Lab{}, &NotValidError{}
	}

	if name == "" {
		return Lab{}, errors.New("invalid name")
	}

	response, err := c.request("POST", MGMT_ENDPOINT_PATH+"labs", `{"name":"`+name+`"}`, nil, nil)

	if err != nil {
		return Lab{}, errors.Wrap(err, "error during add lab request")
	}
	if response.StatusCode() != 201 { //TODO: right error code?
		return Lab{}, getHttpError(response)
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

	response, err := c.request("DELETE", MGMT_ENDPOINT_PATH+"labs/"+strconv.Itoa(id), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 204 {
		return getHttpError(response)
	}
	return nil
}

/*
AddAgentToLab adds an Agent to a Lab.
*/
func (c *ManagementClient) AddAgentToLab(labId, agentId int) error {
	if !c.isValid() {
		return &NotValidError{}
	}

	response, err := c.request("PUT", MGMT_ENDPOINT_PATH+"labs/"+strconv.Itoa(labId)+"/agent/"+strconv.Itoa(agentId), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 200 {
		return getHttpError(response)
	}
	return nil
}

/*
RemoveAgentFromLab removes an Agent from a Lab.
*/
func (c *ManagementClient) RemoveAgentFromLab(labId, agentId int) error {
	if !c.isValid() {
		return &NotValidError{}
	}

	response, err := c.request("DELETE", MGMT_ENDPOINT_PATH+"labs/"+strconv.Itoa(labId)+"/agent/"+strconv.Itoa(agentId), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 204 {
		return getHttpError(response)
	}
	return nil
}

/*
SetLabPower activates or deactivates a lab.
*/
func (c *ManagementClient) SetLabPower(labId int, power bool) error {
	if !c.isValid() {
		return &NotValidError{}
	}

	var labPowerState string

	if power {
		labPowerState = "on"
	} else {
		labPowerState = "off"
	}

	response, err := c.request("PUT", MGMT_ENDPOINT_PATH+"labs/"+strconv.Itoa(labId)+"/power/"+labPowerState, "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 200 {
		return getHttpError(response)
	}

	return nil
}

/*
ENGINES
*/

/*
GetEngines returns a list of all engines.
*/
func (c *ManagementClient) GetEngines() (Engines, error) {
	if !c.isValid() {
		return nil, &NotValidError{}
	}

	response, err := c.request("GET", MGMT_ENDPOINT_PATH+"engines", "", nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error during get labs request")
	}
	if response.StatusCode() != 200 {
		return nil, getHttpError(response)
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

	response, err := c.request("GET", MGMT_ENDPOINT_PATH+"engines/"+strconv.Itoa(id), "", nil, nil)
	if err != nil {
		return Engine{}, errors.Wrap(err, "error during get labs request")
	}
	if response.StatusCode() != 200 {
		return Engine{}, getHttpError(response)
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
func (c *ManagementClient) CreateEngine(name, engineId string) (Engine, error) {
	if !c.isValid() {
		return Engine{}, &NotValidError{}
	}
	if name == "" {
		return Engine{}, errors.New("invalid name")
	}

	//TODO: engine id should be removed! it should always be auto generated!
	if engineId == "" {
		engineId = "auto"
	}

	type requestParams struct {
		Name     string `json:"name"`
		EngineId string `json:"engine_id"`
	}

	params := requestParams{name, engineId}
	jsonString, err := json.Marshal(params)
	if err != nil {
		return Engine{}, errors.Wrap(err, "error during marshal")
	}

	response, err := c.request("POST", MGMT_ENDPOINT_PATH+"engines", string(jsonString), nil, nil)
	if err != nil {
		return Engine{}, errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 201 {
		return Engine{}, getHttpError(response)
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

	response, err := c.request("DELETE", MGMT_ENDPOINT_PATH+"engines/"+strconv.Itoa(id), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 204 {
		return getHttpError(response)
	}
	return nil
}

/*
AddUserToEngine adds an User to an Engine
*/
func (c *ManagementClient) AddUserToEngine(engineId, userId int) error {
	if !c.isValid() {
		return &NotValidError{}
	}

	response, err := c.request("PUT", MGMT_ENDPOINT_PATH+"engines/"+strconv.Itoa(engineId)+"/user/"+strconv.Itoa(userId), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 200 {
		return getHttpError(response)
	}

	return nil
}

/*
RemoveUserFromEngine removes an User from an Engine.
*/
func (c *ManagementClient) RemoveUserFromEngine(engineId, userId int) error {
	if !c.isValid() {
		return &NotValidError{}
	}

	response, err := c.request("DELETE", MGMT_ENDPOINT_PATH+"engines/"+strconv.Itoa(engineId)+"/user/"+strconv.Itoa(userId), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 204 {
		return getHttpError(response)
	}

	return nil
}

/*
AddEndpointToEngine adds an Endpoint to an Engine.
*/
func (c *ManagementClient) AddEndpointToEngine(engineId, endpointId int) error {
	if !c.isValid() {
		return &NotValidError{}
	}

	response, err := c.request("PUT", MGMT_ENDPOINT_PATH+"engines/"+strconv.Itoa(engineId)+"/endpoint/"+strconv.Itoa(endpointId), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 200 {
		return getHttpError(response)
	}

	return nil
}

/*
RemoveEndpointFromEngine removes an Endpoint from an Engine.
*/
func (c *ManagementClient) RemoveEndpointFromEngine(engineId, endpointId int) error {
	if !c.isValid() {
		return &NotValidError{}
	}

	response, err := c.request("DELETE", MGMT_ENDPOINT_PATH+"engines/"+strconv.Itoa(engineId)+"/endpoint/"+strconv.Itoa(endpointId), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 204 {
		return getHttpError(response)
	}
	return nil
}

/*
AGENTS
*/

/*
GetAgents returns a list of all agents.
*/
func (c *ManagementClient) GetAgents() (Agents, error) {
	if !c.isValid() {
		return nil, &NotValidError{}
	}

	response, err := c.request("GET", MGMT_ENDPOINT_PATH+"agents", "", nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error during get labs request")
	}
	if response.StatusCode() != 200 {
		return nil, getHttpError(response)
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

	response, err := c.request("GET", MGMT_ENDPOINT_PATH+"agents/"+strconv.Itoa(id), "", nil, nil)
	if err != nil {
		return Agent{}, errors.Wrap(err, "error during get labs request")
	}
	if response.StatusCode() != 200 {
		return Agent{}, getHttpError(response)
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
	if !c.isValid() {
		return Agent{}, &NotValidError{}
	}

	if name == "" {
		return Agent{}, errors.New("invalid name")
	}

	if dataDir == "" {
		dataDir = "."
	}

	type requestParams struct {
		Name    string `json:"name"`
		DataDir string `json:"data_dir"`
	}

	params := requestParams{name, dataDir}
	jsonString, err := json.Marshal(params)
	if err != nil {
		return Agent{}, errors.Wrap(err, "error during marshal")
	}

	response, err := c.request("POST", MGMT_ENDPOINT_PATH+"agents", string(jsonString), nil, nil)
	if err != nil {
		return Agent{}, errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 201 {
		return Agent{}, getHttpError(response)
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

	response, err := c.request("DELETE", MGMT_ENDPOINT_PATH+"agents/"+strconv.Itoa(id), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 204 {
		return getHttpError(response)
	}
	return nil
}

/*
AddEngineToAgent adds an Engine to an Agent.
*/
func (c *ManagementClient) AddEngineToAgent(agentId, engineId int) error {
	if !c.isValid() {
		return &NotValidError{}
	}

	response, err := c.request("PUT", MGMT_ENDPOINT_PATH+"agents/"+strconv.Itoa(agentId)+"/engine/"+strconv.Itoa(engineId), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 200 {
		return getHttpError(response)
	}

	return nil
}

/*
RemoveEngineFromAgent removes an Engine from an Agent.
*/
func (c *ManagementClient) RemoveEngineFromAgent(agentId, engineId int) error {
	if !c.isValid() {
		return &NotValidError{}
	}

	response, err := c.request("DELETE", MGMT_ENDPOINT_PATH+"agents/"+strconv.Itoa(agentId)+"/engine/"+strconv.Itoa(engineId), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}

	if response.StatusCode() != 204 {
		return getHttpError(response)
	}

	return nil
}

/*
AddSelectorToAgent adds a Selector to an Agent.
*/
func (c *ManagementClient) AddSelectorToAgent(agentId, selectorId int) (Agent, error) {
	//TODO: Not implemented yet!!! Selectors are not implemented yet, see SELECTORS section
	return Agent{}, errors.New("Not implemented yet!")
}

/*
RemoveSelectorFromAgent removes a Selector from an Agent.
*/
func (c *ManagementClient) RemoveSelectorFromAgent(agentId, selectorId int) (Agent, error) {
	//TODO: Not implemented yet!!! Selectors are not implemented yet, see SELECTORS section
	return Agent{}, errors.New("Not implemented yet!")
}

/*
ENDPOINTS
*/

/*
GetEndpoints returns a list of all endpoints.
*/
func (c *ManagementClient) GetEndpoints() (Endpoints, error) {
	if !c.isValid() {
		return nil, &NotValidError{}
	}

	response, err := c.request("GET", MGMT_ENDPOINT_PATH+"endpoints", "", nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error during get labs request")
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
GetEndpoint returns the endpoint with the given id.
*/
func (c *ManagementClient) GetEndpoint(id int) (Endpoint, error) {
	if !c.isValid() {
		return Endpoint{}, &NotValidError{}
	}

	response, err := c.request("GET", MGMT_ENDPOINT_PATH+"endpoints/"+strconv.Itoa(id), "", nil, nil)
	if err != nil {
		return Endpoint{}, errors.Wrap(err, "error during get labs request")
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
CreateEndpoint creates a new endpoint.
*/
func (c *ManagementClient) CreateEndpoint(name, address, protocol string) (Endpoint, error) {
	if !c.isValid() {
		return Endpoint{}, &NotValidError{}
	}

	if name == "" {
		return Endpoint{}, errors.New("invalid name")
	}

	if protocol == "" {
		protocol = "udpv4"
	}

	if address == "" {
		return Endpoint{}, errors.New("invalid address")
	}

	type requestParams struct {
		Name     string `json:"name"`
		Address  string `json:"address"`
		Protocol string `json:"protocol"`
	}

	params := requestParams{name, address, protocol}
	jsonString, err := json.Marshal(params)
	if err != nil {
		return Endpoint{}, errors.Wrap(err, "error during marshal")
	}

	response, err := c.request("POST", MGMT_ENDPOINT_PATH+"endpoints", string(jsonString), nil, nil)

	if err != nil {
		return Endpoint{}, errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 201 {
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
DeleteEndpoint the Lab with the given id.
*/
func (c *ManagementClient) DeleteEndpoint(id int) error {
	if !c.isValid() {
		return &NotValidError{}
	}

	response, err := c.request("DELETE", MGMT_ENDPOINT_PATH+"endpoints/"+strconv.Itoa(id), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 204 {
		return getHttpError(response)
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

	response, err := c.request("GET", MGMT_ENDPOINT_PATH+"recordings", "", nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error during get labs request")
	}
	if response.StatusCode() != 200 {
		return nil, getHttpError(response)
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
UploadRecordFile uploads the given record data to the api and saves it as a .snmprec file at the given remote path inside of the data dir.
*/
func (c *ManagementClient) UploadRecordFileString(recordContents *string, remotePath string) error {
	headerMap := make(map[string]string)
	headerMap["Content-Type"] = "text/plain"
	response, err := c.request("POST", MGMT_ENDPOINT_PATH+"recordings/"+remotePath, *recordContents, headerMap, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 204 {
		return getHttpError(response)
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
	response, err := c.request("DELETE", MGMT_ENDPOINT_PATH+"recordings/"+remotePath, "", headerMap, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 204 {
		return getHttpError(response)
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
	response, err := c.request("GET", MGMT_ENDPOINT_PATH+"recordings/"+remotePath, "", headerMap, nil)
	if err != nil {
		return "", errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 200 {
		return "", getHttpError(response)
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
	if !c.isValid() {
		return User{}, &NotValidError{}
	}

	if name == "" {
		return User{}, errors.New("invalid name")
	}

	if user == "" {
		return User{}, errors.New("invalid user")
	}

	if authProto == "" {
		authProto = "none"
	}
	if privProto == "" {
		privProto = "none"
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
		User:      user,
		Name:      name,
		AuthProto: authProto,
		PrivProto: privProto,
	}

	if authKey == "" {
		params.AuthKey = nil
	} else {
		params.AuthKey = &authKey
	}

	if privKey == "" {
		params.PrivKey = nil
	} else {
		params.PrivKey = &privKey
	}

	jsonString, err := json.Marshal(params)
	if err != nil {
		return User{}, errors.Wrap(err, "error during marshal")
	}

	response, err := c.request("POST", MGMT_ENDPOINT_PATH+"users", string(jsonString), nil, nil)
	if err != nil {
		return User{}, errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 201 {
		return User{}, getHttpError(response)
	}

	var newUser User
	err = json.Unmarshal(response.Body(), &newUser)
	if err != nil {
		return User{}, errors.Wrap(err, "error during unmarshalling http response")
	}
	return newUser, nil
}

/*
GetUsers returns a list of all users.
*/
func (c *ManagementClient) GetUsers() (Users, error) {
	if !c.isValid() {
		return nil, &NotValidError{}
	}

	response, err := c.request("GET", MGMT_ENDPOINT_PATH+"users", "", nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error during get labs request")
	}
	if response.StatusCode() != 200 {
		return nil, getHttpError(response)
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

	response, err := c.request("GET", MGMT_ENDPOINT_PATH+"users/"+strconv.Itoa(id), "", nil, nil)
	if err != nil {
		return User{}, errors.Wrap(err, "error during get labs request")
	}
	if response.StatusCode() != 200 {
		return User{}, getHttpError(response)
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

	response, err := c.request("DELETE", MGMT_ENDPOINT_PATH+"users/"+strconv.Itoa(id), "", nil, nil)
	if err != nil {
		return errors.Wrap(err, "error during request")
	}
	if response.StatusCode() != 204 {
		return getHttpError(response)
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
	return Selector{}, errors.New("Not implemented yet!")
}

/*
GetSelectors returns a list of all selectors.
*/
func (c *ManagementClient) GetSelectors() (Selectors, error) {
	//TODO: Not implemented yet!!! This also isnt implemented in the api so we have to wait until its ready.
	return nil, errors.New("Not implemented yet!")
}

/*
GetSelector returns the selector with the given id.
*/
func (c *ManagementClient) GetSelector(id int) (Selector, error) {
	//TODO: Not implemented yet!!! This also isnt implemented in the api so we have to wait until its ready.
	return Selector{}, errors.New("Not implemented yet!")
}

/*
DeleteSelector deletes the selector with the given id.
*/
func (c *ManagementClient) DeleteSelector(id int) error {
	//TODO: Not implemented yet!!! This also isnt implemented in the api so we have to wait until its ready.
	return errors.New("Not implemented yet!")
}
