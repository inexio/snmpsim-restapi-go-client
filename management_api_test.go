package snmpsimclient

import (
	"github.com/soniah/gosnmp"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestManagementClient_buildUpSetupAndTestIt(t *testing.T) {
	community := "public"
	//	Agent 1
	//Agent
	agentName1 := "test-buildUpSetupAndTestIt-agent1"
	agentDataDir1 := configManagementTest.RootDataDir + "test-buildUpSetupAndTestIt-agent1"
	//Endpoint
	endpointName1 := "test-buildUpSetupAndTestIt-endpoint1"
	address1 := configManagementTest.Agent1.EndpointAddress + ":" + strconv.Itoa(configManagementTest.Agent1.EndpointPort[0])
	//User
	name1 := "test-buildUpSetupAndTestIt-user1"
	userIdentifier1 := "test-buildUpSetupAndTestIt-user1"
	authKey1 := "0x50dd4d3ec79a1cf4dfa5fee9f76b0847647fcf74"
	authProto1 := "sha"
	privKey1 := "0x50dd4d3ec79a1cf4dfa5fee9f76b0847"
	privProto1 := "des"
	//engine
	engineName1 := "test-buildUpSetupAndTestIt-engine1"
	engineId1 := "0102030405070809"
	//Record File:
	localRecordFilePath1 := configManagementTest.TestDataDir + "snmprecs/TestManagementClient_buildUpSetupAndTestIt/agent1/" + community + ".snmprec"
	remoteRecordFilePath1 := agentDataDir1 + "/" + community + ".snmprec"

	// Agent 2
	//Agent
	agentName2 := "test-buildUpSetupAndTestIt-agent1"
	agentDataDir2 := configManagementTest.RootDataDir + "test-buildUpSetupAndTestIt-agent2"
	//Endpoint
	endpointName2 := "api_test_endpoint2"
	address2 := configManagementTest.Agent2.EndpointAddress + ":" + strconv.Itoa(configManagementTest.Agent2.EndpointPort[0])
	//User
	name2 := "test-buildUpSetupAndTestIt-user2"
	userIdentifier2 := "test-buildUpSetupAndTestIt-user2"
	//Engine
	engineName2 := "test-buildUpSetupAndTestIt-engine2"
	engineId2 := "010203040507080A"
	//Record File
	localRecordFilePath2 := configManagementTest.TestDataDir + "snmprecs/TestManagementClient_buildUpSetupAndTestIt/agent2/" + community + ".snmprec"
	remoteRecordFilePath2 := agentDataDir2 + "/" + community + ".snmprec"

	//Create a new api client
	client, err := NewManagementClient(configManagementTest.Http.BaseUrl)
	if !assert.NoError(t, err, "error while creating a new api client") {
		return
	}
	//Set configManagementTest.Http.AuthUsername and password
	if configManagementTest.Http.AuthUsername != "" && configManagementTest.Http.AuthPassword != "" {
		err = client.SetUsernameAndPassword(configManagementTest.Http.AuthUsername, configManagementTest.Http.AuthPassword)
		if !assert.NoError(t, err, "error while creating a new api client") {
			return
		}
	}

	//Create a new lab
	lab, err := createLabAndCheckForSuccess(t, client, "TestManagementClient_buildUpSetupAndTestIt")
	if err != nil {
		return
	}
	//Clean up: delete lab
	defer func() {
		err = deleteLabAndCheckForSuccess(t, client, lab)
		assert.NoError(t, err, "error during delete lab")
	}()

	//Record file agent 1
	//TODO: remove this when its possible to overwrite files
	err = client.DeleteRecordFile(remoteRecordFilePath1)
	if err != nil {
		if err, ok := err.(HttpError); assert.True(t, ok, "unknown error returned while deleting record file") {
			if !assert.True(t, err.StatusCode == 404, "http error code for deleting record file is not 404! error: "+err.Error()) {
				return
			}
		} else {
			return
		}
	}

	err = uploadRecordFileAndCheckForSuccess(t, client, localRecordFilePath1, remoteRecordFilePath1)
	if err != nil {
		return
	}
	//Clean Up
	defer func() {
		err = deleteRecordFileAndCheckForSuccess(t, client, remoteRecordFilePath1)
		assert.NoError(t, err, "error during delete record file")
	}()

	//Record file agent 2
	//TODO: remove this when its possible to overwrite files
	err = client.DeleteRecordFile(remoteRecordFilePath2)
	if err != nil {
		if err, ok := err.(HttpError); assert.True(t, ok, "unknown error returned while deleting record file") {
			if !assert.True(t, err.StatusCode == 404, "http error code for deleting record file is not 404! error: "+err.Error()) {
				return
			}
		} else {
			return
		}
	}

	err = uploadRecordFileAndCheckForSuccess(t, client, localRecordFilePath2, remoteRecordFilePath2)
	if err != nil {
		return
	}
	//Clean Up
	defer func() {
		err = deleteRecordFileAndCheckForSuccess(t, client, remoteRecordFilePath2)
		assert.NoError(t, err, "error during delete record file")
	}()

	/*--------------------
			Agent 1
	  --------------------*/

	//Create an engine1
	engine1, err := createEngineAndCheckForSuccess(t, client, engineName1, engineId1)
	if err != nil {
		return
	}
	//Cleanup: delete engine1
	defer func() {
		err = deleteEngineAndCheckForSuccess(t, client, engine1)
		assert.NoError(t, err, "error during delete engine")
	}()

	//Create endpoint1
	endpoint1, err := createEndpointAndCheckForSuccess(t, client, endpointName1, address1, configManagementTest.Protocol)
	if err != nil {
		return
	}
	//Cleanup: delete endpoint1
	defer func() {
		err = deleteEndpointAndCheckForSuccess(t, client, endpoint1)
		assert.NoError(t, err, "error during delete endpoint")
	}()

	//Create user1
	user1, err := createUserAndCheckForSuccess(t, client, userIdentifier1, name1, authKey1, authProto1, privKey1, privProto1)
	if err != nil {
		return
	}
	//Cleanup: delete user1
	defer func() {
		err = deleteUserAndCheckForSuccess(t, client, user1)
		assert.NoError(t, err, "error during delete user")
	}()

	//Add User1 to Engine1
	err = addUserToEngineAndCheckForSuccess(t, client, engine1, user1)
	if err != nil {
		return
	}
	//Cleanup: remove user1 from engine1
	defer func() {
		err = removeUserFromEngineAndCheckForSuccess(t, client, engine1, user1)
		assert.NoError(t, err, "error during remove user from engine")
	}()

	//Add Endpoint1 to Engine1
	err = addEndpointToEngineAndCheckForSuccess(t, client, engine1, endpoint1)
	if err != nil {
		return
	}
	//Cleanup: remove endpoint1 from engine1
	defer func() {
		//Remove endpoint1 from engine1
		err = removeEndpointFromEngineAndCheckForSuccess(t, client, engine1, endpoint1)
		assert.NoError(t, err, "error during remove endpoint from engine")
	}()

	//Create agent1
	agent1, err := createAgentAndCheckForSuccess(t, client, agentName1, agentDataDir1)
	if err != nil {
		return
	}
	//Clean up: delete agent1
	defer func() {
		err = deleteAgentAndCheckForSuccess(t, client, agent1)
		assert.NoError(t, err, "error during delete agent")
	}()

	//Add engine1 to agent1
	err = addEngineToAgentAndCheckForSuccess(t, client, agent1, engine1)
	if err != nil {
		return
	}
	//Cleanup: remove engine1 from agent1
	defer func() {
		err = removeEngineFromAgentAndCheckForSuccess(t, client, agent1, engine1)
		assert.NoError(t, err, "error during remove engine from agent")
	}()

	//Add agent1 to lab
	err = addAgentToLabAndCheckForSuccess(t, client, lab, agent1)
	if err != nil {
		return
	}
	//Cleanup: remove agent1 from lab
	defer func() {
		err = removeAgentFromLabAndCheckForSuccess(t, client, lab, agent1)
		assert.NoError(t, err, "error during remove agent from lab")
	}()

	/*--------------------
			Agent 2
	  --------------------*/

	//Create an engine2
	engine2, err := createEngineAndCheckForSuccess(t, client, engineName2, engineId2)
	if err != nil {
		return
	}
	//Cleanup: delete engine2
	defer func() {
		err = deleteEngineAndCheckForSuccess(t, client, engine2)
		assert.NoError(t, err, "error during delete engine")
	}()

	//Create endpoint2
	endpoint2, err := createEndpointAndCheckForSuccess(t, client, endpointName2, address2, configManagementTest.Protocol)
	if err != nil {
		return
	}
	//Cleanup: delete endpoint2
	defer func() {
		err = deleteEndpointAndCheckForSuccess(t, client, endpoint2)
		assert.NoError(t, err, "error during delete endpoint")
	}()

	//Create user2
	user2, err := createUserAndCheckForSuccess(t, client, userIdentifier2, name2, "", "", "", "")
	if err != nil {
		return
	}
	//Cleanup: delete user2
	defer func() {
		err = deleteUserAndCheckForSuccess(t, client, user2)
		assert.NoError(t, err, "error during delete user")
	}()

	//Add User2 to Engine2
	err = addUserToEngineAndCheckForSuccess(t, client, engine2, user2)
	if err != nil {
		return
	}
	//Cleanup: remove user2 from engine2
	defer func() {
		err = removeUserFromEngineAndCheckForSuccess(t, client, engine2, user2)
		assert.NoError(t, err, "error during remove user from engine")
	}()

	//Add endpoint2 to engine2
	err = addEndpointToEngineAndCheckForSuccess(t, client, engine2, endpoint2)
	if err != nil {
		return
	}
	defer func() {
		err = removeEndpointFromEngineAndCheckForSuccess(t, client, engine2, endpoint2)
		assert.NoError(t, err, "error during remove endpoint from engine")
	}()

	//Create agent2
	agent2, err := createAgentAndCheckForSuccess(t, client, agentName2, agentDataDir2)
	if err != nil {
		return
	}
	//Clean up: delete agent2
	defer func() {
		err = deleteAgentAndCheckForSuccess(t, client, agent2)
		assert.NoError(t, err, "error during delete agent")
	}()

	//Add engine2 to agent2
	err = addEngineToAgentAndCheckForSuccess(t, client, agent2, engine2)
	if err != nil {
		return
	}
	//Cleanup: remove engine2 from agent2
	defer func() {
		err = removeEngineFromAgentAndCheckForSuccess(t, client, agent2, engine2)
		assert.NoError(t, err, "error during remove engine from agent")
	}()

	//Add agent2 to lab
	err = addAgentToLabAndCheckForSuccess(t, client, lab, agent2)
	if err != nil {
		return
	}
	//Cleanup: remove agent2 from lab
	defer func() {
		err = removeAgentFromLabAndCheckForSuccess(t, client, lab, agent2)
		assert.NoError(t, err, "error during remove agent from lab")
	}()

	//Power on lab
	err = setLabPowerAndCheckForSuccess(t, client, lab, true)
	if err != nil {
		return
	}
	//Cleanup: turn lab off
	defer func() {
		err = setLabPowerAndCheckForSuccess(t, client, lab, false)
		assert.NoError(t, err, "error during power off lab")
	}()

	/*--------------------
		  SNMP Requests
	  --------------------*/

	//SNMP Request Agent 1 SNMPv3

	agent1Snmpv3 := &gosnmp.GoSNMP{
		Target:  configManagementTest.Agent1.EndpointAddress,
		Port:    uint16(configManagementTest.Agent1.EndpointPort[0]),
		Timeout: time.Duration(2) * time.Second,
		Version: gosnmp.Version3,
		SecurityParameters: &gosnmp.UsmSecurityParameters{UserName: userIdentifier1,
			AuthenticationProtocol:   gosnmp.SHA,
			AuthenticationPassphrase: authKey1,
			PrivacyProtocol:          gosnmp.DES,
			PrivacyPassphrase:        privKey1,
		},
		MsgFlags:      gosnmp.AuthPriv,
		SecurityModel: gosnmp.UserSecurityModel,
		Transport:     "udp",
		ContextName:   community,
	}

	err = agent1Snmpv3.ConnectIPv4()
	if !assert.NoError(t, err, "error during snmp connect v4") {
		return
	}
	defer func() {
		err = agent1Snmpv3.Conn.Close()
		assert.NoError(t, err, "error during snmp connection close")
	}()

	for i := 1; i <= 12; i++ {
		res, err := agent1Snmpv3.Get([]string{"1.3.6.1.2.1.1.1.0"})
		if err != nil && i < 36 {
			time.Sleep(1 * time.Second)
			continue
		}
		if assert.NoError(t, err, "error during snmp get request") {
			resultByte, ok := res.Variables[0].Value.([]byte)
			if assert.True(t, ok, "failed to convert result to string") {
				resultString := string(resultByte)
				assert.True(t, resultString == "agent1-test-record", "snmpget result is not the expected value! result: "+resultString+" (expected: agent1-test-record)")
			}
		}
		break
	}

	//SNMP Request Agent 1 SNMPv2
	agent1Snmpv2c := &gosnmp.GoSNMP{
		Target:    configManagementTest.Agent1.EndpointAddress,
		Port:      uint16(configManagementTest.Agent1.EndpointPort[0]),
		Timeout:   time.Duration(2) * time.Second,
		Version:   gosnmp.Version2c,
		Community: community,
		Transport: "udp",
	}
	err = agent1Snmpv2c.ConnectIPv4()
	if !assert.NoError(t, err, "error during snmp connect v4") {
		return
	}
	defer func() {
		err = agent1Snmpv2c.Conn.Close()
		assert.NoError(t, err, "error during snmp connection close")
	}()
	res, err := agent1Snmpv2c.Get([]string{"1.3.6.1.2.1.1.1.0"})
	if assert.NoError(t, err, "error during snmp get request") {
		resultByte, ok := res.Variables[0].Value.([]byte)
		if assert.True(t, ok, "failed to convert result to string") {
			resultString := string(resultByte)
			assert.True(t, resultString == "agent1-test-record", "snmpget result is not the expected value! result: "+resultString+" (expected: agent1-test-record)")
		}
	}

	//SNMP Request Agent 2 SNMPv2
	agent2Snmpv2c := &gosnmp.GoSNMP{
		Target:    configManagementTest.Agent2.EndpointAddress,
		Port:      uint16(configManagementTest.Agent2.EndpointPort[0]),
		Timeout:   time.Duration(2) * time.Second,
		Version:   gosnmp.Version2c,
		Community: community,
		Transport: "udp",
	}
	err = agent2Snmpv2c.ConnectIPv4()
	if !assert.NoError(t, err, "error during snmp connect v4") {
		return
	}
	defer func() {
		err = agent2Snmpv2c.Conn.Close()
		assert.NoError(t, err, "error during snmp connection close")
	}()

	res, err = agent2Snmpv2c.Get([]string{"1.3.6.1.2.1.1.1.0"})

	if assert.NoError(t, err, "error during snmp get request") {
		resultByte, ok := res.Variables[0].Value.([]byte)
		if assert.True(t, ok, "failed to convert result to string") {
			resultString := string(resultByte)
			assert.True(t, resultString == "agent2-test-record", "snmpget result is not the expected value! result: "+resultString+" (expected: agent2-test-record)")
		}
	}
}

func TestManagementClient_UploadRecordFileString_DeleteRecordFile(t *testing.T) {

	fileContent := `1.3.6.1.2.1.1.1.0|4|testFile
	1.3.6.1.2.1.1.2.0|6|1.3.6.1.4.1.8072.3.2.10
	1.3.6.1.2.1.1.3.0|67|123999999`

	remotePathFile1 := "test-UploadRecordFileString_DeleteRecordFile/dir1/dir2/public.snmprec"
	remotePathFile2 := "test-UploadRecordFileString_DeleteRecordFile/dir1/public.snmprec"

	//Create a new api client
	client, err := NewManagementClient(configManagementTest.Http.BaseUrl)
	if !assert.NoError(t, err, "error while creating a new api client") {
		return
	}
	//Set configManagementTest.Http.AuthUsername and password
	if configManagementTest.Http.AuthUsername != "" && configManagementTest.Http.AuthPassword != "" {
		err = client.SetUsernameAndPassword(configManagementTest.Http.AuthUsername, configManagementTest.Http.AuthPassword)
		if !assert.NoError(t, err, "error while creating a new api client") {
			return
		}
	}

	//init cleanup
	err = client.DeleteRecordFile(remotePathFile1)
	if err != nil {
		if err, ok := err.(HttpError); assert.True(t, ok, "error while initial delete is not a http error"+err.Error()) {
			if !assert.True(t, err.StatusCode == 404, "cleanup delete error != 404") {
				return
			}
		} else {
			return
		}
	}
	err = client.DeleteRecordFile(remotePathFile2)
	if err != nil {
		if err, ok := err.(HttpError); assert.True(t, ok, "error while initial delete is not a http error"+err.Error()) {
			if !assert.True(t, err.StatusCode == 404, "cleanup delete error != 404") {
				return
			}
		} else {
			return
		}
	}

	err = uploadRecordFileStringAndCheckForSuccess(t, client, &fileContent, remotePathFile1)
	if err != nil {
		return
	}

	defer func() {
		err = deleteRecordFileAndCheckForSuccess(t, client, remotePathFile1)
		assert.NoError(t, err, "error while deleting record file")
	}()

	err = uploadRecordFileStringAndCheckForSuccess(t, client, &fileContent, remotePathFile2)
	if err != nil {
		return
	}

	defer func() {
		err = deleteRecordFileAndCheckForSuccess(t, client, remotePathFile2)
		assert.NoError(t, err, "error while deleting record file")
	}()

	//TODO: this should cause an api error but does not
	/*
		//upload invalid record file
		invalidRecord := "invalid\record\file"
		err = client.UploadRecordFileString(&invalidRecord, "invalid/record/file.snmprec")
		if assert.Error(t, err, "no error when uploading invalid record file") {
			if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
				assert.True(t, err.StatusCode == 400, "error != 400")
			}
		}
	*/
}

func TestManagementClient_Agent_Failures(t *testing.T) {
	//Create a new api client
	client, err := NewManagementClient(configManagementTest.Http.BaseUrl)
	if !assert.NoError(t, err, "error while creating a new api client") {
		return
	}
	//Set configManagementTest.Http.AuthUsername and password
	if configManagementTest.Http.AuthUsername != "" && configManagementTest.Http.AuthPassword != "" {
		err = client.SetUsernameAndPassword(configManagementTest.Http.AuthUsername, configManagementTest.Http.AuthPassword)
		if !assert.NoError(t, err, "error while creating a new api client") {
			return
		}
	}

	//TODO: this should cause an api error but does not
	/*
		//Create Agent with invalid data dir
		_, err = client.CreateAgent("name", "test-CreateAgent_Failure")
		if assert.Error(t, err, "no error when an agent with an invalid data dir was created") {
			if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
				assert.True(t, err.StatusCode == 404, "error != 404")
			}
		}
	*/

	//Get Invalid Agent
	_, err = client.GetAgent(-1)
	if assert.Error(t, err, "no error when trying to get an invalid agent") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 404, "error != 404")
		}
	}

	//Delete invalid agent
	err = client.DeleteAgent(-1)
	if assert.Error(t, err, "no error when a non existent agent was deleted") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 404, "error != 404")
		}
	}

	//create valid agent
	agent, err := createAgentAndCheckForSuccess(t, client, "test-Agent_Failures-agent1", ".")
	if err != nil {
		return
	}
	defer func() { _ = deleteAgentAndCheckForSuccess(t, client, agent) }()

	//TODO: this should cause an api error but does not
	/*
		//add invalid engine to agent
		err = client.AddEngineToAgent(agent.Id, -1)
		if assert.Error(t, err, "no error when an invalid engine id was added to an agent") {
			if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
				assert.True(t, err.StatusCode == 404, "error != 404" , err.Error())
			}
		}
	*/

	//create valid engine
	engine, err := createEngineAndCheckForSuccess(t, client, "test-Agent_Failures-engine1", "010203040507080B")
	if err != nil {
		return
	}
	defer func() { _ = deleteEngineAndCheckForSuccess(t, client, engine) }()

	//add valid engine to invalid agent
	err = client.AddEngineToAgent(-1, engine.Id)
	if assert.Error(t, err, "no error when an existent engine was added to a non existing agent") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 400, "error != 400")
		}
	}

	//add engine to agent
	err = addEngineToAgentAndCheckForSuccess(t, client, agent, engine)
	if err != nil {
		return
	}
	//its removed later, no defer needed

	//add already attached engine to agent
	err = client.AddEngineToAgent(agent.Id, engine.Id)
	if assert.Error(t, err, "no error when an engine was added twice to an agent") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 400, "error != 400")
		}
	}

	//remove non existing engine from non existing agent
	err = client.RemoveEngineFromAgent(-1, -1)
	if assert.Error(t, err, "no error when removing non existing engine from non existing agent") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 404, "error != 404")
		}
	}

	//remove non existing engine from exisiting agent
	err = client.RemoveEngineFromAgent(agent.Id, -1)
	if assert.Error(t, err, "no error when removing non existing engine from existing agent") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 404, "error != 404", err.Error())
		}
	}

	//remove existing engine from non exisiting agent
	err = client.RemoveEngineFromAgent(-1, engine.Id)
	if assert.Error(t, err, "no error when removing existing engine from non existing agent") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 400, "error != 400", err.Error())
		}
	}

	err = removeEngineFromAgentAndCheckForSuccess(t, client, agent, engine)
	if err != nil {
		return
	}
	err = client.RemoveEngineFromAgent(agent.Id, engine.Id)
	if assert.Error(t, err, "no error when removing an engine from an agent that is not attached to the agent") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 404, "error != 404", err.Error())
		}
	}
}

