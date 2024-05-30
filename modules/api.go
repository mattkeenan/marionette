// Package modules contain the implementation of our modules.  Each
// module has a name/type such as "git", "file", etc.  The modules
// each accept an arbitrary set of parameters which are module-specific.
package modules

import (
	"github.com/skx/marionette/config"
	"github.com/skx/marionette/environment"
)

// ModuleConstructor is the signature of a constructor-function.
type ModuleConstructor func(cfg *config.Config, env *environment.Environment) ModuleAPI

// ModuleAPI is the interface which all of our modules must implement.
//
// There are only two methods, one to check if the supplied parameters
// make sense, the other to actually execute the rule.
//
// If a module wishes to setup a variable in the environment then they
// can optionally implement the `ModuleOutput` interface too.
type ModuleAPI interface {

	// Check allows a module to ensures that any mandatory parameters
	// are present, or perform similar setup-work.
	//
	// If no error is returned then the module will be executed later
	// via a call to Execute.
	Check(map[string]interface{}) error

	// Execute runs the module with the given arguments.
	//
	// The return value is true if the module made a change
	// and false otherwise.
	Execute(map[string]interface{}) (bool, error)
}

// ModuleOutput is an optional interface that may be implemented by any of
// our internal modules.
//
// If this interface is implemented it is possible for modules to set
// values in the environment after they've been executed.
type ModuleOutput interface {

	// GetOutputs will return a set of key-value pairs.
	//
	// These will be set in the environment, scoped by the rule-name,
	// if the module is successfully executed.
	GetOutputs() map[string]string
}

// StringParam returns the named parameter, as a string, from the map.
//
// If the parameter was not present an empty array is returned.
func StringParam(vars map[string]interface{}, param string) string {

	// Get the value
	val, ok := vars[param]
	if !ok {
		return ""
	}

	// Can it be cast into a string?
	str, valid := val.(string)
	if valid {
		return str
	}

	// OK not a string parameter
	return ""
}

// ArrayParam returns the named parameter, as an array, from the map.
//
// If the parameter was not present an empty array is returned.
func ArrayParam(vars map[string]interface{}, param string) []string {

	var empty []string

	// Get the value
	val, ok := vars[param]
	if !ok {
		return empty
	}

	// if we are passed a string then return an array with
	// just the one string
	str, valid := val.(string)
	if valid {
		return []string{str}
	}

	// Can it be cast into a string array?
	strs, valid := val.([]string)
	if valid {
		return strs
	}

	// OK not a string parameter
	return empty
}
