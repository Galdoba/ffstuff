package actions

import (
	"reflect"
	"testing"
)

func TestScanDir(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		want1   []string
		wantErr bool
	}{
		// TODO: Add test cases.
		//{name: "simple", args: args{`\\192.168.31.4\buffer\IN\@AMEDIA_IN\`}, want: []string{`amedia_tv_series.xml`}, want1: []string{}, wantErr: false},
		{name: "bad", args: args{`\\192.168.31.4\buffer\IN\@GOBLIN\`}, want: []string{}, want1: []string{`2023.10.13`}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ScanDir(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScanDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ScanDir() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ScanDir() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
