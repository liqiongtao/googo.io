package goo_utils

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

//func TestSHAWithRSA_prod(t *testing.T) {
//	env := "prod"
//
//	privateKey, publicKey, jwksData, err := RSA_SHA256()
//	if err != nil {
//		return
//	}
//
//	var m1 M
//	json.Unmarshal(jwksData, &m1)
//	m1["kid"] = GenIdStr()
//	m := M{
//		"keys": []interface{}{
//			m1,
//		},
//	}
//
//	os.MkdirAll(env, 0755)
//
//	WriteToFile(env+"/private_key.pem", privateKey)
//	WriteToFile(env+"/public_key.pem", publicKey)
//	WriteToFile(env+"/jwks.json", m.Json())
//}

func TestSHAWithRSA_token_prod(t *testing.T) {
	env := "prod"

	privateKey, _ := ioutil.ReadFile(env + "/private_key.pem")
	publicKey, _ := ioutil.ReadFile(env + "/public_key.pem")

	header := map[string]interface{}{
		"kid": "7169331482501058561",
	}

	data := map[string]interface{}{}
	data["scp"] = []string{"durables:registration:create", "durables:products:read", "durables:parties:read"}
	data["firstname"] = ""
	//data["iss"] = "https://championplanet-dev.oss-cn-shenzhen.aliyuncs.com/jwt/"
	data["iss"] = "https://championplanet-pd.oss-cn-shenzhen.aliyuncs.com/jwt/"
	data["lastname"] = ""
	data["phone"] = "18510381580"
	//data["jwks_uri"] = "https://championplanet-dev.oss-cn-shenzhen.aliyuncs.com/jwt/.well-known/jwks.json"
	data["jwks_uri"] = "https://championplanet-pd.oss-cn-shenzhen.aliyuncs.com/jwt/.well-known/jwks.json"
	data["CountryCode"] = "360"
	data["ABONo"] = "750853426"
	data["partyId"] = "3600075085342601"
	data["jti"] = "a2791c5c-e056-49bc-9e2d-d38abeca47cd"
	data["email"] = ""
	data["exp"] = time.Now().Add(100 * time.Hour).Unix()
	data["iat"] = time.Now().Unix()

	token, err := JWTTokenCreate(data, header, privateKey)
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

//func TestSHAWithRSA_test(t *testing.T) {
//	env := "test"
//
//	privateKey, publicKey, jwksData, err := RSA_SHA256()
//	if err != nil {
//		return
//	}
//
//	var m1 M
//	json.Unmarshal(jwksData, &m1)
//	m1["kid"] = GenIdStr()
//	m := M{
//		"keys": []interface{}{
//			m1,
//		},
//	}
//
//	os.MkdirAll(env, 0755)
//
//	WriteToFile(env+"/private_key.pem", privateKey)
//	WriteToFile(env+"/public_key.pem", publicKey)
//	WriteToFile(env+"/jwks.json", m.Json())
//}

func TestSHAWithRSA_token_test(t *testing.T) {
	env := "test"

	privateKey, _ := ioutil.ReadFile(env + "/private_key.pem")
	publicKey, _ := ioutil.ReadFile(env + "/public_key.pem")

	header := map[string]interface{}{
		"kid": "5fe5224e1ba6a55d02e89f6934b45f44",
	}

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

	token, err := JWTTokenCreate(data, header, privateKey)
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
