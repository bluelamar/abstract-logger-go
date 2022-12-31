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
	"bytes"
	"strings"
	"testing"
)

var (
	errMsg = "bad thing happened"
)

func TestNilTargetLocalTime(t *testing.T) {
	logger := New(nil, true).WithTag("status").WithTag("schedule")

	logger.Errorln(errMsg)
}

func TestNilTargetUtcTime(t *testing.T) {
	logger := New(nil, false)

	logger.SetTimeLabels(":lcl:", ":utc:")

	logger.Errorln(errMsg)
}

func TestDefaultLogMessages(t *testing.T) {

	buf := new(bytes.Buffer)
	logger := New(buf, true)

	logger.Errorln(errMsg)
	outMsg := ErrorLabel + ": " + errMsg
	if !strings.Contains(buf.String(), outMsg) {
		t.Fatalf(`Expected logger buffer(%s) equal to msg(%s)`, buf.String(), outMsg)
	}

	logger.Errorf("%d - %s", 1, errMsg)
	outMsg = ErrorLabel + ": 1 - " + errMsg
	if !strings.Contains(buf.String(), outMsg) {
		t.Fatalf(`Expected logger buffer(%s) equal to msg(%s)`, buf.String(), outMsg)
	}

	logger.Warnln(errMsg)
	outMsg = WarnLabel + ": " + errMsg
	if !strings.Contains(buf.String(), outMsg) {
		t.Fatalf(`Expected logger buffer(%s) equal to msg(%s)`, buf.String(), outMsg)
	}

	logger.Warnf("%d - %s", 1, errMsg)
	outMsg = WarnLabel + ": 1 - " + errMsg
	if !strings.Contains(buf.String(), outMsg) {
		t.Fatalf(`Expected logger buffer(%s) equal to msg(%s)`, buf.String(), outMsg)
	}

	logger.Debugln(errMsg)
	outMsg = DebugLabel + ": " + errMsg
	if !strings.Contains(buf.String(), outMsg) {
		t.Fatalf(`Expected logger buffer(%s) equal to msg(%s)`, buf.String(), outMsg)
	}

	logger.Debugf("%d - %s", 1, errMsg)
	outMsg = DebugLabel + ": 1 - " + errMsg
	if !strings.Contains(buf.String(), outMsg) {
		t.Fatalf(`Expected logger buffer(%s) equal to msg(%s)`, buf.String(), outMsg)
	}

	logger.Infoln(errMsg)
	outMsg = InfoLabel + ": " + errMsg
	if !strings.Contains(buf.String(), outMsg) {
		t.Fatalf(`Expected logger buffer(%s) equal to msg(%s)`, buf.String(), outMsg)
	}

	logger.Infof("%d - %s", 1, errMsg)
	outMsg = InfoLabel + ": 1 - " + errMsg
	if !strings.Contains(buf.String(), outMsg) {
		t.Fatalf(`Expected logger buffer(%s) equal to msg(%s)`, buf.String(), outMsg)
	}
}

func TestNewLabelsLogMessages(t *testing.T) {

	labels := []string{"INK", "DOT", "WHY", "EEL", "FUN"}
	tags := []string{"alpha", "beta"}
	tagsInMsg := "[alpha][beta] "

	buf := new(bytes.Buffer)
	logger := New(buf, true)

	for _, tag := range tags {
		logger = logger.WithTag(tag)
	}

	err := logger.SetLogLabels(labels[0], labels[1], labels[2], labels[3], labels[4])
	if err != nil {
		t.Fatalf(`Expected log labels to be valid: err=%v`, err)
	}

	logger.Errorln(errMsg)
	outMsg := labels[3] + ":" + tagsInMsg + errMsg
	if !strings.Contains(buf.String(), outMsg) {
		t.Fatalf(`Expected logger buffer(%s) equal to msg(%s)`, buf.String(), outMsg)
	}

	logger.Errorf("%d - %s", 1, errMsg)
	outMsg = labels[3] + ":" + tagsInMsg + "1 - " + errMsg
	if !strings.Contains(buf.String(), outMsg) {
		t.Fatalf(`Expected logger buffer(%s) equal to msg(%s)`, buf.String(), outMsg)
	}

	logger.Warnln(errMsg)
	outMsg = labels[2] + ":" + tagsInMsg + errMsg
	if !strings.Contains(buf.String(), outMsg) {
		t.Fatalf(`Expected logger buffer(%s) equal to msg(%s)`, buf.String(), outMsg)
	}

	logger.Warnf("%d - %s", 3, errMsg)
	outMsg = labels[2] + ":" + tagsInMsg + "3 - " + errMsg
	if !strings.Contains(buf.String(), outMsg) {
		t.Fatalf(`Expected logger buffer(%s) equal to msg(%s)`, buf.String(), outMsg)
	}

	logger.Debugln(errMsg)
	outMsg = labels[1] + ":" + tagsInMsg + errMsg
	if !strings.Contains(buf.String(), outMsg) {
		t.Fatalf(`Expected logger buffer(%s) equal to msg(%s)`, buf.String(), outMsg)
	}

	logger.Debugf("%d - %s", 4, errMsg)
	outMsg = labels[1] + ":" + tagsInMsg + "4 - " + errMsg
	if !strings.Contains(buf.String(), outMsg) {
		t.Fatalf(`Expected logger buffer(%s) equal to msg(%s)`, buf.String(), outMsg)
	}

	logger.Infoln(errMsg)
	outMsg = labels[0] + ":" + tagsInMsg + errMsg
	if !strings.Contains(buf.String(), outMsg) {
		t.Fatalf(`Expected logger buffer(%s) equal to msg(%s)`, buf.String(), outMsg)
	}

	logger.Infof("%d - %s", 5, errMsg)
	outMsg = labels[0] + ":" + tagsInMsg + "5 - " + errMsg
	if !strings.Contains(buf.String(), outMsg) {
		t.Fatalf(`Expected logger buffer(%s) equal to msg(%s)`, buf.String(), outMsg)
	}
}
