package wallet

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/token"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/paxzhu/go-solana/pkg/config"
)

// WalletManager 管理 Solana 钱包的结构体
type WalletManager struct {
	client     *client.Client
	network    string // "mainnet", "testnet", "devnet", "localhost"
	account    types.Account
	tokenCache map[string]common.PublicKey // 缓存代币地址
}

// NewWalletManager 创建新的钱包管理器
func NewWalletManager(network string) (*WalletManager, error) {
	endpoint, ok := config.NetworkConfig[network]
	if !ok {
		return nil, fmt.Errorf("unsupported network: %s", network)
	}

	return &WalletManager{
		client:     client.NewClient(endpoint),
		network:    network,
		tokenCache: make(map[string]common.PublicKey),
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
func (wm *WalletManager) CheckAmount(ctx context.Context, ticker string) (uint64, error) {
	if wm.account.PublicKey.ToBase58() == "" {
		return 0, errors.New("no account loaded")
	}

	if ticker == "SOL" {
		balance, err := wm.client.GetBalance(
			ctx,
			wm.account.PublicKey.ToBase58(),
		)
		if err != nil {
			return 0, fmt.Errorf("failed to get SOL balance: %w", err)
		}
		return balance, nil
	}

	// 获取代币账户信息
	tokenPubKey, err := wm.getTokenPubKey(ticker)
	if err != nil {
		return 0, err
	}

	tokenAccount, err := wm.client.GetTokenAccountsByOwnerByMint(
		ctx,
		wm.account.PublicKey.ToBase58(),
		&client.GetTokenAccountsConfig{
			Mint: tokenPubKey.ToBase58(),
		},
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get token balance: %w", err)
	}

	if len(tokenAccount.Value) == 0 {
		return 0, nil
	}

	return tokenAccount.Value[0].Account.Data.Parsed.Info.TokenAmount.Amount, nil
}

// Transfer 转账功能
func (wm *WalletManager) Transfer(ctx context.Context, ticker string, to string, amount uint64) error {
	if ticker == "SOL" {
		tx, err := types.NewTransaction(types.NewTransactionParam{
			Message: types.NewMessage(types.NewMessageParam{
				FeePayer: wm.account.PublicKey,
				Instructions: []types.Instruction{
					types.NewInstruction(types.SystemProgramID, []byte{2, 0, 0, 0}, // Transfer instruction
						types.NewAccountMeta(wm.account.PublicKey, true, true),
						types.NewAccountMeta(common.PublicKeyFromString(to), false, true),
					),
				},
			}),
			Signers: []types.Account{wm.account},
		})
		if err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		sig, err := wm.client.SendTransaction(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to send transaction: %w", err)
		}

		return wm.client.WaitForConfirmation(ctx, sig, nil)
	}

	// Token transfer implementation
	tokenPubKey, err := wm.getTokenPubKey(ticker)
	if err != nil {
		return err
	}

	tx, err := token.NewTransferInstruction(
		amount,
		wm.account.PublicKey,
		common.PublicKeyFromString(to),
		tokenPubKey,
		[]types.Account{wm.account},
	).Build()
	if err != nil {
		return fmt.Errorf("failed to build token transfer: %w", err)
	}

	sig, err := wm.client.SendTransaction(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to send token transaction: %w", err)
	}

	return wm.client.WaitForConfirmation(ctx, sig, nil)
}

// 内部辅助方法：获取代币公钥
func (wm *WalletManager) getTokenPubKey(ticker string) (common.PublicKey, error) {
	if pubKey, ok := wm.tokenCache[ticker]; ok {
		return pubKey, nil
	}

	// 这里应该实现从代币符号到代币地址的映射
	// 实际使用时需要维护一个代币地址映射表
	return common.PublicKey{}, fmt.Errorf("token not supported: %s", ticker)
}
