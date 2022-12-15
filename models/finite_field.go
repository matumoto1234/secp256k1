package models

import (
	"math/big"
	"sync"
)

type FiniteFieldPrime big.Int

var (
	once sync.Once
	prime *big.Int
)

func (f FiniteFieldPrime) Raw() *big.Int {
	return prime
}

func NewFiniteFieldPrime() FiniteFieldPrime {
	once.Do(func(){
		prime = big.NewInt(0)
	})

	return FiniteFieldPrime(*prime)
}

type FiniteField struct {
	Value *big.Int
	Prime *big.Int
}

// Neg() : sets z to -x and returns z
func (f *FiniteField) Neg(x *FiniteField) *FiniteField {
	f.Value.Neg(x.Value)
	f.Value.Mod(f.Value, f.Prime)
	return f
}

// Add() : sets z to the sum x+y and returns z
func (f *FiniteField) Add(x, y *FiniteField) *FiniteField {
	f.Value.Add(x.Value, y.Value)
	f.Value.Mod(f.Value, f.Prime)
	return f
}

// Sub() : sets z to the difference x-y and returns z
func (f *FiniteField) Sub(x, y *FiniteField) *FiniteField {
	f.Value.Sub(x.Value, y.Value)
	f.Value.Mod(f.Value, f.Prime)
	return f
}

// Mul() : sets z to the product x*y and returns z
func (f *FiniteField) Mul(x, y *FiniteField) *FiniteField {
	f.Value.Mul(x.Value, y.Value)
	f.Value.Mod(f.Value, f.Prime)
	return f
}

// Div() : sets z to the quotient x/y and returns z
func (f *FiniteField) Div(x, y *FiniteField) *FiniteField {
	inv := new(big.Int).ModInverse(y.Value, y.Prime)
	f.Value.Mul(x.Value, inv)
	f.Value.Mod(f.Value, f.Prime)
	return f
}

func (f FiniteField) String() string {
	return "value:" + f.Value.String()
}

func (f FiniteField) Equals(x *FiniteField) bool {
	return f.Value.Cmp(x.Value) == 0
}

// NewFiniteField() : constructor of FiniteField
func NewFiniteField(value *big.Int, prime big.Int) *FiniteField {
	return &FiniteField{
		Value: new(big.Int).Mod(value, &prime),
		Prime: &prime,
	}
}

// NewFiniteFieldWithSECP256K1() : constructor of FiniteField with secp256k1 parameter
func NewFiniteFieldWithSECP256K1(value big.Int) *FiniteField {
	two := big.NewInt(2)
	prime := new(big.Int).Exp(two, big.NewInt(256), nil)
	prime.Sub(prime, new(big.Int).Exp(two, big.NewInt(32), nil))
	prime.Sub(prime, new(big.Int).Exp(two, big.NewInt(9), nil))
	prime.Sub(prime, new(big.Int).Exp(two, big.NewInt(8), nil))
	prime.Sub(prime, new(big.Int).Exp(two, big.NewInt(7), nil))
	prime.Sub(prime, new(big.Int).Exp(two, big.NewInt(6), nil))
	prime.Sub(prime, new(big.Int).Exp(two, big.NewInt(4), nil))
	prime.Sub(prime, big.NewInt(1))

	return NewFiniteField(&value, *prime)
}
