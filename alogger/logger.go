// Copyright 2023, Initialize All Once Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package alogger

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	InfoLabel  = "INFO"
	DebugLabel = "DEBUG"
	WarnLabel  = "WARNING"
	ErrorLabel = "ERROR"
	FatalLabel = "FATAL"
)

type LoggerI interface {

	// SetLogLabels allows user to over-ride default labels for the associated log messages.
	SetLogLabels(infoLabel, debugLabel, warnLabel, errorLabel, fatalLabel string)

	Infof(format string, args ...interface{})
	Infoln(args ...interface{})

	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})

	Warnf(format string, args ...interface{})
	Warnln(args ...interface{})

	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})

	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})
}

type defaultLogger struct {
	infoLbl  string
	debugLbl string
	warnLbl  string
	errorLbl string
	fatalLbl string
	logger   *log.Logger
}

func New(target io.Writer, localTime bool) LoggerI {

	if target == nil {
		target = os.Stderr
	}

	logFlags := log.LstdFlags
	if !localTime {
		logFlags |= log.LUTC
	}

	l := log.New(target, "", logFlags)

	return &defaultLogger{
		infoLbl:  InfoLabel,
		debugLbl: DebugLabel,
		warnLbl:  WarnLabel,
		errorLbl: ErrorLabel,
		fatalLbl: FatalLabel,
		logger:   l,
	}
}

// SetLogLabels allows user to over-ride default labels for the associated log messages.
func (d *defaultLogger) SetLogLabels(infoLabel, debugLabel, warnLabel, errorLabel, fatalLabel string) {
	d.infoLbl = infoLabel
	d.debugLbl = debugLabel
	d.warnLbl = warnLabel
	d.errorLbl = errorLabel
	d.fatalLbl = fatalLabel
}

func (d *defaultLogger) Infof(format string, args ...interface{}) {
	d.print(d.infoLbl, fmt.Sprintf(format, args...))
}

func (d *defaultLogger) Infoln(args ...interface{}) {
	d.print(d.infoLbl, fmt.Sprintln(args...))
}

func (d *defaultLogger) Debugf(format string, args ...interface{}) {
	d.print(d.debugLbl, fmt.Sprintf(format, args...))
}

func (d *defaultLogger) Debugln(args ...interface{}) {
	d.print(d.debugLbl, fmt.Sprintln(args...))
}

func (d *defaultLogger) Warnf(format string, args ...interface{}) {
	d.print(d.warnLbl, fmt.Sprintf(format, args...))
}

func (d *defaultLogger) Warnln(args ...interface{}) {
	d.print(d.warnLbl, fmt.Sprintln(args...))
}

func (d *defaultLogger) Errorf(format string, args ...interface{}) {
	d.print(d.errorLbl, fmt.Sprintf(format, args...))
}

func (d *defaultLogger) Errorln(args ...interface{}) {
	d.print(d.errorLbl, fmt.Sprintln(args...))
}

func (d *defaultLogger) Fatalf(format string, args ...interface{}) {
	d.print(d.fatalLbl, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (d *defaultLogger) Fatalln(args ...interface{}) {
	d.print(d.fatalLbl, fmt.Sprintln(args...))
	os.Exit(1)
}

func (d *defaultLogger) print(label, msg string) {
	d.logger.Println(label + ": " + msg)
}
