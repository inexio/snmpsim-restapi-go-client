package snmpsimclient

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"
)

type agentData struct {
	EndpointAddress string `mapstructure:"endpointAddress"`
	EndpointPort    []int  `mapstructure:"endpointPort"`
}

type httpConfig struct {
	BaseUrl      string `mapstructure:"baseUrl"`
	AuthUsername string `mapstructure:"authUsername"`
	AuthPassword string `mapstructure:"authPassword"`
}

type configMetricsApiTest struct {
	Http        httpConfig `mapstructure:"http"`
	Protocol    string     `mapstructure:"protocol"`
	Agent1      agentData  `mapstructure:"agent1"`
	RootDataDir string     `mapstructure:"rootDataDir"`
	TestDataDir string
}

var configMetricsTest configMetricsApiTest

type configManagementApiTest struct {
	Http        httpConfig `mapstructure:"http"`
	Protocol    string     `mapstructure:"protocol"`
	Agent1      agentData  `mapstructure:"agent1"`
	Agent2      agentData  `mapstructure:"agent2"`
	RootDataDir string     `mapstructure:"rootDataDir"`
	TestDataDir string
	TestTagId   int
	TestTagName string `mapstructure:"testTag"`
}

var configManagementTest configManagementApiTest

func init() {
	_, currFilename, _, _ := runtime.Caller(0)
	testDataDir := path.Dir(currFilename) + "/test-data/"

	//management config
	viperManagement := viper.New()
	viperManagement.SetConfigFile(testDataDir + "management-api-test-config.yaml")
	viperManagement.SetEnvPrefix("snmpsim_management_api_test")
	replacer := strings.NewReplacer(".", "_")
	viperManagement.SetEnvKeyReplacer(replacer)
	viperManagement.AutomaticEnv()

	err := viperManagement.ReadInConfig()
	if err != nil {
		fmt.Println("Failed to read management config file!", err)
		os.Exit(2)
	}
	err = viperManagement.Unmarshal(&configManagementTest)
	if err != nil {
		fmt.Println("Failed to unmarshal config file!")
		os.Exit(2)
	}

	configManagementTest.TestDataDir = testDataDir

	//metrics config
	viperMetrics := viper.New()
	viperMetrics.SetConfigFile(testDataDir + "metrics-api-test-config.yaml")
	viperMetrics.SetEnvPrefix("snmpsim_metrics_api_test")
	viperMetrics.SetEnvKeyReplacer(replacer)
	viperMetrics.AutomaticEnv()

	err = viperMetrics.ReadInConfig()
	if err != nil {
		fmt.Println("Failed to metrics read config file!", err)
		os.Exit(2)
	}

	err = viperMetrics.Unmarshal(&configMetricsTest)
	if err != nil {
		fmt.Println("Failed to unmarshal config file!")
		os.Exit(2)
	}

	configMetricsTest.TestDataDir = testDataDir

	//tags
	//Create a new api client
	client, err := NewManagementClient(configManagementTest.Http.BaseUrl)
	if err != nil {
		fmt.Println("Failed to create a new Management Client during init()!", err)
		os.Exit(2)
	}
	//Set configManagementTest.Http.AuthUsername and password
	if configManagementTest.Http.AuthUsername != "" && configManagementTest.Http.AuthPassword != "" {
		err = client.SetUsernameAndPassword(configManagementTest.Http.AuthUsername, configManagementTest.Http.AuthPassword)
		if err != nil {
			fmt.Println("Failed to set http auth username and password during init()!", err)
			os.Exit(2)
		}
	}

	tagFilter := make(map[string]string)
	tagFilter["name"] = configManagementTest.TestTagName
	tags, err := client.GetTags(tagFilter)
	if err != nil {
		fmt.Println("Error while trying to get the object tag for the test!", err)
		os.Exit(2)
	}

	switch tagsCount := len(tags); tagsCount {
	case 1:
		//this is the usual case
		configManagementTest.TestTagId = tags[0].Id

		//delete objects with this text
		_, err = client.DeleteAllObjectsWithTag(configManagementTest.TestTagId)
		if err != nil {
			fmt.Println("Error while trying to delete all objects tagged with tag-id", configManagementTest.TestTagId)
			os.Exit(2)
		}
	case 0:
		//tag not found, this might be legit when the test runs for the first time or the snmpsim db was recreated
		fmt.Println("No tag for management api test found -> creating a new tag!")
		tag, err := client.CreateTag(tagFilter["name"], "tag for tagging objects used in test cases in the management api test")
		if err != nil {
			fmt.Println("Failed to create a new tag for management api test!")
			os.Exit(2)
		}
		fmt.Println("Successfully created new tag for management api test!")
		configManagementTest.TestTagId = tag.Id
	default:
		//more than one tag found, something is wrong
		fmt.Println("There was more than one tag found for the management api test!")
		os.Exit(2)
	}
}

