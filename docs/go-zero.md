# go-zero 学习笔记

## 环境安装

### goctl

goctl 是 Go-Zero 的内置脚手架，是提升开发效率的一大利器，可以一键生成代码、文档、部署 k8s yaml、dockerfile 等。

运行命令安装 goctl

```shell
go install github.com/zeromicro/go-zero/tools/goctl@latest
```

验证安装

```shell
goctl --version
```

### protoc

protoc 是一个用于生成代码的工具，它可以根据 proto 文件生成C++、Java、Python、Go、PHP 等多重语言的代码，而 gRPC 的代码生成还依赖 protoc-gen-go，protoc-gen-go-grpc 插件来配合生成 Go 语言的 gRPC 代码。

运行命令安装 protoc

```shell
go env -w GOPROXY=https://goproxy.cn,direct
goctl env check --install --verbose --force
```

### goctl 集成编译环境插件安装

goctl-intellij 是 Go-Zero api 描述语言的 intellij 编辑器插件，支持 api 描述语言高亮、语法检测、快速提示、创建模板特性。

https://go-zero.dev/docs/tasks/installation/goctl-intellij

goctl vscode 编辑器插件可以安装在 1.46.0+ 版本的 Visual Studio Code 上，首先请确保你的 Visual Studio Code 版本符合要求，并已安装 goctl 命令行工具。如果尚未安装 Visual Studio Code，请安装并打开 Visual Studio Code。

https://go-zero.dev/docs/tasks/installation/goctl-vscode

### Air 热加载工具

Air 是一个实时热加载工具，能够在代码修改后自动重新编译并重启 Go 应用，显著提升开发效率。

```shell
go install github.com/air-verse/air@latest
```

运行命令启动 Air

```shell
air
```

Air 支持在项目根目录创建配置文件(.air.toml 文件)

```toml
# .air.toml
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/main ."
bin = "tmp/main"
# 包含文件
include_ext = ["go","yaml","yml"]
# 排除文件
exclude_dir = ["assets","tmp","vendor"]
```

---

## go-zero DSL(领域特性语言)

Go-Zero 框架支持两种 DSL：api 与 protobuf。

两种 DSL 都能实现微服务开发。这两种 DSL 各自拥有独特的优势，能够满足不同场景下的微服务开发需求。`api` DSL 专注于 HTTP 服务，适用于构建对外暴露的 RESTful API 微服务；而 `protobuf` 则用于构建高性能的 gRPC 服务，常用于微服务间的内部通信，以实现高效的数据交换。

### api

api 是 Go-Zero 自研的领域特性语言（下文称 api 语言 或 api 描述语言），旨在实现人性化的基础描述语言，作为生成 HTTP 服务最基本的描述语言。
**特点：**
* **简洁高效：** api 语言设计简洁，语法类似 Go 结构体，易于学习和使用。
* **专注于 HTTP 服务：** 主要用于定义 HTTP 接口，配合 `goctl` 工具可以快速生成 HTTP 服务代码。
* **约定优于配置：** 通过清晰的 api 文件定义，实现代码生成、路由注册、参数校验等自动化，减少手动配置。
**适用场景：**
* 构建 RESTful API 服务。
* 快速开发对外暴露的 HTTP 接口。

api 领域特性语言包含语法版本，info 块，结构体声明，服务描述等几大块语法组成，其中结构体和 Golang 结构体 语法几乎一样，只是移出了 struct 关键字。

https://go-zero.dev/docs/tasks/dsl/api

### protobuf

Protocol buffers 是 Google 的语言中立、平台中立、可扩展的结构化数据序列化机制——像 XML，但更小、更快、更简单。您定义了一次数据的结构化方式，然后您可以使用特殊生成的源代码轻松地将结构化数据写入各种数据流并使用各种语言从中读取结构化数据。
**特点：**
* **跨语言、跨平台：** 支持多种编程语言，适用于构建多语言服务。
* **高效序列化：** 数据传输效率高，编码后体积小，适合高性能的 RPC 通信。
* **严格的契约定义：** 通过 `.proto` 文件强制定义服务接口和数据结构，保证服务间的兼容性。
**适用场景：**
* 构建高性能的 gRPC 服务。
* 微服务之间进行内部通信。
* 对数据传输效率和跨语言兼容性有较高要求的场景。

