package wallet

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/paxzhu/go-solana/pkg/config"
)

const SOL_MINT_ADDR = "So11111111111111111111111111111111111111112"

// WalletManager 管理 Solana 钱包的结构体
type WalletManager struct {
	client  *client.Client
	network string // "mainnet", "testnet", "devnet", "localhost"
	account types.Account
	// tokenCache map[string]common.PublicKey // 缓存代币地址
}

// 添加 Jupiter API 相关常量
const (
	JupiterQuoteAPI    = "https://quote-api.jup.ag/v6/quote"
	JupiterSwapAPI     = "https://quote-api.jup.ag/v6/swap"
	defaultSlippageBps = 100 // 1% slippage
)

// QuoteResponse Jupiter 报价响应
type QuoteResponse struct {
	InputMint  string `json:"inputMint"`
	OutputMint string `json:"outputMint"`
	Amount     string `json:"amount"`
	OutAmount  string `json:"outAmount"`
	SwapData   []byte `json:"swapData"`
}

// NewWalletManager 创建新的钱包管理器
func NewWalletManager(network string) (*WalletManager, error) {
	endpoint, ok := config.NetworkConfig[network]
	if !ok {
		return nil, fmt.Errorf("unsupported network: %s", network)
	}

	return &WalletManager{
		client:  client.NewClient(endpoint),
		network: network,
		// tokenCache: make(map[string]common.PublicKey),
	}, nil
}

// CreateAccount 创建新账户
func (wm *WalletManager) CreateAccount() (string, error) {
	account := types.NewAccount()
	wm.account = account

	// 如果是测试网络，自动请求空投
	if wm.network == "testnet" || wm.network == "devnet" {
		_, err := wm.client.RequestAirdrop(
			context.Background(),
			account.PublicKey.ToBase58(),
			1000000000, // 1 SOL
		)
		if err != nil {
			return "", fmt.Errorf("airdrop failed: %w", err)
		}
	}

	return account.PublicKey.ToBase58(), nil
}

// LoadAccount 从私钥文件加载账户
func (wm *WalletManager) LoadAccount(keyPath string) error {
	data, err := os.ReadFile(keyPath)
	if err != nil {
		return fmt.Errorf("failed to read key file: %w", err)
	}

	var privateKey []byte
	err = json.Unmarshal(data, &privateKey)
	if err != nil {
		return fmt.Errorf("failed to unmarshal private key: %w", err)
	}

	account, err := types.AccountFromBytes(privateKey)
	if err != nil {
		return fmt.Errorf("failed to create account from bytes: %w", err)
	}

	wm.account = account
	return nil
}

// CheckAmount 检查指定代币的余额
func (wm *WalletManager) CheckAmount(ctx context.Context, mintAddr string) (uint64, error) {
	if wm.account.PublicKey.ToBase58() == "" {
		return 0, errors.New("no account loaded")
	}

	// mintAddr为空，则直接获取SOL余额
	if mintAddr == SOL_MINT_ADDR {
		balance, err := wm.client.GetBalance(
			ctx,
			wm.account.PublicKey.ToBase58(),
		)
		if err != nil {
			return 0, fmt.Errorf("failed to get SOL balance: %w", err)
		}
		return balance, nil
	}

	// 获取SPL代币账户信息
	tokenAccounts, err := wm.client.GetTokenAccountsByOwnerByMint(ctx, wm.account.PublicKey.ToBase58(), mintAddr)
	if err != nil {
		return 0, fmt.Errorf("failed to get token balance: %w", err)
	}

	if len(tokenAccounts) == 0 {
		return 0, nil
	}

	return tokenAccounts[0].Amount, nil
}

// Transfer 转账功能
func (wm *WalletManager) TransferSOL(ctx context.Context, toAddress string, amount uint64) error {
	senderPrivateKey := wm.account.PrivateKey
	senderPubKey := wm.account.PublicKey
	receiverPubKey := common.PublicKeyFromString(toAddress)

	// 获取最近的区块哈希（用于构建交易）
	recentBlockhashResponse, err := wm.client.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("获取最近区块哈希失败: %v", err)
	}
	// 合约的 ProgramID
	programID := common.PublicKeyFromString() // 替换为你的合约 Program ID

	// 构建交易消息
	message := types.NewMessage(types.NewMessageParams{
		FeePayer: 			senderPubKey,
		RecentBlockhash: 	recentBlockhashResponse.Blockhash,
		Instructions: []types.Instruction{

		},
	})}
	)
}