/*
HELPER FUNCTIONS FOR PERFORMING API CALLS AND CHECKING IF THEY WHERE SUCCESSFUL
*/

/*
AGENTS
*/

func createAgentAndCheckForSuccess(t *testing.T, client *ManagementClient, name, dataDir string) (Agent, error) {
	agent, err := client.CreateAgentWithTag(name, dataDir, configManagementTest.TestTagId)
	if !assert.NoError(t, err, "error while creating a new agent") {
		return Agent{}, err
	}

	//Test if lab was created
	agents, err := client.GetAgents(nil)
	if !assert.NoError(t, err, "error during GetAgents()") {
		return Agent{}, err
	}

	if !assert.True(t, agentExists(agent, agents), "created agent was not found in list of agents") {
		return Agent{}, errors.New("assertion failed")
	}
	return agent, nil
}

func deleteAgentAndCheckForSuccess(t *testing.T, client *ManagementClient, agent Agent) error {
	err := client.DeleteAgent(agent.Id)
	assert.NoError(t, err, "error while deleting lab")

	agents, err := client.GetAgents(nil)
	if !assert.NoError(t, err, "error during GetAgents()") {
		return err
	}
	if !assert.False(t, agentExists(agent, agents), "deleted agent was found in list of agents") {
		return errors.New("assertion failed")
	}
	return nil
}

func addEngineToAgentAndCheckForSuccess(t *testing.T, client *ManagementClient, agent Agent, engine Engine) error {
	err := client.AddEngineToAgent(agent.Id, engine.Id)
	if !assert.NoError(t, err, "error while adding engine to agent") {
		return err
	}

	//Test if engine was added to the agent
	agent, err = client.GetAgent(agent.Id)
	if !assert.NoError(t, err, "error while get agent") {
		return err
	}
	if !assert.True(t, isEngineInAgent(agent, engine), "engine was successfully added to an agent, but cannot be found in the agents engine list") {
		return errors.New("assertion failed")
	}
	return nil
}

func removeEngineFromAgentAndCheckForSuccess(t *testing.T, client *ManagementClient, agent Agent, engine Engine) error {
	//Remove engine from agent
	err := client.RemoveEngineFromAgent(agent.Id, engine.Id)
	if !assert.NoError(t, err, "error while deleting engine") {
		return err
	}
	agent, err = client.GetAgent(agent.Id)
	if !assert.NoError(t, err, "error while get agent") {
		return err
	}
	if !assert.False(t, isEngineInAgent(agent, engine), "engine is still in agent after successfully removing it") {
		return errors.New("assertion failed")
	}
	return nil
}

func addTagToAgentAndCheckForSuccess(t *testing.T, client *ManagementClient, agent Agent, tag Tag) error {
	err := client.AddTagToAgent(agent.Id, tag.Id)
	if !assert.NoError(t, err, "error while adding tag to agent") {
		return err
	}

	//Test if tag was added to the agent
	agent, err = client.GetAgent(agent.Id)
	if !assert.NoError(t, err, "error while get agent") {
		return err
	}
	if !assert.True(t, tagExists(tag, agent.Tags), "tag was successfully added to an agent, but cannot be found in the agents tag list") {
		return errors.New("assertion failed")
	}
	return nil
}

