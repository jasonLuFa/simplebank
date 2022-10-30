package util

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// generates a random integer between min and max
func randFloats(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
// generate a random string of length n
func RandomString(n int) string{
	var sb strings.Builder

	for i :=0 ; i < n; i++{
		c := alphabet[rand.Intn(len(alphabet))]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string{
	return RandomString(6)
}

func RandomMoney() string{
	return fmt.Sprintf("%.2f",randFloats(1000,10000))
}

func StringToFloat64(s string) float64{
	result,_ := strconv.ParseFloat(s, 64) 
	return result
}

func RandomCurrency() string{
	currencies := []string{USD, EUR, CAD}
	return currencies[rand.Intn(len(currencies))]
}