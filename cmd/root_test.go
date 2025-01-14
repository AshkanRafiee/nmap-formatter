/*
Copyright © 2021 vdjagilev

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package cmd

import (
	_ "embed"
	"errors"
	"log"
	"os"
	"path"
	"testing"

	"github.com/spf13/cobra"
	"github.com/vdjagilev/nmap-formatter/formatter"
)

func Test_validate(t *testing.T) {
	type args struct {
		config formatter.Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		before  func()
		after   func()
	}{
		{
			name:    "Wrong output format",
			args:    args{config: formatter.Config{OutputFormat: formatter.OutputFormat("test")}},
			wantErr: true,
			before:  func() {},
			after:   func() {},
		},
		{
			name: "Missing input file",
			args: args{
				config: formatter.Config{
					OutputFormat: formatter.CSVOutput,
				},
			},
			wantErr: true,
			before:  func() {},
			after:   func() {},
		},
		{
			name: "Successful validation",
			args: args{
				config: formatter.Config{
					OutputFormat: formatter.CSVOutput,
					InputFile:    formatter.InputFile(path.Join(os.TempDir(), "formatter_cmd_valid_2")),
				},
			},
			wantErr: false,
			before: func() {
				os.Create(path.Join(os.TempDir(), "formatter_cmd_valid_2"))
			},
			after: func() {
				os.Remove(path.Join(os.TempDir(), "formatter_cmd_valid_2"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()
			defer tt.after()
			if err := validate(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_arguments(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "No XML path argument provided",
			args:    args{},
			wantErr: true,
		},
		{
			name: "No Output format argument provided",
			args: args{
				args: []string{"file.xml"},
			},
			wantErr: true,
		},
		{
			name: "Version argument provided",
			args: args{
				args: []string{"version"},
			},
			wantErr: false,
		},
		{
			name: "2 arguments provided",
			args: args{
				args: []string{
					"file.xml",
					"html",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := arguments(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("arguments() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_run(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	before := func(file string, testWorkflow formatter.Workflow, testConfig formatter.Config) {
		os.Create(file)
		workflow = testWorkflow
		config = testConfig
		config.InputFile = formatter.InputFile(file)
	}
	after := func(file string) {
		os.Remove(file)
		workflow = nil
		config = formatter.Config{}
	}
	tests := []struct {
		name      string
		input     string
		workflow  formatter.Workflow
		config    formatter.Config
		args      args
		runBefore bool
		wantErr   bool
	}{
		{
			name: "Fails validation during the run (no settings at all, will fail)",
			args: args{},
			config: formatter.Config{
				ShowVersion: false,
			},
			wantErr: true,
		},
		{
			name:      "Workflow execution fails",
			input:     path.Join(os.TempDir(), "formatter_cmd_run_1"),
			runBefore: true,
			workflow: &testWorkflow{
				executeResult: errors.New("Bad failure"),
			},
			config: formatter.Config{
				OutputFormat: "csv",
				ShowVersion:  false,
			},
			args:    args{},
			wantErr: true,
		},
		{
			name:      "Shows version using flag",
			runBefore: true,
			config: formatter.Config{
				ShowVersion: true,
			},
			args:    args{},
			wantErr: false,
		},
		{
			name:      "Successful workflow execution",
			input:     path.Join(os.TempDir(), "formatter_cmd_run_2"),
			runBefore: true,
			workflow:  &testWorkflow{},
			config: formatter.Config{
				OutputFormat: "html",
				ShowVersion:  false, // false by default
			},
			args:    args{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.runBefore {
				before(tt.input, tt.workflow, tt.config)
				defer after(tt.input)
			}
			if err := run(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type testWorkflow struct {
	executeResult error
}

func (w *testWorkflow) Execute() (err error) {
	return w.executeResult
}

func (w *testWorkflow) SetConfig(c *formatter.Config) {
	log.Println("testWorkflow -> SetConfig")
}

func Test_shouldShowVersion(t *testing.T) {
	type args struct {
		c    *formatter.Config
		args []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Don't show version (html)",
			args: args{
				c: &formatter.Config{
					ShowVersion: false,
				},
				args: []string{
					"html",
					"path/to/file.xml",
				},
			},
			want: false,
		},
		{
			name: "Don't show version (arguments are used incorrectly)",
			args: args{
				c: &formatter.Config{
					ShowVersion: false,
				},
				args: []string{
					"version",
					"path/to/file.xml",
				},
			},
			want: false,
		},
		{
			name: "Show version (flag)",
			args: args{
				c: &formatter.Config{
					ShowVersion: true,
				},
				args: []string{},
			},
			want: true,
		},
		{
			name: "Show version (argument)",
			args: args{
				c: &formatter.Config{
					ShowVersion: false,
				},
				args: []string{
					"version",
				},
			},
			want: true,
		},
		{
			name: "Show version (both)",
			args: args{
				c: &formatter.Config{
					ShowVersion: true,
				},
				args: []string{
					"version",
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldShowVersion(tt.args.c, tt.args.args); got != tt.want {
				t.Errorf("shouldShowVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
