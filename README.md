# sysproxy-go

A Go library for reading, setting, disabling, and monitoring system proxy settings on Windows, macOS, and Linux.

> **Project role:** This is the system-proxy mechanism library used by [KokoroBox Service](https://github.com/amamiyakokoro/kokorobox-service). KokoroBox Service owns authentication, leases, recovery, guard policy, and user-facing operations; this module only implements platform integration.

Version 2 removes the standalone CLI and HTTP server. KokoroBox Desktop sends system-proxy requests to the authenticated KokoroBox Service API. Integrations can import this library for platform operations, or use KokoroBox Service when they need authentication, leases, guard behavior, recovery, or events.

## Features

- Configure a shared HTTP, HTTPS, and SOCKS proxy
- Configure a PAC URL
- Query or disable the current proxy
- Watch for proxy changes
- Target another desktop session or Windows user

## Install

Add the v2 module to a Go project:

```sh
go get github.com/amamiyakokoro/sysproxy-go/v2@v2.0.0
```

Use the canonical module path:

```go
import "github.com/amamiyakokoro/sysproxy-go/v2/sysproxy"
```

## Go API

```go
package main

import "github.com/amamiyakokoro/sysproxy-go/v2/sysproxy"

func main() {
	if err := sysproxy.SetProxy(&sysproxy.Options{
		Proxy:  "127.0.0.1:7890",
		Bypass: "localhost,127.0.0.1",
	}); err != nil {
		panic(err)
	}
}
```

The main operations are:

```go
sysproxy.SetProxy(options)
sysproxy.SetPac(options)
sysproxy.DisableProxy(options)
sysproxy.QueryProxySettings(options)
sysproxy.WaitProxySettingsChange(ctx, options)
```

Use `sysproxy.OptionsForUser(name)` or `sysproxy.OptionsForProcess(pid)` when operating on another user or desktop session.

## Platform support

| Platform | Backend |
| --- | --- |
| Windows | WinINet or the user registry |
| macOS | `networksetup` |
| Linux | `gsettings` for GNOME-compatible desktops; `kwriteconfig5`/`kwriteconfig6` for KDE |

The module is tested on Linux and compiled for Windows, macOS, and Linux in CI. Version 2 does not publish standalone binaries.

## License

[GNU General Public License v3.0](LICENSE)
