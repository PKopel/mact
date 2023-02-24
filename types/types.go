// +kubebuilder:object:generate=true
package types

import (
	"log"
	"os"

	deepcopy "github.com/barkimedes/go-deepcopy"
	"gopkg.in/yaml.v2"
)

type changeType string

const (
	Add    changeType = "add"
	Remove changeType = "remove"
	Modify changeType = "modify"
)

type HttpVerb string

const (
	GET    HttpVerb = "GET"
	POST   HttpVerb = "POST"
	PUT    HttpVerb = "PUT"
	PATCH  HttpVerb = "PATCH"
	DELETE HttpVerb = "DELETE"
)

type Change struct {
	Type  changeType  `yaml:"type"`
	Field string      `yaml:"field"`
	Value interface{} `yaml:"value,omitempty"`
}

type EndpointConfig struct {
	Path        string   `yaml:"path"`
	Verb        HttpVerb `yaml:"verb"`
	StatusCodes []int    `yaml:"statusCodes"`
	Changes     []Change `yaml:"changes"`
}

type ServiceConfig struct {
	Host          string           `yaml:"host"`
	TrustAllCerts bool             `yaml:"trustAllCerts"`
	Endpoints     []EndpointConfig `yaml:"endpoints"`
}

type MactConfig struct {
	Services []ServiceConfig `yaml:"services"`
}

func ReadConfig(configFile string) MactConfig {
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Error readig config file %v: %v ", configFile, err)
	}
	var config MactConfig
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return config
}

// DeepCopyInto is a deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Change) DeepCopyInto(out *Change) {
	*out = *in
	if in.Value != nil {
		in, out := &in.Value, &out.Value
		var err error
		*out, err = deepcopy.Anything(*in)
		if err != nil {
			panic(err)
		}
	}
}

// DeepCopy is a deepcopy function, copying the receiver, creating a new Change.
func (in *Change) DeepCopy() *Change {
	if in == nil {
		return nil
	}
	out := new(Change)
	in.DeepCopyInto(out)
	return out
}
