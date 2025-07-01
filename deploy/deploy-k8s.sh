#!/bin/bash

# 部署K8s脚本 - HuinongFinancial微服务
# 使用方法：./scripts/deploy-k8s.sh [service_name] [version]
# 示例：./scripts/deploy-k8s.sh appuser v1.0.0
#       ./scripts/deploy-k8s.sh all v1.0.0

set -e

# 颜色定义