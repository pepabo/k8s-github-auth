package server

import "testing"

func Test_isJWT(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "jwt",
			args: args{
				token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dnZWRJbkFzIjoiYWRtaW4iLCJpYXQiOjE0MjI3Nzk2Mzh9.gzSraSYS8EXBxLN_oWnFSRgCzcmJmMjLiuyu5CSpyHI",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "not jwt",
			args: args{
				token: "aaaaaaaaaaaaaaaa",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isJWT(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("isJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isJWT() got = %v, want %v", got, tt.want)
			}
		})
	}
}
