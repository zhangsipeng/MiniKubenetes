package data

import runtimedata "example/Minik8s/pkg/data/RuntimeData"

var credential runtimedata.RuntimeConfig

func InitCredential(c runtimedata.RuntimeConfig) {
	credential = c
}

func GetCredential() runtimedata.RuntimeConfig {
	return credential
}
