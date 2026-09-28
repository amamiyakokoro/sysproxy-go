//go:build !linux && !windows

package sysproxy

import (
	"fmt"
	"runtime"
)

func OptionsForUser(_ string) (*Options, error) {
	return nil, fmt.Errorf("%s does not support targeting a specific user", runtime.GOOS)
}

func OptionsForProcess(_ int) (*Options, error) {
	return nil, fmt.Errorf("%s does not support targeting a specific process", runtime.GOOS)
}
