// Copyright 2016 - 2026 The excelize Authors. All rights reserved. Use of
// this source code is governed by a BSD-style license that can be found in
// the LICENSE file.

//go:build darwin

package excelize

import (
	"syscall"
	"unsafe"
)

// sysSysctl is the darwin syscall number for sysctl(2). The stdlib constant
// syscall.SYS_SYSCTL was removed in recent Go versions, so we hardcode the
// value (202), which is stable on both amd64 and arm64 darwin.

// mibHWMemsize is the sysctl MIB for hw.memsize: CTL_HW (6) + HW_MEMSIZE (24).
var mibHWMemsize = [2]int32{6, 24}

// availableMemoryBytes returns an estimate of available system memory on
// macOS by calling sysctl(hw.memsize) via syscall.Syscall6. It reads the
// total physical RAM and assumes 60 % is currently available as a
// conservative baseline for tuning streaming I/O buffers.
func availableMemoryBytes() uint64 {
	var total uint64
	size := uintptr(unsafe.Sizeof(total))
	_, _, errno := syscall.Syscall6(
		syscall.SYS___SYSCTL,
		uintptr(unsafe.Pointer(&mibHWMemsize[0])), 2,
		uintptr(unsafe.Pointer(&total)),
		uintptr(unsafe.Pointer(&size)),
		0, 0,
	)
	if errno != 0 || total == 0 {
		return autoTuneFallbackMem
	}
	// Assume ~60 % of physical RAM is currently available for use.
	return total * 6 / 10
}
