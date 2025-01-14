package formatter

import (
	"os"
	"path"
	"reflect"
	"testing"
)

func TestMainWorkflow_parse(t *testing.T) {
	tests := []struct {
		name        string
		w           *MainWorkflow
		wantNMAPRun NMAPRun
		wantErr     bool
		fileContent string
		fileName    string
	}{
		{
			name: "Wrong path (file does not exists)",
			w: &MainWorkflow{
				Config: &Config{
					InputFile: InputFile(""),
				},
			},
			wantNMAPRun: NMAPRun{},
			wantErr:     true,
		},
		{
			name: "Non-xml file",
			w: &MainWorkflow{
				Config: &Config{}, // will be set dynamically
			},
			wantNMAPRun: NMAPRun{},
			wantErr:     true,
			fileName:    "main_workflow_parse_2_test",
			fileContent: "[NOT XML file]",
		},
		{
			name: "XML file (empty content)",
			w: &MainWorkflow{
				Config: &Config{},
			},
			wantNMAPRun: NMAPRun{},
			wantErr:     false,
			fileContent: `<?xml version="1.0"?>
			<?xml-stylesheet href="file:///usr/local/bin/../share/nmap/nmap.xsl" type="text/xsl"?>
			<nmaprun></nmaprun>`,
			fileName: "main_workflow_parse_3_test",
		},
		{
			name: "XML file with some matching output",
			w: &MainWorkflow{
				Config: &Config{},
			},
			wantNMAPRun: NMAPRun{
				Scanner: "nmap",
				Version: "5.59BETA3",
				ScanInfo: ScanInfo{
					Services: "1-1000",
				},
			},
			wantErr: false,
			fileContent: `<?xml version="1.0"?>
			<nmaprun scanner="nmap" version="5.59BETA3">
				<scaninfo services="1-1000"/>
			</nmaprun>`,
			fileName: "main_workflow_parse_4_test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fileName != "" {
				name := path.Join(os.TempDir(), tt.fileName)
				// Creating file with test-case content
				err := os.WriteFile(name, []byte(tt.fileContent), os.ModePerm)
				if err != nil {
					t.Errorf("Could not write file, error %v", err)
				}
				// deferring file removal after the test
				defer os.Remove(name)
				tt.w.Config.InputFile = InputFile(name)
			}
			gotNMAPRun, err := tt.w.parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("MainWorkflow.parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNMAPRun, tt.wantNMAPRun) {
				t.Errorf("MainWorkflow.parse() = %+v, want %+v", gotNMAPRun, tt.wantNMAPRun)
			}
		})
	}
}

func TestMainWorkflow_Execute(t *testing.T) {
	tests := []struct {
		name        string
		w           *MainWorkflow
		wantErr     bool
		fileName    string
		fileContent string
		before      func(file string)
	}{
		{
			name: "Cannot open OutputFile for a write",
			w: &MainWorkflow{
				Config: &Config{
					OutputFile: OutputFile(""),
				},
			},
			wantErr: true,
		},
		{
			name: "OutputFile already exists",
			w: &MainWorkflow{
				Config: &Config{},
			},
			wantErr:     true,
			fileName:    "main_workflow_Execute_2_test",
			fileContent: "",
			before: func(file string) {
				// Creating output file
				os.Create(path.Join(os.TempDir(), file+"_output"))
			},
		},
		{
			name: "Parse of the file has failed",
			w: &MainWorkflow{
				Config: &Config{},
			},
			wantErr:     true,
			fileName:    "main_workflow_Execute_3_test",
			fileContent: "[NOT XML file]",
		},
		{
			name: "Formatter is not defined (OutputFormat == nil)",
			w: &MainWorkflow{
				Config: &Config{},
			},
			wantErr:  true,
			fileName: "main_workflow_Execute_4_test",
			fileContent: `<?xml version="1.0"?>
			<nmaprun></nmaprun>`,
		},
		{
			name: "Empty CSV with header",
			w: &MainWorkflow{
				Config: &Config{
					OutputFormat: CSVOutput,
				},
			},
			wantErr:  false,
			fileName: "main_workflow_Execute_5_test",
			fileContent: `<?xml version="1.0"?>
			<nmaprun></nmaprun>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(tt.fileName)
			}
			if tt.fileName != "" {
				name := path.Join(os.TempDir(), tt.fileName)
				err := os.WriteFile(name, []byte(tt.fileContent), os.ModePerm)
				if err != nil {
					t.Errorf("Could not write file, error %v", err)
				}
				defer os.Remove(name)
				defer os.Remove(name + "_output")
				tt.w.Config.InputFile = InputFile(name)
				tt.w.Config.OutputFile = OutputFile(name + "_output")
			}
			if err := tt.w.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("MainWorkflow.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
