<div align="center">

# sysproxy-go

Go library for reading, setting, disabling, and watching system proxy settings.

[License](LICENSE)

</div>

## Features

- Set an HTTP, HTTPS, or SOCKS proxy, or a PAC URL
- Query, disable, and watch proxy settings
- Target another user or desktop session on Windows and Linux

## Supported platforms

| Platform | Backend |
| --- | --- |
| Windows | WinINet or user registry |
| macOS | `networksetup` |
| Linux | GNOME `gsettings` or KDE `kwriteconfig5`/`kwriteconfig6` |

## Get started

```sh
go get github.com/amamiyakokoro/sysproxy-go/v2@v2.0.0
```

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

## Development

Requires Go 1.25+.

```sh
go test ./...
```

## Documentation

- [Go package reference](https://pkg.go.dev/github.com/amamiyakokoro/sysproxy-go/v2/sysproxy)
- [KokoroBox Service](https://github.com/amamiyakokoro/kokorobox-service) uses this library for platform integration and provides authentication, leases, recovery, and guard policy.

## License

Licensed under [GNU GPLv3](LICENSE).
