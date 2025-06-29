#!/bin/bash

# 惠农金融微服务认证架构部署脚本
set -e

echo "🚀 开始部署惠农金融微服务认证架构..."

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查依赖
check_dependencies() {
    echo -e "${YELLOW}🔍 检查部署依赖...${NC}"
    
    if ! command -v kubectl &> /dev/null; then
        echo -e "${RED}❌ kubectl 未安装${NC}"
        exit 1
    fi
    
    if ! command -v goctl &> /dev/null; then
        echo -e "${RED}❌ goctl 未安装${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✅ 依赖检查通过${NC}"
}

# 生成微服务代码
generate_code() {
    echo -e "${YELLOW}🔨 生成微服务代码...${NC}"
    
    # 生成 auth-api 代码
    echo "生成 auth-api 代码..."
    cd app/auth/api && goctl api go -api auth.api -dir . --style goZero
    
    # 生成 auth-rpc 代码
    echo "生成 auth-rpc 代码..."
    cd ../rpc && goctl rpc protoc auth.proto --go_out=. --go-grpc_out=. --zrpc_out=. --style goZero
    
    # 生成 appuser-rpc 代码
    echo "生成 appuser-rpc 代码..."
    cd ../../appuser/rpc && goctl rpc protoc appuser.proto --go_out=. --go-grpc_out=. --zrpc_out=. --style goZero
    
    # 生成 oauser-rpc 代码
    echo "生成 oauser-rpc 代码..."
    cd ../../oauser/rpc && goctl rpc protoc oauser.proto --go_out=. --go-grpc_out=. --zrpc_out=. --style goZero
    
    cd ../../../../
    echo -e "${GREEN}✅ 代码生成完成${NC}"
}

# 构建Docker镜像
build_images() {
    echo -e "${YELLOW}🐳 构建Docker镜像...${NC}"
    
    # 构建 auth-api 镜像
    echo "构建 auth-api 镜像..."
    docker build -t registry.huinong.internal/auth-api:latest app/auth/api/
    
    # 构建 auth-rpc 镜像
    echo "构建 auth-rpc 镜像..."
    docker build -t registry.huinong.internal/auth-rpc:latest app/auth/rpc/
    
    # 构建其他服务镜像...
    # (这里可以继续添加其他服务的镜像构建)
    
    echo -e "${GREEN}✅ 镜像构建完成${NC}"
}

# 部署基础设施
deploy_infrastructure() {
    echo -e "${YELLOW}🏗️ 部署基础设施...${NC}"
    
    # 创建命名空间
    kubectl create namespace huinong-financial || echo "命名空间已存在"
    
    # 部署MySQL
    echo "部署MySQL..."
    kubectl apply -f deploy/mysql/mysql.yaml
    
    # 部署Redis
    echo "部署Redis..."
    kubectl apply -f deploy/redis/redis.yaml
    
    # 等待基础设施就绪
    echo "等待基础设施就绪..."
    kubectl wait --for=condition=ready pod -l app=mysql -n huinong-financial --timeout=300s
    kubectl wait --for=condition=ready pod -l app=redis -n huinong-financial --timeout=300s
    
    echo -e "${GREEN}✅ 基础设施部署完成${NC}"
}

# 部署认证服务
deploy_auth_services() {
    echo -e "${YELLOW}🔐 部署认证服务...${NC}"
    
    # 部署 auth-rpc
    echo "部署 auth-rpc..."
    kubectl apply -f deploy/k8s/auth-rpc.yaml
    
    # 等待 auth-rpc 就绪
    kubectl wait --for=condition=ready pod -l app=auth-rpc -n huinong-financial --timeout=300s
    
    # 部署 auth-api
    echo "部署 auth-api..."
    kubectl apply -f deploy/k8s/auth-api.yaml
    
    # 等待 auth-api 就绪
    kubectl wait --for=condition=ready pod -l app=auth-api -n huinong-financial --timeout=300s
    
    echo -e "${GREEN}✅ 认证服务部署完成${NC}"
}

