package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/paxzhu/go-solana/pkg/wallet"
)

const SOL_MINT_ADDR = "So11111111111111111111111111111111111111112"
const GOAT_MINT_ADDR = "CzLSujWBLFsSjncfkh59rUFqvafWcY5tzedWJSuypump"
const keyPath1 = "assets/wallet_7dEc3i8Niz.json"
const keyPath2 = "assets/wallet_3qQEWctNXM.json"
const singleTransferAmount = 100000000 // 0.1 SOL
const QuoteAPI = "https://quote-api.jup.ag/v6/quote?inputMint=So11111111111111111111111111111111111111112&outputMint=CzLSujWBLFsSjncfkh59rUFqvafWcY5tzedWJSuypump&amount=500000000&slippageBps=100"
const mintAddr = "J2LDFD6Cso2aRq6XhT2hBq1jcgqo9j4bEBM1VCh71ojX"
const mintAuthority = "3qQEWctNXMxHiyoaQwAwn2xBRZGC7c2XeoPuLcZ6ZrUP"
const ata2 = "3zwAJGJTTnSJfVHaLLtubVmpZXDLFY51fhVZjg1HCQxb"
const ata1 = "31Fuuv2ekbt4ATZD1w8x9MkiiEnScWGZkGFDPteHsgXP"

func main() {
	// wallet.GenerateWallets(2)
	// 创建 WalletManager 实例
	wm1, err := wallet.NewWalletManager("devnet")
	if err != nil {
		log.Fatalf("Error creating WalletManager: %v", err)
	}
	wm2, err := wallet.NewWalletManager("devnet")
	if err != nil {
		log.Fatalf("Error creating WalletManager: %v", err)
	}

	// 检查私钥文件是否存在，不存在则创建新账户，存在则加载现有账户
	var publicKey1 string
	if _, err := os.Stat(keyPath1); os.IsNotExist(err) {
		// 创建新账户并显示公钥
		publicKey1, err = wm1.CreateAccount()
		if err != nil {
			log.Fatalf("Error creating Account1: %v", err)
		}
		fmt.Printf("New Account1 created with public key: %s\n", publicKey1)
	} else {
		// 加载现有账户（假设私钥存储在文件中）
		publicKey1, err = wm1.LoadAccount(keyPath1) // 替换为实际路径
		if err != nil {
			log.Fatalf("Error loading Account1: %v", err)
		}
		fmt.Printf("Loaded Account1 with public key: %s\n", publicKey1)
	}

	var publicKey2 string
	if _, err := os.Stat(keyPath2); os.IsNotExist(err) {
		// 创建新账户并显示公钥
		publicKey2, err = wm2.CreateAccount()
		if err != nil {
			log.Fatalf("Error creating Account2: %v", err)
		}
		fmt.Printf("New Account2 created with public key: %s\n", publicKey2)
	} else {
		// 加载现有账户（假设私钥存储在文件中）
		publicKey2, err = wm2.LoadAccount(keyPath2) // 替换为实际路径
		if err != nil {
			log.Fatalf("Error loading Account2: %v", err)
		}
		fmt.Printf("Loaded Account2 with public key: %s\n", publicKey2)
	}

	// 检查余额是否为空，如果是则请求空投
	balance1, err := wm1.CheckAmount(context.Background(), SOL_MINT_ADDR)
	if err != nil {
		log.Fatalf("Error checking balance: %v", err)
	}
	fmt.Printf("Balance for Account1 %s: %d lamports\n", publicKey1, balance1)
	if balance1 == 0 {
		// 请求空投以获得测试 SOL
		isSuccessed := wm1.RequestAirdrop(publicKey1, 1000000000) // 1 SOL
		if isSuccessed {
			fmt.Println("Airdrop successful for Account1!")
		} else {
			fmt.Println("Airdrop failed! for Account1")
		}
	}

	balance2, err := wm2.CheckAmount(context.Background(), SOL_MINT_ADDR)
	if err != nil {
		log.Fatalf("Error checking balance: %v", err)
	}
	fmt.Printf("Balance for Account2 %s: %d lamports\n", publicKey2, balance2)
	if balance2 == 0 {
		isSuccessed := wm2.RequestAirdrop(publicKey2, 1000000000) // 1 SOL
		if isSuccessed {
			fmt.Println("Airdrop successful for Account2!")
		} else {
			fmt.Println("Airdrop failed! for Account2")
		}
	}

	// 账户1向账户2转账
	txhash1, err := wm1.TransferSOL(context.Background(), publicKey2, singleTransferAmount)
	if err != nil {
		log.Fatalf("[Account1 transfer SOL to Account2 failed] - Error transferring SOL: %v", err)
	}
	fmt.Printf("[Account1 transfer SOL to Account2 successful] - txhash1: %s\n", txhash1)

	// 账户2向账户1转账
	txhash2, err := wm2.TransferSOL(context.Background(), publicKey1, singleTransferAmount)
	if err != nil {
		log.Fatalf("[Account2 transfer SOL to Account1 failed] - Error transferring SOL: %v", err)
	}
	fmt.Printf("[Account2 transfer SOL to Account1 successful] - txhash2: %s\n", txhash2)

	// wallet.CreateMint(wm1.Account, wm2.Account)
	// ata1, err := wm1.CreateTokenAccount(context.Background(), mintAddr)
	// if err != nil {
	// 	log.Fatalf("Error creating token account: %v", err)
	// }
	// fmt.Println("CreateTokenAccount successful!ata1:", ata1)
	// ata2, err := wm2.CreateTokenAccount(context.Background(), mintAddr)
	// if err != nil {
	// 	log.Fatalf("Error creating token account: %v", err)
	// }
	// fmt.Println("CreateTokenAccount successful!ata2:", ata2)
	_, err = wm1.TransferTokensChecked(mintAddr, wm2.Account, ata1, ata2, singleTransferAmount, 8)
	if err != nil {
		log.Fatalf("[Account1 transfer token to Account2 failed] - Error transferring token: %v", err)
	}
	// 检查目标账户余额
	// balance, err = wm1.CheckAmount(context.Background(), SOL_MINT_ADDR)
	// if err != nil {
	// 	log.Fatalf("Error checking balance: %v", err)
	// }
	// fmt.Printf("Balance for account %s: %d lamports\n", toPublicKey, balance)

	// ata, err := wm1.CreateTokenAccount(context.Background(), GOAT_MINT_ADDR)
	// if err != nil {
	// 	log.Fatalf("Error creating token account: %v", err)
	// }
	// fmt.Println("CreateTokenAccount successful!ata:", ata)
	// // 测试 Jupiter Swap
	// // wm.GetQuote(QuoteAPI)
	// wm1.Sell(context.Background(), GOAT_MINT_ADDR, 0.5)
}

// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"

// 	"github.com/blocto/solana-go-sdk/client"
// 	"github.com/blocto/solana-go-sdk/rpc"
// )

// func getTokenBalance(walletAddress string, tokenMintAddress string) (uint64, error) {
// 	// 创建一个 Solana 客户端
// 	c := client.NewClient(rpc.DevnetRPCEndpoint)
// 	// c := client.NewClient(rpc.LocalnetRPCEndpoint)
// 	if tokenMintAddress == "" {
// 		balance, err := c.GetBalance(
// 			context.TODO(),
// 			walletAddress,
// 		)
// 		if err != nil {
// 			return 0, fmt.Errorf("failed to get SOL balance: %w", err)
// 		}
// 		return balance, nil
// 	}

// 	// 使用 GetTokenAccountsByOwnerByMint 获取代币账户
// 	tokenAccounts, err := c.GetTokenAccountsByOwnerByMint(
// 		context.TODO(),
// 		walletAddress,
// 		tokenMintAddress,
// 	)
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to get token accounts: %w", err)
// 	}
// 	// 打印 tokenAccounts 的详细信息
// 	for i, account := range tokenAccounts {
// 		fmt.Printf("Account %d:\n", i+1)
// 		fmt.Printf("  PublicKey: %s\n", account.PublicKey.ToBase58())
// 		fmt.Printf("  Mint: %s\n", account.Mint.ToBase58())
// 		fmt.Printf("  Owner: %s\n", account.Owner.ToBase58())
// 		fmt.Printf("  Token Amount: %d\n", account.Amount)
// 		fmt.Printf("  TokenAccountState: %d\n", account.State)
// 	}
// 	return 0, nil
// 	// // 遍历代币账户获取余额
// 	// for _, account := range tokenAccounts {
// 	// 	accountInfo, err := c.GetAccountInfo(context.TODO(), account.Pubkey)
// 	// 	if err != nil {
// 	// 		return 0, fmt.Errorf("failed to get account info: %w", err)
// 	// 	}

// 	// 	parsedAccount, err := token.AccountFromData(accountInfo.Data)
// 	// 	if err != nil {
// 	// 		return 0, fmt.Errorf("failed to parse account data: %w", err)
// 	// 	}

// 	// 	return parsedAccount.Amount, nil
// 	// }

// 	// return 0, fmt.Errorf("no token account found for the specified mint")
// }

// func main() {
// 	walletAddress := "657u8g2j83MmSd7sbxkmJLD6onbqp86SJPoQbNWNeToe"
// 	tokenMintAddress := "Cr7Q5ttDHLj64ASZiaDzyEAWyMooLBecT1YwEQscUw2k"

// 	balance, err := getTokenBalance(walletAddress, tokenMintAddress)
// 	if err != nil {
// 		log.Fatalf("Error getting token balance: %v", err)
// 	}

// 	fmt.Printf("Token balance: %d\n", balance)

// 	balance, err = getTokenBalance(walletAddress, "")
// 	if err != nil {
// 		log.Fatalf("Error getting token balance: %v", err)
// 	}
// 	fmt.Printf("Token balance: %d\n", balance)
// }
