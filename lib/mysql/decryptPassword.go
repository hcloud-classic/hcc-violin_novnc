package mysql

import (
	"hcc/violin-novnc/action/grpc/client"
	"hcc/violin-novnc/lib/config"
	"hcc/violin-novnc/lib/rsautil"
	"io/ioutil"
)

func getDecryptPassword() (string, error) {
	privKeyData, err := ioutil.ReadFile(config.Rsakey.PrivateKeyFile)
	if err != nil {
		return "", err
	}

	privKey, err := rsautil.BytesToPrivateKey(privKeyData)
	if err != nil {
		return "", err
	}

	encryptedPassword, err := client.RC.GetMYSQLDEncryptedPassword()
	if err != nil {
		return "", err
	}

	decryptedPassword, err := rsautil.DecryptWithPrivateKey(encryptedPassword, privKey)
	if err != nil {
		return "", err
	}

	return string(decryptedPassword), nil
}
