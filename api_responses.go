package snmpsimclient

/*
Labs is an array of Labs.
*/
type Labs []Lab

/*
Lab - Group of SNMP agents belonging to the same virtual laboratory. Some operations can be applied to them all at once.
*/
type Lab struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Agents Agents `json:"agents"`
	Power  string `json:"power"`
	Tags   Tags   `json:"tags"`
}

/*
Engines is an array of engines.
*/
type Engines []Engine

/*
Engine - Represents a unique, independent and fully operational SNMP engine, though not yet attached to any transport endpoints.
*/
type Engine struct {
	Id        int       `json:"id"`
	EngineId  string    `json:"engine_id"`
	Name      string    `json:"name"`
	Agents    Agents    `json:"agents"`
	Endpoints Endpoints `json:"endpoints"`
	Users     Users     `json:"users"`
	Tags      Tags      `json:"tags"`
}

/*
Agents is an array of agents.
*/
type Agents []Agent

/*
Agent - Represents SNMP agent. Consists of SNMP engine and transport endpoints it binds.
*/
type Agent struct {
	Id        int       `json:"id"`
	Engines   Engines   `json:"engines"`
	Name      string    `json:"name"`
	Endpoints Endpoints `json:"endpoints"`
	Labs      Labs      `json:"labs"`
	Selectors Selectors `json:"selectors"`
	DataDir   string    `json:"data_dir"`
	Tags      Tags      `json:"tags"`
}

/*
Endpoints is an array of endpoints.
*/
type Endpoints []Endpoint

/*
Endpoint - SNMP transport endpoint object. Each SNMP engine can bind one or more transport endpoints. Each transport endpoint can only be bound by one SNMP engine.
*/
type Endpoint struct {
	Id       int     `json:"id"`
	Engines  Engines `json:"engines"`
	Name     string  `json:"name"`
	Protocol string  `json:"protocol"`
	Address  string  `json:"address"`
	Tags     Tags    `json:"tags"`
}

/*
Recordings is an array of recordings.
*/
type Recordings []Recording

/*
Recording - Represents a single simulation data file residing by path under simulation data root.
*/
type Recording struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

/*
Users is an array of users.
*/
type Users []User

/*
User - SNMPv3 USM user object. Contains SNMPv3 credentials grouped by user name.
*/
type User struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	AuthKey   string  `json:"auth_key"`
	AuthProto string  `json:"auth_proto"`
	Engines   Engines `json:"engines"`
	PrivKey   string  `json:"priv_key"`
	PrivProto string  `json:"priv_proto"`
	User      string  `json:"user"`
	Tags      Tags    `json:"tags"`
}

//TODO: not implemented int the api yet, there is only one default selector for snmpv2c and one for snmpv3.
// We have to wait with the implementation until its implemented in the api.
// The strucuture is like this in the api doc but it might be different in the real api, this has to be checked first before it can be used.

/*
Selectors is an array of selectors.
*/
type Selectors []Selector

/*
Selector - Each selector should end up being a path to a simulation data file relative to the command responder's data directory.
The value of the selector can be static or, more likely, it contains templates that are expanded at run time. Each template can expand into some property of the inbound request.
Known templates include:
    ${context-engine-id}
    ${context-name}
    ${endpoint-id}
    ${source-address}
*/
type Selector struct {
	Id       int    `json:"id"`
	Comment  string `json:"comment"`
	Template string `json:"template"`
	Tags     Tags   `json:"tags"`
}

/*
ErrorResponse contains error information.
*/
type ErrorResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

/*
ProcessMetrics - SNMP simulator system is composed of many running processes. This object describes common properties of a process.

	Cmdline   *string           `json:"cmdline"`
	Uptime    *int              `json:"uptime"`
	Owner     *string           `json:"owner"`

	LifeCycle  	*ProcessLifeCycle	`json:"lifecycle"`
*/
type ProcessMetrics struct {
	Id             int          `json:"id"`
	Path           string       `json:"path"`
	Runtime        int          `json:"runtime"`
	Cpu            int          `json:"cpu"`
	Memory         int          `json:"memory"`
	Files          int          `json:"files"`
	Exits          int          `json:"exits"`
	Changes        int          `json:"changes"`
	UpdateInterval int          `json:"update_interval"`
	LastUpdate     string       `json:"last_update"`
	ConsolePages   ConsolePages `json:"console_pages"`
	Supervisor     Supervisor   `json:"supervisor"`
}

/*
ProcessLifeCycle - How this process has being doing.
*/
type ProcessLifeCycle struct {
	Exits    *int `json:"exits"`
	Restarts *int `json:"restarts"`
}

/*
ProcessesMetrics is an array of ProcessMetrics.
*/
type ProcessesMetrics []ProcessMetrics

/*
PacketMetrics - Transport endpoint related activity. Includes raw network packet counts as well as SNMP messages failed to get processed at the later stages.
*/
type PacketMetrics struct {
	FirstHit        *int   `json:"first_hit"`
	LastHit         *int   `json:"last_hit"`
	Total           *int64 `json:"total"`
	ParseFailures   *int64 `json:"parse_failures"`
	AuthFailures    *int64 `json:"auth_failures"`
	ContextFailures *int64 `json:"context_failures"`
}

/*
MessageMetrics - SNMP message level metrics.
*/
type MessageMetrics struct {
	FirstHit   *int       `json:"first_hit"`
	LastHit    *int       `json:"last_hit"`
	Pdus       *int64     `json:"pdus"`
	VarBinds   *int64     `json:"var_binds"`
	Failures   *int64     `json:"failures"`
	Variations Variations `json:"variations"`
}

/*
Variation - Variation module metrics.
*/
type Variation struct {
	FirstHit *int    `json:"first_hit"`
	LastHit  *int    `json:"last_hit"`
	Total    *int64  `json:"total"`
	Name     *string `json:"name"`
	Failures *int64  `json:"failures"`
}

/*
Variations is an array of variations.
*/
type Variations []Variation

/*
Tag - tags a collection of SNMP simulator control plane resources
*/
type Tag struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Agents      Agents    `json:"agents"`
	Endpoints   Endpoints `json:"endpoints"`
	Engines     Engines   `json:"engines"`
	Labs        Labs      `json:"labs"`
	Selectors   Selectors `json:"selectors"`
	Users       Users     `json:"users"`
}

/*
Tags is an array of tags.
*/
type Tags []Tag

/*
Console - contains information regarding the console
*/
type Console struct {
	Id        int    `json:"id"`
	Timestamp string `json:"timestamp"`
	Text      string `json:"text"`
}

/*
Consoles is an array of Consoles
*/
type Consoles []Console

/*
ConsolePages - contains information regarding snmpsim consoles
*/
type ConsolePages struct {
	Count      int    `json:"count"`
	LastUpdate string `json:"last_update"`
}

/*
Supervisor - contains information regarding the supervisor
*/
type Supervisor struct {
	Hostname string `json:"hostname"`
	WatchDir string `json:"watch_dir"`
}
