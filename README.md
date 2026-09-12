# sysproxy-go

A small Go library and command-line tool for reading, setting, disabling, and monitoring system proxy settings on Windows, macOS, and Linux.

> **Fork status:** This is a maintained fork of [UruhaLushia/sysproxy-go](https://github.com/UruhaLushia/sysproxy-go), adapted for [KokoroBox-Desktop](https://github.com/amamiyakokoro/KokoroBox-Desktop). It may diverge from upstream, and its releases prioritize the platforms and architectures supported by KokoroBox-Desktop.

## Features

- Configure a shared HTTP, HTTPS, and SOCKS proxy
- Configure a PAC URL
- Query or disable the current proxy
- Watch for proxy changes
- Guard a proxy configuration and restore it after external changes
- Target another desktop session or Windows user
- Optionally expose the operations through a local HTTP service

## Install

Download a prebuilt binary from [GitHub Releases](https://github.com/amamiyakokoro/sysproxy-go/releases), or build the fork from source:

```sh
git clone https://github.com/amamiyakokoro/sysproxy-go.git
cd sysproxy-go
go build -trimpath -o sysproxy .
```

The module retains the upstream import path for compatibility:

```go
import "github.com/UruhaLushia/sysproxy-go/sysproxy"
```

## CLI

```sh
# Set a proxy
sysproxy proxy --server 127.0.0.1:7890 --bypass "localhost,127.0.0.1"

# Wait for the proxy port before changing the system setting
sysproxy proxy --server 127.0.0.1:7890 --wait-server

# Set a PAC URL
sysproxy pac --url http://127.0.0.1:7890/proxy.pac

# Inspect, disable, watch, or guard the setting
sysproxy status
sysproxy disable
sysproxy watch
sysproxy guard --server 127.0.0.1:7890
```

Run `sysproxy --help` or `sysproxy <command> --help` for the complete option reference.

## Go API

```go
package main

import "github.com/UruhaLushia/sysproxy-go/sysproxy"

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

The x64 binaries require an x86-64-v3 processor. Other Go-supported targets may still build from source, but they are not published by this fork.

## Optional HTTP service

The service is excluded from default builds. Enable it with the `sysproxy_server` build tag:

```sh
go build -tags sysproxy_server -o sysproxy .
sysproxy server --network tcp --listen 127.0.0.1:9090
```

It provides `/ping`, `/status`, `/proxy`, `/pac`, `/disable`, and the `/events` server-sent event stream.

## License

[GNU General Public License v3.0](LICENSE)
