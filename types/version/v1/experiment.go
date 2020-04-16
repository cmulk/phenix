package v1

import (
	"fmt"
	"os"
	"time"
)

type Schedule map[string]string

type ExperimentSpec struct {
	ExperimentName string        `json:"experimentName" yaml:"experimentName" structs:"experimentName"`
	Topology       *TopologySpec `json:"topology" yaml:"topology"`
	Scenario       *ScenarioSpec `json:"scenario" yaml:"scenario"`
	Schedules      Schedule      `json:"schedules" yaml:"schedules"`
	VLANMin        int           `json:"vlanMin" yaml:"vlanMin" structs:"vlanMin"`
	VLANMax        int           `json:"vlanMax" yaml:"vlanMax" structs:"vlanMax"`
	RunLocal       bool          `json:"runLocal" yaml:"runLocal" structs:"runLocal"`
}

type ExperimentStatus struct {
	StartTime time.Time `json:"startTime" yaml:"startTime"`
}

func (this *ExperimentSpec) SetDefaults() {
	this.Topology.SetDefaults()
}

func (this ExperimentSpec) VerifyScenario() error {
	if this.Scenario == nil {
		return nil
	}

	hosts := make(map[string]struct{})

	for _, node := range this.Topology.Nodes {
		hosts[node.General.Hostname] = struct{}{}
	}

	for _, app := range this.Scenario.Apps.Host {
		for _, host := range app.Hosts {
			if _, ok := hosts[host.Hostname]; !ok {
				return fmt.Errorf("host %s in app %s not in topology", host.Hostname, app.Name)
			}
		}
	}

	return nil
}

func (this ExperimentSpec) SnapshotName(node string) string {
	hostname, _ := os.Hostname()

	return fmt.Sprintf("%s_%s_%s_snapshot", hostname, this.ExperimentName, node)
}
