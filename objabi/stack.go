// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package objabi

import (
	"github.com/inovacc/module/internal/cfg"
)

const StackNosplitBase = 800

func StackNosplit(race bool) int {
	// This arithmetic must match that in runtime/stack.go:stackNosplit.
	return StackNosplitBase * stackGuardMultiplier(race)
}

// stackGuardMultiplier returns a multiplier to apply to the default
// stack guard size. Larger multipliers are used for non-optimized
// builds that have larger stack frames or for specific targets.
func stackGuardMultiplier(race bool) int {
	// This arithmetic must match that in runtime/internal/sys/consts.go:StackGuardMultiplier.
	n := 1
	// On AIX, a larger stack is needed for syscalls.
	if cfg.GOOS == "aix" {
		n += 1
	}
	// The race build also needs more stack.
	if race {
		n += 1
	}
	return n
}