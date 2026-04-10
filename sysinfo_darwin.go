// Copyright 2016 - 2026 The excelize Authors. All rights reserved. Use of
// this source code is governed by a BSD-style license that can be found in
// the LICENSE file.

//go:build darwin

package excelize

import "golang.org/x/sys/unix"

// availableMemoryBytes returns an estimate of available system memory on
// macOS. It reads total physical RAM via the hw.memsize sysctl using the
// golang.org/x/sys/unix package (no subprocess, no /proc) and assumes 60 %
// is currently available — a conservative baseline for tuning I/O buffers.
func availableMemoryBytes() uint64 {
	total, err := unix.SysctlUint64("hw.memsize")
	if err != nil {
		return autoTuneFallbackMem
	}
	// Assume ~60 % of physical RAM is available for use.
	return total * 6 / 10
}
