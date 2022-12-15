package models

import (
	"math/big"
	"testing"
)

func Test_FiniteField_Neg(t *testing.T) {
	type arg struct {
		a *FiniteField
	}

	prime := big.NewInt(223)

	tests := []struct {
		name string
		arg  arg
		want FiniteField
	}{
		{
			name: "-(0) = 0",
			arg: arg{
				a: NewFiniteField(big.NewInt(0), *prime),
			},
			want: *NewFiniteField(big.NewInt(0), *prime),
		},
		{
			name: "-(1) = -1",
			arg: arg{
				a: NewFiniteField(big.NewInt(1), *prime),
			},
			want: *NewFiniteField(big.NewInt(-1), *prime),
		},
		{
			name: "-(3) = -3",
			arg: arg{
				a: NewFiniteField(big.NewInt(3), *prime),
			},
			want: *NewFiniteField(big.NewInt(-3), *prime),
		},
		{
			name: "-prime = 0",
			arg: arg{
				a: NewFiniteField(prime, *prime),
			},
			want: *NewFiniteField(big.NewInt(0), *prime),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := new(FiniteField).Neg(tt.arg.a); !got.Equals(&tt.want) {
				t.Errorf("%v : Add() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_FiniteField_Add(t *testing.T) {
	type args struct {
		a *FiniteField
		b *FiniteField
	}

	two := big.NewInt(2)
	prime := big.NewInt(0).Exp(two, big.NewInt(256), nil)
	prime.Sub(prime, big.NewInt(0).Exp(two, big.NewInt(32), nil))
	prime.Sub(prime, big.NewInt(0).Exp(two, big.NewInt(9), nil))
	prime.Sub(prime, big.NewInt(0).Exp(two, big.NewInt(8), nil))
	prime.Sub(prime, big.NewInt(0).Exp(two, big.NewInt(7), nil))
	prime.Sub(prime, big.NewInt(0).Exp(two, big.NewInt(6), nil))
	prime.Sub(prime, big.NewInt(0).Exp(two, big.NewInt(4), nil))
	prime.Sub(prime, big.NewInt(1))

	tests := []struct {
		name string
		args args
		want FiniteField
	}{
		{
			name: "0 + 0 = 0",
			args: args{
				NewFiniteField(big.NewInt(0), *prime),
				NewFiniteField(big.NewInt(0), *prime),
			},
			want: *NewFiniteField(big.NewInt(0), *prime),
		},
		{
			name: "0 + 1 = 1",
			args: args{
				NewFiniteField(big.NewInt(0), *prime),
				NewFiniteField(big.NewInt(1), *prime),
			},
			want: *NewFiniteField(big.NewInt(1), *prime),
		},
		{
			name: "1 + 0 = 0",
			args: args{
				NewFiniteField(big.NewInt(1), *prime),
				NewFiniteField(big.NewInt(0), *prime),
			},
			want: *NewFiniteField(big.NewInt(1), *prime),
		},
		{
			name: "10 + 10 = 20",
			args: args{
				NewFiniteField(big.NewInt(10), *prime),
				NewFiniteField(big.NewInt(10), *prime),
			},
			want: *NewFiniteField(big.NewInt(20), *prime),
		},
		{
			name: "prime + 1 = 1",
			args: args{
				NewFiniteField(prime, *prime),
				NewFiniteField(big.NewInt(1), *prime),
			},
			want: *NewFiniteField(big.NewInt(1), *prime),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFiniteField(big.NewInt(0), *prime).Add(tt.args.a, tt.args.b); !got.Equals(&tt.want) {
				t.Errorf("%v : Add() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_FiniteField_Sub(t *testing.T) {
	type args struct {
		a *FiniteField
		b *FiniteField
	}

	two := big.NewInt(2)
	prime := big.NewInt(0).Exp(two, big.NewInt(256), nil)
	prime.Sub(prime, big.NewInt(0).Exp(two, big.NewInt(32), nil))
	prime.Sub(prime, big.NewInt(0).Exp(two, big.NewInt(9), nil))
	prime.Sub(prime, big.NewInt(0).Exp(two, big.NewInt(8), nil))
	prime.Sub(prime, big.NewInt(0).Exp(two, big.NewInt(7), nil))
	prime.Sub(prime, big.NewInt(0).Exp(two, big.NewInt(6), nil))
	prime.Sub(prime, big.NewInt(0).Exp(two, big.NewInt(4), nil))
	prime.Sub(prime, big.NewInt(1))

	tests := []struct {
		name string
		args args
		want FiniteField
	}{
		{
			name: "0 - 0 = 0",
			args: args{
				a: NewFiniteField(big.NewInt(0), *prime),
				b: NewFiniteField(big.NewInt(0), *prime),
			},
			want: *NewFiniteField(big.NewInt(0), *prime),
		},
		{
			name: "1 - 0 = 1",
			args: args{
				a: NewFiniteField(big.NewInt(1), *prime),
				b: NewFiniteField(big.NewInt(0), *prime),
			},
			want: *NewFiniteField(big.NewInt(1), *prime),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFiniteField(big.NewInt(0), *prime).Sub(tt.args.a, tt.args.b); !got.Equals(&tt.want) {
				t.Errorf("%v : Sub() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_FiniteField_Mul(t *testing.T) {
	type args struct {
		a *FiniteField
		b *FiniteField
	}

	two := big.NewInt(2)
	prime := big.NewInt(0).Exp(two, big.NewInt(256), nil)
	prime.Sub(prime, big.NewInt(0).Exp(two, big.NewInt(32), nil))
	prime.Sub(prime, big.NewInt(0).Exp(two, big.NewInt(9), nil))
	prime.Sub(prime, big.NewInt(0).Exp(two, big.NewInt(8), nil))
	prime.Sub(prime, big.NewInt(0).Exp(two, big.NewInt(7), nil))
	prime.Sub(prime, big.NewInt(0).Exp(two, big.NewInt(6), nil))
	prime.Sub(prime, big.NewInt(0).Exp(two, big.NewInt(4), nil))
	prime.Sub(prime, big.NewInt(1))

	tests := []struct {
		name string
		args args
		want FiniteField
	}{
		{
			name: "0 * 0 = 0",
			args: args{
				a: NewFiniteField(big.NewInt(0), *prime),
				b: NewFiniteField(big.NewInt(0), *prime),
			},
			want: *NewFiniteField(big.NewInt(0), *prime),
		},
		{
			name: "1 * 0 = 1",
			args: args{
				a: NewFiniteField(big.NewInt(1), *prime),
				b: NewFiniteField(big.NewInt(0), *prime),
			},
			want: *NewFiniteField(big.NewInt(0), *prime),
		},
		{
			name: "0 * 1 = 0",
			args: args{
				a: NewFiniteField(big.NewInt(0), *prime),
				b: NewFiniteField(big.NewInt(1), *prime),
			},
			want: *NewFiniteField(big.NewInt(0), *prime),
		},
		{
			name: "2 * 1 = 2",
			args: args{
				a: NewFiniteField(big.NewInt(2), *prime),
				b: NewFiniteField(big.NewInt(1), *prime),
			},
			want: *NewFiniteField(big.NewInt(2), *prime),
		},
		{
			name: "prime * 2 = 0",
			args: args{
				a: NewFiniteField(prime, *prime),
				b: NewFiniteField(big.NewInt(2), *prime),
			},
			want: *NewFiniteField(big.NewInt(0), *prime),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFiniteField(big.NewInt(0), *prime).Mul(tt.args.a, tt.args.b); !got.Equals(&tt.want) {
				t.Errorf("%v : Sub() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_FiniteField_Div(t *testing.T) {
	type args struct {
		a *FiniteField
		b *FiniteField
	}

	two := big.NewInt(2)
	prime := big.NewInt(0).Exp(two, big.NewInt(256), nil)
	prime.Sub(prime, big.NewInt(0).Exp(two, big.NewInt(32), nil))
	prime.Sub(prime, big.NewInt(0).Exp(two, big.NewInt(9), nil))
	prime.Sub(prime, big.NewInt(0).Exp(two, big.NewInt(8), nil))
	prime.Sub(prime, big.NewInt(0).Exp(two, big.NewInt(7), nil))
	prime.Sub(prime, big.NewInt(0).Exp(two, big.NewInt(6), nil))
	prime.Sub(prime, big.NewInt(0).Exp(two, big.NewInt(4), nil))
	prime.Sub(prime, big.NewInt(1))

	inv2, _ := new(big.Int).SetString("57896044618658097711785492504343953926634992332820282019728792003954417335832", 10)

	tests := []struct {
		name string
		args args
		want FiniteField
	}{
		{
			name: "0 / 1 = 0",
			args: args{
				a: NewFiniteField(big.NewInt(0), *prime),
				b: NewFiniteField(big.NewInt(1), *prime),
			},
			want: *NewFiniteField(big.NewInt(0), *prime),
		},
		{
			name: "1 / 1 = 1",
			args: args{
				a: NewFiniteField(big.NewInt(1), *prime),
				b: NewFiniteField(big.NewInt(1), *prime),
			},
			want: *NewFiniteField(big.NewInt(1), *prime),
		},
		{
			name: "1 / 2 = 57896044618658097711785492504343953926634992332820282019728792003954417335832",
			args: args{
				a: NewFiniteField(big.NewInt(1), *prime),
				b: NewFiniteField(big.NewInt(2), *prime),
			},
			want: *NewFiniteField(inv2, *prime),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFiniteField(big.NewInt(0), *prime).Div(tt.args.a, tt.args.b); !got.Equals(&tt.want) {
				t.Errorf("%v : Sub() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
