package snmpsimclient

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"strings"
	"os"
	"path"
	"runtime"
	"testing"
)

type agentData struct {
	EndpointAddress string `mapstructure:"endpointAddress"`
	EndpointPort    []int  `mapstructure:"endpointPort"`
}

type httpConfig struct {
	BaseUrl         string    `mapstructure:"baseUrl"`
	AuthUsername 	string    `mapstructure:"authUsername"`
	AuthPassword 	string    `mapstructure:"authPassword"`
}

type configMetricsApiTest struct {
	Http            httpConfig	`mapstructure:"http"`
	Protocol        string    	`mapstructure:"protocol"`
	Agent1          agentData 	`mapstructure:"agent1"`
	RootDataDir     string    	`mapstructure:"rootDataDir"`
	TestDataDir     string
}

var configMetricsTest configMetricsApiTest

type configManagementApiTest struct {
	Http            httpConfig	`mapstructure:"http"`
	Protocol        string    	`mapstructure:"protocol"`
	Agent1          agentData 	`mapstructure:"agent1"`
	Agent2          agentData 	`mapstructure:"agent2"`
	RootDataDir     string    	`mapstructure:"rootDataDir"`
	TestDataDir     string
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
}

/*
HELPER FUNCTIONS FOR PERFORMING API CALLS AND CHECKING IF THEY WHERE SUCCESSFUL
*/

/*
AGENTS
*/

func createAgentAndCheckForSuccess(t *testing.T, client *ManagementClient, name, dataDir string) (Agent, error) {
	agent, err := client.CreateAgent(name, dataDir)
	if !assert.NoError(t, err, "error while creating a new agent") {
		return Agent{}, err
	}

	//Test if lab was created
	agents, err := client.GetAgents()
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

	agents, err := client.GetAgents()
	if !assert.NoError(t, err, "error during GetAgents()") {
		return err
	}
	if !assert.False(t, agentExists(agent, agents), "created agent was found in list of agents") {
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

/*
LABS
*/

func createLabAndCheckForSuccess(t *testing.T, client *ManagementClient, name string) (Lab, error) {
	lab, err := client.CreateLab(name)
	if !assert.NoError(t, err, "error while creating a new lab") {
		return Lab{}, err
	}

	//Test if lab was created
	labs, err := client.GetLabs()
	if !assert.NoError(t, err, "error during GetLabs()") {
		return Lab{}, err
	}
	if !assert.True(t, labsExists(lab, labs), "created lab was not found in list of labs") {
		return Lab{}, errors.New("assertion failed")
	}
	return lab, nil
}

func deleteLabAndCheckForSuccess(t *testing.T, client *ManagementClient, lab Lab) error {
	err := client.DeleteLab(lab.Id)
	assert.NoError(t, err, "error while deleting lab")

	//Test if lab was deleted
	labs, err := client.GetLabs()
	if !assert.NoError(t, err, "error during GetLabs()") {
		return err
	}
	if !assert.False(t, labsExists(lab, labs), "created lab was not found in list of labs") {
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

/*
ENGINES
*/

func createEngineAndCheckForSuccess(t *testing.T, client *ManagementClient, name, engineId string) (Engine, error) {
	//Create an engine
	engine, err := client.CreateEngine(name, engineId)
	if !assert.NoError(t, err, "error while creating a new api engine") {
		return Engine{}, err
	}

	//Test if engine was created
	engines, err := client.GetEngines()
	if !assert.NoError(t, err, "error during GetEngines()") {
		return Engine{}, err
	}
	if !assert.True(t, engineExist(engine, engines), "created engine was not found in list of engines") {
		return Engine{}, errors.New("assertion failed")
	}
	return engine, nil
}

func deleteEngineAndCheckForSuccess(t *testing.T, client *ManagementClient, engine Engine) error {
	err := client.DeleteEngine(engine.Id)
	assert.NoError(t, err, "error while deleting engine")

	//Test if engine was deleted
	engines, err := client.GetEngines()
	if !assert.NoError(t, err, "error during GetEngines()") {
		return err
	}
	if !assert.False(t, engineExist(engine, engines), "created engine was not found in list of engines") {
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

/*
ENDPOINTS
*/

func createEndpointAndCheckForSuccess(t *testing.T, client *ManagementClient, name, address, domain string) (Endpoint, error) {
	endpoint, err := client.CreateEndpoint(name, address, domain)
	if !assert.NoError(t, err, "error while creating a new endpoint") {
		return Endpoint{}, err
	}

	//Test if endpoint was created
	endpoints, err := client.GetEndpoints()
	if !assert.NoError(t, err, "error during GetEndpoints()") {
		return Endpoint{}, err
	}

	if !assert.True(t, endpointExist(endpoint, endpoints), "created endpoint was not found in list of endpoints") {
		return Endpoint{}, errors.New("assertion failed")
	}
	return endpoint, nil
}

func deleteEndpointAndCheckForSuccess(t *testing.T, client *ManagementClient, endpoint Endpoint) error {
	err := client.DeleteEndpoint(endpoint.Id)
	assert.NoError(t, err, "error while deleting endpoint")

	//Test if endpoint was deleted
	endpoints, err := client.GetEndpoints()
	if !assert.NoError(t, err, "error during GetEndpoints()") {
		return err
	}
	if !assert.False(t, endpointExist(endpoint, endpoints), "created endpoint was not found in list of endpoints") {
		return err
	}
	return nil
}

/*
USERS
*/

func createUserAndCheckForSuccess(t *testing.T, client *ManagementClient, userIdentifier, name, authKey, authProto, privKey, privProto string) (User, error) {
	user, err := client.CreateUser(userIdentifier, name, authKey, authProto, privKey, privProto)
	if !assert.NoError(t, err, "error while creating a new user") {
		return User{}, err
	}

	//Test if user was created
	users, err := client.GetUsers()
	if !assert.NoError(t, err, "error during GetUsers()") {
		return User{}, err
	}

	if !assert.True(t, userExist(user, users), "created user was not found in list of users") {
		return User{}, errors.New("assertion failed")
	}
	return user, nil
}

func deleteUserAndCheckForSuccess(t *testing.T, client *ManagementClient, user User) error {
	err := client.DeleteUser(user.Id)
	assert.NoError(t, err, "error while deleting user")

	//Test if user was deleted
	users, err := client.GetUsers()
	if !assert.NoError(t, err, "error during GetUsers()") {
		return err
	}
	if !assert.False(t, userExist(user, users), "deleted user was not found in list of users") {
		return err
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

func labsExists(lab Lab, labs Labs) bool {
	labWasCreated := false
	for _, currLab := range labs {
		if currLab.Id == lab.Id {
			labWasCreated = true
			break
		}
	}
	return labWasCreated
}

func engineExist(engine Engine, engines Engines) bool {
	engineWasCreated := false
	for _, currEngine := range engines {
		if currEngine.Id == engine.Id {
			engineWasCreated = true
			break
		}
	}
	return engineWasCreated
}

func endpointExist(endpoint Endpoint, endpoints Endpoints) bool {
	endpointWasCreated := false
	for _, currEndpoint := range endpoints {
		if currEndpoint.Id == endpoint.Id {
			endpointWasCreated = true
			break
		}
	}
	return endpointWasCreated
}
func userExist(user User, users Users) bool {
	userWasCreated := false
	for _, currUser := range users {
		if currUser.Id == user.Id {
			userWasCreated = true
			break
		}
	}
	return userWasCreated
}
