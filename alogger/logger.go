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
	"errors"
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

	LocalTimeLabel = "[lcl]"
	UtcTimeLabel   = "[utc]"
)

type LoggerI interface {

	// SetLogLabels allows user to over-ride default severity strings for the associated log messages.
	// Empty string for label is an error.
	SetLogLabels(infoLabel, debugLabel, warnLabel, errorLabel, fatalLabel string) error

	// SetTimeLabels allows user to over-ride default date/time labels for the log messages.
	SetTimeLabels(localTimeLabel, utcTimeLabel string)

	// WithTag allows user to specify a tag that will be prefixed to all log messages.
	// Multiple tags may be added to the logger.
	// For example, tag could represent a component called "status", inwhich case the
	// string "[status]" would be prefixed to the message after the severity label.
	WithTag(tag string) LoggerI

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
	dtLocal  bool
	dtLabel  string
	tags     string
}

func New(target io.Writer, localTime bool) LoggerI {

	if target == nil {
		target = os.Stderr
	}

	logFlags := log.LstdFlags
	dtLbl := LocalTimeLabel
	if !localTime {
		logFlags |= log.LUTC
		dtLbl = UtcTimeLabel
	}

	l := log.New(target, "", logFlags)

	return &defaultLogger{
		infoLbl:  InfoLabel,
		debugLbl: DebugLabel,
		warnLbl:  WarnLabel,
		errorLbl: ErrorLabel,
		fatalLbl: FatalLabel,
		logger:   l,
		dtLocal:  localTime,
		dtLabel:  dtLbl,
		tags:     "",
	}
}

// SetLogLabels allows user to over-ride default labels for the associated log messages.
// Empty string for label is an error.
func (d *defaultLogger) SetLogLabels(infoLabel, debugLabel, warnLabel, errorLabel, fatalLabel string) error {
	if infoLabel == "" {
		return errors.New("invalid info label specified")
	}
	d.infoLbl = infoLabel

	if debugLabel == "" {
		return errors.New("invalid debug label specified")
	}
	d.debugLbl = debugLabel

	if warnLabel == "" {
		return errors.New("invalid warning label specified")
	}
	d.warnLbl = warnLabel

	if errorLabel == "" {
		return errors.New("invalid error label specified")
	}
	d.errorLbl = errorLabel

	if fatalLabel == "" {
		return errors.New("invalid fatal label specified")
	}
	d.fatalLbl = fatalLabel

	return nil
}

// SetTimeLabels allows user to over-ride default date/time labels for the log messages.
func (d *defaultLogger) SetTimeLabels(localTimeLabel, utcTimeLabel string) {
	if d.dtLocal {
		d.dtLabel = localTimeLabel
	} else {
		d.dtLabel = utcTimeLabel
	}
}

// WithTag allows user to specify a tag that will be prefixed to all log messages.
// Multiple tags may be added to the logger.
func (d *defaultLogger) WithTag(tag string) LoggerI {
	if tag == "" {
		return d
	}

	if d.tags == "" {
		d.tags = "[" + tag + "]"
	} else {
		d.tags = fmt.Sprintf("%s[%s]", d.tags, tag)
	}

	return d
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
	if d.tags == "" {
		d.logger.Println(d.dtLabel + " " + label + ": " + msg)
	} else {
		d.logger.Println(d.dtLabel + " " + label + ":" + d.tags + " " + msg)
	}
}