goctl 根据 proto 生成 gRPC 代码时：
1. service 中的 Message（入参&出参） 必须要在 main proto 文件中，不支持引入的文件
2. 引入的 Message 只能嵌套在 main proto 中的 Message 中使用
3. goctl 生成 gRPC 代码时，不会生成被引入的 proto 文件的 Go 代码，需要自行预先生成

https://go-zero.dev/docs/tasks/dsl/proto

### 选择建议

在选择使用 `api` 或 `protobuf` 进行微服务开发时，可以根据以下几点进行权衡：

*   **对外接口优先（HTTP/RESTful）：** 如果你的微服务主要职责是向外部系统（如前端应用、第三方服务）提供 HTTP/RESTful API 接口，并且追求开发效率和简洁性，那么 `api` DSL 是一个非常好的选择。它能够快速生成符合 HTTP 规范的服务代码，方便前后端分离开发。

*   **内部通信优先（高性能/跨语言）：** 如果你的微服务主要处理服务间的内部通信，对数据传输效率、序列化性能和跨语言兼容性有较高要求，那么 `protobuf` 结合 gRPC 是更优的选择。它能提供更快的通信速度和更小的数据包大小，非常适合构建大规模的分布式系统。

*   **混合使用：** 在实际项目中，你可能会发现两种 DSL 混合使用的场景。例如，对外提供 HTTP 接口的服务可以使用 `api`，而内部服务之间进行高效率通信则使用 `protobuf`。Go-Zero 框架的灵活性允许你在同一个项目中同时使用这两种 DSL。

---

## Go-Zero 项目结构

Go-Zero 采用约定式项目布局，不同服务类型 (API/RPC) 结构略有差异，但核心目录保持一致。

### API 服务典型目录结构

```shell
go-zero-demo
├── api                 # API 定义文件 (.api)
├── etc                 # 配置文件 (YAML)
|   └── service.yaml
├── internal            # 业务逻辑代码 (goctl 生成 + 手动编写)
|   ├── config          # 配置结构体 (goctl 生成)
|   ├── handler         # HTTP 路由处理函数 (goctl 生成)
|   ├── logic           # 业务逻辑处理 (手动编写)
|   ├── middleware      # 中间件 (手动编写)
|   ├── svc             # 依赖资源池 (数据库/缓存/队列等)
|   └── types           # 请求/响应数据结构体 (goctl 生成)
├── model               # 数据库模型 (goctl 生成)
|   └── user.go         # 用户模型 (goctl 生成)
├── scripts             # 脚本
├── go.mod              # 依赖管理
├── go.sum              # 依赖管理
├── service-name.go     # 服务入口
└── README.md           # 项目说明
```

### RPC 服务典型目录结构

```shell
go-zero-demo
├── proto               # protobuf 定义文件 (.proto)
├── etc                 # 配置文件 (YAML)
|   └── service.yaml
├── internal            # 业务逻辑代码 (goctl 生成 + 手动编写)
|   ├── config          # 配置结构体 (goctl 生成)
|   ├── logic           # 业务逻辑处理 (手动编写)
|   ├── svc             # 依赖资源池 (数据库/缓存/队列等)
|   └── types           # 请求/响应数据结构体 (goctl 生成)
├── model               # 数据库模型 (goctl 生成)
|   └── user.go         # 用户模型
├── scripts             # 脚本
├── go.mod              # 依赖管理
├── go.sum              # 依赖管理
├── service-name.go     # 服务入口
└── README.md           # 项目说明
```

### 关键目录详解

#### api 与 protobuf 目录

api 与 protobuf 目录是 Go-Zero 项目中最重要的目录，它们分别用于定义 API 和 RPC 服务。

##### api 目录

api 目录存放 `.api` 文件，定义 HTTP 接口的路由、请求/响应数据结构体。

示例语法：

```api
get /user/info (UserInfoReq) returns (UserInfoResp)
```

##### protobuf 目录

