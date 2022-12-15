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

	prime, _ := new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)

	a := models.NewFiniteField(big.NewInt(0), *prime)
	b := models.NewFiniteField(big.NewInt(7), *prime)
	order, _ := new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)

	secp256k1 := models.NewEllipticCurve(
		a,
		b,
		*order,
	)

	x, _ := new(big.Int).SetString("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", 16)
	y, _ := new(big.Int).SetString("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8", 16)

	// 生成点G
	G := models.NewEllipticCurvePoint(
		models.NewFiniteField(x, *prime),
		models.NewFiniteField(y, *prime),
		false,
		secp256k1,
	)

	// 秘密鍵
	n, _ := new(big.Int).SetString("83ecb3984a4f9ff03e84d5f9c0d7f888a81833643047acc58eb6431e01d9bac8", 16)
	sec := models.NewFiniteField(n, *secp256k1.Order)
	fmt.Println("sec:", sec)
	fmt.Println()

	// 公開鍵
	pub := new(models.EllipticCurvePoint).MulByScalar(G, sec, secp256k1)

	// sign
	fmt.Println("sign...")

	// t = (h(m)+r*sec)/k
	z := messageToFiniteField("hello", *secp256k1.Order)
	fmt.Println("z:", z)

	k := newRandomFiniteField(*secp256k1.Order)
	fmt.Println("k:", k)

	Q := new(models.EllipticCurvePoint).MulByScalar(G, k, secp256k1)
	fmt.Println("Q:", Q)

	r := models.NewFiniteField(Q.X.Value, *secp256k1.Order)
	rs := new(models.FiniteField).Mul(r, sec)
	t := new(models.FiniteField).Add(z, rs)

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

	z = messageToFiniteField("hello", *secp256k1.Order)

	// w = 1 / t
	w := new(models.FiniteField).Div(models.NewFiniteField(big.NewInt(1), *secp256k1.Order), &sign.t)

	u1 := new(models.FiniteField).Mul(z, w)

	fmt.Println("z:", z)
	fmt.Println("w:", w)
	u2 := new(models.FiniteField).Mul(&sign.r, w)
	Pu1 := new(models.EllipticCurvePoint).MulByScalar(G, u1, secp256k1)
	Su2 := new(models.EllipticCurvePoint).MulByScalar(pub, u2, secp256k1)
	fmt.Println("u1:", u1)
	fmt.Println("u2:", u2)
	fmt.Println("Pu1:", Pu1)
	fmt.Println("Su2:", Su2)

	fmt.Println(reflect.DeepEqual(Pu1, Su2))

	Q = new(models.EllipticCurvePoint).Add(Pu1, Su2, secp256k1)
	fmt.Println()
	if Q.IsZero {
		fmt.Println("Q is zero")
	} else {
		fmt.Println("r:", sign.r)
		fmt.Println("x:", Q.X)
		fmt.Println("verify:", sign.r.Equals(Q.X))
	}
}
