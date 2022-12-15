package models

import (
	"fmt"
	"math/big"
)

type EllipticCurvePoint struct {
	X      *FiniteField
	Y      *FiniteField
	IsZero bool
}

type EllipticCurve struct {
	A     *FiniteField
	B     *FiniteField
	Order *big.Int // 位数
}

// DeepCopy() : ppの値を代入する
func (p *EllipticCurvePoint) DeepCopy(pp *EllipticCurvePoint) *EllipticCurvePoint {
	var cloneX, cloneY FiniteField
	if pp.X != nil {
		cloneX = *pp.X
	}
	if pp.Y != nil {
		cloneY = *pp.Y
	}
	p.X = &cloneX
	p.Y = &cloneY
	p.IsZero = pp.IsZero
	return p
}

// Equals() : 2つの楕円曲線上の点が等しいかどうかを判定する
func (p *EllipticCurvePoint) Equals(pp *EllipticCurvePoint) bool {
	if p.IsZero && pp.IsZero {
		return true
	}
	if p.IsZero != pp.IsZero {
		return false
	}
	return p.X.Equals(pp.X) && p.Y.Equals(pp.Y)
}

// Add() : sets z to the sum x+y on ec and returns z
func (z *EllipticCurvePoint) Add(x, y *EllipticCurvePoint, ec *EllipticCurve) *EllipticCurvePoint {
	if x.IsZero {
		z.DeepCopy(y)
		return z
	}
	if y.IsZero {
		z.DeepCopy(x)
		return z
	}

	x1 := x.X
	y1 := x.Y
	x2 := y.X
	y2 := y.Y

	L := NewFiniteField(big.NewInt(0), *x1.Prime)

	if x1.Equals(x2) {
		// P + (-P) = 0
		if y1.Equals(NewFiniteField(y2.Value, *y2.Prime).Neg(y2)) {
			z.IsZero = true
			return z
		}

		// L = (3 * x1^2 + a) / (2 * y1)
		x1Square := NewFiniteField(big.NewInt(0), *x1.Prime).Mul(x1, x1)
		L.Add(x1Square, x1Square)
		L.Add(L, x1Square)
		L.Add(L, ec.A)
		L.Div(L, NewFiniteField(big.NewInt(0), *x1.Prime).Add(y1, y1))
	} else {
		// L = (y2 - y1) / (x2 - x1)
		L.Sub(y2, y1)
		L.Div(L, NewFiniteField(big.NewInt(0), *x1.Prime).Sub(x2, x1))
	}

	// x3 = L^2 - x1 - x2
	x3 := NewFiniteField(big.NewInt(0), *x1.Prime).Mul(L, L)
	x3.Sub(x3, x1)
	x3.Sub(x3, x2)

	// y3 = L * (x1 - x3) - y1
	y3 := NewFiniteField(big.NewInt(0), *x1.Prime).Sub(x1, x3)
	y3.Mul(y3, L)
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
		z.IsZero = true
		return z
	}

	if p.IsZero {
		z.DeepCopy(p)
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

	z.DeepCopy(sum)
	z.IsZero = false
	return z
}

// IsValid : Whether it satisfies y^2 = x^3 + ax + b
func (p *EllipticCurvePoint) isValid(ec *EllipticCurve) bool {
	if p.IsZero {
		return true
	}

	if p.X.Prime.Cmp(p.Y.Prime) != 0 {
		return false
	}

	prime := p.X.Prime

	// y * y == (x * x + a) * x + b
	lhs := NewFiniteField(big.NewInt(0), *prime).Mul(p.Y, p.Y)
	rhs := NewFiniteField(big.NewInt(0), *prime).Mul(p.X, p.X)
	rhs.Add(rhs, ec.A)
	rhs.Mul(rhs, p.X)
	rhs.Add(rhs, ec.B)

	return lhs.Equals(rhs)
}

func (p *EllipticCurvePoint) String() string {
	if p.IsZero {
		return "x:zero y:zero"
	}
	return "x:" + p.X.String() + " y:" + p.Y.String()
}

func NewEllipticCurvePoint(x, y *FiniteField, isZero bool, ec *EllipticCurve) *EllipticCurvePoint {
	ecp := &EllipticCurvePoint{
		X:      x,
		Y:      y,
		IsZero: isZero,
	}

	if !ecp.isValid(ec) {
		panic(fmt.Sprintf("NewEllipticCurvePoint() : Invalid elliptic curve point : %v on elliptic curve : %v", ecp, ec))
	}

	return ecp
}

func NewEllipticCurve(a, b *FiniteField, order big.Int) *EllipticCurve {
	return &EllipticCurve{
		A:     a,
		B:     b,
		Order: &order,
	}
}