protobuf 目录存放 `.proto` 文件，定义 gRPC 服务的方法和消息格式。

示例语法：

```protobuf
service UserService {
    rpc GetUserInfo(GetUserInfoReq) returns (GetUserInfoResp);
}
```

#### etc 目录

etc 目录存放服务配置文件，包含服务配置文件 (YAML)，支持动态加载。

示例语法：

```yaml
Name: hello-api
Host: 0.0.0.0
Port: 8888
Database:
  Host: 127.0.0.1
  Port: 3306
  User: root
  Password: 123456
  Name: hello
  MaxIdleConns: 10
```

#### internal 目录

##### config 目录

config 目录存放配置结构体，包含配置结构体 (goctl 生成)，自动生成配置结构体与 `etc/*.yaml` 配置文件映射。

示例语法：

```go
type Config struct {
    rest.RestConf
    // 其他配置
    Database struct {
      Host      string 
      Port      int   
      User      string
      Password  string 
      Name      string 
      MaxIdleConns int
    }
}
```

##### logic 目录

业务逻辑处理目录，包含业务逻辑处理函数 (手动编写)，业务逻辑处理函数是业务逻辑的实现。

示例语法：

```go
type UserLogic struct {
	svcCtx *svc.servicecontext
}

func NewUserLogic(ctx context.Context, svcCtx *svc.servicecontext) *UserLogic {
	return &UserLogic{
		svcCtx: svcCtx,
	}
}
```

##### handler 目录

handler 目录存放 HTTP 路由处理函数，包含 HTTP 路由处理函数 (goctl 生成)，HTTP 路由处理函数是 HTTP 请求的入口。

示例语法：

```go
func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	logic.NewUserLogic(r.Context(), svc.Newservicecontext(c)).GetUserInfo()
}
```

##### middleware 目录

middleware 目录存放中间件，包含中间件 (手动编写)，中间件是 HTTP 请求的拦截器。

示例语法：

```go
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 认证逻辑
		next(w, r)
	}
}
```

##### svc 目录

svc 目录存放依赖资源池，包含依赖资源池 (数据库/缓存/队列等)，依赖资源池是业务逻辑处理的依赖。

示例语法：

```go
type servicecontext struct {
	Config config.Config
	UserModel model.UserModel
}

func Newservicecontext(c config.Config) *servicecontext {
	return &servicecontext{
		Config: c,
	}
}
```

##### types 目录

types 目录存放请求/响应数据结构体，包含请求/响应数据结构体 (goctl 生成)，请求/响应数据结构体是 HTTP 请求/响应的参数和返回值。

示例语法：

```go
type UserInfoReq struct {
	Id int64 `path:"id"`
}

type UserInfoResp struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
```

#### model 目录

model 目录存放数据库模型，包含数据库模型 (goctl model 生成)，数据库模型是数据库的映射。

goctl model 命令可以生成数据库模型，支持多种数据库，如 MySQL、PostgreSQL、SQLite 等。

```shell
goctl model mysql dsn -c -src ./sql/user.sql -dir ./model
```

生成的 model 示例语法：

```go
type User struct {
	Id   int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Name string `gorm:"column:name"`
}
```

---

## 代码生成逻辑

### 生成流程

1.定义接口：编写 `.api` 或 `.proto` 文件，定义接口。
2. 生成骨架代码：
    - API 服务：生成 `api` 目录下的 `*.api` 文件，包含接口定义、请求/响应结构体。
    ```shell
    goctl api go -api *.api -dir .
    ```
    - RPC 服务：生成 `proto` 目录下的 `*.proto` 文件，包含服务定义、消息定义。
    ```shell
    goctl rpc proto *.proto --go_out=. --go-grpc_out=. --zrpc_out=.
    ```
3. 填充业务逻辑：在 `internal/logic` 下实现业务逻辑。

自动生成示例：

```api
syntax = "v1"

type Request {
	Name string `path:"name,options=you|me"`
}

type Response {
	Message string `json:"message"`
}

service hello-api {
	@handler HelloHandler
	get /from/:name (Request) returns (Response)
}
```

