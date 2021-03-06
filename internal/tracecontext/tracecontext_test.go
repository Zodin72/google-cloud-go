// Copyright 2014 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tracecontext

import (
	"reflect"
	"testing"
)

var validData = []byte{0, 0, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 1, 97, 98, 99, 100, 101, 102, 103, 104, 2, 1}

func TestDecode(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		wantTraceID [2]uint64
		wantSpanID  uint64
		wantOpts    byte
		wantOk      bool
	}{
		{
			name:        "nil data",
			data:        nil,
			wantTraceID: [2]uint64{},
			wantSpanID:  0,
			wantOpts:    0,
			wantOk:      false,
		},
		{
			name:        "short data",
			data:        []byte{0, 0, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77},
			wantTraceID: [2]uint64{},
			wantSpanID:  0,
			wantOpts:    0,
			wantOk:      false,
		},
		{
			name:        "wrong field number",
			data:        []byte{0, 1, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77},
			wantTraceID: [2]uint64{},
			wantSpanID:  0,
			wantOpts:    0,
			wantOk:      false,
		},
		{
			name:        "valid data",
			data:        validData,
			wantTraceID: [2]uint64{0x4F4E4D4C4B4A4948, 0x4746454443424140},
			wantSpanID:  0x6867666564636261,
			wantOpts:    1,
			wantOk:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTraceID, gotSpanID, gotOpts, gotOk := Decode(tt.data)
			if !reflect.DeepEqual(gotTraceID, tt.wantTraceID) {
				t.Errorf("Decode() gotTraceID = %v, want %v", gotTraceID, tt.wantTraceID)
			}
			if gotSpanID != tt.wantSpanID {
				t.Errorf("Decode() gotSpanID = %v, want %v", gotSpanID, tt.wantSpanID)
			}
			if gotOpts != tt.wantOpts {
				t.Errorf("Decode() gotOpts = %v, want %v", gotOpts, tt.wantOpts)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Decode() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	tests := []struct {
		name     string
		dst      []byte
		traceID  [2]uint64
		spanID   uint64
		opts     byte
		wantN    int
		wantData []byte
	}{
		{
			name:     "short data",
			dst:      make([]byte, 0),
			traceID:  [2]uint64{5714589967255750984, 5135868584551137600},
			spanID:   0x6867666564636261,
			opts:     1,
			wantN:    -1,
			wantData: make([]byte, 0),
		},
		{
			name:     "valid data",
			dst:      make([]byte, totalLen),
			traceID:  [2]uint64{5714589967255750984, 5135868584551137600},
			spanID:   0x6867666564636261,
			opts:     1,
			wantN:    totalLen,
			wantData: validData,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotN := Encode(tt.dst, tt.traceID, tt.spanID, tt.opts)
			if gotN != tt.wantN {
				t.Errorf("n = %v, want %v", gotN, tt.wantN)
			}
			if gotData := tt.dst; !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("dst = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

func BenchmarkDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Decode(validData)
	}
}

func BenchmarkEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Encode(validData, [2]uint64{1, 1}, 1, 1)
	}
}
