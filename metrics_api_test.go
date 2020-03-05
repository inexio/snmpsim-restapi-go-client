package snmpsimclient

import (
	"github.com/pkg/errors"
	"github.com/soniah/gosnmp"
	"github.com/stretchr/testify/assert"
	"regexp"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestMetricsClient_BuildUpSetupAndTestMetrics(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestMetricsClient_BuildUpSetupAndTestMetrics in short mode")
	}
	community := "public"
	//	Agent 1
	//Agent
	agentName1 := "test-buildUpSetupAndTestMetrics-agent1"
	agentDataDir1 := configMetricsTest.RootDataDir + "test-buildUpSetupAndTestMetrics-agent1"
	//Endpoint
	endpointName1 := "test-buildUpSetupAndTestMetrics-endpoint1"
	address1 := configMetricsTest.Agent1.EndpointAddress + ":" + strconv.Itoa(configMetricsTest.Agent1.EndpointPort[0]) //only start requests
	//Endpoint
	endpointName2 := "test-buildUpSetupAndTestMetrics-endpoint2"
	address2 := configMetricsTest.Agent1.EndpointAddress + ":" + strconv.Itoa(configMetricsTest.Agent1.EndpointPort[1]) //many requests

	//User
	name1 := "test-buildUpSetupAndTestMetrics"
	userIdentifier1 := "test-buildUpSetupAndTestMetrics"
	//engine
	engineName1 := "test-buildUpSetupAndTestMetrics-engine1"
	engineId1 := "0102030405070809"
	//Record File:
	localRecordFilePath1 := configMetricsTest.TestDataDir + "snmprecs/TestMetricsClient_BuildUpSetupAndTestMetrics/" + community + ".snmprec"
	remoteRecordFilePath1 := agentDataDir1 + "/" + community + ".snmprec"

	//Create a new api client
	managementClient, err := NewManagementClient(configManagementTest.Http.BaseUrl)
	if !assert.NoError(t, err, "error while creating a new api client") {
		return
	}
	//Set configMetricsTest.Http.AuthUsername and password
	if configManagementTest.Http.AuthUsername != "" && configManagementTest.Http.AuthPassword != "" {
		err = managementClient.SetUsernameAndPassword(configManagementTest.Http.AuthUsername, configManagementTest.Http.AuthPassword)
		if !assert.NoError(t, err, "error while creating a new api client") {
			return
		}
	}

	//Record file
	//TODO: remove this when its possible to overwrite files
	err = managementClient.DeleteRecordFile(remoteRecordFilePath1)
	if err != nil {
		if err, ok := err.(HttpError); assert.True(t, ok, "unknown error returned while deleting record file") {
			if !assert.True(t, err.StatusCode == 404, "http error code for deleting record file is not 404! error: "+err.Error()) {
				return
			}
		} else {
			return
		}
	}

	err = uploadRecordFileAndCheckForSuccess(t, managementClient, localRecordFilePath1, remoteRecordFilePath1)
	if err != nil {
		return
	}
	//Clean Up
	defer func() {
		err = deleteRecordFileAndCheckForSuccess(t, managementClient, remoteRecordFilePath1)
		assert.NoError(t, err, "error during delete record file")
	}()

	//Create a new lab
	lab, err := createLabAndCheckForSuccess(t, managementClient, "TestMetricsClient_BuildUpSetupAndTestMetrics")
	if err != nil {
		return
	}
	//Clean up: delete lab
	defer func() {
		err = deleteLabAndCheckForSuccess(t, managementClient, lab)
		assert.NoError(t, err, "error during delete lab")
	}()

	//Create an engine1
	engine1, err := createEngineAndCheckForSuccess(t, managementClient, engineName1, engineId1)
	if err != nil {
		return
	}
	//Cleanup: delete engine1
	defer func() {
		err = deleteEngineAndCheckForSuccess(t, managementClient, engine1)
		assert.NoError(t, err, "error during delete engine")
	}()

	//Create endpoint1
	endpoint1, err := createEndpointAndCheckForSuccess(t, managementClient, endpointName1, address1, configMetricsTest.Protocol)
	if err != nil {
		return
	}
	//Cleanup: delete endpoint1
	defer func() {
		err = deleteEndpointAndCheckForSuccess(t, managementClient, endpoint1)
		assert.NoError(t, err, "error during delete endpoint")
	}()

	//Create endpoint2
	endpoint2, err := createEndpointAndCheckForSuccess(t, managementClient, endpointName2, address2, configMetricsTest.Protocol)
	if err != nil {
		return
	}
	//Cleanup: delete endpoint2
	defer func() {
		err = deleteEndpointAndCheckForSuccess(t, managementClient, endpoint2)
		assert.NoError(t, err, "error during delete endpoint")
	}()

	//Create user1
	user1, err := createUserAndCheckForSuccess(t, managementClient, userIdentifier1, name1, "", "", "", "")
	if err != nil {
		return
	}
	//Cleanup: delete user1
	defer func() {
		err = deleteUserAndCheckForSuccess(t, managementClient, user1)
		assert.NoError(t, err, "error during delete user")
	}()

	//Add User1 to Engine1
	err = addUserToEngineAndCheckForSuccess(t, managementClient, engine1, user1)
	if err != nil {
		return
	}
	//Cleanup: remove user1 from engine1
	defer func() {
		err = removeUserFromEngineAndCheckForSuccess(t, managementClient, engine1, user1)
		assert.NoError(t, err, "error during remove user from engine")
	}()

	//Add endpoint1 to engine1
	err = addEndpointToEngineAndCheckForSuccess(t, managementClient, engine1, endpoint1)
	if err != nil {
		return
	}
	defer func() {
		err = removeEndpointFromEngineAndCheckForSuccess(t, managementClient, engine1, endpoint1)
		assert.NoError(t, err, "error during remove endpoint from engine")
	}()
	//Add endpoint2 to engine1
	err = addEndpointToEngineAndCheckForSuccess(t, managementClient, engine1, endpoint2)
	if err != nil {
		return
	}
	defer func() {
		err = removeEndpointFromEngineAndCheckForSuccess(t, managementClient, engine1, endpoint2)
		assert.NoError(t, err, "error during remove endpoint from engine")
	}()

	//Create agent1
	agent1, err := createAgentAndCheckForSuccess(t, managementClient, agentName1, agentDataDir1)
	if err != nil {
		return
	}
	//Clean up: delete agent1
	defer func() {
		err = deleteAgentAndCheckForSuccess(t, managementClient, agent1)
		assert.NoError(t, err, "error during delete agent")
	}()

	//Add engine1 to agent1
	err = addEngineToAgentAndCheckForSuccess(t, managementClient, agent1, engine1)
	if err != nil {
		return
	}
	//Cleanup: remove engine1 from agent1
	defer func() {
		err = removeEngineFromAgentAndCheckForSuccess(t, managementClient, agent1, engine1)
		assert.NoError(t, err, "error during remove engine from agent")
	}()

	//Add agent1 to lab
	err = addAgentToLabAndCheckForSuccess(t, managementClient, lab, agent1)
	if err != nil {
		return
	}
	//Cleanup: remove agent1 from lab
	defer func() {
		err = removeAgentFromLabAndCheckForSuccess(t, managementClient, lab, agent1)
		assert.NoError(t, err, "error during remove agent from lab")
	}()

	//Power on lab
	err = setLabPowerAndCheckForSuccess(t, managementClient, lab, true)
	if err != nil {
		return
	}
	//Cleanup: turn lab off
	defer func() {
		err = setLabPowerAndCheckForSuccess(t, managementClient, lab, false)
		assert.NoError(t, err, "error during power off lab")
	}()

	agent1Snmpv2c := &gosnmp.GoSNMP{
		Target:    configMetricsTest.Agent1.EndpointAddress,
		Port:      uint16(configMetricsTest.Agent1.EndpointPort[0]),
		Timeout:   time.Duration(2) * time.Second,
		Version:   gosnmp.Version2c,
		Community: "public",
		Transport: "udp",
	}
	err = agent1Snmpv2c.ConnectIPv4()
	if !assert.NoError(t, err, "error during snmp connect, cannot send request!") {
		return
	}
	defer func() {
		err = agent1Snmpv2c.Conn.Close()
		assert.NoError(t, err, "error during snmp connection close")
	}()

	for i := 0; ; i++ {
		_, err := agent1Snmpv2c.Get([]string{"1.3.6.1.2.1.1.1.0"})
		if err != nil && i < 36 {
			time.Sleep(1 * time.Second)
			continue
		}
		if !assert.NoError(t, err, "cannot succeed initial snmp request") {
			return
		}
		break
	}

	//Send SNMP Request

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go requestSender(t, &wg, community)
	}

	wg.Wait()

	//Now collect metrics

	//TODO test metric failures

	//Create a new api client
	metricsClient, err := NewMetricsClient(configMetricsTest.Http.BaseUrl)
	if !assert.NoError(t, err, "error while creating a new api client") {
		return
	}
	//Set configMetricsTest.Http.AuthUsername and password
	if configMetricsTest.Http.AuthUsername != "" && configMetricsTest.Http.AuthPassword != "" {
		err = metricsClient.SetUsernameAndPassword(configMetricsTest.Http.AuthUsername, configMetricsTest.Http.AuthPassword)
		if !assert.NoError(t, err, "error while creating a new api client") {
			return
		}
	}

	packets, err := metricsClient.GetPackets(nil)
	if assert.NoError(t, err, "error during GetPackets") {
		assert.True(t, *packets.Total > 10000, "less than 10000 packets, packets: "+strconv.FormatInt(*packets.Total, 10))
		assert.True(t, *packets.AuthFailures == 0, "more than 0 auth failures, auth failures: "+strconv.FormatInt(*packets.AuthFailures, 10))
		assert.True(t, *packets.ParseFailures == 0, "more than 0 parse failures, parse failures: "+strconv.FormatInt(*packets.ParseFailures, 10))
		assert.True(t, *packets.ContextFailures == 0, "more than 0 context failures, context failures: "+strconv.FormatInt(*packets.ContextFailures, 10))
	}

	filters := make(map[string]string)
	filters["local_address"] = address1
	packetsAddr1, errAddr1 := metricsClient.GetPackets(filters)
	if assert.NoError(t, errAddr1, "error during GetPackets for addr1") {
		assert.True(t, *packetsAddr1.Total <= *packets.Total, "filtered packet requests returned more total requests than request for all packets")
		assert.True(t, *packetsAddr1.AuthFailures <= *packets.AuthFailures, "filtered packet requests returned more auth failures than request for all packets")
		assert.True(t, *packetsAddr1.ParseFailures <= *packets.ParseFailures, "more than 0 parse failures, parse failures: "+strconv.FormatInt(*packets.ParseFailures, 10))
		assert.True(t, *packetsAddr1.ContextFailures <= *packets.ContextFailures, "more than 0 context failures, context failures: "+strconv.FormatInt(*packets.ContextFailures, 10))
	}
	filters["local_address"] = address2
	packetsAddr2, errAddr2 := metricsClient.GetPackets(filters)
	if assert.NoError(t, errAddr2, "error during GetPackets for addr2") {
		assert.True(t, *packetsAddr2.Total <= *packets.Total, "filtered packet requests returned more total requests than request for all packets")
		assert.True(t, *packetsAddr2.AuthFailures <= *packets.AuthFailures, "filtered packet requests returned more auth failures than request for all packets")
		assert.True(t, *packetsAddr2.ParseFailures <= *packets.ParseFailures, "filtered packet requests returned more parse failures than request for all packets")
		assert.True(t, *packetsAddr2.ContextFailures <= *packets.ContextFailures, "filtered packet requests returned more context failures than request for all packets")
	}
	if errAddr1 != nil && errAddr2 != nil {
		assert.True(t, *packetsAddr1.Total <= *packetsAddr2.Total, "endpoint with less incoming packet requests returned more total packets than endpoint with more requests")
	}

	messages, err := metricsClient.GetMessages(nil)
	if assert.NoError(t, err, "error during GetMessages") {
		assert.True(t, *messages.Pdus > 10000 && *messages.VarBinds > 10000, "there are less than 10000 pdus and var binds, pdus: "+strconv.FormatInt(*messages.Pdus, 10)+", var binds: "+strconv.FormatInt(*messages.VarBinds, 10))
		assert.True(t, *messages.Failures == 0, "there are failures, failures: "+strconv.FormatInt(*messages.Failures, 10))
	}

	filters = make(map[string]string)
	filters["local_address"] = address1
	messagesAddr1, errAddr1 := metricsClient.GetMessages(filters)
	if assert.NoError(t, errAddr1, "error during GetMessages for addr1") {
		assert.True(t, *messagesAddr1.Pdus <= *messages.Pdus, "filtered message requests returned more pdus than request for all packets")
		assert.True(t, *messagesAddr1.VarBinds <= *messages.VarBinds, "filtered message requests returned more var binds than request for all packets")
		assert.True(t, *messagesAddr1.Failures <= *messages.Failures, "filtered message requests returned more failures than request for all packets")
	}

	filters["local_address"] = address2
	messagesAddr2, errAddr2 := metricsClient.GetMessages(filters)
	if assert.NoError(t, errAddr2, "error during GetMessages for addr1") {
		assert.True(t, *messagesAddr2.Pdus <= *messages.Pdus, "filtered message requests returned more pdus than request for all packets")
		assert.True(t, *messagesAddr2.VarBinds <= *messages.VarBinds, "filtered message requests returned more var binds than request for all packets")
		assert.True(t, *messagesAddr2.Failures <= *messages.Failures, "filtered message requests returned more failures than request for all packets")
	}

	if errAddr1 != nil && errAddr2 != nil {
		assert.True(t, *messagesAddr1.Pdus <= *messagesAddr2.Pdus, "endpoint with less incoming message requests returned more pdus than endpoint with more requests")
		assert.True(t, *messagesAddr1.VarBinds <= *messagesAddr2.VarBinds, "endpoint with less incoming message requests returned more var binds than endpoint with more requests")
	}

	messageFilters, err := metricsClient.GetMessageFilters()
	if assert.NoError(t, err, "error during GetMessageFilters") {
		assert.True(t, len(messageFilters) > 0, "no message filters found")
	}
	possibleValuesForFilter, err := metricsClient.GetPossibleValuesForMessageFilter("local_address")
	if assert.NoError(t, err, "error during GetPossibleValuesForMessageFilter") {
		assert.True(t, len(possibleValuesForFilter) >= 2, "less than 2 values for message filter 'local_address' found")
	}

	packetFilters, err := metricsClient.GetPacketFilters()
	if assert.NoError(t, err, "error during GetPacketFilters") {
		assert.True(t, len(packetFilters) > 0, "no packet filters found")
	}
	possibleValuesForFilter, err = metricsClient.GetPossibleValuesForPacketFilter("local_address")
	if assert.NoError(t, err, "error during GetPossibleValuesForPacketFilter") {
		assert.True(t, len(possibleValuesForFilter) >= 2, "less than 2 values for packet filter 'local_address' found")
	}
}

