package experiment

import (
	"fmt"
	"strings"
)

const (
	set         = "set"
	pushValue   = "pushValue"
	setValue    = "setValue"
	removeValue = "removeValue"
)

type Experiment struct {
	Parameter string `yaml:"parameter"`
	Value     string `yaml:"value"`
	Key       string `yaml:"key"`
	Action    string `yaml:"action""`
	Applied   bool
}

func (e *Experiment) Execute(commands []string) ([]string, error) {
	if e.Action == set {
		return e.executeSet(commands)
	}
	if e.Action == pushValue {
		return e.executePushValue(commands)
	}

	if e.Action == setValue {
		return e.executeSetValue(commands)
	}

	if e.Action == removeValue {
		return e.executeRemoveValue(commands)
	}

	return nil, fmt.Errorf("invalid action: %s", e.Action)
}

func (e *Experiment) executeSet(commands []string) ([]string, error) {
	for _, c := range commands {
		parts := strings.Split(c, "=")
		if parts[0] == e.Parameter {
			return nil, fmt.Errorf("executing set: cannot add %s as it already exists", e.Parameter)
		}
	}

	c := fmt.Sprintf("%s=%s", e.Parameter, e.Value)
	commands = append(commands, c)

	return commands, nil
}

func (e *Experiment) executePushValue(commands []string) ([]string, error) {
	for idx, c := range commands {
		parts := strings.Split(c, "=")
		if parts[0] == e.Parameter {
			args := strings.Split(parts[1], ",")
			for _, a := range args {
				if a == e.Value {
					return nil, fmt.Errorf("cannot add %s as it already exists", c)
				}
			}

			args = append(args, e.Value)
			arg := strings.Join(args, ",")
			command := strings.Join([]string{parts[0], arg}, "=")
			commands[idx] = command
			return commands, nil
		}
	}

	return nil, fmt.Errorf("parameter %s not found among current args", e.Parameter)
}

func (e *Experiment) executeSetValue(commands []string) ([]string, error) {
	for idx, c := range commands {
		parts := strings.Split(c, "=")
		if parts[0] == e.Parameter {
			newVal := strings.Join([]string{parts[0], e.Value}, "=")
			commands[idx] = newVal
			return commands, nil
		}
	}

	return nil, fmt.Errorf("parameter %s not found among current args", e.Parameter)
}

func (e *Experiment) executeRemoveValue(commands []string) ([]string, error) {
	for idx, c := range commands {
		parts := strings.Split(c, "=")
		if parts[0] == e.Parameter {
			args := strings.Split(parts[1], ",")
			for i, a := range args {
				if a == e.Value {
					remove(args, i)
					break
				}
			}

			arg := strings.Join(args, ",")
			command := strings.Join([]string{parts[0], arg}, "=")
			commands[idx] = command
			return commands, nil
		}
	}

	return nil, fmt.Errorf("parameter %s not found among current args", e.Parameter)
}

func (e *Experiment) CheckIfApplied(cmdArgs []string) bool {
	if e.Action == pushValue || e.Action == setValue || e.Action == set {
		return e.checkValue(cmdArgs)
	}

	if e.Action == removeValue {
		return e.checkRemoval(cmdArgs)
	}

	return false
}

func (e *Experiment) checkValue(cmdArgs []string) bool {
	for _, c := range cmdArgs {
		parts := strings.Split(c, "=")
		if len(parts) > 0 {
			if parts[0] == e.Parameter && parts[1] == e.Value {
				return true
			}
		}
	}
	return false
}

// check if value still exists in specified parameter
// example: --authorization-mode=Node,RBAC
func (e *Experiment) checkRemoval(cmdArgs []string) bool {
	for _, c := range cmdArgs {
		parts := strings.Split(c, "=")
		if len(parts) > 0 {
			if parts[0] == e.Parameter {
				args := strings.Split(parts[1], ",")
				for _, a := range args {
					if a == e.Value {
						return false
					}
				}
				return true
			}
		}
	}
	return true
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}
