package internal

import (
	"fmt"
)

type Depth int

const (
	DepthService  = 0
	DepthMethod   = 1
	DepthArgument = 2
	//DepthArgumentType = 3
	//DepthAll          = 4
)

func DepthFromString(s string) (Depth, error) {
	switch s {
	case "", "service":
		return DepthService, nil
	case "method", "operation":
		return DepthMethod, nil
	case "argument":
		return DepthArgument, nil
	//case "argument_type", "argumentType":
	//	return DepthArgumentType, nil
	//case "all":
	//	return DepthAll, nil
	default:
		return 0, fmt.Errorf("unknown depth value %q", s)
	}
}
