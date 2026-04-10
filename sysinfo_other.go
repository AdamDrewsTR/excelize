// Copyright 2016 - 2026 The excelize Authors. All rights reserved. Use of
// this source code is governed by a BSD-style license that can be found in
// the LICENSE file.

//go:build !linux && !darwin && !windows

package excelize

// availableMemoryBytes returns the fallback memory estimate on platforms
// where we have no OS-level memory API available.
func availableMemoryBytes() int64 { return autoTuneFallbackMem }
