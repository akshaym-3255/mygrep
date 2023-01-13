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

	tests := []struct {
		name    string
		fields  fields
		want    []string
		WantErr error
	}{
		{
			name: "file present and more than one match",
			fields: fields{
				Path:    "../testdata/test.txt",
				Pattern: "abc",
			},
			want:    []string{"abcd", "abcd", "abc", "abc"},
			WantErr: nil,
		},
		{
			name: "file present and 0 match",
			fields: fields{
				Path:    "../testdata/test.txt",
				Pattern: "abcdef",
			},

			want:    nil,
			WantErr: nil,
		},
		{
			name: "file present and exact 1 match",
			fields: fields{
				Path:    "../testdata/test.txt",
				Pattern: "Abcd",
			},
			want:    []string{"Abcd"},
			WantErr: nil,
		},
		{
			name: "file not present return error",
			fields: fields{
				Path:    "notPresent.txt",
				Pattern: "Abcd",
			},
			want:    nil,
			WantErr: errors.New("open notPresent.txt: no such file or directory"),
		},
		{
			name: "test case insensitive match",
			fields: fields{
				CaseInSensitive: true,
				Path:            "../testdata/test.txt",
				Pattern:         "ABC",
			},
			want:    []string{"abcd", "abcd", "Abcd", "abc", "abc"},
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
			matchedLines, err := g.matchPatternInFile()
			if err != nil {
				fmt.Println(err.Error())
			} else {
				assert.Equal(t, tt.want, matchedLines)
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
		reader io.Reader
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name: "test from stdin 1 match ",
			fields: fields{
				ReadFromStdIn: true,
				Pattern:       "abc",
			},
			args: args{
				reader: strings.NewReader("abc"),
			},
			want: []string{"abc"},
		},
		{
			name: "from stdin more than one match",
			fields: fields{
				ReadFromStdIn: true,
				Pattern:       "abc",
			},
			args: args{
				reader: strings.NewReader("abc\nabcd"),
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

			got := g.matchPatternInStdIn(tt.args.reader)
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
				OutputFile: "../testdata/outFile.txt",
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

func TestGrep_matchPatternInDir(t *testing.T) {
	type fields struct {
		CaseInSensitive bool
		Recursive       bool
		ReadFromStdIn   bool
		Path            string
		Pattern         string
		OutputFile      string
	}

	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr error
	}{
		{
			name: "match pattern in directory successful execution",
			fields: fields{
				Recursive: true,
				Path:      "../testdata/",
				Pattern:   "abc",
			},
			want:    []string{"../testdata/outFile.txt   abc", "../testdata/outFile.txt   abcd", "../testdata/test.txt   abcd", "../testdata/test.txt   abcd", "../testdata/test.txt   abc", "../testdata/test.txt   abc"},
			wantErr: nil,
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
			matchedLines, _ := g.matchPatternInDir()
			assert.Equal(t, tt.want, matchedLines)
		})
	}
}