func TestManagementClient_Lab_Failures(t *testing.T) {
	//Create a new api client
	client, err := NewManagementClient(configManagementTest.Http.BaseUrl)
	if !assert.NoError(t, err, "error while creating a new api client") {
		return
	}
	//Set configManagementTest.Http.AuthUsername and password
	if configManagementTest.Http.AuthUsername != "" && configManagementTest.Http.AuthPassword != "" {
		err = client.SetUsernameAndPassword(configManagementTest.Http.AuthUsername, configManagementTest.Http.AuthPassword)
		if !assert.NoError(t, err, "error while creating a new api client") {
			return
		}
	}

	//TODO: this should cause an api error but does not

	//Get Invalid Lab
	_, err = client.GetLab(-1)
	if assert.Error(t, err, "no error when trying to get an invalid lab") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 404, "error != 404")
		}
	}

	//Delete invalid lab
	err = client.DeleteLab(-1)
	if assert.Error(t, err, "no error when a non existent lab was deleted") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 404, "error != 404")
		}
	}

	//create valid lab
	lab, err := createLabAndCheckForSuccess(t, client, "test-Lab_Failures-lab1")
	if err != nil {
		return
	}
	defer func() { _ = deleteLabAndCheckForSuccess(t, client, lab) }()

	//TODO: this should cause an api error but does not
	/*
		//add invalid agent to lab
		err = client.AddAgentToLab(lab.Id, -1)
		if assert.Error(t, err, "no error when an invalid agent id was added to an lab") {
			if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
				assert.True(t, err.StatusCode == 404, "error != 404", err.Error())
			}
		}
	*/

	//create valid agent
	agent, err := createAgentAndCheckForSuccess(t, client, "test-Lab_Failures-agent1", ".")
	if err != nil {
		return
	}
	defer func() { _ = deleteAgentAndCheckForSuccess(t, client, agent) }()

	//add valid agent to invalid lab
	err = client.AddAgentToLab(-1, agent.Id)
	if assert.Error(t, err, "no error when an existent agent was added to a non existing lab") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 400, "error != 400")
		}
	}

	//add agent to lab
	err = addAgentToLabAndCheckForSuccess(t, client, lab, agent)
	if err != nil {
		return
	}
	//its removed later, no defer needed

	//add already attached agent to lab
	err = client.AddAgentToLab(lab.Id, agent.Id)
	if assert.Error(t, err, "no error when an agent was added twice to an lab") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 400, "error != 400")
		}
	}

	//remove non existing agent from non existing lab
	err = client.RemoveAgentFromLab(-1, -1)
	if assert.Error(t, err, "no error when removing non existing agent from non existing lab") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 404, "error != 404")
		}
	}

	//remove non existing agent from exisiting lab
	err = client.RemoveAgentFromLab(lab.Id, -1)
	if assert.Error(t, err, "no error when removing non existing agent from existing lab") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 404, "error != 404", err.Error())
		}
	}

	//remove existing agent from non exisiting lab
	err = client.RemoveAgentFromLab(-1, agent.Id)
	if assert.Error(t, err, "no error when removing existing agent from non existing lab") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 400, "error != 400", err.Error())
		}
	}

	err = removeAgentFromLabAndCheckForSuccess(t, client, lab, agent)
	if err != nil {
		return
	}
	err = client.RemoveAgentFromLab(lab.Id, agent.Id)
	if assert.Error(t, err, "no error when removing an agent from an lab that is not attached to the lab") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 404, "error != 404", err.Error())
		}
	}
}

