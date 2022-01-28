package entity

import "time"

const BootstrapV4 = "Bootstrap V4"
const BootstrapV3 = "Bootstrap V3"

type HealthCheck struct {
	AppName        string    `json:"appName"`
	AppVersion     string    `json:"appVersion"`
	ClusterName    string    `json:"clusterName"`
	ClusterVersion string    `json:"clusterVersion"`
	DeployedAt     time.Time `json:"deployedAt"`
	Git            struct {
		Hash   string `json:"hash"`
		Branch string `json:"branch"`
		Url    string `json:"url"`
	} `json:"git"`
}

type Infos struct {
	ProdBlue, ProdGreen, PreprodBlue, PreprodGreen *HealthCheck
}

type State struct {
	Selected int
	Screen   string
}
