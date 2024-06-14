package main

import (
	"flag"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kpetremann/salt-exporter/internal/metrics"
	"github.com/kpetremann/salt-exporter/pkg/listener"
	"github.com/spf13/viper"
)

func TestReadConfigFlagOnly(t *testing.T) {
	tests := []struct {
		name  string
		flags []string
		want  Config
	}{
		{
			name: "simple config, flags only",
			flags: []string{
				"-host=127.0.0.1",
				"-port=8080",
			},
			want: Config{
				LogLevel:      defaultLogLevel,
				ListenAddress: "127.0.0.1",
				ListenPort:    8080,
				IPCFile:       listener.DefaultIPCFilepath,
				PKIDir:        listener.DefaultPKIDirpath,
				TLS: struct {
					Enabled     bool
					Key         string
					Certificate string
				}{
					Enabled:     false,
					Key:         "",
					Certificate: "",
				},
				Metrics: metrics.Config{
					HealthMinions: true,
					Global: struct {
						Filters struct {
							IgnoreTest bool `mapstructure:"ignore-test"`
							IgnoreMock bool `mapstructure:"ignore-mock"`
						}
					}{
						Filters: struct {
							IgnoreTest bool `mapstructure:"ignore-test"`
							IgnoreMock bool `mapstructure:"ignore-mock"`
						}{
							IgnoreTest: false,
							IgnoreMock: false,
						},
					},
					SaltNewJobTotal: struct{ Enabled bool }{
						Enabled: true,
					},
					SaltExpectedResponsesTotal: struct{ Enabled bool }{
						Enabled: true,
					},
					SaltFunctionResponsesTotal: struct {
						Enabled        bool
						AddMinionLabel bool `mapstructure:"add-minion-label"`
					}{
						Enabled:        true,
						AddMinionLabel: false,
					},
					SaltScheduledJobReturnTotal: struct {
						Enabled        bool
						AddMinionLabel bool `mapstructure:"add-minion-label"`
					}{
						Enabled:        true,
						AddMinionLabel: false,
					},
					SaltResponsesTotal: struct{ Enabled bool }{
						Enabled: true,
					},
					SaltFunctionStatus: struct {
						Enabled bool
						Filters struct {
							Functions []string
							States    []string
						}
					}{
						Enabled: true,
						Filters: struct {
							Functions []string
							States    []string
						}{
							Functions: []string{
								"state.highstate",
							},
							States: []string{
								"highstate",
							},
						},
					},
				},
			},
		},
		{
			name: "advanced config, flags only",
			flags: []string{
				"-host=127.0.0.1",
				"-port=8080",
				"-ipc-file=/dev/null",
				"-health-minions=false",
				"-health-functions-filter=test.sls",
				"-health-states-filter=nop",
				"-ignore-test",
				"-ignore-mock",
				"-tls",
				"-tls-cert=./cert",
				"-tls-key=./key",
			},
			want: Config{
				LogLevel:      defaultLogLevel,
				ListenAddress: "127.0.0.1",
				ListenPort:    8080,
				IPCFile:       "/dev/null",
				PKIDir:        "/etc/salt/pki/master",
				TLS: struct {
					Enabled     bool
					Key         string
					Certificate string
				}{
					Enabled:     true,
					Key:         "./key",
					Certificate: "./cert",
				},
				Metrics: metrics.Config{
					HealthMinions: false,
					Global: struct {
						Filters struct {
							IgnoreTest bool `mapstructure:"ignore-test"`
							IgnoreMock bool `mapstructure:"ignore-mock"`
						}
					}{
						Filters: struct {
							IgnoreTest bool `mapstructure:"ignore-test"`
							IgnoreMock bool `mapstructure:"ignore-mock"`
						}{
							IgnoreTest: true,
							IgnoreMock: true,
						},
					},
					SaltNewJobTotal: struct{ Enabled bool }{
						Enabled: true,
					},
					SaltExpectedResponsesTotal: struct{ Enabled bool }{
						Enabled: true,
					},
					SaltFunctionResponsesTotal: struct {
						Enabled        bool
						AddMinionLabel bool `mapstructure:"add-minion-label"`
					}{
						Enabled:        true,
						AddMinionLabel: false,
					},
					SaltScheduledJobReturnTotal: struct {
						Enabled        bool
						AddMinionLabel bool `mapstructure:"add-minion-label"`
					}{
						Enabled:        true,
						AddMinionLabel: false,
					},
					SaltResponsesTotal: struct{ Enabled bool }{
						Enabled: false,
					},
					SaltFunctionStatus: struct {
						Enabled bool
						Filters struct {
							Functions []string
							States    []string
						}
					}{
						Enabled: false,
						Filters: struct {
							Functions []string
							States    []string
						}{
							Functions: []string{
								"test.sls",
							},
							States: []string{
								"nop",
							},
						},
					},
				},
			},
		},
	}

	name := os.Args[0]
	backupArgs := os.Args
	backupCommandLine := flag.CommandLine
	defer func() {
		flag.CommandLine = backupCommandLine
		os.Args = backupArgs
		viper.Reset()
	}()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Args = append([]string{name}, test.flags...)
			flag.CommandLine = flag.NewFlagSet(name, flag.ContinueOnError)
			viper.Reset()

			cfg, err := ReadConfig()

			if diff := cmp.Diff(cfg, test.want); diff != "" {
				t.Errorf("Mismatch for '%s' test:\n%s", test.name, diff)
			}

			if err != nil {
				t.Errorf("Unexpected error for '%s': '%s'", test.name, err)
			}
		})
	}
}

