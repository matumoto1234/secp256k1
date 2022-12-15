package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
	"reflect"

	"github.com/matumoto1234/secp256k1/models"
)

func messageToFiniteField(message string, prime big.Int) *models.FiniteField {
	hash := sha256.Sum256([]byte(message))
	value := big.NewInt(0).SetBytes(hash[:])
	return models.NewFiniteField(value, prime)
}

func newRandomFiniteField(prime big.Int) *models.FiniteField {
	n, err := rand.Int(rand.Reader, &prime)
	if err != nil {
		log.Fatal(err)
	}
	return models.NewFiniteField(n, prime)
}

type Signature struct {
	r models.FiniteField
	t models.FiniteField
}

func main() {
	// 鍵生成
	fmt.Println("generate key...")

	x, _ := big.NewInt(0).SetString("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", 16)
	y, _ := big.NewInt(0).SetString("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8", 16)

	// 生成点G
	G := models.NewEllipticCurveWithSECP256K1(
		*models.NewFiniteFieldWithSECP256K1(*x),
		*models.NewFiniteFieldWithSECP256K1(*y),
	)

	n, _ := new(big.Int).SetString("83ecb3984a4f9ff03e84d5f9c0d7f888a81833643047acc58eb6431e01d9bac8", 16)

	// 秘密鍵
	sec := models.NewFiniteField(n, *G.Order)
	fmt.Println("sec:", sec)
	fmt.Println()

	// 公開鍵
	pub := models.NewEmptyEllipticCurve(*G.Order).MulByScalar(*G, *sec)

	// sign
	fmt.Println("sign...")

	// t = (h(m)+r*s)/k
	z := messageToFiniteField("hello", *G.Order)
	fmt.Println("z:", z)

	k := newRandomFiniteField(*G.Order)
	fmt.Println("k:", k)

	Q := models.NewEmptyEllipticCurve(*G.Order).MulByScalar(*G, *k)
	fmt.Println("Q:", Q)
	r := models.NewFiniteField(Q.X.Value, *G.A.Prime)
	rs := models.NewFiniteField(big.NewInt(0), *G.Order).Mul(r, sec)
	t := models.NewFiniteField(big.NewInt(0), *G.Order).Add(z, rs)

	t.Div(t, k)

	sign := Signature{}
	sign.r = *r
	sign.t = *t
	fmt.Println("sign.r:", sign.r)
	fmt.Println("sign.t:", sign.t)

	// verify
	fmt.Println()
	fmt.Println("verify...")
	fmt.Println()

	z = messageToFiniteField("hello", *G.Order)
	w := models.NewFiniteField(big.NewInt(0), *G.Order).Div(models.NewFiniteField(big.NewInt(1), *G.Order), &sign.t)
	u1 := models.NewFiniteField(big.NewInt(0), *G.Order).Mul(z, w)

	fmt.Println("z:", z)
	fmt.Println("w:", w)
	u2 := models.NewFiniteField(big.NewInt(0), *G.Order).Mul(&sign.r, w)
	Pu1 := models.NewEmptyEllipticCurve(*G.Order).MulByScalar(*G, *u1)
	Su2 := models.NewEmptyEllipticCurve(*G.Order).MulByScalar(pub, *u2)
	fmt.Println("u1:", u1)
	fmt.Println("u2:", u2)
	fmt.Println("Pu1:", Pu1)
	fmt.Println("Su2:", Su2)

	fmt.Println(reflect.DeepEqual(Pu1, Su2))

	Q = models.NewEmptyEllipticCurve(*G.Order).Add(Pu1, Su2)
	fmt.Println()
	if Q.IsZero {
		fmt.Println("Q is zero")
	} else {
		fmt.Println("r:", sign.r)
		fmt.Println("x:", Q.X)
		fmt.Println("verify:", sign.r.Equals(Q.X))
	}
}
