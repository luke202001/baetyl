package config

import (
	"encoding/json"
	"github.com/baetyl/baetyl-core/common"
	"github.com/baetyl/baetyl-core/models"
	"github.com/baetyl/baetyl-go/log"
	"github.com/baetyl/baetyl/protocol/http"
	"time"
)

// Config config
type Config struct {
	APIServer APIServer `yaml:"apiServer" json:"apiServer" default:"{}"`
	Agent  AgentConfig `yaml:"agent" json:"agent"`
	Logger log.Config  `yaml:"logger" json:"logger"`
}

type APIServer struct {
	InCluster  bool   `yaml:"inCluster" json:"inCluster" default:"false"`
	ConfigPath string `yaml:"configPath" json:"configPath" default:"config.yaml"`
}

type AgentConfig struct {
	Remote struct {
		HTTP   *http.ClientInfo `yaml:"http" json:"http" default:"{}"`
		Report struct {
			URL      string        `yaml:"url" json:"url" default:"/v3/edge/info"`
			Topic    string        `yaml:"topic" json:"topic" default:"$baidu/iot/edge/%s/core/forward"`
			Interval time.Duration `yaml:"interval" json:"interval" default:"20s"`
		} `yaml:"report" json:"report"`
		Desire struct {
			URL   string `yaml:"url" json:"url"`
			Topic string `yaml:"topic" json:"topic" default:"$baidu/iot/edge/%s/core/backward"`
		} `yaml:"desire" json:"desire"`
	} `yaml:"remote" json:"remote"`
}

type BackwardInfo struct {
	Delta    map[string]interface{} `yaml:"delta,omitempty" json:"delta,omitempty"`
	Metadata map[string]interface{} `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

type ForwardInfo struct {
	Metadata map[string]string `yaml:"metadata" json:"metadata" default:"{}"`
	//Status     *Inspect          `yaml:"status" json:"status"`          // node update
	Apps map[string]string `yaml:"apps" json:"apps" default:"{}"` // shadow update
}

type ApplicationResource struct {
	Type    string      `yaml:"type" json:"type"`
	Name    string      `yaml:"name" json:"name"`
	Version string      `yaml:"version" json:"version"`
	Value   models.Application `yaml:"value" json:"value"`
}

type ConfigurationResource struct {
	Type    string        `yaml:"type" json:"type"`
	Name    string        `yaml:"name" json:"name"`
	Version string        `yaml:"version" json:"version"`
	Value   models.Configuration `yaml:"value" json:"value"`
}

type DesireRequest struct {
	Resources []*BaseResource `yaml:"resources" json:"resources"`
}

type DesireResponse struct {
	Resources []*Resource `yaml:"resources" json:"resources"`
}

type VolumeDevice struct {
	DevicePath string `json:"devicePath,omitempty"`
}

type BaseResource struct {
	Type    common.Resource `yaml:"type,omitempty" json:"type,omitempty"`
	Name    string          `yaml:"name,omitempty" json:"name,omitempty"`
	Version string          `yaml:"version,omitempty" json:"version,omitempty"`
}

type Resource struct {
	BaseResource `yaml:",inline" json:",inline"`
	Data         []byte      `yaml:"data,omitempty" json:"data,omitempty"`
	Value        interface{} `yaml:"value,omitempty" json:"value,omitempty"`
}

func (r *Resource) GetApplication() *models.Application {
	if r.Type == common.Application {
		return r.Value.(*models.Application)
	}
	return nil
}

func (r *Resource) GetConfiguration() *models.Configuration {
	if r.Type == common.Configuration {
		return r.Value.(*models.Configuration)
	}
	return nil
}

func (r *Resource) UnmarshalJSON(b []byte) error {
	var base BaseResource
	err := json.Unmarshal(b, &base)
	if err != nil {
		return err
	}
	switch base.Type {
	case common.Application:
		var app ApplicationResource
		err := json.Unmarshal(b, &app)
		if err != nil {
			return err
		}
		r.Value = &app.Value
	case common.Configuration:
		var config ConfigurationResource
		err := json.Unmarshal(b, &config)
		if err != nil {
			return err
		}
		r.Value = &config.Value
	}
	r.Data = b
	r.BaseResource = base
	return nil
}

type StorageObject struct {
	Md5         string `json:"md5,omitempty" yaml:"md5"`
	URL         string `json:"url,omitempty" yaml:"url"`
	Compression string `json:"compression,omitempty" yaml:"compression"`
}
