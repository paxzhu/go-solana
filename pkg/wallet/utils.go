package wallet

import (
	"fmt"
	"os"

	"github.com/blocto/solana-go-sdk/types"
)

func generateWallets(numWallets int, filename string) error {
	if filename == "" {
		filename = "wallets.csv"
	}
	// 创建或打开文件
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("无法创建文件: %v", err)
	}
	defer file.Close()

	// 写入文件头
	_, err = file.WriteString("PublicKey,PrivateKey\n")
	if err != nil {
		return fmt.Errorf("无法写入文件: %v", err)
	}

	// 生成钱包并写入文件
	for i := 0; i < numWallets; i++ {
		wallet := types.NewAccount() // 生成新钱包
		publicKey := wallet.PublicKey.ToBase58()
		privateKey := wallet.PrivateKey

		// 将公私钥对写入文件
		_, err = file.WriteString(fmt.Sprintf("%s,%v\n", publicKey, privateKey))
		if err != nil {
			return fmt.Errorf("无法写入文件: %v", err)
		}

		fmt.Printf("Wallet %d: Public Key: %s\n", i+1, publicKey)
	}

	return nil
}
