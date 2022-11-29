package service

import (
	"testing"
)

func Test_buildURL(t *testing.T) {
	type args struct {
		urlT    string
		address addressInfo
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "apply templates",
			args: args{
				urlT: "{{.Address}}#{{.Host}}#{{.Port}}",
				address: addressInfo{
					Address: "address",
					Host:    "host",
					Port:    "8080",
				},
			},
			want:    "address#host#8080",
			wantErr: false,
		},
		{
			name: "no templates",
			args: args{
				urlT: "no-templates",
				address: addressInfo{
					Address: "address",
					Host:    "host",
					Port:    "8080",
				},
			},
			want:    "no-templates",
			wantErr: false,
		},
		{
			name: "invalid templates",
			args: args{
				urlT: "{{.Address}}#{{.Host}}#{{.Unknown}}",
				address: addressInfo{
					Address: "address",
					Host:    "host",
					Port:    "8080",
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildURL(tt.args.urlT, tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("buildURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildCommand(t *testing.T) {
	type args struct {
		commandT []string
		address  addressInfo
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "apply templates",
			args: args{
				commandT: []string{"abc", "{{.Address}}", "{{.Host}}", "{{.Port}}"},
				address: addressInfo{
					Address: "address",
					Host:    "host",
					Port:    "8080",
				},
			},
			want:    "abc address host 8080",
			wantErr: false,
		},
		{
			name: "invalid templates",
			args: args{
				commandT: []string{"abc", "{{.Unknown}}"},
				address: addressInfo{
					Address: "address",
					Host:    "host",
					Port:    "8080",
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "no templates",
			args: args{
				commandT: []string{"abc", "def"},
				address: addressInfo{
					Address: "address",
					Host:    "host",
					Port:    "8080",
				},
			},
			want:    "abc def",
			wantErr: false,
		},
		{
			name: "empty command",
			args: args{
				commandT: []string{},
				address: addressInfo{
					Address: "address",
					Host:    "host",
					Port:    "8080",
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildCommand(tt.args.commandT, tt.args.address)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("buildCommand() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if got.String() != tt.want {
				t.Errorf("buildCommand() got = %v, want %v", got, tt.want)
			}
		})
	}
}