# 部署业务服务
deploy_business_services() {
    echo -e "${YELLOW}💼 部署业务服务...${NC}"
    
    # 部署用户服务
    echo "部署用户服务..."
    kubectl apply -f deploy/k8s/appuser-api.yaml
    kubectl apply -f deploy/k8s/appuser-rpc.yaml
    kubectl apply -f deploy/k8s/oauser-api.yaml
    kubectl apply -f deploy/k8s/oauser-rpc.yaml
    
    # 部署业务服务
    echo "部署业务服务..."
    kubectl apply -f deploy/k8s/loan-api.yaml
    kubectl apply -f deploy/k8s/loan-rpc.yaml
    kubectl apply -f deploy/k8s/risk-api.yaml
    kubectl apply -f deploy/k8s/risk-rpc.yaml
    
    echo -e "${GREEN}✅ 业务服务部署完成${NC}"
}

# 配置网关
configure_gateway() {
    echo -e "${YELLOW}🌐 配置Higress网关...${NC}"
    
    # 应用网关配置
    kubectl apply -f deploy/higress/auth-config.yaml
    
    # 等待Ingress就绪
    sleep 10
    
    echo -e "${GREEN}✅ 网关配置完成${NC}"
}

# 验证部署
verify_deployment() {
    echo -e "${YELLOW}🔍 验证部署状态...${NC}"
    
    # 检查所有Pod状态
    echo "检查Pod状态..."
    kubectl get pods -n huinong-financial
    
    # 检查服务状态
    echo "检查服务状态..."
    kubectl get svc -n huinong-financial
    
    # 检查Ingress状态
    echo "检查Ingress状态..."
    kubectl get ingress
    
    # 测试认证接口
    echo "测试认证接口..."
    GATEWAY_IP=$(kubectl get ingress huinong-financial-gateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
    if [ -z "$GATEWAY_IP" ]; then
        GATEWAY_IP="api.huinong.internal"
    fi
    
    echo "网关地址: $GATEWAY_IP"
    
    # 健康检查
    echo "执行健康检查..."
    curl -f http://$GATEWAY_IP/health || echo "健康检查失败，请检查服务状态"
    
    echo -e "${GREEN}✅ 部署验证完成${NC}"
}

# 显示访问信息
show_access_info() {
    echo -e "${GREEN}🎉 部署成功！${NC}"
    echo ""
    echo "访问信息："
    echo "- API网关: http://api.huinong.internal"
    echo "- 认证接口: http://api.huinong.internal/api/v1/auth"
    echo "- Consul UI: consul.huinong.internal:32130/ui/"
    echo ""
    echo "测试命令："
    echo "# 用户注册"
    echo 'curl -X POST http://api.huinong.internal/api/v1/auth/register \'
    echo '  -H "Content-Type: application/json" \'
    echo '  -d '"'"'{"account":"13800138000","password":"123456","name":"测试用户","age":25,"type":"appuser"}'"'"
    echo ""
    echo "# 用户登录"
    echo 'curl -X POST http://api.huinong.internal/api/v1/auth/login \'
    echo '  -H "Content-Type: application/json" \'
    echo '  -d '"'"'{"account":"13800138000","password":"123456","type":"appuser"}'"'"
    echo ""
    echo "详细文档请查看: docs/认证鉴权架构设计.md"
}

# 主执行流程
main() {
    case "$1" in
        "check")
            check_dependencies
            ;;
        "code")
            generate_code
            ;;
        "build")
            build_images
            ;;
        "infra")
            deploy_infrastructure
            ;;
        "auth")
            deploy_auth_services
            ;;
        "business")
            deploy_business_services
            ;;
        "gateway")
            configure_gateway
            ;;
        "verify")
            verify_deployment
            ;;
        "all"|"")
            check_dependencies
            generate_code
            build_images
            deploy_infrastructure
            deploy_auth_services
            deploy_business_services
            configure_gateway
            verify_deployment
            show_access_info
            ;;
        *)
            echo "用法: $0 [check|code|build|infra|auth|business|gateway|verify|all]"
            echo ""
            echo "  check     - 检查部署依赖"
            echo "  code      - 生成微服务代码"
            echo "  build     - 构建Docker镜像"
            echo "  infra     - 部署基础设施"
            echo "  auth      - 部署认证服务"
            echo "  business  - 部署业务服务"
            echo "  gateway   - 配置网关"
            echo "  verify    - 验证部署"
            echo "  all       - 执行完整部署流程（默认）"
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@" 