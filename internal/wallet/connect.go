package wallet

import (
	"context"
	"fmt"

	// 更新导入路径
	"github.com/blocto/solana-go-sdk/client"
)

const DevnetRPCEndpoint = "https://api.devnet.solana.com"

// WalletConnection 管理钱包连接的结构体
type WalletConnection struct {
	client *client.Client
}

// NewConnection 创建新的钱包连接
func NewConnection(endpoint string) *WalletConnection {
	if endpoint == "" {
		endpoint = DevnetRPCEndpoint
	}
	return &WalletConnection{
		client: client.NewClient(endpoint),
	}
}

// GetWalletInfo 获取钱包信息
func (wc *WalletConnection) GetWalletInfo(publicKey string) error {
	// 直接使用字符串形式的公钥
	// 获取账户信息
	accountInfo, err := wc.client.GetAccountInfo(context.Background(), publicKey)
	if err != nil {
		return fmt.Errorf("获取账户信息失败: %v", err)
	}

	// 获取余额
	balance, err := wc.client.GetBalance(context.Background(), publicKey)
	if err != nil {
		return fmt.Errorf("获取余额失败: %v", err)
	}

	fmt.Printf("Account Info: %+v\n", accountInfo)
	fmt.Printf("Balance: %d lamports (%.9f SOL)\n", balance, float64(balance)/1e9)

	return nil
}
