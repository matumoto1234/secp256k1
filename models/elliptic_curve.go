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

// deepCopy() : ppの値をpに代入する
func (p *EllipticCurvePoint) deepCopy(pp *EllipticCurvePoint) *EllipticCurvePoint {
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

// equals() : 2つの楕円曲線上の点が等しいかどうかを判定する
func (p *EllipticCurvePoint) equals(pp *EllipticCurvePoint) bool {
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

func ToEllipticCurvePoint(x, y, prime *big.Int) *EllipticCurvePoint {
	// (0, 0) is zero
	if x.Sign() == 0 && y.Sign() == 0 {
		return NewEllipticCurvePoint(nil, nil, true)
	}

	return NewEllipticCurvePoint(
		NewFiniteField(x, *prime),
		NewFiniteField(y, *prime),
		false,
	)
}

func NewEllipticCurvePoint(x, y *FiniteField, isZero bool) *EllipticCurvePoint {
	if !isZero && x.Prime.Cmp(y.Prime) != 0 {
		panic("NewEllipticCurvePoint(): the primes of x and y are not same")
	}

	return &EllipticCurvePoint{
		X:      x,
		Y:      y,
		IsZero: isZero,
	}
}

type EllipticCurve struct {
	a       *FiniteField
	b       *FiniteField
	prime   *big.Int
	g       *EllipticCurvePoint
	bigSize int
	name    string
	order   *big.Int // 位数
}

func panicIfNotOnCurveP(ec *EllipticCurve, p *EllipticCurvePoint) {
	if p.IsZero {
		return
	}

	if !ec.IsOnCurveP(p) {
		panic("attempted operation on invalid point")
	}
}

func (ec *EllipticCurve) Params() *elliptic.CurveParams {
	return &elliptic.CurveParams{
		P:       ec.prime,
		N:       ec.order,
		B:       ec.b.Value,
		Gx:      ec.g.X.Value,
		Gy:      ec.g.Y.Value,
		BitSize: ec.bigSize,
		Name:    ec.name,
	}
}

func (ec *EllipticCurve) IsOnCurve(x, y *big.Int) bool {
	p := ToEllipticCurvePoint(x, y, ec.prime)
	return ec.IsOnCurveP(p)
}

func (ec *EllipticCurve) IsOnCurveP(p *EllipticCurvePoint) bool {
	// y * y == (x * x + a) * x + b
	lhs := new(FiniteField).Mul(p.Y, p.Y)
	rhs := new(FiniteField).Mul(p.X, p.X)
	rhs.Add(rhs, ec.a)
	rhs.Mul(rhs, p.X)
	rhs.Add(rhs, ec.b)

	return lhs.Equals(rhs)
}

func (ec *EllipticCurve) Add(x1, y1, x2, y2 *big.Int) (*big.Int, *big.Int) {
	p1 := ToEllipticCurvePoint(x1, y1, ec.prime)
	p2 := ToEllipticCurvePoint(x2, y2, ec.prime)
	result := ec.AddP(p1, p2)
	return result.X.Value, result.Y.Value
}

func (ec *EllipticCurve) AddP(p1, p2 *EllipticCurvePoint) *EllipticCurvePoint {
	panicIfNotOnCurveP(ec, p1)
	panicIfNotOnCurveP(ec, p2)

	if p1.IsZero {
		return new(EllipticCurvePoint).deepCopy(p2)
	}
	if p2.IsZero {
		return new(EllipticCurvePoint).deepCopy(p2)
	}

	x1 := p1.X
	y1 := p1.Y
	x2 := p2.X
	y2 := p2.Y

	L := NewFiniteField(big.NewInt(0), *ec.prime)

	if x1.Equals(x2) {
		// P + (-P) = 0
		if y1.Equals(new(FiniteField).Neg(y2)) {
			return NewEllipticCurvePoint(nil, nil, true)
		}

		// L = (3 * x1^2 + a) / (2 * y1)
		x1Square := new(FiniteField).Mul(x1, x1)
		L.Add(x1Square, x1Square)
		L.Add(L, x1Square)
		L.Add(L, ec.a)
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

	return NewEllipticCurvePoint(x3, y3, false)
}

func (ec *EllipticCurve) Double(x, y *big.Int) (*big.Int, *big.Int) {
	p := ToEllipticCurvePoint(x, y, ec.prime)
	result := ec.DoubleP(p)
	return result.X.Value, result.Y.Value
}

func (ec *EllipticCurve) DoubleP(p *EllipticCurvePoint) *EllipticCurvePoint {
	return ec.AddP(p, p)
}

func (ec *EllipticCurve) ScalarMult(x, y *big.Int, k []byte) (*big.Int, *big.Int) {
	p := ToEllipticCurvePoint(x, y, ec.prime)
	result := ec.ScalarMultP(p, k)
	return result.X.Value, result.Y.Value
}

// k is big-endian
func (ec *EllipticCurve) ScalarMultP(p *EllipticCurvePoint, k []byte) *EllipticCurvePoint {
	panicIfNotOnCurveP(ec, p)

	if len(k) == 0 { // k == 0
		return NewEllipticCurvePoint(nil, nil, true)
	}

	if p.IsZero {
		return new(EllipticCurvePoint).deepCopy(p)
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
			sum = ec.AddP(sum, sum)
			if rb&byte(1) == 1 {
				sum = ec.AddP(sum, p)
			}
			rb >>= 1
		}
	}

	return sum
}

func (ec *EllipticCurve) ScalarBaseMult(k []byte) (*big.Int, *big.Int) {
	result := ec.ScalarBaseMultP(k)
	return result.X.Value, result.Y.Value
}

func (ec *EllipticCurve) ScalarBaseMultP(k []byte) *EllipticCurvePoint {
	return ec.ScalarMultP(ec.g, k)
}

func NewEllipticCurve(a, b *FiniteField, prime *big.Int, G *EllipticCurvePoint, bitSize int, name string, order *big.Int) *EllipticCurve {
	return &EllipticCurve{
		a:       a,
		b:       b,
		prime:   prime,
		g:       G,
		bigSize: 0,
		name:    "",
		order:   order,
	}
}
