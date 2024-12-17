# Solana 钱包生成器和连接器

## 项目概述

这个项目是一个用 Go 语言编写的 Solana 钱包工具。它可以批量生成 Solana 钱包，并提供与 Solana 网络连接的功能。主要特性包括：

1. 批量生成 Solana 钱包（包括公钥和私钥）
2. 连接到 Solana 网络（默认使用 Devnet）
3. 查询钱包信息（如余额和账户数据）
4. 尝试为新生成的钱包请求 SOL 空投（在 Devnet 上）

## 安装要求

- Go 1.22 或更高版本
- Solana CLI（用于空投功能）

## 安装步骤

1. 克隆仓库：

   - git clone [https://github.com/paxzhu/go-solana.git](https://github.com/yourusername/go-solana.git)
   - cd go-solana
2. 安装依赖：

   - go mod tidy
3. 确保 Solana CLI 已正确安装并配置为使用 Devnet：

   - solana config set --url [https://api.devnet.solana.com](https://api.devnet.solana.com/)

## 使用方法

运行主程序：
   - go run cmd/main/main.go

这将执行以下操作：

1. 生成多个 Solana 钱包
2. 打印每个钱包的公钥和私钥
3. 保存第一个钱包的公钥到文件
4. 尝试为第一个钱包请求 SOL 空投
5. 查询并显示第一个钱包的信息

### 说明

- **cmd/main**: 此目录下的 `main.go` 文件是应用的入口点，所有执行流程从这里开始。
- **internal/wallet**: 这个包包含了与 Solana 钱包交互的所有功能实现，如创建新钱包、连接到现有钱包等。
- **go.mod 和 go.sum**: 这两个文件用于管理 Go 项目的依赖关系，确保构建环境的一致性。

## 注意事项

- 本项目默认使用 Solana Devnet。在生产环境中使用时，请确保更改网络设置。
- 空投功能可能受到 Devnet 的速率限制。如果遇到问题，请稍后再试或使用在线水龙头。
