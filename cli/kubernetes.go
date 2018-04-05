package main

import (
	"github.com/stelligent/config-lint/assertion"
)

// KubernetesLinter lints resources in Kubernets YAML files
type KubernetesLinter struct {
	Filenames   []string
	Log         assertion.LoggingFunction
	ValueSource assertion.ValueSource
}

// KubernetesResourceLoader converts Kubernetes configuration files into a collection of Resource objects
type KubernetesResourceLoader struct {
	Log assertion.LoggingFunction
}

func getResourceIDFromMetadata(m map[string]interface{}) (string, bool) {
	if metadata, ok := m["metadata"].(map[string]interface{}); ok {
		if name, ok := metadata["name"].(string); ok {
			return name, true
		}
	}
	return "", false
}

// Load converts a text file into a collection of Resource objects
func (l KubernetesResourceLoader) Load(filename string) ([]assertion.Resource, error) {
	resources := make([]assertion.Resource, 0)
	yamlResources, err := loadYAML(filename, l.Log)
	if err != nil {
		return resources, err
	}
	for _, resource := range yamlResources {
		m := resource.(map[string]interface{})
		var resourceID string
		if name, ok := getResourceIDFromMetadata(m); ok {
			resourceID = name
		} else {
			resourceID = getResourceIDFromFilename(filename)
		}
		kr := assertion.Resource{
			ID:         resourceID,
			Type:       m["kind"].(string),
			Properties: m,
			Filename:   filename,
		}
		resources = append(resources, kr)
	}
	return resources, nil
}

// Validate runs validate on a collection of filenames using a RuleSet
func (l KubernetesLinter) Validate(ruleSet assertion.RuleSet, options LinterOptions) (assertion.ValidationReport, error) {
	loader := KubernetesResourceLoader{Log: l.Log}
	f := FileLinter{Filenames: l.Filenames, Log: l.Log, ValueSource: l.ValueSource, Loader: loader}
	return f.ValidateFiles(ruleSet, options)
}

// Search evaluates a JMESPath expression against the resources in a collection of filenames
func (l KubernetesLinter) Search(ruleSet assertion.RuleSet, searchExpression string) {
	loader := KubernetesResourceLoader{Log: l.Log}
	f := FileLinter{Filenames: l.Filenames, Log: l.Log, ValueSource: l.ValueSource, Loader: loader}
	f.SearchFiles(ruleSet, searchExpression)
}
