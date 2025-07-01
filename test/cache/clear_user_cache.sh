#!/bin/bash

# Redis配置 - 根据您的实际情况修改
REDIS_HOST="127.0.0.1"
REDIS_PORT="6379"
REDIS_PASSWORD="ChinaSkills@"
REDIS_DB="0"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

show_usage() {
    echo -e "${BLUE}用法: $0 <操作> [参数]${NC}"
    echo ""
    echo "操作:"
    echo "  all                    - 清理所有用户缓存"
    echo "  phone <手机号>         - 清理指定手机号的缓存"
    echo "  id <用户ID>            - 清理指定用户ID的缓存"
    echo "  list                   - 列出所有用户缓存键"
    echo "  test                   - 测试Redis连接"
    echo ""
    echo "示例:"
    echo "  $0 all"
    echo "  $0 phone 13452552490"
    echo "  $0 id 1001"
    echo "  $0 list"
}

test_redis() {
    echo -e "${BLUE}测试Redis连接...${NC}"
    if redis-cli -h $REDIS_HOST -p $REDIS_PORT -n $REDIS_DB -a $REDIS_PASSWORD ping > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Redis连接成功${NC}"
        return 0
    else
        echo -e "${RED}❌ Redis连接失败${NC}"
        echo "请检查Redis配置: $REDIS_HOST:$REDIS_PORT DB:$REDIS_DB"
        return 1
    fi
}

clear_all_cache() {
    echo -e "${YELLOW}🧹 清理所有用户缓存...${NC}"
    
    # 清理AppUser缓存
    appuser_count=$(redis-cli -h $REDIS_HOST -p $REDIS_PORT -n $REDIS_DB -a $REDIS_PASSWORD --scan --pattern "cache:appUsers:*" | wc -l)
    if [ $appuser_count -gt 0 ]; then
        redis-cli -h $REDIS_HOST -p $REDIS_PORT -n $REDIS_DB -a $REDIS_PASSWORD --scan --pattern "cache:appUsers:*" | xargs redis-cli -h $REDIS_HOST -p $REDIS_PORT -n $REDIS_DB -a $REDIS_PASSWORD del > /dev/null 2>&1
        echo -e "${GREEN}删除AppUser缓存: $appuser_count 个键${NC}"
    else
        echo "AppUser缓存: 0 个键"
    fi
    
    # 清理OaUser缓存
    oauser_count=$(redis-cli -h $REDIS_HOST -p $REDIS_PORT -n $REDIS_DB -a $REDIS_PASSWORD --scan --pattern "cache:oaUsers:*" | wc -l)
    if [ $oauser_count -gt 0 ]; then
        redis-cli -h $REDIS_HOST -p $REDIS_PORT -n $REDIS_DB -a $REDIS_PA SSWORD --scan --pattern "cache:oaUsers:*" | xargs redis-cli -h $REDIS_HOST -p $REDIS_PORT -n $REDIS_DB -a $REDIS_PASSWORD del > /dev/null 2>&1
        echo -e "${GREEN}删除OaUser缓存: $oauser_count 个键${NC}"
    else
        echo "OaUser缓存: 0 个键"
    fi
    
    total=$((appuser_count + oauser_count))
    echo -e "${GREEN}✅ 总共删除了 $total 个缓存键${NC}"
}

