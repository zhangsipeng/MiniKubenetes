package data

import "crypto/rsa"

var serverKey *rsa.PrivateKey
var authKey *rsa.PrivateKey

func GetServerKey() *rsa.PrivateKey {
	return serverKey
}

func SetServerKey(key *rsa.PrivateKey) {
	serverKey = key
}

func GetAuthKey() *rsa.PrivateKey {
	return authKey
}

func SetAuthKey(key *rsa.PrivateKey) {
	authKey = key
}
