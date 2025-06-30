# loanproduct 服务文档

本目录包含 loanproduct 服务的API和RPC文档。

## 文档列表

### API文档 (Swagger)

### RPC文档 (Protocol Buffer)
- [loanproduct-rpc-proto.md](./loanproduct-rpc-proto.md) - Proto文件说明

## 使用说明

### 查看Swagger文档
可以使用以下方式查看Swagger文档：
1. 将JSON/YAML文件导入到 [Swagger Editor](https://editor.swagger.io/)
2. 使用IDE插件（如VSCode的OpenAPI扩展）
3. 使用swagger-ui等工具

### 更新文档
当修改了API或RPC定义文件后，运行以下命令更新文档：
```bash
./scripts/gen-code.sh loanproduct docs
```

生成时间: 2025-06-30 10:30:07
