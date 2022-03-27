package pkg

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// ParamType defines CLI param type.
type ParamType string

const (
	ParamTypeFlag ParamType = "flag"
	ParamTypeArg  ParamType = "argument"
)

var (
	errEmpty = errors.New("empty")
)

// GetStringFlag returns CLI string flag value.
func GetStringFlag(cmd *cobra.Command, flagName string, isOptional bool) (*string, error) {
	if !shouldHandleFlag(cmd, flagName, isOptional) {
		return nil, nil
	}

	v, err := cmd.Flags().GetString(flagName)
	if err != nil {
		return nil, BuildParamErr(flagName, ParamTypeFlag, err)
	}

	if v == "" {
		return nil, BuildParamErr(flagName, ParamTypeFlag, errEmpty)
	}

	return &v, nil
}

// GetUintFlag returns CLI uint flag value.
func GetUintFlag(cmd *cobra.Command, flagName string, isOptional bool) (*uint, error) {
	if !shouldHandleFlag(cmd, flagName, isOptional) {
		return nil, nil
	}

	v, err := cmd.Flags().GetUint(flagName)
	if err != nil {
		return nil, BuildParamErr(flagName, ParamTypeFlag, err)
	}

	return &v, nil
}

// GetUintArg returns CLI uint arg value.
func GetUintArg(argName, argValue string) (uint, error) {
	v, err := strconv.ParseUint(argValue, 10, 16)
	if err != nil {
		return 0, BuildParamErr(argName, ParamTypeArg, err)
	}

	return uint(v), nil
}

// GetBoolFlag returns CLI bool arg value.
func GetBoolFlag(cmd *cobra.Command, flagName string, isOptional bool) (*bool, error) {
	if !shouldHandleFlag(cmd, flagName, isOptional) {
		return nil, nil
	}

	v, err := cmd.Flags().GetBool(flagName)
	if err != nil {
		return nil, BuildParamErr(flagName, ParamTypeArg, err)
	}

	return &v, nil
}

// BuildParamErr builds a human-readable CLI params parsing error.
func BuildParamErr(pName string, pType ParamType, err error) error {
	return fmt.Errorf("parsing %q %s parameter: %w", pName, string(pType), err)
}

func shouldHandleFlag(cmd *cobra.Command, flagName string, isOptional bool) bool {
	if !isOptional {
		return true
	}

	f := cmd.Flags().Lookup(flagName)
	if f == nil || !f.Changed {
		return false
	}

	return true
}