func removeTagFromAgentAndCheckForSuccess(t *testing.T, client *ManagementClient, agent Agent, tag Tag) error {
	err := client.RemoveTagFromAgent(agent.Id, tag.Id)
	if !assert.NoError(t, err, "error while removing tag from agent") {
		return err
	}
	agent, err = client.GetAgent(agent.Id)
	if !assert.NoError(t, err, "error while get agent") {
		return err
	}
	if !assert.False(t, tagExists(tag, agent.Tags), "agent is still tagged after successfully removing tag") {
		return errors.New("assertion failed")
	}
	return nil
}

/*
LABS
*/

func createLabAndCheckForSuccess(t *testing.T, client *ManagementClient, name string) (Lab, error) {
	lab, err := client.CreateLabWithTag(name, configManagementTest.TestTagId)
	if !assert.NoError(t, err, "error while creating a new lab") {
		return Lab{}, err
	}

	//Test if lab was created
	labs, err := client.GetLabs(nil)
	if !assert.NoError(t, err, "error during GetLabs()") {
		return Lab{}, err
	}
	if !assert.True(t, labExists(lab, labs), "created lab was not found in list of labs") {
		return Lab{}, errors.New("assertion failed")
	}
	return lab, nil
}

func deleteLabAndCheckForSuccess(t *testing.T, client *ManagementClient, lab Lab) error {
	err := client.DeleteLab(lab.Id)
	assert.NoError(t, err, "error while deleting lab")

	//Test if lab was deleted
	labs, err := client.GetLabs(nil)
	if !assert.NoError(t, err, "error during GetLabs()") {
		return err
	}
	if !assert.False(t, labExists(lab, labs), "deleted lab was found in list of labs") {
		return err
	}
	return nil
}

func addAgentToLabAndCheckForSuccess(t *testing.T, client *ManagementClient, lab Lab, agent Agent) error {
	err := client.AddAgentToLab(lab.Id, agent.Id)
	if !assert.NoError(t, err, "error while adding agent to lab") {
		return err
	}

	//Test if agent was added to the lab
	lab, err = client.GetLab(lab.Id)
	if !assert.NoError(t, err, "error while get lab") {
		return err
	}
	if !assert.True(t, isAgentInLab(lab, agent), "agent was successfully added to an lab, but cannot be found in the labs agent list") {
		return errors.New("assertion failed")
	}
	return nil
}

func removeAgentFromLabAndCheckForSuccess(t *testing.T, client *ManagementClient, lab Lab, agent Agent) error {
	//Remove agent from lab
	err := client.RemoveAgentFromLab(lab.Id, agent.Id)
	if !assert.NoError(t, err, "error while deleting agent") {
		return err
	}
	lab, err = client.GetLab(lab.Id)
	if !assert.NoError(t, err, "error while get lab") {
		return err
	}
	if !assert.False(t, isAgentInLab(lab, agent), "agent is still in lab after successfully removing it") {
		return errors.New("assertion failed")
	}
	return nil
}

func setLabPowerAndCheckForSuccess(t *testing.T, client *ManagementClient, lab Lab, powerBool bool) error {
	power := "on"
	if powerBool == false {
		power = "off"
	}
	err := client.SetLabPower(lab.Id, powerBool)
	if !assert.NoError(t, err, "error while turning power on for lab") {
		return err
	}
	lab, err = client.GetLab(lab.Id)
	if !assert.NoError(t, err, "error while get lab") {
		return err
	}
	if !assert.True(t, lab.Power == power, `lab power is not "`+power+`" after successfully setting it to "`+power+`"`) {
		return errors.New("assertion failed")
	}
	return nil
}

