package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"elizebch/elizeutils"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
)

const (
	wallFileName string = "Elize.wallet"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

var userWallet *wallet

func Wallet() *wallet {
	if userWallet == nil {
		userWallet = &wallet{}
		if hasWalletFile() {
			userWallet.restoreWallet()
		} else {
			userWallet.createKey()
			persistKey(userWallet.privateKey)
		}
		userWallet.getAddress()
	}
	return userWallet
}

func hasWalletFile() bool {
	_, err := os.Stat(wallFileName)
	return os.IsExist(err)
}

func persistKey(key *ecdsa.PrivateKey) {
	keyAsByte, err := x509.MarshalECPrivateKey(key)
	elizeutils.Errchk(err)
	os.WriteFile(wallFileName, keyAsByte, 0644)
}

func (w *wallet) restoreWallet() {
	bytes, err := os.ReadFile(wallFileName)
	elizeutils.Errchk(err)
	privateKey, err := x509.ParseECPrivateKey(bytes)
	elizeutils.Errchk(err)
	w.privateKey = privateKey
}

func (w *wallet) createKey() {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	elizeutils.Errchk(err)
	w.privateKey = privateKey
}

func (w *wallet) getAddress() {
	w.Address = fmt.Sprintf("%x", append(w.privateKey.X.Bytes(), w.privateKey.Y.Bytes()...))
}

func Sign(w *wallet, TxID string) string {
	TxAsByte, err := hex.DecodeString(TxID)
	elizeutils.Errchk(err)
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, TxAsByte)
	elizeutils.Errchk(err)
	return fmt.Sprintf("%x", append(r.Bytes(), s.Bytes()...))
}

func Verify(signature, address, hashedTxID string) bool {
	x, y, err := restoreBigInt(address)
	elizeutils.Errchk(err)
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	r, s, err := restoreBigInt(signature)
	elizeutils.Errchk(err)
	return ecdsa.Verify(&publicKey, []byte(hashedTxID), r, s)
}

func restoreBigInt(hexaPayload string) (*big.Int, *big.Int, error) {
	bytes, err := hex.DecodeString(hexaPayload)
	if err != nil {
		return nil, nil, err
	}
	bigA, bigB := big.Int{}, big.Int{}
	firstHalf := bytes[:len(bytes)/2]
	lastHalf := bytes[len(bytes)/2:]
	return bigA.SetBytes(firstHalf), bigB.SetBytes(lastHalf), nil
}