输出：
- 自动生成 `handler/userinfohandler.go` (路由绑定、参数校验、返回值处理)
- 自动生成 `logic/userinfologic.go` (空逻辑待实现)
- 自动生成 `types/types.go` (请求/响应数据结构体)

---

## 设计哲学分析

### 分层清晰

- `handler`：仅做参数绑定/返回值处理
- `logic`：纯业务逻辑处理
- `svc`：统一管理依赖资源池

### 生成与手写分离

- 生成与手写分离，生成代码仅做参数绑定/返回值处理，不允许修改生成代码
- 手写代码仅做业务逻辑处理，不会被覆盖

### 约定优于配置

Go-Zero 框架遵循约定优于配置的设计哲学，通过约定好的目录结构和文件命名，实现代码生成、路由注册、参数校验等自动化，减少手动配置。

---

## Go-Zero 配置加载

Go-Zero 框架支持多种配置加载方式，包括：

- *.json
- *.toml
- *.yaml
- *.yml

此外，Go-Zero 还支持从环境变量中加载配置，环境变量的优先级高于配置文件。

### 定义 Config 结构体

```go
type Config struct {
  Name string
  Host string `json:"default=0.0.0.0"`
  Port int
}
```

如上，我们定义了一个 Config 结构体。

### 定义配置路径

```go
var f = flag.String("f", "config.yaml", "config file")
```

我们一般希望可以在启动的时候指定配置文件的路径，所以我们定一个 flag 用于接受配置文件的路径。

### 编写配置文件

我们使用 yaml 格式当做实例，生成 config.yaml 文件。写入如下内容

```yaml
Name: hello
Host: 0.0.0.0
Port: 8080
```

### 加载配置

```go
flag.Parse()
var c Config
conf.MustLoad(*f, &c)
fmt.Println(c.Name)
```

---

## Go-Zero 数据库操作 (ORM 与 Model)

Go-Zero 框架支持多种数据库操作，包括：

- MySQL
- PostgreSQL
- SQLite
- MongoDB

并自动生成 model 文件和 ORM 操作文件。

### 创建 sql 文件

创建一个 `sql` 文件夹，里面存放 sql 文件，用于生成 model。

实例创建一个 `sql/user.sql` 文件，写入如下内容

```sql
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `age` int(11) NOT NULL,
  `email` varchar(255) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 创建 api 文件

创建一个 `api/user.api` 文件，写入如下内容

```api
type (
	User {
		Id        int64  `json:"id"`
		Name      string `json:"name"`
		Age       int64  `json:"age"`
		Email     string `json:"email"`
		CreatedAt int64  `json:"created_at"`
		UpdatedAt int64  `json:"updated_at"`
	}
	UserCreateReq {
		Name  string `json:"name"`
		Age   int64  `json:"age"`
		Email string `json:"email"`
	}
	UserUpdateReq {
		Id    int64  `json:"id"`
		Name  string `json:"name"`
		Age   int64  `json:"age"`
		Email string `json:"email"`
	}
	UserQueryReq {
		Id    int64  `json:"id"`
		Email string `json:"email"`
	}
	Response {
		Code    int64  `json:"code"`
		Message string `json:"message"`
		Data    User   `json:"data"`
	}
)

service user-api {
	@handler createUser
	post /user (UserCreateReq) returns (Response)

	@handler updateUser
	put /user (UserUpdateReq) returns (Response)

	@handler queryUser
	get /user (UserQueryReq) returns (Response)
}
```

### 生成 api 代码

生成 api 代码，生成 api 代码时，会自动生成 config、svc、model、handler、logic 等文件。

```shell
goctl api go -api api/user.api -dir .
```

输出：

```shell
etc/user-api.yaml exists, ignored generation
internal/config/config.go exists, ignored generation
user.go exists, ignored generation
internal/svc/servicecontext.go exists, ignored generation
internal/handler/createuserhandler.go exists, ignored generation
internal/handler/updateuserhandler.go exists, ignored generation
internal/handler/queryuserhandler.go exists, ignored generation
internal/logic/createuserlogic.go exists, ignored generation
internal/logic/updateuserlogic.go exists, ignored generation
internal/logic/queryuserlogic.go exists, ignored generation
Done.
```

### 生成 model 代码

生成 model 代码，生成 model 代码时，会自动生成 model 文件。

```shell
goctl model mysql ddl -src sql/user.sql -dir internal/model -c=false
```

输出：

```shell
usermodel_gen.go exists, ignored generation
usermodel.go exists, ignored generation
vars.go exists, ignored generation
Done.
```

### 手动完善代码

#### 完善 config 配置

在 `etc/user-api.yaml` 文件中添加数据库配置

```yaml
Mysql:
  DataSource: root:zero@tcp(127.0.0.1:3306)/zero?charset=utf8mb4&parseTime=True&loc=Local