clear_cache_by_phone() {
    local phone=$1
    echo -e "${YELLOW}🧹 清理手机号 $phone 的缓存...${NC}"
    
    keys=(
        "cache:appUsers:phone:$phone"
        "cache:oaUsers:phone:$phone"
    )
    
    deleted=0
    for key in "${keys[@]}"; do
        if redis-cli -h $REDIS_HOST -p $REDIS_PORT -n $REDIS_DB -a $REDIS_PASSWORD exists "$key" | grep -q "1"; then
            redis-cli -h $REDIS_HOST -p $REDIS_PORT -n $REDIS_DB -a $REDIS_PASSWORD del "$key" > /dev/null
            echo -e "${GREEN}删除缓存键: $key${NC}"
            deleted=$((deleted + 1))
        fi
    done
    
    # 查找可能的ID缓存（通过模式匹配）
    appuser_id_keys=$(redis-cli -h $REDIS_HOST -p $REDIS_PORT -n $REDIS_DB -a $REDIS_PASSWORD --scan --pattern "cache:appUsers:id:*")
    oauser_id_keys=$(redis-cli -h $REDIS_HOST -p $REDIS_PORT -n $REDIS_DB -a $REDIS_PASSWORD --scan --pattern "cache:oaUsers:id:*")
    
    # 这里只是删除所有ID缓存，因为shell脚本不容易解析JSON内容
    # 在生产环境中，建议使用更精确的方法
    
    echo -e "${GREEN}✅ 为手机号 $phone 删除了 $deleted 个明确的缓存键${NC}"
    echo -e "${YELLOW}💡 建议运行 '$0 all' 来确保清理所有相关缓存${NC}"
}

clear_cache_by_id() {
    local id=$1
    echo -e "${YELLOW}🧹 清理用户ID $id 的缓存...${NC}"
    
    keys=(
        "cache:appUsers:id:$id"
        "cache:oaUsers:id:$id"
    )
    
    deleted=0
    for key in "${keys[@]}"; do
        if redis-cli -h $REDIS_HOST -p $REDIS_PORT -n $REDIS_DB -a $REDIS_PASSWORD exists "$key" | grep -q "1"; then
            redis-cli -h $REDIS_HOST -p $REDIS_PORT -n $REDIS_DB -a $REDIS_PASSWORD del "$key" > /dev/null
            echo -e "${GREEN}删除缓存键: $key${NC}"
            deleted=$((deleted + 1))
        fi
    done
    
    echo -e "${GREEN}✅ 为用户ID $id 删除了 $deleted 个缓存键${NC}"
}

list_cache_keys() {
    echo -e "${BLUE}📋 列出所有用户缓存键...${NC}"
    
    echo -e "\n${YELLOW}=== AppUser 缓存 ===${NC}"
    appuser_keys=$(redis-cli -h $REDIS_HOST -p $REDIS_PORT -n $REDIS_DB -a $REDIS_PASSWORD --scan --pattern "cache:appUsers:*")
    if [ -z "$appuser_keys" ]; then
        echo "无缓存键"
    else
        echo "$appuser_keys" | while IFS= read -r key; do
            if [ -n "$key" ]; then
                echo "  $key"
            fi
        done
    fi
    
    echo -e "\n${YELLOW}=== OaUser 缓存 ===${NC}"
    oauser_keys=$(redis-cli -h $REDIS_HOST -p $REDIS_PORT -n $REDIS_DB -a $REDIS_PASSWORD --scan --pattern "cache:oaUsers:*")
    if [ -z "$oauser_keys" ]; then
        echo "无缓存键"
    else
        echo "$oauser_keys" | while IFS= read -r key; do
            if [ -n "$key" ]; then
                echo "  $key"
            fi
        done
    fi
}

# 主逻辑
if [ $# -eq 0 ]; then
    show_usage
    exit 1
fi

# 测试Redis连接
if ! test_redis; then
    exit 1
fi

case "$1" in
    "all")
        clear_all_cache
        ;;
    "phone")
        if [ -z "$2" ]; then
            echo -e "${RED}❌ 请提供手机号${NC}"
            exit 1
        fi
        clear_cache_by_phone "$2"
        ;;
    "id")
        if [ -z "$2" ]; then
            echo -e "${RED}❌ 请提供用户ID${NC}"
            exit 1
        fi
        clear_cache_by_id "$2"
        ;;
    "list")
        list_cache_keys
        ;;
    "test")
        echo -e "${GREEN}✅ Redis连接正常${NC}"
        ;;
    *)
        echo -e "${RED}❌ 未知操作: $1${NC}"
        show_usage
        exit 1
        ;;
esac 