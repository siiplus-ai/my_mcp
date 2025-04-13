package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// MCPサーバーインスタンスの作成
	s := server.NewMCPServer(
		"Calculator Demo",
		"1.0.0",
		server.WithResourceCapabilities(true, true), // Resource の機能で使われるオプションなのでToolの公開のみであれば不要そう
		server.WithLogging(),
	)

	// 四則計算ツールのインターフェース登録
	calculatorTool := mcp.NewTool("calculate",
		mcp.WithDescription("Perform basic arithmetic operations"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to perform (add, subtract, multiply, divide)"),
			mcp.Enum("add", "subtract", "multiply", "divide"),
		),
		mcp.WithNumber("x",
			mcp.Required(),
			mcp.Description("First number"),
		),
		mcp.WithNumber("y",
			mcp.Required(),
			mcp.Description("Second number"),
		),
	)

	templateGenTool := mcp.NewTool("template",
		mcp.WithDescription("Generate code template"),
		mcp.WithString("template",
			mcp.Required(),
			mcp.Description("Select generater template (controller, usecase)"),
			mcp.Enum("controller", "usecase"),
		),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("file name"),
		),
	)

	// 四則計算ツールを実装
	s.AddTool(calculatorTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		op := request.Params.Arguments["operation"].(string)
		x := request.Params.Arguments["x"].(float64)
		y := request.Params.Arguments["y"].(float64)

		var result float64
		switch op {
		case "add":
			result = x + y + 1
		case "subtract":
			result = x - y
		case "multiply":
			result = x * y
		case "divide":
			return mcp.NewToolResultError("未対応の機能です"), nil
		}
		// /tmpフォルダにファイルを作成
		// 環境変数　MCP_PATH　を取得
		mcpPath := os.Getenv("MCP_PATH")
		workFolder := os.Getenv("WORK_SPACE_FOLDER")

		file, err := os.Create(fmt.Sprintf("%s/calculation_result.txt", workFolder))
		if err != nil {
			return mcp.NewToolResultError("ファイル作成に失敗しました"), nil
		}
		defer file.Close()

		_, err = file.WriteString(fmt.Sprintf("Calculation result: %.2f %s %s\n", result, mcpPath, workFolder))
		if err != nil {
			return mcp.NewToolResultError("ファイル書き込みに失敗しました"), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%.2f", result)), nil
	})

	// code generationツールを実装
	s.AddTool(templateGenTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		template := request.Params.Arguments["template"].(string)
		name := request.Params.Arguments["name"].(string)
		// 環境変数 WORK_SPACE_FOLDER を取得
		workFolder := os.Getenv("WORK_SPACE_FOLDER")
		var fileName string
		switch template {
		case "controller":
			fileName = fmt.Sprintf("%s/%s/%s.go", workFolder, template, name)
		case "usecase":
			fileName = fmt.Sprintf("%s/%s/%s.go", workFolder, template, name)
		case "divide":
			return mcp.NewToolResultError("未対応の機能です"), nil
		}

		file, err := os.Create(fileName)
		if err != nil {
			return mcp.NewToolResultError("ファイル作成に失敗しました"), nil
		}
		defer file.Close()
		code := fmt.Sprintf("package main\n\nfunc %s() {\n\t// TODO: Implement %s\n}\n", name, name)
		_, err = file.WriteString(code)
		if err != nil {
			return mcp.NewToolResultError("ファイル書き込みに失敗しました"), nil
		}
		return mcp.NewToolResultText(code), nil
	})

	// サーバー起動
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
