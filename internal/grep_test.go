package internal

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGrep_matchPatternInFile(t *testing.T) {
	type fields struct {
		CaseInSensitive bool
		Recursive       bool
		ReadFromStdIn   bool
		Path            string
		Pattern         string
		OutputFile      string
	}
	type args struct {
		path    string
		pattern string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		WantErr error
	}{
		{
			name:   "file present and more than one match",
			fields: fields{},
			args: args{
				"test.txt",
				"abc",
			},
			want:    4,
			WantErr: nil,
		},
		{
			name:   "file present and 0 match",
			fields: fields{},
			args: args{
				"test.txt",
				"abcdef",
			},
			want:    0,
			WantErr: nil,
		},
		{
			name:   "file present and exact 1 match",
			fields: fields{},
			args: args{
				"test.txt",
				"Abcd",
			},
			want:    1,
			WantErr: nil,
		},
		{
			name:   "file not present return error",
			fields: fields{},
			args: args{
				"notpresent.txt",
				"Abcd",
			},
			want:    1,
			WantErr: errors.New("open notpresent.txt: no such file or directory"),
		},
		{
			name:   "test case insensitive match",
			fields: fields{CaseInSensitive: true},
			args: args{
				"test.txt",
				"ABC",
			},
			want:    5,
			WantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Grep{
				CaseInSensitive: tt.fields.CaseInSensitive,
				Recursive:       tt.fields.Recursive,
				ReadFromStdIn:   tt.fields.ReadFromStdIn,
				Path:            tt.fields.Path,
				Pattern:         tt.fields.Pattern,
				OutputFile:      tt.fields.OutputFile,
			}
			matchedLines, err := g.matchPatternInFile(tt.args.path, tt.args.pattern)
			if err != nil {
				assert.EqualError(t, err, tt.WantErr.Error())
			} else {
				assert.Equal(t, tt.want, len(matchedLines))
			}
		})
	}
}

func TestGrep_matchPatternInStdIn(t *testing.T) {
	type fields struct {
		CaseInSensitive bool
		Recursive       bool
		ReadFromStdIn   bool
		Path            string
		Pattern         string
		OutputFile      string
	}
	type args struct {
		pattern string
		reader  io.Reader
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name:   "test from stdin 1 match ",
			fields: fields{ReadFromStdIn: true},
			args: args{
				pattern: "abc",
				reader:  strings.NewReader("abc"),
			},
			want: []string{"abc"},
		},
		{
			name:   "from stdin more than one match",
			fields: fields{ReadFromStdIn: true},
			args: args{
				pattern: "abc",
				reader:  strings.NewReader("abc\nabcd"),
			},
			want: []string{"abc", "abcd"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Grep{
				CaseInSensitive: tt.fields.CaseInSensitive,
				Recursive:       tt.fields.Recursive,
				ReadFromStdIn:   tt.fields.ReadFromStdIn,
				Path:            tt.fields.Path,
				Pattern:         tt.fields.Pattern,
				OutputFile:      tt.fields.OutputFile,
			}

			got := g.matchPatternInStdIn(tt.args.pattern, tt.args.reader)
			assert.Equal(t, tt.want, got)

		})
	}
}

func TestGrep_WriteOutput(t *testing.T) {
	type fields struct {
		CaseInSensitive bool
		Recursive       bool
		ReadFromStdIn   bool
		Path            string
		Pattern         string
		OutputFile      string
	}
	type args struct {
		matchedLines []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test write to file",
			fields: fields{
				OutputFile: "outFile.txt",
			},
			args: args{matchedLines: []string{"abc", "abcd"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Grep{
				CaseInSensitive: tt.fields.CaseInSensitive,
				Recursive:       tt.fields.Recursive,
				ReadFromStdIn:   tt.fields.ReadFromStdIn,
				Path:            tt.fields.Path,
				Pattern:         tt.fields.Pattern,
				OutputFile:      tt.fields.OutputFile,
			}
			g.WriteOutput(tt.args.matchedLines)
			content, _ := os.ReadFile(tt.fields.OutputFile)
			con := string(content)
			lines := strings.Split(con, "\n")
			fmt.Println(lines)
			assert.Equal(t, 2, len(lines)-1)
			assert.Equal(t, []string{"abc", "abcd"}, tt.args.matchedLines)

		})
	}
}
