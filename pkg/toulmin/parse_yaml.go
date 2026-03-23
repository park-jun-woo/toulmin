//ff:func feature=graph type=parser control=sequence
//ff:what ParseYAML — parses YAML file into GraphDef
package toulmin

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ParseYAML reads a YAML file and returns a GraphDef.
func ParseYAML(path string) (GraphDef, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return GraphDef{}, fmt.Errorf("toulmin: read yaml: %w", err)
	}
	var def GraphDef
	if err := yaml.Unmarshal(data, &def); err != nil {
		return GraphDef{}, fmt.Errorf("toulmin: parse yaml: %w", err)
	}
	return def, nil
}
