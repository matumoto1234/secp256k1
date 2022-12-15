package models

import (
	"math/big"
	"testing"
)

func Test_EllipticCurve_Add(t *testing.T) {
	type args struct {
		a EllipticCurve
		b EllipticCurve
	}

	secp256k1Order, _ := new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)

	tests := []struct {
		name  string
		args  args
		order *big.Int
		want  EllipticCurve
	}{
		{
			name: "P + (-P) = 0",
			args: args{
				a: *NewEllipticCurveWithSECP256K1(
					*NewFiniteFieldWithSECP256K1(*big.NewInt(2)),
					*NewFiniteFieldWithSECP256K1(*big.NewInt(3)),
				),
				b: *NewEllipticCurveWithSECP256K1(
					*NewFiniteFieldWithSECP256K1(*big.NewInt(2)),
					*NewFiniteFieldWithSECP256K1(*big.NewInt(-3)),
				),
			},
			order: secp256k1Order,
			want:  *NewEmptyEllipticCurve(*secp256k1Order),
		},
		{
			name: "(170, 142) + (60, 139) = (220, 181)",
			args: args{
				a: *NewEllipticCurve(
					*NewFiniteField(big.NewInt(170), *big.NewInt(223)),
					*NewFiniteField(big.NewInt(142), *big.NewInt(223)),
					*big.NewInt(223),
				),
				b: *NewEllipticCurve(
					*NewFiniteField(big.NewInt(60), *big.NewInt(223)),
					*NewFiniteField(big.NewInt(139), *big.NewInt(223)),
					*big.NewInt(223),
				),
			},
			order: big.NewInt(223),
			want: *NewEllipticCurve(
				*NewFiniteField(big.NewInt(220), *big.NewInt(223)),
				*NewFiniteField(big.NewInt(181), *big.NewInt(223)),
				*big.NewInt(223),
			),
		},
		{
			name: "(192, 105) * (192, 105) = (49, 71)",
			args: args{
				a: *NewEllipticCurve(
					*NewFiniteField(big.NewInt(192), *big.NewInt(223)),
					*NewFiniteField(big.NewInt(105), *big.NewInt(223)),
					*big.NewInt(223),
				),
				b: *NewEllipticCurve(
					*NewFiniteField(big.NewInt(192), *big.NewInt(223)),
					*NewFiniteField(big.NewInt(105), *big.NewInt(223)),
					*big.NewInt(223),
				),
			},
			order: big.NewInt(223),
			want: *NewEllipticCurve(
				*NewFiniteField(big.NewInt(49), *big.NewInt(223)),
				*NewFiniteField(big.NewInt(71), *big.NewInt(223)),
				*big.NewInt(223),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEmptyEllipticCurve(*tt.order)
			if got := e.Add(tt.args.a, tt.args.b); !got.Equals(tt.want) {
				t.Errorf("%v : EllipticCurve.MulByScalar() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_EllipticCurve_MulByScalar(t *testing.T) {
	type args struct {
		a EllipticCurve
		n FiniteField
	}

	tests := []struct {
		name  string
		args  args
		order *big.Int
		want  EllipticCurve
	}{
		{
			name: "2 * (192, 105) = (49, 71)",
			args: args{
				a: *NewEllipticCurve(
					*NewFiniteField(big.NewInt(192), *big.NewInt(223)),
					*NewFiniteField(big.NewInt(105), *big.NewInt(223)),
					*big.NewInt(223),
				),
				n: *NewFiniteField(big.NewInt(2), *big.NewInt(223)),
			},
			order: big.NewInt(223),
			want: *NewEllipticCurve(
				*NewFiniteField(big.NewInt(49), *big.NewInt(223)),
				*NewFiniteField(big.NewInt(71), *big.NewInt(223)),
				*big.NewInt(223),
			),
		},
		{
			name: "4 * (47, 71) = (194, 51)",
			args: args{
				a: *NewEllipticCurve(
					*NewFiniteField(big.NewInt(47), *big.NewInt(223)),
					*NewFiniteField(big.NewInt(71), *big.NewInt(223)),
					*big.NewInt(223),
				),
				n: *NewFiniteField(big.NewInt(4), *big.NewInt(223)),
			},
			order: big.NewInt(223),
			want: *NewEllipticCurve(
				*NewFiniteField(big.NewInt(194), *big.NewInt(223)),
				*NewFiniteField(big.NewInt(51), *big.NewInt(223)),
				*big.NewInt(223),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEmptyEllipticCurve(*tt.order)
			if got := e.MulByScalar(tt.args.a, tt.args.n); !got.Equals(tt.want) {
				t.Errorf("%v : EllipticCurve.MulByScalar() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
