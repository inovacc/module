// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Action graph execution methods related to coverage.

package work

import (
	"github.com/inovacc/module/internal/base"
	"github.com/inovacc/module/internal/cfg"
	"github.com/inovacc/module/str"
)

// CovData invokes "go tool covdata" with the specified arguments
// as part of the execution of action 'a'.
func (b *Builder) CovData(a *Action, cmdargs ...any) ([]byte, error) {
	cmdline := str.StringList(cmdargs...)
	args := append([]string{}, cfg.BuildToolexec...)
	args = append(args, base.Tool("covdata"))
	args = append(args, cmdline...)
	return b.Shell(a).runOut(a.Objdir, nil, args)
}
