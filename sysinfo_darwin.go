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
// macOS. It reads total physical RAM via sysctl(hw.memsize) and then sums
// vm.page_free_count and vm.page_speculative_count — two per-page counters
// that are readable without root and reflect current memory pressure. The
// result is multiplied by hw.pagesize to get bytes.
//
// "Speculative" pages are pre-faulted but immediately reclaimable, so they
// should be counted as free. This matches what Activity Monitor and vm_stat
// report as available before compressor/inactive pages are considered.
//
// Falls back to total × 0.6 if any sysctl call fails.
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

	pageSize, err := syscall.SysctlUint32("hw.pagesize")
	if err != nil || pageSize == 0 {
		pageSize = 4096
	}

	freePages, err := syscall.SysctlUint32("vm.page_free_count")
	if err != nil {
		return int64(total * 6 / 10)
	}
	specPages, err := syscall.SysctlUint32("vm.page_speculative_count")
	if err != nil {
		specPages = 0 // treat as unavailable, not fatal
	}

	avail := uint64(freePages+specPages) * uint64(pageSize)
	if avail == 0 || avail >= total {
		return int64(total * 6 / 10)
	}
	return int64(avail)
}
