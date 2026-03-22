//ff:func feature=engine type=engine control=sequence dimension=1
//ff:what LoadGraphYAML — parses YAML file and builds a live Graph
package toulmin

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadGraphYAML reads a YAML file and builds a live *Graph using the provided function and backing registries.
// Combines YAML parsing and LoadGraph into a single call.
func LoadGraphYAML(path string, functions map[string]any, backings map[string]any) (*Graph, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("toulmin: read yaml: %w", err)
	}
	var def GraphDef
	if err := yaml.Unmarshal(data, &def); err != nil {
		return nil, fmt.Errorf("toulmin: parse yaml: %w", err)
	}
	return LoadGraph(def, functions, backings)
}
