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
					`c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\app\audlite\Dom_C.wav`,
					`c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\app\audlite\Dom_L.wav`,
					`c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\app\audlite\Dom_R.wav`,
					`c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\app\audlite\Dom_Lfe.wav`,
					`c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\app\audlite\Dom_Rs.wav`,
					`c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\app\audlite\Dom_Ls.wav`,
				},
			},
		},
		{
			name: "test 1",
			args: args{
				[]string{
					`c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\app\audlite\Archipelago_HD_Trailer.20.RUS.clear.mp4`,
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
		`c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\app\audlite\Dom_C.wav`,
		`c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\app\audlite\Dom_L.wav`,
		`c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\app\audlite\Dom_R.wav`,
		`c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\app\audlite\Dom_Lfe.wav`,
		`c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\app\audlite\Dom_Rs.wav`,
		`c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\app\audlite\Dom_Ls.wav`,
	}
	mp, _ := NewMerge(input...)
	fmt.Println("------------")
	mp.SetOriginalFPS(fpsToFloat("24/1"))
	//mp.SetOriginalDuration("00:27:25.14")
	mp.SetTargetName("out")
	mp.tgDuration = 3600.0
	prompt, err := mp.Prompt()
	fmt.Println(err)
	fmt.Println(prompt)
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