```

在 `internal/config/config.go` 文件中添加数据库配置

```go
type Config struct {
	rest.RestConf
	Mysql struct {
		DataSource string
	}
}
```

### 修改 svc 上下文 

在 `internal/svc/servicecontext.go` 文件中添加数据库模型

```go
type servicecontext struct {
	Config config.Config
	UserModel model.UserModel
}

func Newservicecontext(c config.Config) *servicecontext {
	return &servicecontext{
		Config: c,
		UserModel: model.NewUserModel(sqlx.NewMysql(c.Mysql.DataSource)),
	}
}
```

### 修改 logic 逻辑

在 `internal/logic/createuserlogic.go` 文件中添加创建用户逻辑

```go
func (l *CreateUserLogic) CreateUser(req *types.UserCreateReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	user := model.User{
		Name: req.Name,
		Age:  req.Age,
		Email: req.Email,
	}
	res, err := l.svcCtx.UserModel.Insert(l.ctx, &user)
	if err != nil {
		return nil, err
	}
	user.Id, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}
	resp = &types.Response{
		Code: 0,
		Message: "success",
		Data: types.User{
			Id: user.Id,
			Name: user.Name,
			Age: user.Age,
			Email: user.Email,
			CreatedAt: user.CreatedAt.UnixMilli(),
			UpdatedAt: user.UpdatedAt.UnixMilli(),
		},
	}
	return resp, nil
}
```

在 `internal/logic/queryuserlogic.go` 文件中添加查询用户逻辑

```go
func (l *QueryUserLogic) QueryUser(req *types.UserQueryReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}
	resp = &types.Response{
		Code: 0,
		Message: "success",
		Data: types.User{
			Id: user.Id,
			Name: user.Name,
			Age: user.Age,
			Email: user.Email,
			CreatedAt: user.CreatedAt.UnixMilli(),
			UpdatedAt: user.UpdatedAt.UnixMilli(),
		},
	}
	return resp, nil
}
```

在 `internal/logic/updateuserlogic.go` 文件中添加更新用户逻辑

```go
func (l *UpdateUserLogic) UpdateUser(req *types.UserUpdateReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}
	user.Name = req.Name
	user.Age = req.Age
	user.Email = req.Email
	err = l.svcCtx.UserModel.Update(l.ctx, &user)
	if err != nil {
		return nil, err
	}
	resp = &types.Response{
		Code: 0,
		Message: "success",
		Data: types.User{
			Id: user.Id,
			Name: user.Name,
			Age: user.Age,
			Email: user.Email,
			CreatedAt: user.CreatedAt.UnixMilli(),
			UpdatedAt: user.UpdatedAt.UnixMilli(),
		},
	}
	return resp, nil
}
```

### 迁移数据库

Go-Zero 框架代码没有支持迁移数据库，迁移数据库需要手动执行。

### 运行项目

```shell
go run user.go
```

### 测试

```shell
curl -X POST http://localhost:8888/user -H "Content-Type: application/json" -d '{"name": "test", "age": 18, "email": "test@test.com"}'
```

---

## Go-Zero 缓存配置 (Redis)

### 添加 Redis 配置

编辑 `etc/user-api.yaml` 文件，添加 Redis 缓存配置

```yaml
Redis:
  Host: localhost:6379
  Type: node
  Pass: "ChinaSkills@"