func TestConfigFileOnly(t *testing.T) {
	name := os.Args[0]
	backupArgs := os.Args
	backupCommandLine := flag.CommandLine
	defer func() {
		flag.CommandLine = backupCommandLine
		os.Args = backupArgs
		viper.Reset()
	}()

	flags := []string{
		"-config-file=config_test.yml",
	}

	os.Args = append([]string{name}, flags...)
	flag.CommandLine = flag.NewFlagSet(name, flag.ContinueOnError)
	viper.Reset()

	cfg, err := ReadConfig()

	want := Config{
		LogLevel:      "info",
		ListenAddress: "127.0.0.1",
		ListenPort:    2113,
		IPCFile:       "/dev/null",
		PKIDir:        "/tmp/pki",
		TLS: struct {
			Enabled     bool
			Key         string
			Certificate string
		}{
			Enabled:     true,
			Key:         "/path/to/key",
			Certificate: "/path/to/certificate",
		},
		Metrics: metrics.Config{
			HealthMinions: true,
			Global: struct {
				Filters struct {
					IgnoreTest bool `mapstructure:"ignore-test"`
					IgnoreMock bool `mapstructure:"ignore-mock"`
				}
			}{
				Filters: struct {
					IgnoreTest bool `mapstructure:"ignore-test"`
					IgnoreMock bool `mapstructure:"ignore-mock"`
				}{
					IgnoreTest: true,
					IgnoreMock: false,
				},
			},
			SaltNewJobTotal: struct{ Enabled bool }{
				Enabled: true,
			},
			SaltExpectedResponsesTotal: struct{ Enabled bool }{
				Enabled: true,
			},
			SaltFunctionResponsesTotal: struct {
				Enabled        bool
				AddMinionLabel bool `mapstructure:"add-minion-label"`
			}{
				Enabled:        true,
				AddMinionLabel: true,
			},
			SaltScheduledJobReturnTotal: struct {
				Enabled        bool
				AddMinionLabel bool `mapstructure:"add-minion-label"`
			}{
				Enabled:        true,
				AddMinionLabel: true,
			},
			SaltResponsesTotal: struct{ Enabled bool }{
				Enabled: true,
			},
			SaltFunctionStatus: struct {
				Enabled bool
				Filters struct {
					Functions []string
					States    []string
				}
			}{
				Enabled: true,
				Filters: struct {
					Functions []string
					States    []string
				}{
					Functions: []string{
						"state.sls",
					},
					States: []string{
						"test",
					},
				},
			},
		},
	}

	if diff := cmp.Diff(cfg, want); diff != "" {
		t.Errorf("Mismatch:\n%s", diff)
	}

	if err != nil {
		t.Errorf("Unexpected error: '%s'", err)
	}
}

func TestConfigFileWithFlags(t *testing.T) {
	name := os.Args[0]
	backupArgs := os.Args
	backupCommandLine := flag.CommandLine
	defer func() {
		flag.CommandLine = backupCommandLine
		os.Args = backupArgs
		viper.Reset()
	}()

	flags := []string{
		"-config-file=config_test.yml",
		"-host=127.0.0.1",
		"-port=8080",
		"-health-minions=false",
		"-health-functions-filter=test.sls",
		"-health-states-filter=nop",
		"-ignore-mock",
		"-ipc-file=/somewhere",
	}

	os.Args = append([]string{name}, flags...)
	flag.CommandLine = flag.NewFlagSet(name, flag.ContinueOnError)
	viper.Reset()

	cfg, err := ReadConfig()
	want := Config{
		LogLevel:      "info",
		ListenAddress: "127.0.0.1",
		ListenPort:    8080,
		IPCFile:       "/somewhere",
		PKIDir:        "/tmp/pki",
		TLS: struct {
			Enabled     bool
			Key         string
			Certificate string
		}{
			Enabled:     true,
			Key:         "/path/to/key",
			Certificate: "/path/to/certificate",
		},
		Metrics: metrics.Config{
			HealthMinions: false,
			Global: struct {
				Filters struct {
					IgnoreTest bool `mapstructure:"ignore-test"`
					IgnoreMock bool `mapstructure:"ignore-mock"`
				}
			}{
				Filters: struct {
					IgnoreTest bool `mapstructure:"ignore-test"`
					IgnoreMock bool `mapstructure:"ignore-mock"`
				}{
					IgnoreTest: true,
					IgnoreMock: true,
				},
			},
			SaltNewJobTotal: struct{ Enabled bool }{
				Enabled: true,
			},
			SaltExpectedResponsesTotal: struct{ Enabled bool }{
				Enabled: true,
			},
			SaltFunctionResponsesTotal: struct {
				Enabled        bool
				AddMinionLabel bool `mapstructure:"add-minion-label"`
			}{
				Enabled:        true,
				AddMinionLabel: true,
			},
			SaltScheduledJobReturnTotal: struct {
				Enabled        bool
				AddMinionLabel bool `mapstructure:"add-minion-label"`
			}{
				Enabled:        true,
				AddMinionLabel: true,
			},
			SaltResponsesTotal: struct{ Enabled bool }{
				Enabled: true,
			},
			SaltFunctionStatus: struct {
				Enabled bool
				Filters struct {
					Functions []string
					States    []string
				}
			}{
				Enabled: true,
				Filters: struct {
					Functions []string
					States    []string
				}{
					Functions: []string{
						"test.sls",
					},
					States: []string{
						"nop",
					},
				},
			},
		},
	}

	if diff := cmp.Diff(cfg, want); diff != "" {
		t.Errorf("Mismatch:\n%s", diff)
	}

	if err != nil {
		t.Errorf("Unexpected error: '%s'", err)
	}
}
