package validators

import "testing"

func TestVerifyPassword(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test №1 Valid Password",
			args: args{s: "123456Password"},
			want: true,
		},
		{
			name: "Test №2 Invalid Password - no letter",
			args: args{s: "12345678910"},
			want: false,
		},
		{
			name: "Test №3 Invalid Password - no number",
			args: args{s: "Hellopassword"},
			want: false,
		},
		{
			name: "Test №4 Invalid Password - no upper",
			args: args{s: "hellopassword1234567"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VerifyPassword(tt.args.s); got != tt.want {
				t.Errorf("VerifyPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
