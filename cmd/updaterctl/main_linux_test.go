//go:build linux

package main

import "testing"

func TestDefaultExeNameLinux(t *testing.T) {
	if got := defaultExeName(); got != "agent" {
		t.Fatalf("defaultExeName() = %q, want %q", got, "agent")
	}
}