func addTagToLabAndCheckForSuccess(t *testing.T, client *ManagementClient, lab Lab, tag Tag) error {
	err := client.AddTagToLab(lab.Id, tag.Id)
	if !assert.NoError(t, err, "error while adding tag to lab") {
		return err
	}

	//Test if tag was added to the lab
	lab, err = client.GetLab(lab.Id)
	if !assert.NoError(t, err, "error while get lab") {
		return err
	}
	if !assert.True(t, tagExists(tag, lab.Tags), "tag was successfully added to an lab, but cannot be found in the labs tag list") {
		return errors.New("assertion failed")
	}
	return nil
}

func removeTagFromLabAndCheckForSuccess(t *testing.T, client *ManagementClient, lab Lab, tag Tag) error {
	err := client.RemoveTagFromLab(lab.Id, tag.Id)
	if !assert.NoError(t, err, "error while removing tag from lab") {
		return err
	}
	lab, err = client.GetLab(lab.Id)
	if !assert.NoError(t, err, "error while get lab") {
		return err
	}
	if !assert.False(t, tagExists(tag, lab.Tags), "lab is still tagged after successfully removing tag") {
		return errors.New("assertion failed")
	}
	return nil
}

/*
ENGINES
*/
func createEngineAndCheckForSuccess(t *testing.T, client *ManagementClient, name, engineId string) (Engine, error) {
	//Create an engine
	engine, err := client.CreateEngineWithTag(name, engineId, configManagementTest.TestTagId)
	if !assert.NoError(t, err, "error while creating a new api engine") {
		return Engine{}, err
	}

	//Test if engine was created
	engines, err := client.GetEngines(nil)
	if !assert.NoError(t, err, "error during GetEngines()") {
		return Engine{}, err
	}
	if !assert.True(t, engineExists(engine, engines), "created engine was not found in list of engines") {
		return Engine{}, errors.New("assertion failed")
	}
	return engine, nil
}

func deleteEngineAndCheckForSuccess(t *testing.T, client *ManagementClient, engine Engine) error {
	err := client.DeleteEngine(engine.Id)
	assert.NoError(t, err, "error while deleting engine")

	//Test if engine was deleted
	engines, err := client.GetEngines(nil)
	if !assert.NoError(t, err, "error during GetEngines()") {
		return err
	}
	if !assert.False(t, engineExists(engine, engines), "created engine was not found in list of engines") {
		return err
	}
	return nil
}

func addUserToEngineAndCheckForSuccess(t *testing.T, client *ManagementClient, engine Engine, user User) error {
	err := client.AddUserToEngine(engine.Id, user.Id)
	if !assert.NoError(t, err, "error while adding user to engine") {
		return err
	}

	//Test if user was added to the engine
	engine, err = client.GetEngine(engine.Id)
	if !assert.NoError(t, err, "error while get engine") {
		return err
	}
	if !assert.True(t, isUserInEngine(engine, user), "user was successfully added to an engine, but cannot be found in the engines user list") {
		return errors.New("assertion failed")
	}
	return nil
}

func removeUserFromEngineAndCheckForSuccess(t *testing.T, client *ManagementClient, engine Engine, user User) error {
	//Remove user from engine
	err := client.RemoveUserFromEngine(engine.Id, user.Id)
	if !assert.NoError(t, err, "error while deleting user") {
		return err
	}
	engine, err = client.GetEngine(engine.Id)
	if !assert.NoError(t, err, "error while get engine") {
		return err
	}
	if !assert.False(t, isUserInEngine(engine, user), "user is still in engine after successfully removing it") {
		return errors.New("assertion failed")
	}
	return nil
}

func addEndpointToEngineAndCheckForSuccess(t *testing.T, client *ManagementClient, engine Engine, endpoint Endpoint) error {
	err := client.AddEndpointToEngine(engine.Id, endpoint.Id)
	if !assert.NoError(t, err, "error while adding endpoint to engine") {
		return err
	}

	//Test if endpoint was added to the engine
	engine, err = client.GetEngine(engine.Id)
	if !assert.NoError(t, err, "error while get engine") {
		return err
	}
	if !assert.True(t, isEndpointInEngine(engine, endpoint), "endpoint was successfully added to an engine, but cannot be found in the engines endpoint list") {
		return errors.New("assertion failed")
	}
	return nil
}

