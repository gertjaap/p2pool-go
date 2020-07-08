package util

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

func GetMyPublicIP() (net.IP, error) {
	url := "https://api.ipify.org?format=text" // we are using a pulib IP API, we're using ipify here, below are some others
	// https://www.ipify.org
	// http://myexternalip.com
	// http://api.ident.me
	// http://whatismyipaddress.com/api
	fmt.Printf("Getting IP address from  ipify ...\n")
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return net.ParseIP(string(ip)), nil
}

func Sha256d(b []byte) []byte {
	h := sha256.Sum256(b)
	h = sha256.Sum256(h[:])
	return h[:]
}

func GetRandomId() *chainhash.Hash {
	idBytes := make([]byte, 32)
	rand.Read(idBytes)
	id, _ := chainhash.NewHash(idBytes)
	return id
}
