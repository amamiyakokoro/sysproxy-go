//go:build !darwin && !linux && !windows

package sysproxy

import "fmt"

func DisableProxy(_ *Options) error {
	return fmt.Errorf("unsupported operating system")
}

func SetProxy(_ *Options) error {
	return fmt.Errorf("unsupported operating system")
}

func SetPac(_ *Options) error {
	return fmt.Errorf("unsupported operating system")
}

func QueryProxySettings(_ *Options) (*ProxyConfig, error) {
	return nil, fmt.Errorf("unsupported operating system")
}
