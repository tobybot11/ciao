/*
// Copyright (c) 2016 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
*/

package types

import (
	"errors"
	"time"

	"github.com/01org/ciao/ciao-storage"
	"github.com/01org/ciao/openstack/block"
	"github.com/01org/ciao/payloads"
)

// SourceType contains the valid values of the storage source.
type SourceType string

const (
	// ImageService indicates the source comes from the image service.
	ImageService SourceType = "image"

	// VolumeService indicates the source comes from the volume service.
	VolumeService SourceType = "volume"

	// Empty indicates that there is no source for the storage source
	Empty SourceType = "empty"
)

// StorageResource defines a storage resource for a workload.
// TBD: should the workload support multiple of these?
type StorageResource struct {
	// ID indicates a volumeID. If ID is blank, then it needs to be created.
	ID string

	// Bootable indicates whether should the resource be used for booting
	Bootable bool

	// Persistent indicates whether the storage is temporary
	// TBD: do we bother to save info about temp storage?
	//      does it count against quota?
	Persistent bool

	// Size is the size of the storage to be created if new.
	Size int

	// ImageType indicates whether we are making a new resource
	// based on an image or existing volume.
	// Needed only for new storage.
	SourceType SourceType

	// SourceID represents the ID of either the image or the volume
	// that the storage resource is based on.
	SourceID string
}

// Workload contains resource and configuration information for a user
// workload.
type Workload struct {
	ID          string                       `json:"id"`
	Description string                       `json:"description"`
	FWType      string                       `json:"-"`
	VMType      payloads.Hypervisor          `json:"-"`
	ImageID     string                       `json:"-"`
	ImageName   string                       `json:"-"`
	Config      string                       `json:"-"`
	Defaults    []payloads.RequestedResource `json:"-"`
	Storage     *StorageResource             `json:"-"`
}

// Instance contains information about an instance of a workload.
type Instance struct {
	ID          string              `json:"instance_id"`
	TenantID    string              `json:"tenant_id"`
	State       string              `json:"instance_state"`
	WorkloadID  string              `json:"workload_id"`
	NodeID      string              `json:"node_id"`
	MACAddress  string              `json:"mac_address"`
	IPAddress   string              `json:"ip_address"`
	SSHIP       string              `json:"ssh_ip"`
	SSHPort     int                 `json:"ssh_port"`
	CNCI        bool                `json:"-"`
	Usage       map[string]int      `json:"-"`
	Attachments []StorageAttachment `json:"-"`
	CreateTime  time.Time           `json:"-"`
}

// SortedInstancesByID implements sort.Interface for Instance by ID string
type SortedInstancesByID []*Instance

func (s SortedInstancesByID) Len() int           { return len(s) }
func (s SortedInstancesByID) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SortedInstancesByID) Less(i, j int) bool { return s[i].ID < s[j].ID }

// SortedComputeNodesByID implements sort.Interface for Node by ID string
type SortedComputeNodesByID []CiaoComputeNode

func (s SortedComputeNodesByID) Len() int           { return len(s) }
func (s SortedComputeNodesByID) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SortedComputeNodesByID) Less(i, j int) bool { return s[i].ID < s[j].ID }

// Tenant contains information about a tenant or project.
type Tenant struct {
	ID        string
	Name      string
	CNCIID    string
	CNCIMAC   string
	CNCIIP    string
	Resources []*Resource
}

// Resource contains quota or limit information on a resource type.
type Resource struct {
	Rname string
	Rtype int
	Limit int
	Usage int
}

// OverLimit calculates whether a request will put a tenant over it's limit.
func (r *Resource) OverLimit(request int) bool {
	if r.Limit > 0 && r.Usage+request > r.Limit {
		return true
	}
	return false
}

// LogEntry stores information about events.
type LogEntry struct {
	Timestamp time.Time `json:"time_stamp"`
	TenantID  string    `json:"tenant_id"`
	EventType string    `json:"type"`
	Message   string    `json:"message"`
}

// NodeStats stores statistics for individual nodes in the cluster.
type NodeStats struct {
	NodeID          string    `json:"node_id"`
	Timestamp       time.Time `json:"time_stamp"`
	Load            int       `json:"load"`
	MemTotalMB      int       `json:"mem_total_mb"`
	MemAvailableMB  int       `json:"mem_available_mb"`
	DiskTotalMB     int       `json:"disk_total_mb"`
	DiskAvailableMB int       `json:"disk_available_mb"`
	CpusOnline      int       `json:"cpus_online"`
}