func removeEndpointFromEngineAndCheckForSuccess(t *testing.T, client *ManagementClient, engine Engine, endpoint Endpoint) error {
	//Remove endpoint from engine
	err := client.RemoveEndpointFromEngine(engine.Id, endpoint.Id)
	if !assert.NoError(t, err, "error while deleting endpoint") {
		return err
	}
	engine, err = client.GetEngine(engine.Id)
	if !assert.NoError(t, err, "error while get engine") {
		return err
	}
	if !assert.False(t, isEndpointInEngine(engine, endpoint), "endpoint is still in engine after successfully removing it") {
		return errors.New("assertion failed")
	}
	return nil
}

func addTagToEngineAndCheckForSuccess(t *testing.T, client *ManagementClient, engine Engine, tag Tag) error {
	err := client.AddTagToEngine(engine.Id, tag.Id)
	if !assert.NoError(t, err, "error while adding tag to engine") {
		return err
	}

	//Test if tag was added to the engine
	engine, err = client.GetEngine(engine.Id)
	if !assert.NoError(t, err, "error while get engine") {
		return err
	}
	if !assert.True(t, tagExists(tag, engine.Tags), "tag was successfully added to an engine, but cannot be found in the engines tag list") {
		return errors.New("assertion failed")
	}
	return nil
}

func removeTagFromEngineAndCheckForSuccess(t *testing.T, client *ManagementClient, engine Engine, tag Tag) error {
	err := client.RemoveTagFromEngine(engine.Id, tag.Id)
	if !assert.NoError(t, err, "error while removing tag from engine") {
		return err
	}
	engine, err = client.GetEngine(engine.Id)
	if !assert.NoError(t, err, "error while get engine") {
		return err
	}
	if !assert.False(t, tagExists(tag, engine.Tags), "engine is still tagged after successfully removing tag") {
		return errors.New("assertion failed")
	}
	return nil
}

/*
ENDPOINTS
*/

func createEndpointAndCheckForSuccess(t *testing.T, client *ManagementClient, name, address, domain string) (Endpoint, error) {
	endpoint, err := client.CreateEndpointWithTag(name, address, domain, configManagementTest.TestTagId)
	if !assert.NoError(t, err, "error while creating a new endpoint") {
		return Endpoint{}, err
	}

	//Test if endpoint was created
	endpoints, err := client.GetEndpoints(nil)
	if !assert.NoError(t, err, "error during GetEndpoints()") {
		return Endpoint{}, err
	}

	if !assert.True(t, endpointExists(endpoint, endpoints), "created endpoint was not found in list of endpoints") {
		return Endpoint{}, errors.New("assertion failed")
	}
	return endpoint, nil
}

func deleteEndpointAndCheckForSuccess(t *testing.T, client *ManagementClient, endpoint Endpoint) error {
	err := client.DeleteEndpoint(endpoint.Id)
	assert.NoError(t, err, "error while deleting endpoint")

	//Test if endpoint was deleted
	endpoints, err := client.GetEndpoints(nil)
	if !assert.NoError(t, err, "error during GetEndpoints()") {
		return err
	}
	if !assert.False(t, endpointExists(endpoint, endpoints), "deleted endpoint was found in list of endpoints") {
		return err
	}
	return nil
}

func addTagToEndpointAndCheckForSuccess(t *testing.T, client *ManagementClient, endpoint Endpoint, tag Tag) error {
	err := client.AddTagToEndpoint(endpoint.Id, tag.Id)
	if !assert.NoError(t, err, "error while adding tag to endpoint") {
		return err
	}
	//Test if tag was added to the endpoint
	endpoint, err = client.GetEndpoint(endpoint.Id)
	if !assert.NoError(t, err, "error while get endpoint") {
		return err
	}
	if !assert.True(t, tagExists(tag, endpoint.Tags), "tag was successfully added to an endpoint, but cannot be found in the endpoints tag list") {
		return errors.New("assertion failed")
	}
	return nil
}

