package models

import (
	"errors"
	"math/big"
)

type EllipticCurvePoint struct {
	X      *FiniteField
	Y      *FiniteField
	IsZero bool
}

type EllipticCurve struct {
	A      *FiniteField
	B      *FiniteField
	Order  *big.Int // 位数
}

// Clone() : 楕円曲線上の点を複製する
func (ecp *EllipticCurvePoint) Clone() *EllipticCurvePoint {
	cloneX := *ecp.X
	cloneY := *ecp.Y
	return &EllipticCurvePoint{
		X:      &cloneX,
		Y:      &cloneY,
		IsZero: ecp.IsZero,
	}
}

// Equals() : 2つの楕円曲線上の点が等しいかどうかを判定する
func (ecp *EllipticCurvePoint) Equals(p *EllipticCurvePoint) bool {
	return ecp.X.Equals(p.X) && ecp.Y.Equals(p.Y) && ecp.IsZero == p.IsZero
}

// Add() : sets z to the sum x+y on ec and returns z
func (z *EllipticCurvePoint) Add(x, y *EllipticCurvePoint, ec *EllipticCurve) *EllipticCurvePoint {
	if x.IsZero {
		z = y.Clone()
		return z
	}
	if y.IsZero {
		z = x.Clone()
		return z
	}

	x1 := x.X
	y1 := x.Y
	x2 := y.X
	y2 := y.Y

	L := *NewFiniteField(big.NewInt(0), *x1.Prime)

	if x1.Equals(x2) {
		// P + (-P) = 0
		if y1.Equals(NewFiniteField(y2.Value, *y2.Prime).Neg(y2)) {
			z.X = NewFiniteField(big.NewInt(0), *x1.Prime)
			z.Y = NewFiniteField(big.NewInt(0), *x1.Prime)
			z.IsZero = true
			return z
		}

		// L = (3 * x1^2 + a) / (2 * y1)
		x1Square := NewFiniteField(big.NewInt(0), *x1.Prime).Mul(x1, x1)
		L.Add(x1Square, x1Square)
		L.Add(&L, x1Square)
		L.Add(&L, ec.A)
		L.Div(&L, NewFiniteField(big.NewInt(0), *x1.Prime).Add(y1, y1))
	} else {
		// L = (y2 - y1) / (x2 - x1)
		L.Sub(y2, y1)
		L.Div(&L, NewFiniteField(big.NewInt(0), *x1.Prime).Sub(x2, x1))
	}

	// x3 = L^2 - x1 - x2
	x3 := NewFiniteField(big.NewInt(0), *x1.Prime).Mul(&L, &L)
	x3.Sub(x3, x1)
	x3.Sub(x3, x2)

	// y3 = L * (x1 - x3) - y1
	y3 := NewFiniteField(big.NewInt(0), *x1.Prime).Sub(x1, x3)
	y3.Mul(y3, &L)
	y3.Sub(y3, y1)

	z.X = x3
	z.Y = y3
	z.IsZero = false
	return z
}

// MulByScalar() : sets z to the p*x on ec and returns z
func (z *EllipticCurvePoint) MulByScalar(p *EllipticCurvePoint, x *FiniteField, ec *EllipticCurve) *EllipticCurvePoint {
	// x == 0
	if x.Value.Cmp(big.NewInt(0)) == 0 {
		z.X = NewFiniteField(big.NewInt(0), *x.Prime)
		z.Y = NewFiniteField(big.NewInt(0), *x.Prime)
		z.IsZero = true
		return z
	}

	if p.IsZero {
		z = p.Clone()
		return z
	}

	sum := new(EllipticCurvePoint)
	sum.IsZero = true
	for i := x.Value.BitLen() - 1; i >= 0; i-- {
		sum.Add(sum, sum, ec)
		if x.Value.Bit(i) == 1 {
			sum.Add(sum, p, ec)
		}
	}

	z = sum.Clone()
	z.IsZero = false
	return z
}

// IsValid : Whether it satisfies y^2 = x^3 + ax + b
func (ecp *EllipticCurvePoint) IsValid(ec *EllipticCurve) bool {
	if ecp.IsZero {
		return true
	}

	if ecp.X.Prime.Cmp(ecp.Y.Prime) != 0 {
		return false
	}

	p := ecp.X.Prime

	// y * y == (x * x + a) * x + b
	lhs := NewFiniteField(big.NewInt(0), *p).Mul(ecp.Y, ecp.Y)
	rhs := NewFiniteField(big.NewInt(0), *p).Mul(ecp.X, ecp.X)
	rhs.Add(rhs, ec.A)
	rhs.Mul(rhs, ecp.X)
	rhs.Add(rhs, ec.B)

	return lhs.Equals(rhs)
}

func (ecp *EllipticCurvePoint) String() string {
	return "x:" + ecp.X.String() + " y:" + ecp.Y.String()
}

func NewEllipticCurvePoint(x, y *FiniteField, isZero bool, ec *EllipticCurve) (*EllipticCurvePoint, error) {
	ecp := &EllipticCurvePoint{
		X:      x,
		Y:      y,
		IsZero: isZero,
	}

	if !ecp.IsValid(ec) {
		return nil, errors.New("Invalid EllipticCurvePoint")
	}

	return ecp, nil
}

func NewEllipticCurve(a, b *FiniteField, order big.Int) *EllipticCurve {
	return &EllipticCurve{
		A:      a,
		B:      b,
		Order:  &order,
	}
}