```

编辑 `internal/config/config.go` 文件，添加 Redis 配置

```go
type Config struct {
	rest.RestConf
	Redis redis.RedisConf
}
```

### 修改 svc 上下文

编辑 `internal/svc/servicecontext.go` 文件，添加 Redis 配置

```go
type servicecontext struct {
	Config config.Config
	UserModel model.UserModel
	Redis     *redis.Redis
}

func Newservicecontext(c config.Config) *servicecontext {
	return &servicecontext{
		Config: c,
		UserModel: model.NewUserModel(sqlx.NewMysql(c.Mysql.DataSource)),
		Redis:     c.Redis.NewRedis(),
	}
}
```

### 修改 api 文件

编辑 `api/user.api` 文件，添加缓存接口

```api
// 缓存接口
type (
	UserCache {
		Id        int64  `json:"id"`
		Name      string `json:"name"`
		Age       int64  `json:"age"`
		Email     string `json:"email"`
		CreatedAt int64  `json:"created_at"`
		UpdatedAt int64  `json:"updated_at"`
	}
	UserCacheCreateReq {
		Id    int64  `json:"id"`
		Name  string `json:"name"`
		Age   int64  `json:"age"`
		Email string `json:"email"`
	}
	UserCacheUpdateReq {
		Id    int64  `json:"id"`
		Name  string `json:"name"`
		Age   int64  `json:"age"`
		Email string `json:"email"`
	}
	UserCacheQueryReq {
		Id int64 `json:"id"`
	}
	ResponseCache {
		Code    int64     `json:"code"`
		Message string    `json:"message"`
		Data    UserCache `json:"data"`
	}
)

