# sysproxy-go

A Go library for reading, setting, disabling, and monitoring system proxy settings on Windows, macOS, and Linux. A legacy standalone command-line interface remains available during the transition to a library-only module.

> **Fork status:** This is the system-proxy mechanism library used by [KokoroBox Service](https://github.com/amamiyakokoro/kokorobox-service). KokoroBox Service owns authentication, leases, recovery, guard policy, and user-facing operations; this module only implements platform integration.

## Standalone deprecation

The `sysproxy` executable remains deprecated. Its `guard` command and optional `sysproxy_server` HTTP server have been removed. Other CLI commands remain available during the transition to a library-only module.

KokoroBox Desktop sends system-proxy requests to the authenticated KokoroBox Service API. New integrations should import the `sysproxy` package for platform operations, or use KokoroBox Service when they need authentication, leases, guard behavior, recovery, or events.

## Features

- Configure a shared HTTP, HTTPS, and SOCKS proxy
- Configure a PAC URL
- Query or disable the current proxy
- Watch for proxy changes
- Target another desktop session or Windows user

## Install

Download a prebuilt binary from [GitHub Releases](https://github.com/amamiyakokoro/sysproxy-go/releases), or build from source:

```sh
git clone https://github.com/amamiyakokoro/sysproxy-go.git
cd sysproxy-go
go build -trimpath -o sysproxy .
```

Use the canonical module path:

```go
import "github.com/amamiyakokoro/sysproxy-go/sysproxy"
```

## Deprecated CLI

```sh
# Set a proxy
sysproxy proxy --server 127.0.0.1:7890 --bypass "localhost,127.0.0.1"

# Wait for the proxy port before changing the system setting
sysproxy proxy --server 127.0.0.1:7890 --wait-server

# Set a PAC URL
sysproxy pac --url http://127.0.0.1:7890/proxy.pac

# Inspect, disable, or watch the setting
sysproxy status
sysproxy disable
sysproxy watch
```

Run `sysproxy --help` or `sysproxy <command> --help` for the complete option reference. The command prints a deprecation notice directing KokoroBox users to KokoroBox Service.

## Go API

```go
package main

import "github.com/amamiyakokoro/sysproxy-go/sysproxy"

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

Release builds follow KokoroBox-Desktop's supported matrix:

| Platform | Architectures |
| --- | --- |
| Windows | `amd64-v3`, `arm64` |
| macOS | `amd64-v3`, `arm64` |
| Linux | `amd64-v3`, `arm64` |

The x64 binaries require an x86-64-v3 processor. Other Go-supported targets may still build from source, but they are not published by this project.

## License

[GNU General Public License v3.0](LICENSE)
