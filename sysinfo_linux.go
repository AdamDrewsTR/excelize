// Copyright 2016 - 2026 The excelize Authors. All rights reserved. Use of
// this source code is governed by a BSD-style license that can be found in
// the LICENSE file.

//go:build linux

package excelize

import "golang.org/x/sys/unix"

// availableMemoryBytes returns the available system memory in bytes on Linux
// using unix.Sysinfo from golang.org/x/sys/unix. Freeram * Unit gives the
// available physical RAM without spawning a subprocess or parsing /proc.
func availableMemoryBytes() uint64 {
	var info unix.Sysinfo_t
	if err := unix.Sysinfo(&info); err != nil {
		return autoTuneFallbackMem
	}
	return info.Freeram * uint64(info.Unit)
}
