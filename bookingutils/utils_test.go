package bookingutils

import (
	"reflect"
	"testing"
)

func TestStringSliceToInt(t *testing.T) {
	type args struct {
		strings []string
	}
	tests := []struct {
		name               string
		args               args
		wantConvertedUints []uint
		wantErr            bool
	}{
		{
			"Error", args{strings: []string{"", "a", "b"}}, []uint{}, true,
		},
		{
			"converted correctly", args{strings: []string{"0", "01", "123"}}, []uint{0, 1, 123}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConvertedUints, err := StringSliceToInt(tt.args.strings)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringSliceToInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotConvertedUints, tt.wantConvertedUints) {
				t.Errorf("StringSliceToInt() = %v, want %v", gotConvertedUints, tt.wantConvertedUints)
			}
		})
	}
}

func TestStringToUint(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    uint
		wantErr bool
	}{
		{"error negative number", args{input: "1.1"}, 0, true},
		{"error invalid number", args{input: "abc"}, 0, true},
		{"valid converting", args{input: "3"}, 3, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StringToUint(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringToUint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StringToUint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrToFloat(t *testing.T) {
	type args struct {
		strValue string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{"error not a number", args{strValue: "x1"}, 0, true},
		{"error nil number", args{strValue: ""}, 0, true},
		{"valid conversion from string to float", args{strValue: "-2.54"}, -2.54, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StrToFloat(tt.args.strValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrToFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StrToFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUintToString(t *testing.T) {
	type args struct {
		uintValue uint
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"convert to String", args{uintValue: 999}, "999"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UintToString(tt.args.uintValue); got != tt.want {
				t.Errorf("UintToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatFloatToAmmount(t *testing.T) {
	type args struct {
		value float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"format a long float to decimal", args{value: 9.9564}, "9.96"},
		{"format an integer to decimal", args{value: 3}, "3.00"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatFloatToAmmount(tt.args.value); got != tt.want {
				t.Errorf("FormatFloatToAmmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlmostZero(t *testing.T) {
	type args struct {
		val float64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"should be calculated as zero positive", args{val: 0.05}, true},
		{"should be calculated as zero negative", args{val: -0.05}, true},
		{"should be not as zero negative", args{val: -0.059}, false},
		{"should be not as zero positive", args{val: 0.059}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AlmostZero(tt.args.val); got != tt.want {
				t.Errorf("AlmostZero() = %v, want %v", got, tt.want)
			}
		})
	}
}
