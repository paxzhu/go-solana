// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"

// 	"github.com/paxzhu/go-solana/pkg/wallet"

// 	"github.com/blocto/solana-go-sdk/types"
// )

// func main() {
// 	// 创建 WalletManager 实例
// 	wm, err := wallet.NewWalletManager("devnet")
// 	if err != nil {
// 		log.Fatalf("Error creating WalletManager: %v", err)
// 	}

// 	// 创建新账户并显示公钥
// 	publicKey, err := wm.CreateAccount()
// 	if err != nil {
// 		log.Fatalf("Error creating account: %v", err)
// 	}
// 	fmt.Printf("New account created with public key: %s\n", publicKey)

// 	// 请求空投以获得测试 SOL
// 	txSig, err := wm.client.RequestAirdrop(context.Background(), publicKey, 1000000000) // 1 SOL
// 	if err != nil {
// 		log.Fatalf("Error requesting airdrop: %v", err)
// 	}
// 	fmt.Printf("Airdrop transaction signature: %s\n", txSig)

// 	// 检查账户余额
// 	balance, err := wm.CheckAmount(context.Background(), "SOL")
// 	if err != nil {
// 		log.Fatalf("Error checking balance: %v", err)
// 	}
// 	fmt.Printf("Balance for account %s: %d SOL\n", publicKey, balance)

// 	// 创建另一个账户以进行转账
// 	toAccount := types.NewAccount()
// 	toPublicKey := toAccount.PublicKey.ToBase58()
// 	fmt.Printf("Transfering to new account with public key: %s\n", toPublicKey)

// 	// 转账 0.5 SOL
// 	err = wm.Transfer(context.Background(), "SOL", toPublicKey, 500000000) // 0.5 SOL
// 	if err != nil {
// 		log.Fatalf("Error transferring SOL: %v", err)
// 	}
// 	fmt.Println("Transfer successful!")

// 	// 检查目标账户余额
// 	toBalance, err := wm.CheckAmount(context.Background(), "SOL")
// 	if err != nil {
// 		log.Fatalf("Error checking recipient balance: %v", err)
// 	}
// 	fmt.Printf("Balance for recipient account %s: %d SOL\n", toPublicKey, toBalance)

// 	// 加载现有账户（假设私钥存储在文件中）
// 	err = wm.LoadAccount("path/to/private/key.json") // 替换为实际路径
// 	if err != nil {
// 		log.Fatalf("Error loading account: %v", err)
// 	}
// 	fmt.Printf("Loaded account with public key: %s\n", wm.account.PublicKey.ToBase58())

// 	// 再次检查余额
// 	newBalance, err := wm.CheckAmount(context.Background(), "SOL")
// 	if err != nil {
// 		log.Fatalf("Error checking loaded account balance: %v", err)
// 	}
// 	fmt.Printf("Balance for loaded account: %d SOL\n", newBalance)
// }

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/rpc"
)

func getTokenBalance(walletAddress string, tokenMintAddress string) (uint64, error) {
	// 创建一个 Solana 客户端
	c := client.NewClient(rpc.DevnetRPCEndpoint)
	// c := client.NewClient(rpc.LocalnetRPCEndpoint)
	if tokenMintAddress == "" {
		balance, err := c.GetBalance(
			context.TODO(),
			walletAddress,
		)
		if err != nil {
			return 0, fmt.Errorf("failed to get SOL balance: %w", err)
		}
		return balance, nil
	}

	// 使用 GetTokenAccountsByOwnerByMint 获取代币账户
	tokenAccounts, err := c.GetTokenAccountsByOwnerByMint(
		context.TODO(),
		walletAddress,
		tokenMintAddress,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get token accounts: %w", err)
	}
	// 打印 tokenAccounts 的详细信息
	for i, account := range tokenAccounts {
		fmt.Printf("Account %d:\n", i+1)
		fmt.Printf("  PublicKey: %s\n", account.PublicKey.ToBase58())
		fmt.Printf("  Mint: %s\n", account.Mint.ToBase58())
		fmt.Printf("  Owner: %s\n", account.Owner.ToBase58())
		fmt.Printf("  Token Amount: %d\n", account.Amount)
		fmt.Printf("  TokenAccountState: %d\n", account.State)
	}
	return 0, nil
	// // 遍历代币账户获取余额
	// for _, account := range tokenAccounts {
	// 	accountInfo, err := c.GetAccountInfo(context.TODO(), account.Pubkey)
	// 	if err != nil {
	// 		return 0, fmt.Errorf("failed to get account info: %w", err)
	// 	}

	// 	parsedAccount, err := token.AccountFromData(accountInfo.Data)
	// 	if err != nil {
	// 		return 0, fmt.Errorf("failed to parse account data: %w", err)
	// 	}

	// 	return parsedAccount.Amount, nil
	// }

	// return 0, fmt.Errorf("no token account found for the specified mint")
}

func main() {
	walletAddress := "657u8g2j83MmSd7sbxkmJLD6onbqp86SJPoQbNWNeToe"
	tokenMintAddress := "Cr7Q5ttDHLj64ASZiaDzyEAWyMooLBecT1YwEQscUw2k"

	balance, err := getTokenBalance(walletAddress, tokenMintAddress)
	if err != nil {
		log.Fatalf("Error getting token balance: %v", err)
	}

	fmt.Printf("Token balance: %d\n", balance)

	balance, err = getTokenBalance(walletAddress, "")
	if err != nil {
		log.Fatalf("Error getting token balance: %v", err)
	}
	fmt.Printf("Token balance: %d\n", balance)
}
