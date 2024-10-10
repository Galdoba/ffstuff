package target

import "testing"

func Test_autoRename(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "test with number",
			args:    args{`c:\Users\pemaltynov\Programs\grabber\copy_(6)_Zastryavshie_vne_zhizni_s01_TRL_HD.mp4`},
			want:    `c:\Users\pemaltynov\Programs\grabber\copy_(7)_Zastryavshie_vne_zhizni_s01_TRL_HD.mp4`,
			wantErr: false,
		},
		{
			name:    "test with no number",
			args:    args{`c:\Users\pemaltynov\Programs\grabber\Zastryavshie_vne_zhizni_s01_TRL_HD.mp4`},
			want:    `c:\Users\pemaltynov\Programs\grabber\copy_(1)_Zastryavshie_vne_zhizni_s01_TRL_HD.mp4`,
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := autoRename(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("autoRename() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("autoRename() = %v, want %v", got, tt.want)
			}
		})
	}
}