// NodeSummary contains summary information for all nodes in the cluster.
type NodeSummary struct {
	NodeID                string `json:"node_id"`
	TotalInstances        int    `json:"total_instances"`
	TotalRunningInstances int    `json:"total_running_instances"`
	TotalPendingInstances int    `json:"total_pending_instances"`
	TotalPausedInstances  int    `json:"total_paused_instances"`
}

// TenantCNCI contains information about the CNCI instance for a tenant.
type TenantCNCI struct {
	TenantID   string   `json:"tenant_id"`
	IPAddress  string   `json:"ip_address"`
	MACAddress string   `json:"mac_address"`
	InstanceID string   `json:"instance_id"`
	Subnets    []string `json:"subnets"`
}

// FrameStat contains tracing information per node.
type FrameStat struct {
	ID               string  `json:"node_id"`
	TotalElapsedTime float64 `json:"total_elapsed_time"`
	ControllerTime   float64 `json:"total_controller_time"`
	LauncherTime     float64 `json:"total_launcher_time"`
	SchedulerTime    float64 `json:"total_scheduler_time"`
}

// BatchFrameStat contains tracing information for a group of start requests
// by label.
type BatchFrameStat struct {
	NumInstances             int     `json:"num_instances"`
	TotalElapsed             float64 `json:"total_elapsed"`
	AverageElapsed           float64 `json:"average_elapsed"`
	AverageControllerElapsed float64 `json:"average_controller_elapsed"`
	AverageLauncherElapsed   float64 `json:"average_launcher_elapsed"`
	AverageSchedulerElapsed  float64 `json:"average_scheduler_elapsed"`
	VarianceController       float64 `json:"controller_variance"`
	VarianceLauncher         float64 `json:"launcher_variance"`
	VarianceScheduler        float64 `json:"scheduler_variance"`
}

// BatchFrameSummary provides summary information on tracing per label.
type BatchFrameSummary struct {
	BatchID      string `json:"batch_id"`
	NumInstances int    `json:"num_instances"`
}

// Node contains information about a physical node in the cluster.
type Node struct {
	ID       string `json:"node_id"`
	IPAddr   string `json:"ip_address"`
	Hostname string `json:"hostname"`
}

// BlockState represents the state of the block device in the controller
// datastore. This is a subset of the openstack status type.
type BlockState string

const (
	// Available means that the volume is ok for attaching.
	Available BlockState = BlockState(block.Available)

	// Attaching means that the volume is in the process
	// of attaching to an instance.
	Attaching BlockState = BlockState(block.Attaching)

	// InUse means that the volume has been successfully
	// attached to an instance.
	InUse BlockState = BlockState(block.InUse)

	// Detaching means that the volume is in process
	// of detaching.
	Detaching BlockState = "detaching"
)

// BlockData respresents the attributes of this block device.
// TBD - do we really need to store this as actual data,
// or can we use a set of interfaces to get the info?
type BlockData struct {
	storage.BlockDevice
	TenantID    string     // the tenant who owns this volume
	Size        int        // size in GB
	State       BlockState // status of
	CreateTime  time.Time  // when we created the volume
	Name        string     // a human readable name for this volume
	Description string     // some text to describe this volume.
}

// StorageAttachment represents a link between a block device and
// an instance.
type StorageAttachment struct {
	ID         string // a uuid
	InstanceID string // the instance this volume is attached to
	BlockID    string // the ID of the block device
}

