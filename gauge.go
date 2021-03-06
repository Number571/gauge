// Copyright 2015 ThoughtWorks, Inc.

// This file is part of Gauge.

// Gauge is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// Gauge is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with Gauge.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"github.com/getgauge/gauge/cmd"
	"github.com/getgauge/gauge/logger"
	"os"
	"runtime/debug"
)

func main() {
	defer recoverPanic()
	if err := cmd.Parse(); err != nil {
		logger.Info(true, err.Error())
		os.Exit(1)
	}
}

func recoverPanic() {
	if r := recover(); r != nil {
		logger.Fatalf(true, "Panicing : %v\n%s", r, string(debug.Stack()))
	}
}
