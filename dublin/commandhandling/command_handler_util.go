package commandhandling

import (
	"reflect"
)

// CommandHandlerUtil is a utility to use in composition with a CommandHandler
type CommandHandlerUtil struct {

}

// getMethodName returns the string name of a Command
func (sch CommandHandlerUtil) getMethodName(cmd Command) string {
	var commandFullName string

	if reflect.ValueOf(cmd).Kind() == reflect.Ptr {
		commandFullName = "Handle" + reflect.TypeOf(cmd).Elem().Name()
	} else {
		commandFullName = "Handle" + reflect.TypeOf(cmd).Name()
	}

	return commandFullName
}
