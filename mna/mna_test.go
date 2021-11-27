package mna

import (
	"strings"
	"testing"
)

func FilterByPrefix(prefix string) FilterPhoneFunc {
	return func(phone string) bool {
		return strings.HasPrefix(phone, prefix)
	}
}

func FilterBySubstring(substr string) FilterPhoneFunc {
	return func(phone string) bool {
		return strings.Contains(phone, substr)
	}
}

func FilterBySuffix(suffix string) FilterPhoneFunc {
	return func(phone string) bool {
		return strings.HasSuffix(phone, suffix)
	}
}

func OperatorsListFilter(ops ...Operator) FilterOperatorFunc {
	return func(op Operator) bool {
		for _, operator := range ops {
			if op == operator {
				return true
			}
		}

		return false
	}
}


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
		f2           FilterOperatorFunc
		f1          FilterPhoneFunc
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
				f2:           OperatorsListFilter(Tigo, Vodacom),
				f1:          FilterBySuffix("799"),
			},
			want:    Tigo,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAndFilter(tt.args.phoneNumber, tt.args.f1, tt.args.f2)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAndFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetAndFilter() got = %v, want %v", got, tt.want)
			}
		})
	}
}