func TestManagementClient_Engine_Failures(t *testing.T) {
	//Create a new api client
	client, err := NewManagementClient(configManagementTest.Http.BaseUrl)
	if !assert.NoError(t, err, "error while creating a new api client") {
		return
	}
	//Set configManagementTest.Http.AuthUsername and password
	if configManagementTest.Http.AuthUsername != "" && configManagementTest.Http.AuthPassword != "" {
		err = client.SetUsernameAndPassword(configManagementTest.Http.AuthUsername, configManagementTest.Http.AuthPassword)
		if !assert.NoError(t, err, "error while creating a new api client") {
			return
		}
	}

	//TODO: this should cause an api error but does not
	/*
		//Create Engine with invalid params
		_, err = client.CreateEngine("name", "this is not a valid engine id")
		if assert.Error(t, err, "no error when an engine with an invalid engine id was created") {
			if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
				assert.True(t, err.StatusCode == 404, "error != 404")
			}
		}
	*/

	//Get Invalid Engine
	_, err = client.GetEngine(-1)
	if assert.Error(t, err, "no error when trying to get an invalid engine") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 404, "error != 404")
		}
	}

	//Delete invalid engine
	err = client.DeleteEngine(-1)
	if assert.Error(t, err, "no error when a non existent engine was deleted") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 404, "error != 404")
		}
	}

	//create valid engine
	engine, err := createEngineAndCheckForSuccess(t, client, "test-Engine_Failures-engine1", "010203040507080B")
	if err != nil {
		return
	}
	defer func() { _ = deleteEngineAndCheckForSuccess(t, client, engine) }()

	//TODO: this should cause an api error but does not
	/*
		_, err = client.CreateEngine("test-Engine_Failures-engine1", "010203040507080B")
		if assert.Error(t, err, "no error when an engine was created twice") {
			if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
				assert.True(t, err.StatusCode == 404, "error != 404")
			}
		}
	*/

	//TODO: this should cause an api error but does not
	/*
		//add invalid user to engine
		err = client.AddUserToEngine(engine.Id, -1)
		if assert.Error(t, err, "no error when an invalid user id was added to an engine") {
			if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
				assert.True(t, err.StatusCode == 404, "error != 404" , err.Error())
			}
		}
	*/

	//create valid user
	user, err := createUserAndCheckForSuccess(t, client, "test-Engine_Failures-user1", "test-Engine_Failures-user1", "", "", "", "")
	if err != nil {
		return
	}
	defer func() { _ = deleteUserAndCheckForSuccess(t, client, user) }()

	//add valid user to invalid engine
	err = client.AddUserToEngine(-1, user.Id)
	if assert.Error(t, err, "no error when an existent user was added to a non existing engine") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 400, "error != 400")
		}
	}

	//add user to engine
	err = addUserToEngineAndCheckForSuccess(t, client, engine, user)
	if err != nil {
		return
	}
	//its removed later, no defer needed

	//add already attached user to engine
	err = client.AddUserToEngine(engine.Id, user.Id)
	if assert.Error(t, err, "no error when an user was added twice to an engine") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 400, "error != 400")
		}
	}

	//remove non existing user from non existing engine
	err = client.RemoveUserFromEngine(-1, -1)
	if assert.Error(t, err, "no error when removing non existing user from non existing engine") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 404, "error != 404")
		}
	}

	//remove non existing user from exisiting engine
	err = client.RemoveUserFromEngine(engine.Id, -1)
	if assert.Error(t, err, "no error when removing non existing user from existing engine") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 404, "error != 404", err.Error())
		}
	}

	//remove existing user from non exisiting engine
	err = client.RemoveUserFromEngine(-1, user.Id)
	if assert.Error(t, err, "no error when removing existing user from non existing engine") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 400, "error != 400", err.Error())
		}
	}

	err = removeUserFromEngineAndCheckForSuccess(t, client, engine, user)
	if err != nil {
		return
	}
	err = client.RemoveUserFromEngine(engine.Id, user.Id)
	if assert.Error(t, err, "no error when removing an user from an engine that is not attached to the engine") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 404, "error != 404", err.Error())
		}
	}

	//TODO: this should cause an api error but does not
	/*
		//add invalid endpoint to engine
		err = client.AddEndpointToEngine(engine.Id, -1)
		if assert.Error(t, err, "no error when an invalid endpoint id was added to an engine") {
			if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
				assert.True(t, err.StatusCode == 404, "error != 404" , err.Error())
			}
		}
	*/

	//create valid endpoint
	endpoint, err := createEndpointAndCheckForSuccess(t, client, "test-Engine_Failures-endpoint1", configManagementTest.Agent1.EndpointAddress+":"+strconv.Itoa(configManagementTest.Agent1.EndpointPort[0]), configManagementTest.Protocol)
	if err != nil {
		return
	}
	defer func() { _ = deleteEndpointAndCheckForSuccess(t, client, endpoint) }()

	//add valid endpoint to invalid engine
	err = client.AddEndpointToEngine(-1, endpoint.Id)
	if assert.Error(t, err, "no error when an existent endpoint was added to a non existing engine") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 400, "error != 400")
		}
	}

	//add endpoint to engine
	err = addEndpointToEngineAndCheckForSuccess(t, client, engine, endpoint)
	if err != nil {
		return
	}
	//its removed later, no defer needed

	//add already attached endpoint to engine
	err = client.AddEndpointToEngine(engine.Id, endpoint.Id)
	if assert.Error(t, err, "no error when an endpoint was added twice to an engine") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 400, "error != 400")
		}
	}

	//remove non existing endpoint from non existing engine
	err = client.RemoveEndpointFromEngine(-1, -1)
	if assert.Error(t, err, "no error when removing non existing endpoint from non existing engine") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 404, "error != 404")
		}
	}

	//remove non existing endpoint from exisiting engine
	err = client.RemoveEndpointFromEngine(engine.Id, -1)
	if assert.Error(t, err, "no error when removing non existing endpoint from existing engine") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 404, "error != 404", err.Error())
		}
	}

	//remove existing endpoint from non exisiting engine
	err = client.RemoveEndpointFromEngine(-1, endpoint.Id)
	if assert.Error(t, err, "no error when removing existing endpoint from non existing engine") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 400, "error != 400", err.Error())
		}
	}

	err = removeEndpointFromEngineAndCheckForSuccess(t, client, engine, endpoint)
	if err != nil {
		return
	}
	err = client.RemoveEndpointFromEngine(engine.Id, endpoint.Id)
	if assert.Error(t, err, "no error when removing an endpoint from an engine that is not attached to the engine") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 404, "error != 404", err.Error())
		}
	}
}

