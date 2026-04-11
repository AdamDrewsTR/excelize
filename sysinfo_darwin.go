// Copyright 2016 - 2026 The excelize Authors. All rights reserved. Use of
// this source code is governed by a BSD-style license that can be found in
// the LICENSE file.

//go:build darwin

package excelize

import (
	"syscall"
	"unsafe"
)

// mibHWMemsize is the sysctl MIB for hw.memsize: CTL_HW (6) + HW_MEMSIZE (24).
var mibHWMemsize = [2]int32{6, 24}

// availableMemoryBytes returns an estimate of available system memory on
// macOS. It reads total physical RAM via sysctl(hw.memsize) and then reads
// kern.memorystatus_level, which is macOS's built-in 0-100 memory-availability
// score. The score already accounts for free, speculative, purgeable, and
// file-backed pages without requiring root privileges.
//
// Available memory is estimated as:
//
//	total x memorystatus_level / 100
//
// Falls back to total x 0.6 if either sysctl is unavailable.
func availableMemoryBytes() int64 {
	// Read total physical RAM.
	var total uint64
	totalSize := uintptr(unsafe.Sizeof(total))
	_, _, errno := syscall.Syscall6(
		syscall.SYS___SYSCTL,
		uintptr(unsafe.Pointer(&mibHWMemsize[0])), 2,
		uintptr(unsafe.Pointer(&total)),
		uintptr(unsafe.Pointer(&totalSize)),
		0, 0,
	)
	if errno != 0 || total == 0 {
		return autoTuneFallbackMem
	}

	// kern.memorystatus_level is a 0-100 score maintained by the jetsam
	// subsystem. 100 means no pressure; values near 0 indicate the kernel is
	// about to start killing processes to reclaim memory.
	level, err := syscall.SysctlUint32("kern.memorystatus_level")
	if err != nil || level == 0 || level > 100 {
		// Sysctl unavailable or value out of range; use a conservative default.
		return int64(total * 6 / 10)
	}

	return int64(total * uint64(level) / 100)
}
