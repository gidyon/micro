package dbutil

import "testing"

func TestParseQuery(t *testing.T) {
	type args struct {
		query     string
		stopWords []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "1", args: args{query: "hello goog friend", stopWords: []string{"good"}}, want: "hello* goog* friend*"},
		{name: "2", args: args{query: "hello good friend", stopWords: []string{"good"}}, want: "hello* friend*"},
		{name: "2", args: args{query: "hello good friend", stopWords: []string{"good "}}, want: "hello* good* friend*"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseQuery(tt.args.query, tt.args.stopWords...); got != tt.want {
				t.Errorf("ParseQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
