package config

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	type args struct {
		configFile string
		config     interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"0", args{"", &struct{}{}}, true},
		{"1", args{"test_files/config1.json", &struct {
			ListenHost   string `json:"listen_host"`
			ListenPort   int64  `json:"listen_port"`
			ReadTimeout  int64  `json:"http_read_timeout_sec"`
			WriteTimeout int64  `json:"http_write_timeout_sec"`

			GracefulShutdownTimeout int64 `json:"graceful_shutdown_timeout_sec"`

			LogLevel   string `json:"log_level"`
			PathToLogs string `json:"path_to_logs"`
		}{}}, false},
		{"2", args{"test_files/config1.json", struct {
			ListenHost   string `json:"listen_host"`
			ListenPort   int64  `json:"listen_port"`
			ReadTimeout  int64  `json:"http_read_timeout_sec"`
			WriteTimeout int64  `json:"http_write_timeout_sec"`

			GracefulShutdownTimeout int64 `json:"graceful_shutdown_timeout_sec"`

			LogLevel   string `json:"log_level"`
			PathToLogs string `json:"path_to_logs"`
		}{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NewConfig(tt.args.configFile, tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
