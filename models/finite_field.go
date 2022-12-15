package models

import (
	"fmt"
	"math/big"
)

type FiniteField struct {
	Value *big.Int
	Prime *big.Int
}

// Neg() : sets z to -x and returns z
func (f *FiniteField) Neg(x *FiniteField) *FiniteField {
	if f.Value == nil {
		f.Value = new(big.Int).Neg(x.Value)
	} else {
		f.Value.Neg(x.Value)
	}

	f.Value.Mod(f.Value, x.Prime)
	f.Prime = x.Prime
	return f
}

// Add() : sets z to the sum x+y and returns z
func (f *FiniteField) Add(x, y *FiniteField) *FiniteField {
	if x.Prime.Cmp(y.Prime) != 0 {
		panic(fmt.Sprintf("Add() : the primes of x and y are not same. x prime: %v y prime: %v", x.Prime, y.Prime))
	}

	if f.Value == nil {
		f.Value = new(big.Int).Add(x.Value, y.Value)
	} else {
		f.Value.Add(x.Value, y.Value)
	}

	f.Value.Mod(f.Value, x.Prime)
	f.Prime = x.Prime
	return f
}

// Sub() : sets z to the difference x-y and returns z
func (f *FiniteField) Sub(x, y *FiniteField) *FiniteField {
	if x.Prime.Cmp(y.Prime) != 0 {
		panic(fmt.Sprintf("Sub() : the primes of x and y are not same. x prime: %v y prime: %v", x.Prime, y.Prime))
	}

	if f.Value == nil {
		f.Value = new(big.Int).Sub(x.Value, y.Value)
	} else {
		f.Value.Sub(x.Value, y.Value)
	}

	f.Value.Mod(f.Value, x.Prime)
	f.Prime = x.Prime
	return f
}

// Mul() : sets z to the product x*y and returns z
func (f *FiniteField) Mul(x, y *FiniteField) *FiniteField {
	if x.Prime.Cmp(y.Prime) != 0 {
		panic(fmt.Sprintf("Mul() : the primes of x and y are not same. x prime: %v y prime: %v", x.Prime, y.Prime))
	}

	if f.Value == nil {
		f.Value = new(big.Int).Mul(x.Value, y.Value)
	} else {
		f.Value.Mul(x.Value, y.Value)
	}

	f.Value.Mod(f.Value, x.Prime)
	f.Prime = x.Prime
	return f
}

// Div() : sets z to the quotient x/y and returns z
func (f *FiniteField) Div(x, y *FiniteField) *FiniteField {
	if x.Prime.Cmp(y.Prime) != 0 {
		panic(fmt.Sprintf("Div() : the primes of x and y are not same. x prime: %v y prime: %v", x.Prime, y.Prime))
	}

	inv := new(big.Int).ModInverse(y.Value, y.Prime)

	if f.Value == nil {
		f.Value = new(big.Int).Mul(x.Value, inv)
	} else {
		f.Value.Mul(x.Value, inv)
	}

	f.Value.Mod(f.Value, x.Prime)
	f.Prime = x.Prime
	return f
}

func (f FiniteField) String() string {
	return f.Value.String()
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
