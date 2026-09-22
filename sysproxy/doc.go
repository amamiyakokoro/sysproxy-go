// Package sysproxy implements platform-specific system-proxy mechanisms.
//
// The package is intentionally policy-free. Callers own authentication,
// privilege management, leases, recovery, guard behavior, and user-facing
// events. KokoroBox integrations provide those concerns through KokoroBox
// Service and use this package only for operating-system access.
package sysproxy
