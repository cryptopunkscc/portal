package main

import (
	"reflect"
	"testing"
)

func Test_parseArgs(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		wantM Make
		wantA []string
	}{
		{
			name:  "no args",
			args:  []string{},
			wantM: All,
			wantA: []string{"linux", "windows"},
		},
		{
			name:  "parse literal",
			args:  []string{"i"},
			wantM: Installer,
		},
		{
			name:  "parse literals",
			args:  []string{"ladpi"},
			wantM: All,
		},
		{
			name:  "parse numeric",
			args:  []string{"62"},
			wantM: All,
		},
		{
			name:  "parse goos",
			args:  []string{"", "a", "b", "c"},
			wantM: All,
			wantA: []string{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotM, gotA := parseArgs(tt.args)
			if gotM != tt.wantM {
				t.Errorf("parseArgs() gotM = %v, want %v", gotM, tt.wantM)
			}
			if !reflect.DeepEqual(gotA, tt.wantA) {
				t.Errorf("parseArgs() gotA = %v, want %v", gotA, tt.wantA)
			}
		})
	}
}
