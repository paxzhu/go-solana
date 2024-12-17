package wallet

import (
	"crypto/ed25519"

	"github.com/mr-tron/base58"
)

// WalletInfo 存储钱包信息的结构体
type WalletInfo struct {
	PublicKey  string
	PrivateKey string
}

// GenerateWallet 生成单个 Solana 钱包
func GenerateWallet() (*WalletInfo, error) {
	// 生成新的公私钥对
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}

	return &WalletInfo{
		PublicKey:  base58.Encode(pubKey),
		PrivateKey: base58.Encode(privKey),
	}, nil
}

// GenerateMultipleWallets 批量生成钱包
func GenerateMultipleWallets(count int) ([]*WalletInfo, error) {
	wallets := make([]*WalletInfo, count)

	for i := 0; i < count; i++ {
		wallet, err := GenerateWallet()
		if err != nil {
			return nil, err
		}
		wallets[i] = wallet
	}

	return wallets, nil
}
