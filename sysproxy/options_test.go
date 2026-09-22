package sysproxy

import "testing"

func TestConcurrentApplyDefaultsToSerialized(t *testing.T) {
	if DefaultConcurrent() {
		t.Fatal("system proxy mutations must be serialized by default")
	}
	if resolveConcurrentApply(nil) {
		t.Fatal("nil options must use the serialized default")
	}
}

func TestConcurrentApplyCanBeExplicitlyEnabled(t *testing.T) {
	enabled := true
	disabled := false
	if !resolveConcurrentApply(&Options{Concurrent: &enabled}) {
		t.Fatal("explicit concurrent mode was ignored")
	}
	if resolveConcurrentApply(&Options{Concurrent: &disabled}) {
		t.Fatal("explicit serialized mode was ignored")
	}
}