func removeTagFromEndpointAndCheckForSuccess(t *testing.T, client *ManagementClient, endpoint Endpoint, tag Tag) error {
	err := client.RemoveTagFromEndpoint(endpoint.Id, tag.Id)
	if !assert.NoError(t, err, "error while removing tag from endpoint") {
		return err
	}
	endpoint, err = client.GetEndpoint(endpoint.Id)
	if !assert.NoError(t, err, "error while get endpoint") {
		return err
	}
	if !assert.False(t, tagExists(tag, endpoint.Tags), "endpoint is still tagged after successfully removing tag") {
		return errors.New("assertion failed")
	}
	return nil
}

/*
USERS
*/

func createUserAndCheckForSuccess(t *testing.T, client *ManagementClient, userIdentifier, name, authKey, authProto, privKey, privProto string) (User, error) {
	user, err := client.CreateUserWithTag(userIdentifier, name, authKey, authProto, privKey, privProto, configManagementTest.TestTagId)
	if !assert.NoError(t, err, "error while creating a new user") {
		return User{}, err
	}

	//Test if user was created
	users, err := client.GetUsers(nil)
	if !assert.NoError(t, err, "error during GetUsers()") {
		return User{}, err
	}

	if !assert.True(t, userExists(user, users), "created user was not found in list of users") {
		return User{}, errors.New("assertion failed")
	}
	return user, nil
}

func deleteUserAndCheckForSuccess(t *testing.T, client *ManagementClient, user User) error {
	err := client.DeleteUser(user.Id)
	assert.NoError(t, err, "error while deleting user")

	//Test if user was deleted
	users, err := client.GetUsers(nil)
	if !assert.NoError(t, err, "error during GetUsers()") {
		return err
	}
	if !assert.False(t, userExists(user, users), "deleted user was not found in list of users") {
		return err
	}
	return nil
}

func addTagToUserAndCheckForSuccess(t *testing.T, client *ManagementClient, user User, tag Tag) error {
	err := client.AddTagToUser(user.Id, tag.Id)
	if !assert.NoError(t, err, "error while adding tag to user") {
		return err
	}

	//Test if tag was added to the user
	user, err = client.GetUser(user.Id)
	if !assert.NoError(t, err, "error while get user") {
		return err
	}
	if !assert.True(t, tagExists(tag, user.Tags), "tag was successfully added to an user, but cannot be found in the users tag list") {
		return errors.New("assertion failed")
	}
	return nil
}

func removeTagFromUserAndCheckForSuccess(t *testing.T, client *ManagementClient, user User, tag Tag) error {
	err := client.RemoveTagFromUser(user.Id, tag.Id)
	if !assert.NoError(t, err, "error while removing tag from user") {
		return err
	}
	user, err = client.GetUser(user.Id)
	if !assert.NoError(t, err, "error while get user") {
		return err
	}
	if !assert.False(t, tagExists(tag, user.Tags), "user is still tagged after successfully removing tag") {
		return errors.New("assertion failed")
	}
	return nil
}

/*
RECORD FILES
*/
func uploadRecordFileAndCheckForSuccess(t *testing.T, client *ManagementClient, localPath, remotePath string) error {
	err := client.UploadRecordFile(localPath, remotePath)
	if !assert.NoError(t, err, "error while uploading record file") {
		return err
	}
	_, err = client.GetRecordFile(remotePath)
	if !assert.NoError(t, err, "error during GetRecordFile()") {
		return err
	}
	return nil
}
func uploadRecordFileStringAndCheckForSuccess(t *testing.T, client *ManagementClient, fileContents *string, remotePath string) error {
	err := client.UploadRecordFileString(fileContents, remotePath)
	if !assert.NoError(t, err, "error while uploading record file") {
		return err
	}
	_, err = client.GetRecordFile(remotePath)
	if !assert.NoError(t, err, "error during GetRecordFile()") {
		return err
	}
	return nil
}
func deleteRecordFileAndCheckForSuccess(t *testing.T, client *ManagementClient, remotePath string) error {
	err := client.DeleteRecordFile(remotePath)
	if !assert.NoError(t, err, "error while uploading record file") {
		return err
	}
	_, err = client.GetRecordFile(remotePath)
	if !assert.Error(t, err, "error during GetRecordFile()") {
		return err
	}
	return nil
}

