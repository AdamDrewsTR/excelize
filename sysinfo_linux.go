// Copyright 2016 - 2026 The excelize Authors. All rights reserved. Use of
// this source code is governed by a BSD-style license that can be found in
// the LICENSE file.

//go:build linux

package excelize

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// availableMemoryBytes returns the available system memory in bytes on Linux
// by reading MemAvailable from /proc/meminfo. This is the correct metric for
// available RAM: unlike MemFree it includes reclaimable page-cache and
// buffer memory, so it accurately reflects what a new allocation could use.
func availableMemoryBytes() int64 {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return autoTuneFallbackMem
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if !strings.HasPrefix(line, "MemAvailable:") {
			continue
		}
		// Format: "MemAvailable:   12345678 kB"
		fields := strings.Fields(line)
		if len(fields) < 2 {
			break
		}
		kb, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			break
		}
		return kb * 1024
	}
	return autoTuneFallbackMem
}