func TestMetricsClient_Processes(t *testing.T) {
	community := "public"
	//	Agent 1
	//Agent
	agentName1 := "test-buildUpSetupAndTestMetrics-agent1"
	agentDataDir1 := configMetricsTest.RootDataDir + "test-buildUpSetupAndTestMetrics-agent1"
	//Endpoint
	endpointName1 := "test-buildUpSetupAndTestMetrics-endpoint1"
	address1 := configMetricsTest.Agent1.EndpointAddress + ":" + strconv.Itoa(configMetricsTest.Agent1.EndpointPort[0]) //only start requests
	//Endpoint
	endpointName2 := "test-buildUpSetupAndTestMetrics-endpoint2"
	address2 := configMetricsTest.Agent1.EndpointAddress + ":" + strconv.Itoa(configMetricsTest.Agent1.EndpointPort[1]) //many requests

	//User
	name1 := "test-buildUpSetupAndTestMetrics"
	userIdentifier1 := "test-buildUpSetupAndTestMetrics"
	//engine
	engineName1 := "test-buildUpSetupAndTestMetrics-engine1"
	engineId1 := "0102030405070809"
	//Record File:
	localRecordFilePath1 := configMetricsTest.TestDataDir + "snmprecs/TestMetricsClient_BuildUpSetupAndTestMetrics/" + community + ".snmprec"
	remoteRecordFilePath1 := agentDataDir1 + "/" + community + ".snmprec"

	//Create a new api client
	managementClient, err := NewManagementClient(configManagementTest.Http.BaseUrl)
	if !assert.NoError(t, err, "error while creating a new api client") {
		return
	}
	//Set configMetricsTest.Http.AuthUsername and password
	if configManagementTest.Http.AuthUsername != "" && configManagementTest.Http.AuthPassword != "" {
		err = managementClient.SetUsernameAndPassword(configManagementTest.Http.AuthUsername, configManagementTest.Http.AuthPassword)
		if !assert.NoError(t, err, "error while creating a new api client") {
			return
		}
	}

	//Record file
	//TODO: remove this when its possible to overwrite files
	err = managementClient.DeleteRecordFile(remoteRecordFilePath1)
	if err != nil {
		if err, ok := err.(HttpError); assert.True(t, ok, "unknown error returned while deleting record file") {
			if !assert.True(t, err.StatusCode == 404, "http error code for deleting record file is not 404! error: "+err.Error()) {
				return
			}
		} else {
			return
		}
	}

	//Create a new api client
	metricsClient, err := NewMetricsClient(configMetricsTest.Http.BaseUrl)
	if !assert.NoError(t, err, "error while creating a new api client") {
		return
	}
	//Set configMetricsTest.Http.AuthUsername and password
	if configMetricsTest.Http.AuthUsername != "" && configMetricsTest.Http.AuthPassword != "" {
		err = metricsClient.SetUsernameAndPassword(configMetricsTest.Http.AuthUsername, configMetricsTest.Http.AuthPassword)
		if !assert.NoError(t, err, "error while creating a new api client") {
			return
		}
	}

	err = uploadRecordFileAndCheckForSuccess(t, managementClient, localRecordFilePath1, remoteRecordFilePath1)
	if err != nil {
		return
	}
	//Clean Up
	defer func() {
		err = deleteRecordFileAndCheckForSuccess(t, managementClient, remoteRecordFilePath1)
		assert.NoError(t, err, "error during delete record file")
	}()

	//Create a new lab
	lab, err := createLabAndCheckForSuccess(t, managementClient, "TestMetricsClient_BuildUpSetupAndTestMetrics")
	if err != nil {
		return
	}
	//Clean up: delete lab
	defer func() {
		err = deleteLabAndCheckForSuccess(t, managementClient, lab)
		assert.NoError(t, err, "error during delete lab")
	}()

	//Create an engine1
	engine1, err := createEngineAndCheckForSuccess(t, managementClient, engineName1, engineId1)
	if err != nil {
		return
	}
	//Cleanup: delete engine1
	defer func() {
		err = deleteEngineAndCheckForSuccess(t, managementClient, engine1)
		assert.NoError(t, err, "error during delete engine")
	}()

	//Create endpoint1
	endpoint1, err := createEndpointAndCheckForSuccess(t, managementClient, endpointName1, address1, configMetricsTest.Protocol)
	if err != nil {
		return
	}
	//Cleanup: delete endpoint1
	defer func() {
		err = deleteEndpointAndCheckForSuccess(t, managementClient, endpoint1)
		assert.NoError(t, err, "error during delete endpoint")
	}()

	//Create endpoint2
	endpoint2, err := createEndpointAndCheckForSuccess(t, managementClient, endpointName2, address2, configMetricsTest.Protocol)
	if err != nil {
		return
	}
	//Cleanup: delete endpoint2
	defer func() {
		err = deleteEndpointAndCheckForSuccess(t, managementClient, endpoint2)
		assert.NoError(t, err, "error during delete endpoint")
	}()

	//Create user1
	user1, err := createUserAndCheckForSuccess(t, managementClient, userIdentifier1, name1, "", "", "", "")
	if err != nil {
		return
	}
	//Cleanup: delete user1
	defer func() {
		err = deleteUserAndCheckForSuccess(t, managementClient, user1)
		assert.NoError(t, err, "error during delete user")
	}()

	//Add User1 to Engine1
	err = addUserToEngineAndCheckForSuccess(t, managementClient, engine1, user1)
	if err != nil {
		return
	}
	//Cleanup: remove user1 from engine1
	defer func() {
		err = removeUserFromEngineAndCheckForSuccess(t, managementClient, engine1, user1)
		assert.NoError(t, err, "error during remove user from engine")
	}()

	//Add endpoint1 to engine1
	err = addEndpointToEngineAndCheckForSuccess(t, managementClient, engine1, endpoint1)
	if err != nil {
		return
	}
	defer func() {
		err = removeEndpointFromEngineAndCheckForSuccess(t, managementClient, engine1, endpoint1)
		assert.NoError(t, err, "error during remove endpoint from engine")
	}()
	//Add endpoint2 to engine1
	err = addEndpointToEngineAndCheckForSuccess(t, managementClient, engine1, endpoint2)
	if err != nil {
		return
	}
	defer func() {
		err = removeEndpointFromEngineAndCheckForSuccess(t, managementClient, engine1, endpoint2)
		assert.NoError(t, err, "error during remove endpoint from engine")
	}()

	//Create agent1
	agent1, err := createAgentAndCheckForSuccess(t, managementClient, agentName1, agentDataDir1)
	if err != nil {
		return
	}
	//Clean up: delete agent1
	defer func() {
		err = deleteAgentAndCheckForSuccess(t, managementClient, agent1)
		assert.NoError(t, err, "error during delete agent")
	}()

	//Add engine1 to agent1
	err = addEngineToAgentAndCheckForSuccess(t, managementClient, agent1, engine1)
	if err != nil {
		return
	}
	//Cleanup: remove engine1 from agent1
	defer func() {
		err = removeEngineFromAgentAndCheckForSuccess(t, managementClient, agent1, engine1)
		assert.NoError(t, err, "error during remove engine from agent")
	}()

	//Add agent1 to lab
	err = addAgentToLabAndCheckForSuccess(t, managementClient, lab, agent1)
	if err != nil {
		return
	}
	//Cleanup: remove agent1 from lab
	defer func() {
		err = removeAgentFromLabAndCheckForSuccess(t, managementClient, lab, agent1)
		assert.NoError(t, err, "error during remove agent from lab")
	}()

	//Power on lab
	err = setLabPowerAndCheckForSuccess(t, managementClient, lab, true)
	if err != nil {
		return
	}
	//Cleanup: turn lab off
	defer func() {
		err = setLabPowerAndCheckForSuccess(t, managementClient, lab, false)
		assert.NoError(t, err, "error during power off lab")
	}()

	//waiting for asynchronous metrics importer
	time.Sleep(20 * time.Second)

	//Test GetProcesses
	processes, err := metricsClient.GetProcesses(nil)
	if assert.NoError(t, err, "error during GetProcesses") {
		assert.True(t, len(processes) > 0, "No proccesses found during GetProcesses")
	}

	//Test GetProcess
	process, err := metricsClient.GetProcess(processes[0].Id)
	if assert.NoError(t, err, "error during GetProcess") {
		assert.NotNil(t, process.Id, "process id is nil")
		assert.True(t, process.Path == "/opt/snmpsim/supervised/snmpsim-run-labs.sh", "process path is incorrect")
		assert.NotNil(t, process.Runtime, "process runtime is nil")
		assert.NotNil(t, process.Cpu, "process cpu is nil")
		assert.NotNil(t, process.Memory, "process memory is nil")
		assert.NotNil(t, process.Files, "process files is nil")
		assert.NotNil(t, process.Exits, "process exits is nil")
		assert.NotNil(t, process.Changes, "process changes is nil")
		assert.NotNil(t, process.UpdateInterval, "process update_interval is nil")
		assert.NotEmpty(t, process.LastUpdate, "process last_update is empty")
		assert.NotNil(t, process.ConsolePages.Count, "process console_pages count is nil")
		assert.NotEmpty(t, process.ConsolePages.LastUpdate, "process console_pages last_update is empty")
		assert.Equal(t, "sim", process.Supervisor.Hostname, "process supervisor hostname does not match with expected")
		assert.Equal(t, "/opt/snmpsim/supervised", process.Supervisor.WatchDir, "process supervisor watch_dir does not match with expected")
	}

	//Test GetProcessEndpoints
	endpoints, err := metricsClient.GetProcessEndpoints(processes[0].Id)
	if assert.NoError(t, err, "error during GetProcessEndpoints") {
		assert.True(t, len(endpoints) > 0, "No endpoints found during GetProcessEndpoints")
	}

	//Test GetProcessEndpoint
	endpoint, err := metricsClient.GetProcessEndpoint(processes[0].Id, endpoints[0].Id)
	if assert.NoError(t, err, "error during GetProcessEndpoint") {
		assert.Equal(t, "udpv4", endpoint.Protocol, "process endpoint protocol does not match with expected")
		match, _ := regexp.MatchString(`[0-9]{3}\.[0-9]{3}\.[0-9]{3}\.[0-9]{3}:[0-9]+`, endpoint.Address)
		assert.True(t, match, "process endpoint address does not match with expected")
		assert.NotEmpty(t, endpoint.Address, "process endpoint address is empty")
	}

	//Test GetProcessConsolePages
	consolePages, err := metricsClient.GetProcessConsolePages(processes[0].Id)
	if assert.NoError(t, err, "error during GetProcessConsoles") {
		assert.True(t, len(consolePages) > 0, "No console pages found during GetProcessConsoles")
	}

	//Test GetProcessConsolePage
	consolePage, err := metricsClient.GetProcessConsolePage(processes[0].Id, consolePages[0].Id)
	if assert.NoError(t, err, "error during GetProcessConsole") {
		assert.NotEmpty(t, consolePage.Timestamp, "process consol_page timestamp is empty")
		match, _ := regexp.MatchString(`[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\+[0-9]{2}:[0-9]{2}`, consolePage.Timestamp)
		assert.True(t, match, "process consol_page does not match expected pattern")
		assert.NotEmpty(t, consolePage.Text, "process consol_page text is empty")
	}
}

