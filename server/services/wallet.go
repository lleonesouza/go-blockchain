package service

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/craton-api/chain/blockchain"
	"github.com/craton-api/chain/server/config"
)

// Wallet
type Wallet struct {
	walletUrl string
}

func MakeWallet() *Wallet {
	env := config.MakeEnv()
	fmt.Println(env.WALLET_URL)
	return &Wallet{walletUrl: env.WALLET_URL}
}

func DeserializeWallets(data []byte) (*blockchain.Wallet, error) {
	wallet := blockchain.Wallet{}
	bytesReader := bytes.NewReader(data)
	gob.Register(elliptic.P256())

	err := gob.NewDecoder(bytesReader).Decode(&wallet)
	if err != nil {
		return nil, err
	} else {
		return &wallet, nil
	}
}

func (w *Wallet) GetWallet(address string) *blockchain.Wallet {
	addrStr := fmt.Sprintf("%s/wallet/%s", w.walletUrl, address)
	sb, _ := Requester(addrStr)
	wa, err := DeserializeWallets([]byte(sb))
	if err != nil {
		fmt.Println(err)
	}
	return wa
}

func (w *Wallet) GetWallets() ([]string, error) {
	addrStr := fmt.Sprintf("%s/wallets", w.walletUrl)
	type ad struct {
		Addresses []string `json:"addresses"`
	}
	var addresses ad

	sb, err := Requester(addrStr)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(sb), &addresses)
	return addresses.Addresses, nil
}

// GetPubKeyHash returns the Address, PubKey and PubKeyHash. In this Order.
func (w *Wallet) GetPubKey(address string) ([]byte, error) {
	addrStr := fmt.Sprintf("%s/wallet/%s", w.walletUrl, address)
	sb, err := Requester(addrStr)
	if err != nil {
		return nil, err
	}

	return []byte(sb), nil
}

// Create a new wallet
func (w *Wallet) Create() (string, error) {
	addrStr := fmt.Sprintf("%s/wallet", w.walletUrl)
	sb, err := HttpClient("post", addrStr, nil)
	if err != nil {
		return "", err
	}
	return sb, nil
}

func Requester(url string) (string, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := netClient.Get(url)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	sb := string(body)
	return sb, nil
}

func HttpClient(method string, url string, body []byte) (string, error) {
	type Address struct {
		Address string
	}

	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	switch method {
	case "get":
		resp, err := netClient.Get(url)
		if err != nil {
			return "", err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		sb := string(body)
		return sb, nil

	case "post":
		resp, err := netClient.Post(url, "application/json", nil)
		if err != nil {
			return "", err
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		var address Address
		json.Unmarshal(body, &address)

		return address.Address, nil

	default:
		return "", errors.New("invalid URL")
	}

}
