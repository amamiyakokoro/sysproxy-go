package sysproxy

import "context"

// Keep the mechanism API consumed by KokoroBox Service explicit. These compile
// time assertions prevent compatibility-only CLI work from becoming the
// package contract again.
var (
	_ func(*Options) error                          = SetProxy
	_ func(*Options) error                          = SetPac
	_ func(*Options) error                          = DisableProxy
	_ func(*Options) (*ProxyConfig, error)          = QueryProxySettings
	_ func(context.Context, *Options) error         = WaitProxySettingsChange
	_ func(context.Context, *Options, func()) error = WaitProxySettingsChangeReady
)

func Example() {
	options := &Options{
		Proxy:  "127.0.0.1:7890",
		Bypass: "localhost,127.0.0.1",
	}

	_ = options
	// Pass options to SetProxy, SetPac, DisableProxy, or QueryProxySettings.
}
