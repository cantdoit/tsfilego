package writertest

import (
	"Golang/internal/file"
	"Golang/internal/writer"
	"os"
	"testing"
)

func TestTsFileWriter_Open(t *testing.T) {
	tests := []struct {
		name        string
		filePath    string
		flags       int
		mode        os.FileMode
		setupFunc   func(filePath string) error
		cleanupFunc func(filePath string)
		wantErr     bool
	}{
		{
			name:     "nonexistent_path",
			filePath: "testfile1.tsfile",
			flags:    os.O_CREATE | os.O_WRONLY,
			mode:     0644,
			setupFunc: func(filePath string) error {
				return nil
			},
			cleanupFunc: func(filePath string) {
				err := os.Remove(filePath)
				if err != nil {
					return
				}
			},
			wantErr: true,
		},
		{
			name:     "file_already_exists",
			filePath: "testfile2.tsfile",
			flags:    os.O_CREATE | os.O_WRONLY,
			mode:     0644,
			setupFunc: func(filePath string) error {
				_, err := os.Create(filePath)
				return err
			},
			cleanupFunc: func(filePath string) {
				err := os.Remove(filePath)
				if err != nil {
					return
				}
			},
			wantErr: true,
		},
		{
			name:     "invalid_path",
			filePath: "/invalid_path/testfile3.tsfile",
			flags:    os.O_CREATE | os.O_WRONLY,
			mode:     0644,
			setupFunc: func(filePath string) error {
				return nil
			},
			cleanupFunc: func(filePath string) {},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.setupFunc(tt.filePath); err != nil {
				t.Fatalf("failed to set up test: %v", err)
			}
			// Cleanup after the test
			defer tt.cleanupFunc(tt.filePath)

			tf := &writer.TsFileWriter{WriteFile: &file.WriteFile{}}

			err := tf.Open(tt.filePath, tt.flags, tt.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
