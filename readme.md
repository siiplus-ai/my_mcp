go mod init my_mcp

docker build -t my_mcp .

touch .vscode/mcp.json

https://mcpgolang.com/quickstart
go get github.com/metoro-io/mcp-golang

go build -o my_mcp main.go

sudo cp -f my_mcp /usr/local/go/bin/