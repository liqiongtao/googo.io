package goo_utils

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt"
	"io/ioutil"
	"log"
	"testing"
	"time"
)

func TestSHAWithRSA(t *testing.T) {
	//privateKey, publicKey, jwksData, err := RSA_SHA256()
	//if err != nil {
	//	return
	//}

	privateKey, _ := ioutil.ReadFile("private_key.pem")
	publicKey, _ := ioutil.ReadFile("public_key.pem")
	//jwksData, _ := ioutil.ReadFile("jwks.json")

	//var m1 M
	//json.Unmarshal(jwksData, &m1)
	//m1["kid"] = GenIdStr()
	//m := M{
	//	"keys": []interface{}{
	//		m1,
	//	},
	//}

	//fmt.Println(m.String())

	//WriteToFile("private_key.pem", privateKey)
	//WriteToFile("public_key.pem", publicKey)
	//WriteToFile("jwks.json", m.Json())

	data := map[string]interface{}{}
	data["scp"] = []string{"durables:registration:create", "durables:products:read", "durables:parties:read"}
	data["firstname"] = ""
	//data["iss"] = "https://championplanet-dev.oss-cn-shenzhen.aliyuncs.com/jwt/"
	data["iss"] = "https://championplanet-dev.amwaynet.com.cn/jwt/"
	data["lastname"] = ""
	data["phone"] = "18510381580"
	//data["jwks_uri"] = "https://championplanet-dev.oss-cn-shenzhen.aliyuncs.com/jwt/.well-known/jwks.json"
	data["jwks_uri"] = "https://championplanet-dev.amwaynet.com.cn/jwt/.well-known/jwks.json"
	data["CountryCode"] = "360"
	data["ABONo"] = "750853426"
	data["partyId"] = "3600075085342601"
	data["jti"] = "a2791c5c-e056-49bc-9e2d-d38abeca47cd"
	data["email"] = ""
	data["exp"] = time.Now().Add(100 * time.Hour).Unix()
	data["iat"] = time.Now().Unix()

	token, err := JWTTokenCreate(data, privateKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(token)

	tt, err := JWT_TokenParse(token, publicKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tt.Claims)
}

func TestAESCBCDecrypt(t *testing.T) {
	// 生成RSA密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	// 提取公钥
	publicKey := &privateKey.PublicKey

	// 创建JWT
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	//claims["user_id"] = 12345
	//claims["username"] = "john_doe"
	//claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // 设置过期时间

	claims["scp"] = []string{"durables:registration:create", "durables:products:read", "durables:parties:read"}
	claims["firstname"] = ""
	claims["iss"] = "https://championplanet-dev.oss-cn-shenzhen.aliyuncs.com/jwt/"
	claims["lastname"] = ""
	claims["phone"] = "18510381580"
	claims["jwks_uri"] = "https://championplanet-dev.oss-cn-shenzhen.aliyuncs.com/jwt/jwks.json"
	claims["CountryCode"] = "360"
	claims["ABONo"] = "750853426"
	claims["partyId"] = "3600075085342601"
	claims["jti"] = "a2791c5c-e056-49bc-9e2d-d38abeca47cd"
	claims["email"] = ""
	claims["exp"] = time.Now().Add(time.Hour * 100).Unix()
	claims["iat"] = time.Now().Add(-time.Hour * 1).Unix()

	// 使用私钥签名JWT
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Signed Token:", signedToken)
	fmt.Println("publicKey:", publicKey)

	// 验证JWT（可选，为了演示完整性）
	parsedToken, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range parsedToken.Claims.(jwt.MapClaims) {
		fmt.Println(k, v)
	}
}
