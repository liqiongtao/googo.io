package goo_utils

import (
	"math/big"
)

func parseBigInt(num string) *big.Int {
	x, ok := new(big.Int).SetString(num, 10)
	if !ok {
		return big.NewInt(0)
	}
	return x
}

// 加
func BigIntAdd(num1 string, num2 string) string {
	x := parseBigInt(num1)
	y := parseBigInt(num2)
	return x.Add(x, y).String()
}

// 减
func BigIntReduce(num1 string, num2 string) string {
	x := parseBigInt(num1)
	y := parseBigInt(num2)
	return x.Sub(x, y).String()
}

// 乘
func BigIntMul(num1 string, num2 string) string {
	x := parseBigInt(num1)
	y := parseBigInt(num2)
	return x.Mul(x, y).String()
}

// 除
func BigIntDiv(num1 string, num2 string) string {
	x := parseBigInt(num1)
	y := parseBigInt(num2)
	if y.Sign() == 0 {
		return "0"
	}
	return x.Div(x, y).String()
}

// 取模
func BigIntMod(num1 string, num2 string) string {
	x := parseBigInt(num1)
	y := parseBigInt(num2)
	if y.Sign() == 0 {
		return "0"
	}
	return x.Mod(x, y).String()
}

// 比大小，大于返回1，等于返回0，小于返回-1
func BigIntCmp(num1 string, num2 string) int {
	x := parseBigInt(num1)
	y := parseBigInt(num2)
	return x.Cmp(y)
}