func TestManagementClient_User_Failures(t *testing.T) {
	//Create a new api client
	client, err := NewManagementClient(configManagementTest.Http.BaseUrl)
	if !assert.NoError(t, err, "error while creating a new api client") {
		return
	}
	//Set configManagementTest.Http.AuthUsername and password
	if configManagementTest.Http.AuthUsername != "" && configManagementTest.Http.AuthPassword != "" {
		err = client.SetUsernameAndPassword(configManagementTest.Http.AuthUsername, configManagementTest.Http.AuthPassword)
		if !assert.NoError(t, err, "error while creating a new api client") {
			return
		}
	}

	//TODO: this should cause an api error but does not
	/*
		//Create User with invalid params
		_, err = client.CreateUser("test-User_Failures-user1", "test-User_Failures-user1")
		if assert.Error(t, err, "no error when an user with invalid params was created") {
			if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
				assert.True(t, err.StatusCode == 404, "error != 404")
			}
		}
	*/

	//Get Invalid User
	_, err = client.GetUser(-1)
	if assert.Error(t, err, "no error when trying to get an invalid user") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 404, "error != 404")
		}
	}

	//Delete invalid user
	err = client.DeleteUser(-1)
	if assert.Error(t, err, "no error when a non existent user was deleted") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 404, "error != 404")
		}
	}

	//create valid user
	user, err := createUserAndCheckForSuccess(t, client, "test-User_Failures-user1", "test-User_Failures-user1", "", "", "", "")
	if err != nil {
		return
	}
	defer func() { _ = deleteUserAndCheckForSuccess(t, client, user) }()

	_, err = client.CreateUser("test-User_Failures-user1", "test-User_Failures-user1", "", "", "", "")
	if assert.Error(t, err, "no error when creating a user twice") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 400, "error != 400")
		}
	}
}