// Buy 市价买入代币
func (wm *WalletManager) Buy(ctx context.Context, tokenTicker string, amountInSOL float64) error {
	// 获取代币地址
	tokenMint, err := wm.getTokenPubKey(tokenTicker)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	// 构建报价请求
	quoteURL := fmt.Sprintf("%s?inputMint=%s&outputMint=%s&amount=%d&slippageBps=%d",
		JupiterQuoteAPI,
		"So11111111111111111111111111111111111111112", // SOL 的代币地址
		tokenMint.ToBase58(),
		uint64(amountInSOL*1e9), // 转换为 lamports
		defaultSlippageBps,
	)

	// 获取报价
	quote, err := wm.getQuote(quoteURL)
	if err != nil {
		return fmt.Errorf("failed to get quote: %w", err)
	}

	// 执行交换
	err = wm.executeSwap(ctx, quote)
	if err != nil {
		return fmt.Errorf("swap failed: %w", err)
	}

	fmt.Printf("Successfully bought %s tokens\n", tokenTicker)
	return nil
}

// Sell 市价卖出代币
func (wm *WalletManager) Sell(ctx context.Context, tokenTicker string, amount float64) error {
	// 获取代币地址
	tokenMint, err := wm.getTokenPubKey(tokenTicker)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	// 检查代币余额
	balance, err := wm.CheckAmount(ctx, tokenTicker)
	if err != nil {
		return fmt.Errorf("failed to check balance: %w", err)
	}

	if float64(balance) < amount {
		return fmt.Errorf("insufficient balance: have %v, want %v", balance, amount)
	}

	// 构建报价请求
	quoteURL := fmt.Sprintf("%s?inputMint=%s&outputMint=%s&amount=%d&slippageBps=%d",
		JupiterQuoteAPI,
		tokenMint.ToBase58(),
		"So11111111111111111111111111111111111111112", // SOL
		uint64(amount),
		defaultSlippageBps,
	)

	// 获取报价
	quote, err := wm.getQuote(quoteURL)
	if err != nil {
		return fmt.Errorf("failed to get quote: %w", err)
	}

	// 执行交换
	err = wm.executeSwap(ctx, quote)
	if err != nil {
		return fmt.Errorf("swap failed: %w", err)
	}

	fmt.Printf("Successfully sold %s tokens\n", tokenTicker)
	return nil
}

// 内部辅助方法

// getQuote 从 Jupiter 获取报价
func (wm *WalletManager) getQuote(quoteURL string) (*QuoteResponse, error) {
	resp, err := http.Get(quoteURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("quote failed: %s", string(body))
	}

	var quote QuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&quote); err != nil {
		return nil, err
	}

	return &quote, nil
}

// executeSwap 执行代币交换
func (wm *WalletManager) executeSwap(ctx context.Context, quote *QuoteResponse) error {
	// 构建交换请求
	swapReq := struct {
		UserPublicKey string `json:"userPublicKey"`
		SwapData      []byte `json:"swapData"`
	}{
		UserPublicKey: wm.account.PublicKey.ToBase58(),
		SwapData:      quote.SwapData,
	}

	// 发送交换请求
	reqBody, err := json.Marshal(swapReq)
	if err != nil {
		return err
	}

	resp, err := http.Post(JupiterSwapAPI, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("swap request failed: %s", string(body))
	}

	// 处理响应，获取交易指令
	var swapResp struct {
		SwapTransaction []byte `json:"swapTransaction"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&swapResp); err != nil {
		return err
	}

	// 签名并发送交易
	tx, err := types.TransactionFromData(swapResp.SwapTransaction)
	if err != nil {
		return err
	}

	// 签名交易
	tx.Sign([]types.Account{wm.account})

	// 发送交易
	sig, err := wm.client.SendTransaction(ctx, tx)
	if err != nil {
		return err
	}

	// 等待交易确认
	return wm.client.WaitForConfirmation(ctx, sig, nil)
}
