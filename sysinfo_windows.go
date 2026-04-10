// Copyright 2016 - 2026 The excelize Authors. All rights reserved. Use of
// this source code is governed by a BSD-style license that can be found in
// the LICENSE file.

//go:build windows

package excelize

import (
	"syscall"
	"unsafe"
)

// memoryStatusEx mirrors the Win32 MEMORYSTATUSEX structure.
type memoryStatusEx struct {
	dwLength                uint32
	dwMemoryLoad            uint32
	ullTotalPhys            uint64
	ullAvailPhys            uint64
	ullTotalPageFile        uint64
	ullAvailPageFile        uint64
	ullTotalVirtual         uint64
	ullAvailVirtual         uint64
	ullAvailExtendedVirtual uint64
}

var (
	modKernel32           = syscall.NewLazyDLL("kernel32.dll")
	procGlobalMemStatusEx = modKernel32.NewProc("GlobalMemoryStatusEx")
)

// availableMemoryBytes returns the available physical memory in bytes on
// Windows by calling GlobalMemoryStatusEx from kernel32.dll.
func availableMemoryBytes() int64 {
	var ms memoryStatusEx
	ms.dwLength = uint32(unsafe.Sizeof(ms))
	ret, _, _ := procGlobalMemStatusEx.Call(uintptr(unsafe.Pointer(&ms)))
	if ret == 0 {
		return autoTuneFallbackMem
	}
	return int64(ms.ullAvailPhys)
}
