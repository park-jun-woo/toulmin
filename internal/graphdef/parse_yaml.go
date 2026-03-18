//ff:func feature=graph type=parser control=iteration dimension=1
//ff:what ParseYAML — parses YAML file into GraphDef
package graphdef

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ParseYAML reads a YAML file and returns a GraphDef.
// Rules with qualifier=0 get default 1.0.
func ParseYAML(path string) (*GraphDef, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var def GraphDef
	if err := yaml.Unmarshal(data, &def); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}
	for i := range def.Rules {
		if def.Rules[i].Qualifier == 0 {
			def.Rules[i].Qualifier = 1.0
		}
	}
	return &def, nil
}
