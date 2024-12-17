package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/paxzhu/go-solana/internal/wallet" // 注意：需要替换 'your-username' 为你的项目名
)

func main() {
	// 生成 5 个钱包
	wallets, err := wallet.GenerateMultipleWallets(5)
	if err != nil {
		log.Fatal("生成钱包失败:", err)
	}

	// 打印生成的钱包信息
	for i, wallet := range wallets {
		fmt.Printf("Wallet %d:\n", i+1)
		fmt.Printf("Public Key: %s\n", wallet.PublicKey)
		fmt.Printf("Private Key: %s\n", wallet.PrivateKey)
		fmt.Println("-------------------")
	}

	// 创建钱包连接（使用 devnet）
	conn := wallet.NewConnection("")

	if len(wallets) > 0 {
		publicKey := wallets[0].PublicKey
		err := os.WriteFile("wallet_pubkey.txt", []byte(publicKey), 0644)
		if err != nil {
			log.Fatal("无法保存公钥:", err)
		}
		fmt.Println("第一个钱包的公钥已保存到 wallet_pubkey.txt")
	}

	cmd := exec.Command("bash", "airdrop.sh")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("执行airdrop脚本失败: %v\n", err)
	}
	fmt.Printf("Airdrop结果: %s\n", output)

	// 连接第一个钱包并获取信息
	if len(wallets) > 0 {
		fmt.Println("空投后的钱包信息:")
		err := conn.GetWalletInfo(wallets[0].PublicKey)
		if err != nil {
			log.Printf("获取钱包信息失败: %v\n", err)
		}
	}

	// if len(wallets) > 0 {
	// 	fmt.Println("尝试获取第一个钱包的信息:")
	// 	err := conn.GetWalletInfo(wallets[0].PublicKey)
	// 	if err != nil {
	// 		log.Printf("获取钱包信息失败: %v\n", err)
	// 	}
	// } else {
	// 	fmt.Println("没有生成任何钱包")
	// }

}
