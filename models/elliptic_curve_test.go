package models

import (
	"math/big"
	"testing"
)

func Test_EllipticCurve_Add(t *testing.T) {
	type args struct {
		x  *EllipticCurvePoint
		y  *EllipticCurvePoint
		ec *EllipticCurve
	}

	prime := big.NewInt(223)

	a := NewFiniteField(big.NewInt(0), *prime)
	b := NewFiniteField(big.NewInt(7), *prime)
	order, _ := new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)

	ec := NewEllipticCurve(
		a,
		b,
		*order,
	)

	tests := []struct {
		name  string
		args  args
		order *big.Int
		want  *EllipticCurvePoint
	}{
		{
			name: "P + (-P) = 0",
			args: args{
				x: NewEllipticCurvePoint(
					NewFiniteField(big.NewInt(170), *prime),
					NewFiniteField(big.NewInt(142), *prime),
					false,
					ec,
				),
				y: NewEllipticCurvePoint(
					NewFiniteField(big.NewInt(170), *prime),
					NewFiniteField(big.NewInt(-142), *prime),
					false,
					ec,
				),
			},
			order: ec.Order,
			want:  NewEllipticCurvePoint(nil, nil, true, ec),
		},
		{
			name: "(170, 142) + (60, 139) = (220, 181)",
			args: args{
				x: NewEllipticCurvePoint(
					NewFiniteField(big.NewInt(170), *big.NewInt(223)),
					NewFiniteField(big.NewInt(142), *big.NewInt(223)),
					false,
					ec,
				),
				y: NewEllipticCurvePoint(
					NewFiniteField(big.NewInt(60), *big.NewInt(223)),
					NewFiniteField(big.NewInt(139), *big.NewInt(223)),
					false,
					ec,
				),
			},
			order: big.NewInt(223),
			want: NewEllipticCurvePoint(
				NewFiniteField(big.NewInt(220), *big.NewInt(223)),
				NewFiniteField(big.NewInt(181), *big.NewInt(223)),
				false,
				ec,
			),
		},
		{
			name: "(192, 105) * (192, 105) = (49, 71)",
			args: args{
				x: NewEllipticCurvePoint(
					NewFiniteField(big.NewInt(192), *big.NewInt(223)),
					NewFiniteField(big.NewInt(105), *big.NewInt(223)),
					false,
					ec,
				),
				y: NewEllipticCurvePoint(
					NewFiniteField(big.NewInt(192), *big.NewInt(223)),
					NewFiniteField(big.NewInt(105), *big.NewInt(223)),
					false,
					ec,
				),
			},
			order: big.NewInt(223),
			want: NewEllipticCurvePoint(
				NewFiniteField(big.NewInt(49), *big.NewInt(223)),
				NewFiniteField(big.NewInt(71), *big.NewInt(223)),
				false,
				ec,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := new(EllipticCurvePoint).Add(tt.args.x, tt.args.y, ec); !got.Equals(tt.want) {
				t.Errorf("%v : EllipticCurve.MulByScalar() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_EllipticCurve_MulByScalar(t *testing.T) {
	type args struct {
		p *EllipticCurvePoint
		x *FiniteField
	}

	prime := big.NewInt(223)

	a := NewFiniteField(big.NewInt(0), *prime)
	b := NewFiniteField(big.NewInt(7), *prime)
	order, _ := new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)

	ec := NewEllipticCurve(
		a,
		b,
		*order,
	)

	tests := []struct {
		name  string
		args  args
		order *big.Int
		want  *EllipticCurvePoint
	}{
		{
			name: "2 * (192, 105) = (49, 71)",
			args: args{
				p: NewEllipticCurvePoint(
					NewFiniteField(big.NewInt(192), *big.NewInt(223)),
					NewFiniteField(big.NewInt(105), *big.NewInt(223)),
					false,
					ec,
				),
				x: NewFiniteField(big.NewInt(2), *big.NewInt(223)),
			},
			order: big.NewInt(223),
			want: NewEllipticCurvePoint(
				NewFiniteField(big.NewInt(49), *big.NewInt(223)),
				NewFiniteField(big.NewInt(71), *big.NewInt(223)),
				false,
				ec,
			),
		},
		{
			name: "4 * (47, 71) = (194, 51)",
			args: args{
				p: NewEllipticCurvePoint(
					NewFiniteField(big.NewInt(47), *big.NewInt(223)),
					NewFiniteField(big.NewInt(71), *big.NewInt(223)),
					false,
					ec,
				),
				x: NewFiniteField(big.NewInt(4), *big.NewInt(223)),
			},
			order: big.NewInt(223),
			want: NewEllipticCurvePoint(
				NewFiniteField(big.NewInt(194), *big.NewInt(223)),
				NewFiniteField(big.NewInt(51), *big.NewInt(223)),
				false,
				ec,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := new(EllipticCurvePoint).MulByScalar(tt.args.p, tt.args.x, ec); !got.Equals(tt.want) {
				t.Errorf("%v : EllipticCurve.MulByScalar() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