service user-api {
	@handler getUserCache
	get /user/cache (UserCacheQueryReq) returns (ResponseCache)

	@handler setUserCache
	post /user/cache (UserCacheCreateReq) returns (ResponseCache)

	@handler updateUserCache
	put /user/cache (UserCacheUpdateReq) returns (ResponseCache)
}
```

### 生成 api 代码

```shell
goctl api go -api api/user.api -dir .
```

输出：

```shell
etc/user-api.yaml exists, ignored generation
internal/config/config.go exists, ignored generation
user.go exists, ignored generation
internal/svc/servicecontext.go exists, ignored generation
internal/handler/createuserhandler.go exists, ignored generation
internal/handler/updateuserhandler.go exists, ignored generation
internal/handler/queryuserhandler.go exists, ignored generation
internal/logic/createuserlogic.go exists, ignored generation
internal/logic/updateuserlogic.go exists, ignored generation
internal/logic/queryuserlogic.go exists, ignored generation
Done.
```

### 修改 logic 逻辑

在 `internal/logic/setusercachelogic.go` 文件中添加设置用户缓存逻辑

```go
func (l *SetUserCacheLogic) SetUserCache(req *types.UserCacheCreateReq) (resp *types.ResponseCache, err error) {
	// 生成缓存键
	cacheKey := fmt.Sprintf("user:cache:%d", req.Id)

	// 创建用户缓存对象
	userCache := &types.UserCache{
		Id:        req.Id,
		Name:      req.Name,
		Age:       req.Age,
		Email:     req.Email,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	// 将用户缓存对象序列化为 JSON
	cacheData, err := json.Marshal(userCache)
	if err != nil {
		logx.Errorf("Failed to marshal user cache data: %v", err)
		return &types.ResponseCache{
			Code:    500,
			Message: "内部服务器错误",
			Data:    types.UserCache{},
		}, nil
	}

	// 设置缓存到 Redis，设置 1 小时过期时间
	err = l.svcCtx.Redis.SetexCtx(l.ctx, cacheKey, string(cacheData), 3600)
	if err != nil {
		logx.Errorf("Failed to set user cache to redis: %v", err)
		return &types.ResponseCache{
			Code:    500,
			Message: "缓存设置失败",
			Data:    types.UserCache{},
		}, nil
	}

	logx.Infof("Successfully set user cache for user ID: %d", req.Id)

	return &types.ResponseCache{
		Code:    200,
		Message: "缓存设置成功",
		Data:    *userCache,
	}, nil
}
```

在 `internal/logic/updateusercachelogic.go` 文件中添加更新用户缓存逻辑

```go
func (l *UpdateUserCacheLogic) UpdateUserCache(req *types.UserCacheUpdateReq) (resp *types.ResponseCache, err error) {
	// 生成缓存键
	cacheKey := fmt.Sprintf("user:cache:%d", req.Id)

	// 先检查缓存是否存在
	existingData, err := l.svcCtx.Redis.GetCtx(l.ctx, cacheKey)
	if err != nil {
		logx.Errorf("Failed to check existing cache: %v", err)
		return &types.ResponseCache{
			Code:    500,
			Message: "检查缓存失败",
			Data:    types.UserCache{},
		}, nil
	}

	var userCache types.UserCache
	if existingData != "" {
		// 如果缓存存在，先解析现有数据
		err = json.Unmarshal([]byte(existingData), &userCache)
		if err != nil {
			logx.Errorf("Failed to unmarshal existing cache data: %v", err)
			return &types.ResponseCache{
				Code:    500,
				Message: "解析现有缓存失败",
				Data:    types.UserCache{},
			}, nil
		}
	} else {
		// 如果缓存不存在，创建新的缓存对象
		userCache.Id = req.Id
		userCache.CreatedAt = time.Now().Unix()
	}

	// 更新缓存数据
	userCache.Name = req.Name
	userCache.Age = req.Age
	userCache.Email = req.Email
	userCache.UpdatedAt = time.Now().Unix()

	// 序列化更新后的数据
	cacheData, err := json.Marshal(userCache)
	if err != nil {
		logx.Errorf("Failed to marshal updated cache data: %v", err)
		return &types.ResponseCache{
			Code:    500,
			Message: "序列化缓存数据失败",
			Data:    types.UserCache{},
		}, nil
	}

	// 更新 Redis 缓存，设置 1 小时过期时间
	err = l.svcCtx.Redis.SetexCtx(l.ctx, cacheKey, string(cacheData), 3600)
	if err != nil {
		logx.Errorf("Failed to update cache in redis: %v", err)
		return &types.ResponseCache{
			Code:    500,
			Message: "更新缓存失败",
			Data:    types.UserCache{},
		}, nil
	}

	logx.Infof("Successfully updated user cache for user ID: %d", req.Id)

	return &types.ResponseCache{
		Code:    200,
		Message: "更新缓存成功",
		Data:    userCache,
	}, nil
}
```

在 `internal/logic/getusercachelogic.go` 文件中添加获取用户缓存逻辑

```go
func (l *GetUserCacheLogic) GetUserCache(req *types.UserCacheQueryReq) (resp *types.ResponseCache, err error) {
	// 生成缓存键
	cacheKey := fmt.Sprintf("user:cache:%d", req.Id)

	// 从 Redis 获取缓存数据
	cacheData, err := l.svcCtx.Redis.GetCtx(l.ctx, cacheKey)
	if err != nil {
		logx.Errorf("Failed to get user cache from redis: %v", err)
		return &types.ResponseCache{
			Code:    500,
			Message: "获取缓存失败",
			Data:    types.UserCache{},
		}, nil
	}

	// 检查缓存是否存在
	if cacheData == "" {
		logx.Infof("User cache not found for user ID: %d", req.Id)
		return &types.ResponseCache{
			Code:    404,
			Message: "缓存不存在",
			Data:    types.UserCache{},
		}, nil
	}

	// 反序列化 JSON 数据
	var userCache types.UserCache
	err = json.Unmarshal([]byte(cacheData), &userCache)
	if err != nil {
		logx.Errorf("Failed to unmarshal user cache data: %v", err)
		return &types.ResponseCache{
			Code:    500,
			Message: "缓存数据解析失败",
			Data:    types.UserCache{},
		}, nil
	}

	logx.Infof("Successfully retrieved user cache for user ID: %d", req.Id)

	return &types.ResponseCache{
		Code:    200,
		Message: "获取缓存成功",
		Data:    userCache,
	}, nil
}
```

### 运行项目

```shell
go run user.go
```

### 测试

```shell
curl -X POST http://localhost:8888/user/cache -H "Content-Type: application/json" -d '{"id": 1, "name": "test", "age": 18, "email": "test@test.com"}'
curl -X GET http://localhost:8888/user/cache -H "Content-Type: application/json" -d '{"id": 1}'
curl -X PUT http://localhost:8888/user/cache -H "Content-Type: application/json" -d '{"id": 1, "name": "test2", "age": 20, "email": "test2@test.com"}'
```

---

## Go-Zero 框架日志处理

logx：是 Go-Zero 的核心日志库，负责实际的日志记录工作，支持多种输出方式和不同级别的日志记录。

logc：是对 logx 的封装，比 logx 多一个 ctx 参数，能够将日志信息与请求上下文关联，便于追踪请求处理过程。

两者在功能上是等效的，以下两种写法效果相同：

```go
logx.WithContext(ctx).Infof("This is a log message")
logc.Info(ctx, "This is a log message")
```

### 日志级别

Go-Zero 框架的日志级别分为以下几种：

- debug：调试级别，用于记录调试信息。
- info：信息级别，用于记录一般信息。
- error：错误级别，用于记录错误信息。
- severe：严重级别，用于记录严重错误信息。

可以通过配置或者代码设置日志级别：

```go
logx.SetLevel(logx.ErrorLevel) // 只记录 error 级别及以上的日志
logc.SetLevel(logc.ErrorLevel) // 只记录 error 级别及以上的日志
```

### 日志配置

Go-Zero 的日志系统通过 `LogConf` 结构体进行配置，可以配置日志的输出方式、输出路径、输出格式等。

```go
type LogConf struct {
	ServiceName 		string	`json:",optional"`	// 服务名称，用于区分不同服务，默认值为当前服务名
	Mode        		string	`json:",default=console,options=[console,file,volume]"`	// 日志输出模式，可选值为 console、file、volume(K8s 容器卷)
	Encoding    		string	`json:",default=json,options=[json,plain]"`	// 日志编码格式，可选值为 json、plain
	Path        		string	`json:",default=logs"`	// 日志文件路径，默认值为 logs，当 Mode 为 file 时有效
	Level       		string	`json:",default=info,options=[debug,info,warn,error,severe]"`	// 日志级别，可选值为 debug、info、warn、error、severe，默认值为 info
	Compress    		bool	`json:",optional"`	// 是否压缩日志文件，可选值为 true、false，默认值为 false
	KeepDays    		int		`json:",optional"`	// 日志文件保留天数，可选值为 0 表示不保留，当 Mode 为 file 时有效，默认值为 0
	StackCooldownMillis int		`json:",default=100"`	// 堆栈信息冷却时间，可选值为 0 表示不冷却，默认值为 100
	Rotation			string	`json:",default=daily,options=[daily,size]"`	// 日志文件轮转方式，可选值为 daily、size，daily 表示按天轮转，size 表示按大小轮转，默认值为 daily
	MaxSize				int		`json:",default=0"`	// 日志文件最大大小，可选值为 0 表示不限制，当 Mode 为 file 时有效，默认值为 0
}
```

---

## Go-Zero 框架内置中间件

Go-Zero 框架内置了多种中间件，可以方便地进行请求处理和响应处理。

```go
import (
	"github.com/tal-tech/go-zero/core/load"
	"github.com/tal-tech/go-zero/core/stat"
	"github.com/tal-tech/go-zero/rest"
	"github.com/tal-tech/go-zero/rest/handler"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func main() {
	server := rest.MustNewServer(rest.RestConf{
		Port: 8080,
	})
	server.Use(
		handler.MaxConns(1000),			// 限制最大连接数
		handler.BreakeHandler(),		// 熔断器
		handler.SheddingHandler(load.NewShedding()),		// 限流器
		handler.PrometheusHandler(),	// 监控指标
		handler.TracingHandler("user-api"),		// 链路追踪
		handler.RecoverHandler(),		// 异常恢复
		handler.StatHandler(stat.NewMetrics("user-api")),			// 统计信息
	)

	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/metrics",
		Handler: metricsHandler,
	})

	server.Start()
}
```