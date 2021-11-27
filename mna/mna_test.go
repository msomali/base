package mna

import "testing"

func TestGet(t *testing.T) {
	type args struct {
		phoneNumber string
	}
	tests := []struct {
		name    string
		args    args
		want    Operator
		wantErr bool
	}{
		{
			name:    "test vodacom number",
			args:    args{
				phoneNumber: "0765999999",
			},
			want:    Vodacom,
			wantErr: false,
		},
		{
			name:    "test tigo number",
			args:    args{
				phoneNumber: "0712999999",
			},
			want:    Tigo,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.args.phoneNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetThenFilter(t *testing.T) {
	type args struct {
		phoneNumber string
		f           FilterOperatorFunc
		f2          FilterPhoneFunc
	}
	tests := []struct {
		name    string
		args    args
		want    Operator
		wantErr bool
	}{
		{
			name:    "test filter with suffix and pass tigo and vodacom numbers only",
			args:    args{
				phoneNumber: "0712915799",
				f:           OperatorsListFilter(Tigo, Vodacom),
				f2:          FilterBySuffix("799"),
			},
			want:    Tigo | Vodacom,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetThenFilter(tt.args.phoneNumber, tt.args.f, tt.args.f2)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetThenFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetThenFilter() got = %v, want %v", got, tt.want)
			}
		})
	}
}