/*
 TAGS
*/
func createTagAndCheckForSuccess(t *testing.T, client *ManagementClient, name, description string) (Tag, error) {
	tag, err := client.CreateTag(name, description)
	if !assert.NoError(t, err, "error while creating a new tag") {
		return Tag{}, err
	}

	//Test if tag was created
	tags, err := client.GetTags(nil)
	if !assert.NoError(t, err, "error during GetTags()") {
		return Tag{}, err
	}
	if !assert.True(t, tagExists(tag, tags), "created tag was not found in list of tags") {
		return Tag{}, errors.New("assertion failed")
	}
	return tag, nil
}

func deleteTagAndCheckForSuccess(t *testing.T, client *ManagementClient, tag Tag) error {
	err := client.DeleteTag(tag.Id)
	assert.NoError(t, err, "error while deleting tag")

	//Test if tag was deleted
	tags, err := client.GetTags(nil)
	if !assert.NoError(t, err, "error during GetTags()") {
		return err
	}
	if !assert.False(t, tagExists(tag, tags), "created tag was not found in list of tags") {
		return err
	}
	return nil
}

/*
HELP FUNCTIONS
*/
func isUserInEngine(engine Engine, user User) bool {
	isUserInEngine := false
	for _, currUser := range engine.Users {
		if currUser.Id == user.Id {
			isUserInEngine = true
			break
		}
	}
	return isUserInEngine
}

func isEndpointInEngine(engine Engine, endpoint Endpoint) bool {
	isEndpointInEngine := false
	for _, currEndpoint := range engine.Endpoints {
		if currEndpoint.Id == endpoint.Id {
			isEndpointInEngine = true
			break
		}
	}
	return isEndpointInEngine
}

func isEngineInAgent(agent Agent, engine Engine) bool {
	isEngineInAgent := false
	for _, currEngine := range agent.Engines {
		if currEngine.Id == engine.Id {
			isEngineInAgent = true
			break
		}
	}
	return isEngineInAgent
}

func isAgentInLab(lab Lab, agent Agent) bool {
	isAgentInLab := false
	for _, currAgent := range lab.Agents {
		if currAgent.Id == agent.Id {
			isAgentInLab = true
			break
		}
	}
	return isAgentInLab
}

func agentExists(agent Agent, agents Agents) bool {
	agentWasCreated := false
	for _, currAgent := range agents {
		if currAgent.Id == agent.Id {
			agentWasCreated = true
			break
		}
	}
	return agentWasCreated
}

func labExists(lab Lab, labs Labs) bool {
	labWasCreated := false
	for _, currLab := range labs {
		if currLab.Id == lab.Id {
			labWasCreated = true
			break
		}
	}
	return labWasCreated
}

func engineExists(engine Engine, engines Engines) bool {
	engineWasCreated := false
	for _, currEngine := range engines {
		if currEngine.Id == engine.Id {
			engineWasCreated = true
			break
		}
	}
	return engineWasCreated
}

func endpointExists(endpoint Endpoint, endpoints Endpoints) bool {
	endpointWasCreated := false
	for _, currEndpoint := range endpoints {
		if currEndpoint.Id == endpoint.Id {
			endpointWasCreated = true
			break
		}
	}
	return endpointWasCreated
}
func userExists(user User, users Users) bool {
	userWasCreated := false
	for _, currUser := range users {
		if currUser.Id == user.Id {
			userWasCreated = true
			break
		}
	}
	return userWasCreated
}
func tagExists(tag Tag, tags Tags) bool {
	tagWasCreated := false
	for _, currTag := range tags {
		if currTag.Id == tag.Id {
			tagWasCreated = true
			break
		}
	}
	return tagWasCreated
}
