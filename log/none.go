// Copyright (c) 2019 Chair of Applied Cryptography, Technische Universität
// Darmstadt, Germany. All rights reserved. This file is part of go-perun. Use
// of this source code is governed by the Apache 2.0 license that can be found
// in the LICENSE file.

package log

import (
	"fmt"
	"os"
)

var exit = os.Exit

type none struct{}

func (none) Printf(string, ...interface{}) {}
func (none) Print(...interface{})          {}
func (none) Println(...interface{})        {}
func (none) Tracef(string, ...interface{}) {}
func (none) Debugf(string, ...interface{}) {}
func (none) Infof(string, ...interface{})  {}
func (none) Warnf(string, ...interface{})  {}
func (none) Errorf(string, ...interface{}) {}
func (none) Trace(...interface{})          {}
func (none) Debug(...interface{})          {}
func (none) Info(...interface{})           {}
func (none) Warn(...interface{})           {}
func (none) Error(...interface{})          {}
func (none) Traceln(...interface{})        {}
func (none) Debugln(...interface{})        {}
func (none) Infoln(...interface{})         {}
func (none) Warnln(...interface{})         {}
func (none) Errorln(...interface{})        {}

func (none) Panic(args ...interface{})                 { panic(fmt.Sprint(args...)) }
func (none) Panicf(format string, args ...interface{}) { panic(fmt.Sprintf(format, args...)) }
func (none) Panicln(args ...interface{})               { panic(fmt.Sprintln(args...)) }

func (none) Fatal(args ...interface{}) {
	fmt.Fprint(os.Stderr, args...)
	exit(1)
}

func (none) Fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	exit(1)
}

func (none) Fatalln(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	exit(1)
}

func (n *none) WithField(key string, value interface{}) Logger {
	return n
}

func (n *none) WithFields(Fields) Logger {
	return n
}

func (n *none) WithError(error) Logger {
	return n
}
