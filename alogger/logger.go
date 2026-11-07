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
	"sync"
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

type LabelType int

const (
	InfoLabelType  LabelType = iota // 0
	DebugLabelType                  // 1
	WarnLabelType                   // 2
	ErrorLabelType                  // 3
	FatalLabelType                  // 4
)

// LoggerI provides leveled logging with configurable labels, time labels, and tags.
// Set methods are expected to mutate the logger.
// With methods are expected to return a new configured logger without changing the original.
// Fatal methods log at fatal severity but may or may not exit according to implementation.
type LoggerI interface {

	// SetLogLabels allows user to override default severity strings for the associated log messages.
	// Empty string for label is an error.
	// An implementation is expected to mutate the logger.
	SetLogLabels(infoLabel, debugLabel, warnLabel, errorLabel, fatalLabel string) error

	// SetTimeLabels allows user to over-ride default date/time labels for the log messages.
	// Default labels are: LocalTimeLabel and UtcTimeLabel.
	// An implementation is expected to mutate the logger.
	SetTimeLabels(localTimeLabel, utcTimeLabel string)

	// WithLogLabels allows user to over-ride default severity strings for the associated log messages.
	// Empty string for label is an error.
	// An implementation Should be Concurrency safe, returning a new LoggerI that will use the labels.
	WithLogLabels(infoLabel, debugLabel, warnLabel, errorLabel, fatalLabel string) (LoggerI, error)

	// WithTimeLabels allows user to over-ride default date/time labels for the log messages.
	// An implementation Should be Concurrency safe, returning a new LoggerI that will use the labels.
	WithTimeLabels(localTimeLabel, utcTimeLabel string) LoggerI

	// WithTag allows user to specify a tag that will be prefixed to all log messages.
	// Multiple tags may be added to the logger.
	// For example, tag could represent a component called "status", in which case the
	// string "[status]" would be prefixed to the message after the severity label.
	// An implementation Should be Concurrency safe, returning a new LoggerI that will use the tag.
	WithTag(tag string) LoggerI

	Infof(format string, args ...interface{})
	Infoln(args ...interface{})

	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})

	Warnf(format string, args ...interface{})
	Warnln(args ...interface{})

	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})

	// Fatalf logs at fatal severity - it may or may not terminate the process according to implementation.
	Fatalf(format string, args ...interface{})
	// Fatalln logs at fatal severity - it may or may not terminate the process according to implementation.
	Fatalln(args ...interface{})
}

type defaultLogger struct {
	mu       *sync.RWMutex
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
		mu:       new(sync.RWMutex),
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

// clone is used by the With methods to make changes to the logger safe.
func (d *defaultLogger) clone() *defaultLogger {
	d.mu.RLock()
	defer d.mu.RUnlock()

	clone := *d
	return &clone
}

// SetLogLabels allows user to over-ride default labels for the associated log messages.
// Empty string for label is an error.
func (d *defaultLogger) setLogLabels(infoLabel, debugLabel, warnLabel, errorLabel, fatalLabel string) error {
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

func (d *defaultLogger) SetLogLabels(infoLabel, debugLabel, warnLabel, errorLabel, fatalLabel string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.setLogLabels(infoLabel, debugLabel, warnLabel, errorLabel, fatalLabel)
}

func (d *defaultLogger) WithLogLabels(infoLabel, debugLabel, warnLabel, errorLabel, fatalLabel string) (LoggerI, error) {
	clone := d.clone()
	return clone, clone.setLogLabels(infoLabel, debugLabel, warnLabel, errorLabel, fatalLabel)
}

// SetTimeLabels allows user to override default date/time labels for the log messages.
func (d *defaultLogger) setTimeLabels(localTimeLabel, utcTimeLabel string) {
	if d.dtLocal {
		d.dtLabel = localTimeLabel
	} else {
		d.dtLabel = utcTimeLabel
	}
}

func (d *defaultLogger) SetTimeLabels(localTimeLabel, utcTimeLabel string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.setTimeLabels(localTimeLabel, utcTimeLabel)
}

func (d *defaultLogger) WithTimeLabels(localTimeLabel, utcTimeLabel string) LoggerI {
	clone := d.clone()
	clone.setTimeLabels(localTimeLabel, utcTimeLabel)

	return clone
}

// WithTag allows user to specify a tag that will be prefixed to all log messages.
// Multiple tags may be added to the logger.
func (d *defaultLogger) WithTag(tag string) LoggerI {
	clone := d.clone()
	if tag == "" {
		return clone
	}

	if clone.tags == "" {
		clone.tags = "[" + tag + "]"
	} else {
		clone.tags = fmt.Sprintf("%s[%s]", clone.tags, tag)
	}

	return clone
}

func (d *defaultLogger) Infof(format string, args ...interface{}) {
	d.print(InfoLabelType, fmt.Sprintf(format, args...))
}

func (d *defaultLogger) Infoln(args ...interface{}) {
	d.print(InfoLabelType, fmt.Sprintln(args...))
}

func (d *defaultLogger) Debugf(format string, args ...interface{}) {
	d.print(DebugLabelType, fmt.Sprintf(format, args...))
}

func (d *defaultLogger) Debugln(args ...interface{}) {
	d.print(DebugLabelType, fmt.Sprintln(args...))
}

func (d *defaultLogger) Warnf(format string, args ...interface{}) {
	d.print(WarnLabelType, fmt.Sprintf(format, args...))
}

func (d *defaultLogger) Warnln(args ...interface{}) {
	d.print(WarnLabelType, fmt.Sprintln(args...))
}

func (d *defaultLogger) Errorf(format string, args ...interface{}) {
	d.print(ErrorLabelType, fmt.Sprintf(format, args...))
}

func (d *defaultLogger) Errorln(args ...interface{}) {
	d.print(ErrorLabelType, fmt.Sprintln(args...))
}

// Fatalf logs at fatal severity - it does not terminate the process.
func (d *defaultLogger) Fatalf(format string, args ...interface{}) {
	d.print(FatalLabelType, fmt.Sprintf(format, args...))
}

// Fatalln logs at fatal severity - it does not terminate the process.
func (d *defaultLogger) Fatalln(args ...interface{}) {
	d.print(FatalLabelType, fmt.Sprintln(args...))
}

func (d *defaultLogger) print(labelType LabelType, msg string) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var label string

	switch labelType {
	case InfoLabelType:
		label = d.infoLbl
	case DebugLabelType:
		label = d.debugLbl
	case WarnLabelType:
		label = d.warnLbl
	case ErrorLabelType:
		label = d.errorLbl
	case FatalLabelType:
		label = d.fatalLbl
	}

	if d.tags == "" {
		d.logger.Print(d.dtLabel + " " + label + ": " + msg)
	} else {
		d.logger.Print(d.dtLabel + " " + label + ":" + d.tags + " " + msg)
	}
}
