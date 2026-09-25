# sysproxy-go

Go library for reading, setting, disabling, and watching system proxy settings on Windows, macOS, and Linux. It provides the platform integration used by [KokoroBox Service](https://github.com/amamiyakokoro/kokorobox-service); authentication, leases, recovery, and guard policy belong to the service.

## Install

```sh
go get github.com/amamiyakokoro/sysproxy-go/v2@v2.0.0
```

## Usage

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

The package also provides `SetPac`, `DisableProxy`, `QueryProxySettings`, and `WaitProxySettingsChange`. Pass `*sysproxy.Options` to each operation; `WaitProxySettingsChange` also takes a context. On Windows and Linux, `OptionsForUser` and `OptionsForProcess` can target another user or session.

## Platforms

| Platform | Backend |
| --- | --- |
| Windows | WinINet or user registry |
| macOS | `networksetup` |
| Linux | GNOME `gsettings` or KDE `kwriteconfig5`/`kwriteconfig6` |

Version 2 is a library only; it does not include a CLI or HTTP server.

## License

[GPL-3.0](LICENSE)
