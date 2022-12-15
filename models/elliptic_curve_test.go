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

	a := NewFiniteField(big.NewInt(0), prime)
	b := NewFiniteField(big.NewInt(7), prime)

	ec := NewEllipticCurve(
		a,
		b,
		prime,
		nil,
		0,
		"test elliptic curve",
		nil,
	)

	tests := []struct {
		name string
		args args
		want *EllipticCurvePoint
	}{
		{
			name: "P + (-P) = 0",
			args: args{
				x: NewEllipticCurvePoint(
					NewFiniteField(big.NewInt(170), prime),
					NewFiniteField(big.NewInt(142), prime),
					false,
				),
				y: NewEllipticCurvePoint(
					NewFiniteField(big.NewInt(170), prime),
					NewFiniteField(big.NewInt(-142), prime),
					false,
				),
			},
			want: NewEllipticCurvePoint(nil, nil, true),
		},
		{
			name: "(170, 142) + (60, 139) = (220, 181)",
			args: args{
				x: NewEllipticCurvePoint(
					NewFiniteField(big.NewInt(170), big.NewInt(223)),
					NewFiniteField(big.NewInt(142), big.NewInt(223)),
					false,
				),
				y: NewEllipticCurvePoint(
					NewFiniteField(big.NewInt(60), big.NewInt(223)),
					NewFiniteField(big.NewInt(139), big.NewInt(223)),
					false,
				),
			},
			want: NewEllipticCurvePoint(
				NewFiniteField(big.NewInt(220), big.NewInt(223)),
				NewFiniteField(big.NewInt(181), big.NewInt(223)),
				false,
			),
		},
		{
			name: "(192, 105) * (192, 105) = (49, 71)",
			args: args{
				x: NewEllipticCurvePoint(
					NewFiniteField(big.NewInt(192), big.NewInt(223)),
					NewFiniteField(big.NewInt(105), big.NewInt(223)),
					false,
				),
				y: NewEllipticCurvePoint(
					NewFiniteField(big.NewInt(192), big.NewInt(223)),
					NewFiniteField(big.NewInt(105), big.NewInt(223)),
					false,
				),
			},
			want: NewEllipticCurvePoint(
				NewFiniteField(big.NewInt(49), big.NewInt(223)),
				NewFiniteField(big.NewInt(71), big.NewInt(223)),
				false,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ec.AddP(tt.args.x, tt.args.y); !got.equals(tt.want) {
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

	a := NewFiniteField(big.NewInt(0), prime)
	b := NewFiniteField(big.NewInt(7), prime)
	ec := NewEllipticCurve(
		a,
		b,
		prime,
		nil,
		0,
		"test elliptic curve",
		nil,
	)

	tests := []struct {
		name string
		args args
		want *EllipticCurvePoint
	}{
		{
			name: "2 * (192, 105) = (49, 71)",
			args: args{
				p: NewEllipticCurvePoint(
					NewFiniteField(big.NewInt(192), big.NewInt(223)),
					NewFiniteField(big.NewInt(105), big.NewInt(223)),
					false,
				),
				x: NewFiniteField(big.NewInt(2), big.NewInt(223)),
			},
			want: NewEllipticCurvePoint(
				NewFiniteField(big.NewInt(49), big.NewInt(223)),
				NewFiniteField(big.NewInt(71), big.NewInt(223)),
				false,
			),
		},
		{
			name: "4 * (47, 71) = (194, 51)",
			args: args{
				p: NewEllipticCurvePoint(
					NewFiniteField(big.NewInt(47), big.NewInt(223)),
					NewFiniteField(big.NewInt(71), big.NewInt(223)),
					false,
				),
				x: NewFiniteField(big.NewInt(4), big.NewInt(223)),
			},
			want: NewEllipticCurvePoint(
				NewFiniteField(big.NewInt(194), big.NewInt(223)),
				NewFiniteField(big.NewInt(51), big.NewInt(223)),
				false,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ec.ScalarMultP(tt.args.p, tt.args.x.Value.Bytes()); !got.equals(tt.want) {
				t.Errorf("%v : EllipticCurve.MulByScalar() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
