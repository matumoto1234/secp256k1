package models

import (
	"crypto/elliptic"
	"math/big"
	"math/bits"
)

type EllipticCurvePoint struct {
	X      *FiniteField
	Y      *FiniteField
	IsZero bool
}

// DeepCopy() : ppの値をpに代入する
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

func (p *EllipticCurvePoint) String() string {
	if p.IsZero {
		return "(zero,zero)"
	}
	return "(" + p.X.String() + "," + p.Y.String() + ")"
}

type EllipticCurve struct {
	A       *FiniteField
	B       *FiniteField
	Prime   *big.Int
	G       *EllipticCurvePoint
	BitSize int
	Name    string
	Order   *big.Int // 位数
}

func (ec *EllipticCurve) Params() *elliptic.CurveParams {
	return &elliptic.CurveParams{
		P:       ec.Prime,
		N:       ec.Order,
		B:       ec.B.Value,
		Gx:      ec.G.X.Value,
		Gy:      ec.G.Y.Value,
		BitSize: ec.BitSize,
		Name:    ec.Name,
	}
}

func (ec *EllipticCurve) IsOnCurve(p *EllipticCurvePoint) bool {
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

func (ec *EllipticCurve) Add(p1, p2 *EllipticCurvePoint) *EllipticCurvePoint {
	p3 := new(EllipticCurvePoint)

	if p1.IsZero {
		p3.DeepCopy(p2)
		return p3
	}
	if p2.IsZero {
		p3.DeepCopy(p1)
		return p3
	}

	x1 := p1.X
	y1 := p1.Y
	x2 := p2.X
	y2 := p2.Y

	L := NewFiniteField(big.NewInt(0), *ec.Prime)

	if x1.Equals(x2) {
		// P + (-P) = 0
		if y1.Equals(new(FiniteField).Neg(y2)) {
			p3.IsZero = true
			return p3
		}

		// L = (3 * x1^2 + a) / (2 * y1)
		x1Square := new(FiniteField).Mul(x1, x1)
		L.Add(x1Square, x1Square)
		L.Add(L, x1Square)
		L.Add(L, ec.A)
		L.Div(L, new(FiniteField).Add(y1, y1))
	} else {
		// L = (y2 - y1) / (x2 - x1)
		L.Sub(y2, y1)
		L.Div(L, new(FiniteField).Sub(x2, x1))
	}

	// x3 = L^2 - x1 - x2
	x3 := new(FiniteField).Mul(L, L)
	x3.Sub(x3, x1)
	x3.Sub(x3, x2)

	// y3 = L * (x1 - x3) - y1
	y3 := new(FiniteField).Sub(x1, x3)
	y3.Mul(y3, L)
	y3.Sub(y3, y1)

	p3.X = x3
	p3.Y = y3
	p3.IsZero = false
	return p3
}

func (ec *EllipticCurve) Double(p *EllipticCurvePoint) *EllipticCurvePoint {
	return ec.Add(p, p)
}

// e.g. k = [10001000, 10001111] -> 0b1000100010001111
func (ec *EllipticCurve) ScalarMult(p *EllipticCurvePoint, k []byte) *EllipticCurvePoint {
	pk := new(EllipticCurvePoint)

	if len(k) == 0 { // k == 0
		pk.IsZero = true
		return pk
	}

	if p.IsZero {
		pk.DeepCopy(p)
		return pk
	}

	// 繰り返し2乗法の応用
	// P + P = 2P
	// 2P + 2P = 4P
	// 4P + 4P = 8P
	// ... を用いて、k倍したP を Θ(log n) で求める
	sum := NewEllipticCurvePoint(nil, nil, true)
	for _, b := range k {
		rb := bits.Reverse8(b)
		for i := 0; i < 8; i++ {
			sum = ec.Add(sum, sum)
			if rb&byte(1) == 1 {
				sum = ec.Add(sum, p)
			}
			rb >>= 1
		}
	}

	pk.DeepCopy(sum)
	pk.IsZero = false
	return pk
}

func (ec *EllipticCurve) ScalarBaseMult(k []byte) *EllipticCurvePoint {
	return ec.ScalarMult(ec.G, k)
}

func NewEllipticCurvePoint(x, y *FiniteField, isZero bool) *EllipticCurvePoint {
	p := &EllipticCurvePoint{
		X:      x,
		Y:      y,
		IsZero: isZero,
	}

	return p
}

// if !ec.IsOnCurve(p) {
// 	panic(fmt.Sprintf("NewEllipticCurvePoint() : Invalid elliptic curve point : %v on elliptic curve : %v", p, ec))
// }


func NewEllipticCurve(a, b *FiniteField, prime *big.Int, G *EllipticCurvePoint, bitSize int, name string, order *big.Int) *EllipticCurve {
	return &EllipticCurve{
		A:     a,
		B:     b,
		Prime: prime,
		G: G,
		BitSize: 0,
		Name:    "",
		Order:   order,
	}
}
