package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"

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
	r *models.FiniteField
	t *models.FiniteField
}

func generateSecp256k1() *models.EllipticCurve {
	// 素数
	prime, _ := new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)

	// 楕円曲線のパラメータ
	a := models.NewFiniteField(big.NewInt(0), *prime)
	b := models.NewFiniteField(big.NewInt(7), *prime)

	// 位数
	order, _ := new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)

	// 生成点
	gx, _ := new(big.Int).SetString("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", 16)
	gy, _ := new(big.Int).SetString("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8", 16)
	G := models.NewEllipticCurvePoint(
		models.NewFiniteField(gx, *prime),
		models.NewFiniteField(gy, *prime),
		false,
	)

	return models.NewEllipticCurve(
		a,
		b,
		prime,
		G,
		256,
		"secp256k1",
		order,
	)
}

func generateKey(ec *models.EllipticCurve) (*models.FiniteField, *models.EllipticCurvePoint) {
	// 鍵生成
	fmt.Println("generate key...")
	// 秘密鍵
	n, _ := new(big.Int).SetString("83ecb3984a4f9ff03e84d5f9c0d7f888a81833643047acc58eb6431e01d9bac8", 16)
	sec := models.NewFiniteField(n, *ec.Order)
	fmt.Println("sec:", sec)
	fmt.Println()

	// 公開鍵
	pub := ec.ScalarBaseMult(sec.Value.Bytes())

	return sec, pub
}

func sign(ec *models.EllipticCurve, msg string, sec *models.FiniteField) Signature {
	// sign
	fmt.Println("sign...")

	// t = (h(m)+r*sec)/k
	z := messageToFiniteField(msg, *ec.Order)
	fmt.Println("z:", z)

	k := newRandomFiniteField(*ec.Order)
	fmt.Println("k:", k)

	Q := ec.ScalarBaseMult(k.Value.Bytes())
	fmt.Println("Q:", Q)

	r := models.NewFiniteField(Q.X.Value, *ec.Order)
	rs := new(models.FiniteField).Mul(r, sec)
	t := new(models.FiniteField).Add(z, rs)

	t.Div(t, k)

	sign := Signature{}
	sign.r = r
	sign.t = t
	fmt.Println("sign.r:", sign.r)
	fmt.Println("sign.t:", sign.t)
	return sign
}

func verify(ec *models.EllipticCurve, msg string, sign Signature, pub *models.EllipticCurvePoint) bool {
	// verify
	fmt.Println()
	fmt.Println("verify...")
	fmt.Println()

	z := messageToFiniteField(msg, *ec.Order)

	// w = 1 / t
	w := new(models.FiniteField).Div(models.NewFiniteField(big.NewInt(1), *ec.Order), sign.t)

	u1 := new(models.FiniteField).Mul(z, w)

	fmt.Println("z:", z)
	fmt.Println("w:", w)
	u2 := new(models.FiniteField).Mul(sign.r, w)
	Pu1 := ec.ScalarBaseMult(u1.Value.Bytes())
	Su2 := ec.ScalarMult(pub, u2.Value.Bytes())
	fmt.Println("u1:", u1)
	fmt.Println("u2:", u2)
	fmt.Println("Pu1:", Pu1)
	fmt.Println("Su2:", Su2)

	Q := ec.Add(Pu1, Su2)
	fmt.Println()
	if Q.IsZero {
		fmt.Println("Q is zero")
		return false
	} else {
		fmt.Println("r:", sign.r)
		fmt.Println("x:", Q.X)
		return sign.r.Equals(Q.X)
	}
}

func main() {
	secp256k1 := generateSecp256k1()
	sec, pub := generateKey(secp256k1)

	msg := "hello"
	signature := sign(secp256k1, msg, sec)
	fmt.Println("verify:", verify(secp256k1, msg, signature, pub))
}
