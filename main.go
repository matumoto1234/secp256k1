package main

import (
	"crypto/elliptic"
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

func generateSecp256k1() elliptic.Curve {
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

func generateKey(ec elliptic.Curve) (*models.FiniteField, *models.EllipticCurvePoint) {
	// 鍵生成
	fmt.Println("generate key...")
	// 秘密鍵
	n, _ := new(big.Int).SetString("83ecb3984a4f9ff03e84d5f9c0d7f888a81833643047acc58eb6431e01d9bac8", 16)
	priv := models.NewFiniteField(n, *ec.Params().N)
	fmt.Println("priv:", priv)
	fmt.Println()

	// 公開鍵
	x, y := ec.ScalarBaseMult(priv.Value.Bytes())
	pub := models.ToEllipticCurvePoint(x, y, ec.Params().P)

	return priv, pub
}

type signature struct {
	r *models.FiniteField
	t *models.FiniteField
}

func sign(ec elliptic.Curve, msg string, priv *models.FiniteField) signature {
	// sign
	fmt.Println("sign...")

	// t = (h(m)+r*priv)/k
	z := messageToFiniteField(msg, *ec.Params().N)
	fmt.Println("z:", z)

	k := newRandomFiniteField(*ec.Params().N)
	fmt.Println("k:", k)

	x, y := ec.ScalarBaseMult(k.Value.Bytes())
	Q := models.ToEllipticCurvePoint(x, y, ec.Params().P)
	fmt.Println("Q:", Q)

	r := models.NewFiniteField(Q.X.Value, *ec.Params().N)
	rs := new(models.FiniteField).Mul(r, priv)
	t := new(models.FiniteField).Add(z, rs)

	t.Div(t, k)

	sign := signature{}
	sign.r = r
	sign.t = t
	fmt.Println("sign.r:", sign.r)
	fmt.Println("sign.t:", sign.t)
	return sign
}

func verify(ec elliptic.Curve, msg string, sign signature, pub *models.EllipticCurvePoint) bool {
	// verify
	fmt.Println()
	fmt.Println("verify...")
	fmt.Println()

	z := messageToFiniteField(msg, *ec.Params().N)

	// w = 1 / t
	w := new(models.FiniteField).Div(models.NewFiniteField(big.NewInt(1), *ec.Params().N), sign.t)

	u1 := new(models.FiniteField).Mul(z, w)

	fmt.Println("z:", z)
	fmt.Println("w:", w)
	u2 := new(models.FiniteField).Mul(sign.r, w)

	x, y := ec.ScalarBaseMult(u1.Value.Bytes())
	Pu1 := models.ToEllipticCurvePoint(x, y, ec.Params().P)
	x, y = ec.ScalarMult(pub.X.Value, pub.Y.Value, u2.Value.Bytes())
	Su2 := models.ToEllipticCurvePoint(x, y, ec.Params().P)
	fmt.Println("u1:", u1)
	fmt.Println("u2:", u2)
	fmt.Println("Pu1:", Pu1)
	fmt.Println("Su2:", Su2)

	x, y = ec.Add(Pu1.X.Value, Pu1.Y.Value, Su2.X.Value, Su2.Y.Value)
	Q := models.ToEllipticCurvePoint(x, y, ec.Params().P)
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
	// ECDSAを行っている
	secp256k1 := generateSecp256k1()
	priv, pub := generateKey(secp256k1)

	msg := "hello"
	signature := sign(secp256k1, msg, priv)
	fmt.Println("verify:", verify(secp256k1, msg, signature, pub))
}
