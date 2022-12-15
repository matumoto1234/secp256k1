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

func toHash(message string, prime big.Int) *models.FiniteField {
	hash := sha256.Sum256([]byte(message))
	value := new(big.Int).SetBytes(hash[:])
	return models.NewFiniteField(value, prime)
}

// generate random number in [1, prime)
func newRandomFiniteField(prime *big.Int) (*models.FiniteField, error) {
	for {
		n, err := rand.Int(rand.Reader, prime)
		if err != nil {
			return nil, err
		}
		if n.Sign() != 0 {
			return models.NewFiniteField(n, *prime), nil
		}
	}
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

// generateKey() : 秘密鍵privと公開鍵pubの鍵生成
// n := 適切に選んだ定数(厳密な秘密鍵)
// priv := 秘密鍵(有限体にしたもの)
// G := 生成点
// pub := k*G
func generateKey(ec elliptic.Curve) (*models.FiniteField, *models.EllipticCurvePoint) {
	// 秘密鍵
	n, _ := new(big.Int).SetString("83ecb3984a4f9ff03e84d5f9c0d7f888a81833643047acc58eb6431e01d9bac8", 16)
	priv := models.NewFiniteField(n, *ec.Params().N)

	// 公開鍵
	x, y := ec.ScalarBaseMult(priv.Value.Bytes())
	pub := models.ToEllipticCurvePoint(x, y, ec.Params().P)

	return priv, pub
}

type signature struct {
	r *models.FiniteField
	t *models.FiniteField
}

// sign() : signature(r, t)なる署名を生成
// G := 生成点
// k := 一時的な秘密鍵(ランダム)
// Q := 一時的な公開鍵(k*G)
// priv := 秘密鍵
// r := 公開鍵Qのx座標
// z := メッセージのハッシュ
// t := (z + s*r) / k を計算した値
func sign(ec elliptic.Curve, msg string, priv *models.FiniteField) (signature, error) {
	// temporary private key
	k, err := newRandomFiniteField(ec.Params().N)
	if err != nil {
		return signature{}, err
	}

	// temporary public key
	x, y := ec.ScalarBaseMult(k.Value.Bytes())
	Q := models.ToEllipticCurvePoint(x, y, ec.Params().N)

	var sign signature
	sign.r = Q.X

	sign.t = new(models.FiniteField).Mul(sign.r, priv)
	z := toHash(msg, *ec.Params().N)
	sign.t.Add(sign.t, z)
	sign.t.Div(sign.t, k)

	return sign, nil
}

// verify() : signature(r, t)の署名検証
// G := 生成点
// z := メッセージのハッシュ
// pub := 公開鍵
// R := (z*G + r*pub) / t
//
//	Rのx座標 == r -> OK
//	Rのx座標 != r -> NG(R == 無限遠点の場合もNG)
func verify(ec elliptic.Curve, msg string, sign signature, pub *models.EllipticCurvePoint) bool {
	// 計算量改善のための式変形
	// R = (z*G + r*pub)/t
	// w := 1/tとして、
	// R = (z*G + r*pub) * w
	// R = ((z*w)*G + (r*w)*pub)を求める

	// w = 1 / t
	one := models.NewFiniteField(big.NewInt(1), *ec.Params().N)
	w := new(models.FiniteField).Div(one, sign.t)

	z := toHash(msg, *ec.Params().N)
	zw := new(models.FiniteField).Mul(z, w)

	x, y := ec.ScalarBaseMult(zw.Value.Bytes())
	zwG := models.ToEllipticCurvePoint(x, y, ec.Params().P)

	rw := new(models.FiniteField).Mul(sign.r, w)

	x, y = ec.ScalarMult(pub.X.Value, pub.Y.Value, rw.Value.Bytes())
	rwpub := models.ToEllipticCurvePoint(x, y, ec.Params().P)

	x, y = ec.Add(zwG.X.Value, zwG.Y.Value, rwpub.X.Value, rwpub.Y.Value)
	R := models.ToEllipticCurvePoint(x, y, ec.Params().P)

	if R.IsZero {
		return false
	}

	return sign.r.Equals(R.X)
}

func main() {
	// ECDSA
	secp256k1 := generateSecp256k1()
	priv, pub := generateKey(secp256k1)

	msg := "hello"

	signature, err := sign(secp256k1, msg, priv)
	if err != nil {
		log.Fatal("sign:", err)
	}

	f := func(msg2 string) {
		isValid := verify(secp256k1, msg2, signature, pub)
		op := "!="
		if isValid {
			op = "=="
		}
		fmt.Printf("%v %v %v\n", msg, op, msg2)
	}

	f(msg)     // OK
	f("hollo") // NG
	f("here")  // NG
}
