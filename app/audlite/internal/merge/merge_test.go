package merge

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewMerge(t *testing.T) {
	type args struct {
		sources []string
	}
	tests := []struct {
		name    string
		args    args
		want    *mergeProc
		wantErr bool
	}{
		{
			name: "test 1",
			args: args{
				[]string{
					`c:\Users\pemaltynov\assets\Bol.L.wav`,
					`c:\Users\pemaltynov\assets\Bol.C.wav`,
					`c:\Users\pemaltynov\assets\Bol.R.wav`,
					`c:\Users\pemaltynov\assets\Bol.Lfe.wav`,
					`c:\Users\pemaltynov\assets\Bol.Ls.wav`,
					`c:\Users\pemaltynov\assets\Bol.Rs.wav`,
				},
			},
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMerge(tt.args.sources...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMerge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMerge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcess(t *testing.T) {
	fmt.Println("test proc")
	input := []string{
		`c:\Users\pemaltynov\assets\Bol.L.wav`,
		`c:\Users\pemaltynov\assets\Bol.R.wav`,
		// `c:\Users\pemaltynov\assets\Bol.C.wav`,
		// `c:\Users\pemaltynov\assets\Bol.Lfe.wav`,
		// `c:\Users\pemaltynov\assets\Bol.Ls.wav`,
		// `c:\Users\pemaltynov\assets\Bol.Rs.wav`,
	}
	mp, _ := NewMerge(input...)
	fmt.Println("------------")
	mp.SetOriginalFPS(fpsToFloat("25/1"))
	mp.SetTargetName(`c:\Users\pemaltynov\assets\Bol.stereo`)
	prompt, err := mp.Prompt()
	fmt.Println(err)
	fmt.Println(prompt)
	fmt.Println("------------")
}

func Test_durationToFl64(t *testing.T) {
	type args struct {
		duration string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name:    "d2f 1",
			args:    args{"02:00:00.000"},
			want:    7200,
			wantErr: false,
		}, // TODO: Add test cases.
		{
			name:    "d2f 2",
			args:    args{"02.016"},
			want:    2.016,
			wantErr: false,
		}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := durationToFl64(tt.args.duration)
			if (err != nil) != tt.wantErr {
				t.Errorf("durationToFl64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("durationToFl64() = %v, want %v", got, tt.want)
			}
			if got == tt.want {
				fmt.Printf("durationToFl64(\"%v\") SUCCESS = got %v, want %v\n", tt.args.duration, got, tt.want)
			}
		})
	}
}