// CiaoComputeTenants represents the unmarshalled version of the contents of a
// /v2.1/tenants response.  It contains information about the tenants in a ciao
// cluster.
type CiaoComputeTenants struct {
	Tenants []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"tenants"`
}

// NewCiaoComputeTenants allocates a CiaoComputeTenants structure.
// It allocates the Tenants slice as well so that the marshalled
// JSON is an empty array and not a nil pointer, as specified by the
// OpenStack APIs.
func NewCiaoComputeTenants() (tenants CiaoComputeTenants) {
	tenants.Tenants = []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}{}
	return
}

// CiaoComputeNode contains status and statistic information for an individual
// node.
type CiaoComputeNode struct {
	ID                    string    `json:"id"`
	Timestamp             time.Time `json:"updated"`
	Status                string    `json:"status"`
	MemTotal              int       `json:"ram_total"`
	MemAvailable          int       `json:"ram_available"`
	DiskTotal             int       `json:"disk_total"`
	DiskAvailable         int       `json:"disk_available"`
	Load                  int       `json:"load"`
	OnlineCPUs            int       `json:"online_cpus"`
	TotalInstances        int       `json:"total_instances"`
	TotalRunningInstances int       `json:"total_running_instances"`
	TotalPendingInstances int       `json:"total_pending_instances"`
	TotalPausedInstances  int       `json:"total_paused_instances"`
}

// CiaoComputeNodes represents the unmarshalled version of the contents of a
// /v2.1/nodes response.  It contains status and statistics information
// for a set of nodes.
type CiaoComputeNodes struct {
	Nodes []CiaoComputeNode `json:"nodes"`
}

// NewCiaoComputeNodes allocates a CiaoComputeNodes structure.
// It allocates the Nodes slice as well so that the marshalled
// JSON is an empty array and not a nil pointer, as specified by the
// OpenStack APIs.
func NewCiaoComputeNodes() (nodes CiaoComputeNodes) {
	nodes.Nodes = []CiaoComputeNode{}
	return
}

// CiaoTenantResources represents the unmarshalled version of the contents of a
// /v2.1/{tenant}/quotas response.  It contains the current resource usage
// information for a tenant.
type CiaoTenantResources struct {
	ID            string    `json:"id"`
	Timestamp     time.Time `json:"updated"`
	InstanceLimit int       `json:"instances_limit"`
	InstanceUsage int       `json:"instances_usage"`
	VCPULimit     int       `json:"cpus_limit"`
	VCPUUsage     int       `json:"cpus_usage"`
	MemLimit      int       `json:"ram_limit"`
	MemUsage      int       `json:"ram_usage"`
	DiskLimit     int       `json:"disk_limit"`
	DiskUsage     int       `json:"disk_usage"`
}

// CiaoUsage contains a snapshot of resource consumption for a tenant.
type CiaoUsage struct {
	VCPU      int       `json:"cpus_usage"`
	Memory    int       `json:"ram_usage"`
	Disk      int       `json:"disk_usage"`
	Timestamp time.Time `json:"timestamp"`
}

// CiaoUsageHistory represents the unmarshalled version of the contents of a
// /v2.1/{tenant}/resources response.  It contains snapshots of usage information
// for a given tenant over a given period of time.
type CiaoUsageHistory struct {
	Usages []CiaoUsage `json:"usage"`
}

// CiaoCNCISubnet contains subnet information for a CNCI.
type CiaoCNCISubnet struct {
	Subnet string `json:"subnet_cidr"`
}

// CiaoCNCI contains information about an individual CNCI.
type CiaoCNCI struct {
	ID        string           `json:"id"`
	TenantID  string           `json:"tenant_id"`
	IPv4      string           `json:"IPv4"`
	Geography string           `json:"geography"`
	Subnets   []CiaoCNCISubnet `json:"subnets"`
}

// CiaoCNCIDetail represents the unmarshalled version of the contents of a
// v2.1/cncis/{cnci}/detail response.  It contains information about a CNCI.
type CiaoCNCIDetail struct {
	CiaoCNCI `json:"cnci"`
}

// CiaoCNCIs represents the unmarshalled version of the contents of a
// v2.1/cncis response.  It contains information about all the CNCIs
// in the ciao cluster.
type CiaoCNCIs struct {
	CNCIs []CiaoCNCI `json:"cncis"`
}

// NewCiaoCNCIs allocates a CiaoCNCIs structure.
// It allocates the CNCIs slice as well so that the marshalled
// JSON is an empty array and not a nil pointer, as specified by the
// OpenStack APIs.
func NewCiaoCNCIs() (cncis CiaoCNCIs) {
	cncis.CNCIs = []CiaoCNCI{}
	return
}

// CiaoServerStats contains status information about a CN or a NN.
type CiaoServerStats struct {
	ID        string    `json:"id"`
	NodeID    string    `json:"node_id"`
	Timestamp time.Time `json:"updated"`
	Status    string    `json:"status"`
	TenantID  string    `json:"tenant_id"`
	IPv4      string    `json:"IPv4"`
	VCPUUsage int       `json:"cpus_usage"`
	MemUsage  int       `json:"ram_usage"`
	DiskUsage int       `json:"disk_usage"`
}

// CiaoServersStats represents the unmarshalled version of the contents of a
// v2.1/nodes/{node}/servers/detail response.  It contains general information
// about a group of instances.
type CiaoServersStats struct {
	TotalServers int               `json:"total_servers"`
	Servers      []CiaoServerStats `json:"servers"`
}

// NewCiaoServersStats allocates a CiaoServersStats structure.
// It allocates the Servers slice as well so that the marshalled
// JSON is an empty array and not a nil pointer, as specified by the
// OpenStack APIs.
func NewCiaoServersStats() (servers CiaoServersStats) {
	servers.Servers = []CiaoServerStats{}
	return
}

// CiaoClusterStatus represents the unmarshalled version of the contents of a
// v2.1/nodes/summary response.  It contains information about the nodes that
// make up a ciao cluster.
type CiaoClusterStatus struct {
	Status struct {
		TotalNodes            int `json:"total_nodes"`
		TotalNodesReady       int `json:"total_nodes_ready"`
		TotalNodesFull        int `json:"total_nodes_full"`
		TotalNodesOffline     int `json:"total_nodes_offline"`
		TotalNodesMaintenance int `json:"total_nodes_maintenance"`
	} `json:"cluster"`
}

// CNCIDetail stores the IPv4 for a CNCI Agent.
type CNCIDetail struct {
	IPv4 string `json:"IPv4"`
}

// CiaoServersAction represents the unmarshalled version of the contents of a
// v2.1/servers/action request.  It contains an action to be performed on
// one or more instances.
type CiaoServersAction struct {
	Action    string   `json:"action"`
	ServerIDs []string `json:"servers"`
}

// CiaoTraceSummary contains information about a specific SSNTP Trace label.
type CiaoTraceSummary struct {
	Label     string `json:"label"`
	Instances int    `json:"instances"`
}

// CiaoTracesSummary represents the unmarshalled version of the response to a
// v2.1/traces request.  It contains a list of all trace labels and the
// number of instances associated with them.
type CiaoTracesSummary struct {
	Summaries []CiaoTraceSummary `json:"summaries"`
}

// CiaoFrameStat contains the elapsed time statistics for a frame.
type CiaoFrameStat struct {
	ID               string  `json:"node_id"`
	TotalElapsedTime float64 `json:"total_elapsed_time"`
	ControllerTime   float64 `json:"total_controller_time"`
	LauncherTime     float64 `json:"total_launcher_time"`
	SchedulerTime    float64 `json:"total_scheduler_time"`
}

// CiaoBatchFrameStat contains frame statisitics for a ciao cluster.
type CiaoBatchFrameStat struct {
	NumInstances             int     `json:"num_instances"`
	TotalElapsed             float64 `json:"total_elapsed"`
	AverageElapsed           float64 `json:"average_elapsed"`
	AverageControllerElapsed float64 `json:"average_controller_elapsed"`
	AverageLauncherElapsed   float64 `json:"average_launcher_elapsed"`
	AverageSchedulerElapsed  float64 `json:"average_scheduler_elapsed"`
	VarianceController       float64 `json:"controller_variance"`
	VarianceLauncher         float64 `json:"launcher_variance"`
	VarianceScheduler        float64 `json:"scheduler_variance"`
}

// CiaoTraceData represents the unmarshalled version of the response to a
// v2.1/traces/{label} request.  It contains statistics computed from the trace
// information of SSNTP commands sent within a ciao cluster.
type CiaoTraceData struct {
	Summary    CiaoBatchFrameStat `json:"summary"`
	FramesStat []CiaoFrameStat    `json:"frames"`
}

// CiaoEvent contains information about an individual event generated
// in a ciao cluster.
type CiaoEvent struct {
	Timestamp time.Time `json:"time_stamp"`
	TenantID  string    `json:"tenant_id"`
	EventType string    `json:"type"`
	Message   string    `json:"message"`
}

// CiaoEvents represents the unmarshalled version of the response to a
// v2.1/{tenant}/event or v2.1/event request.
type CiaoEvents struct {
	Events []CiaoEvent `json:"events"`
}

// NewCiaoEvents allocates a CiaoEvents structure.
// It allocates the Events slice as well so that the marshalled
// JSON is an empty array and not a nil pointer, as specified by the
// OpenStack APIs.
func NewCiaoEvents() (events CiaoEvents) {
	events.Events = []CiaoEvent{}
	return
}

var (
	// ErrQuota is returned when a resource limit is exceeded.
	ErrQuota = errors.New("Over Quota")

	// ErrTenantNotFound is returned when a tenant ID is unknown.
	ErrTenantNotFound = errors.New("Tenant not found")

	// ErrInstanceNotFound is returned when an instance is not found.
	ErrInstanceNotFound = errors.New("Instance not found")

	// ErrInstanceNotAssigned is returned when an instance is not assigned to a node.
	ErrInstanceNotAssigned = errors.New("Cannot perform operation: instance not assigned to Node")
)