func TestManagementClient_Endpoint_Failures(t *testing.T) {
	//Create a new api client
	client, err := NewManagementClient(configManagementTest.Http.BaseUrl)
	if !assert.NoError(t, err, "error while creating a new api client") {
		return
	}
	//Set username and password
	if configManagementTest.Http.AuthUsername != "" && configManagementTest.Http.AuthPassword != "" {
		err = client.SetUsernameAndPassword(configManagementTest.Http.AuthUsername, configManagementTest.Http.AuthPassword)
		if !assert.NoError(t, err, "error while creating a new api client") {
			return
		}
	}

	//Create Endpoint with invalid input
	_, err = client.CreateEndpoint("test-Endpoint_Failures-endpoint1", "noAddress", "no valid protocol")
	if assert.Error(t, err, "no error when an endpoint with invalid params was created") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 400, "error != 400")
		}
	}

	//Get Invalid Endpoint
	_, err = client.GetEndpoint(-1)
	if assert.Error(t, err, "no error when trying to get an invalid endpoint") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 404, "error != 404")
		}
	}

	//Delete invalid endpoint
	err = client.DeleteEndpoint(-1)
	if assert.Error(t, err, "no error when a non existent endpoint was deleted") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 404, "error != 404")
		}
	}

	//create valid endpoint
	endpoint, err := createEndpointAndCheckForSuccess(t, client, "test-Endpoint_Failures-endpoint1", "1.1.1.1:9753", "udpv4")
	if err != nil {
		return
	}
	defer func() { _ = deleteEndpointAndCheckForSuccess(t, client, endpoint) }()

	_, err = client.CreateEndpoint("test-Endpoint_Failures-endpoint1", "1.1.1.1:9753", "udpv4")
	if assert.Error(t, err, "no error when creating a endpoint twice") {
		if err, ok := err.(HttpError); assert.True(t, ok, "error is not a http error", err.Error()) {
			assert.True(t, err.StatusCode == 400, "error != 400")
		}
	}
}