func requestSender(t *testing.T, wg *sync.WaitGroup, snmpCommunity string) {
	defer wg.Done()

	//SNMP Request Agent 1 SNMPv2
	agent1Snmpv2c := &gosnmp.GoSNMP{
		Target:    configMetricsTest.Agent1.EndpointAddress,
		Port:      uint16(configMetricsTest.Agent1.EndpointPort[1]),
		Timeout:   time.Duration(2) * time.Second,
		Version:   gosnmp.Version2c,
		Community: snmpCommunity,
		Transport: "udp",
	}
	err := agent1Snmpv2c.ConnectIPv4()
	if !assert.NoError(t, err, "error during snmp connect, cannot send request!") {
		return
	}
	defer func() {
		err = agent1Snmpv2c.Conn.Close()
		assert.NoError(t, err, "error during snmp connection close")
	}()
	snmpErrCounter := 0

	//send snmp requests for 30s
	for stay, timeout := true, time.After(30*time.Second); stay; {
		select {
		case <-timeout:
			stay = false
		default:
			_, err := agent1Snmpv2c.Get([]string{"1.3.6.1.2.1.1.1.0"})
			if err != nil {
				if snmpErrCounter >= 3 {
					assert.NoError(t, errors.New("snmp get request failed more than 3 times in a row"))
					stay = false
				} else {
					snmpErrCounter++
				}
			} else {
				snmpErrCounter = 0
			}
		}
	}
}
