# EntDomain

[![Go Reference](https://pkg.go.dev/badge/github.com/githonllc/entdomain.svg)](https://pkg.go.dev/github.com/githonllc/entdomain)
[![Go Report Card](https://goreportcard.com/badge/github.com/githonllc/entdomain)](https://goreportcard.com/report/github.com/githonllc/entdomain)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

一个 [Ent](https://entgo.io) 扩展，从带注解的 schema 自动生成 DDD（领域驱动设计）层代码 — 领域模型、仓储接口和服务接口。

## 特性

- **注解驱动** — 使用简洁的构建器标记字段作用域（`DefaultField`、`InputOnlyField`、`OutputOnlyField` 等）
- **完整的 DDD 层** — 生成领域模型、创建/更新请求、响应、查询参数、仓储接口和服务接口
- **类型安全的 ID 系统** — `ID` 接口支持 string 和 int64 底层类型
- **通用基础服务** — `BaseGenericDomainService` 开箱即用的 CRUD、分页和搜索
- **深拷贝与部分更新** — 生成的 `Clone()` 和 `ApplyToDomainModel()` 方法

## 环境要求

- Go 1.23+
- [Ent](https://entgo.io) v0.14+

## 安装

```bash
go get github.com/githonllc/entdomain
```

## 配置

在 `entc.go` 中注册扩展：

```go
//go:build ignore

package main

import (
    "log"

    "entgo.io/ent/entc"
    "entgo.io/ent/entc/gen"
    "github.com/githonllc/entdomain"
)

func main() {
    ext := entdomain.NewExtensionWithOptions(
        entdomain.WithOutputDir("./ent"),
        entdomain.WithPackageName("ent"),
    )

    if err := entc.Generate("./schema", &gen.Config{}, entc.Extensions(ext)); err != nil {
        log.Fatal(err)
    }
}
```

然后运行：

```bash
go generate ./...
```

## 注解构建器

### 基础构建器

```go
entdomain.DefaultField()                      // 所有作用域：创建、更新、查询、响应
entdomain.InputOnlyField()                    // 仅创建和更新（如密码）
entdomain.OutputOnlyField()                   // 仅查询和响应（如时间戳）
entdomain.CreateOnlyField()                   // 创建 + 查询 + 响应（创建后不可变）
entdomain.IdField()                           // 响应 + 查询，只读
entdomain.AuditLogField()                     // 响应 + 查询，只读
entdomain.DomainFieldWithScopes(scopes...)    // 自定义作用域组合
```

### 流式构建器 API

```go
field.String("email").
    Annotations(
        entdomain.DefaultField().
            WithRequired(entdomain.ScopeCreate).
            WithDescription("Email address").
            WithFormat("email").
            AsSearchable().
            AsFilterable().
            AsUniqueLookup(),
    )
```

### 可用方法

| 方法 | 说明 |
|------|------|
| `WithRequired(scope)` | 在指定作用域中标记为必填 |
| `WithDescription(desc)` | 设置字段描述 |
| `WithExample(val)` | 设置示例值 |
| `WithFormat(fmt)` | 设置格式（email, date-time, uuid 等） |
| `WithPattern(regex)` | 设置验证正则 |
| `WithRange(min, max *float64)` | 设置数值范围（nil = 无限制） |
| `WithLength(min, max *int)` | 设置字符串长度（nil = 无限制） |
| `AsSearchable()` | 启用文本搜索 |
| `AsFilterable()` | 启用过滤 |
| `AsSortable()` | 启用排序 |
| `AsUniqueLookup()` | 生成 `FindByX` 方法 |
| `AsRangeLookup()` | 生成 `FindByXRange` 方法 |
| `AsSensitive()` | 标记为敏感字段 |
| `AsReadOnly()` | 标记为只读 |

## Schema 示例

```go
package schema

import (
    "time"

    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "github.com/githonllc/entdomain"
)

type User struct {
    ent.Schema
}

func (User) Fields() []ent.Field {
    return []ent.Field{
        field.Int64("id").
            Unique().
            Immutable().
            Annotations(entdomain.IdField()),

        field.String("name").
            NotEmpty().
            Annotations(
                entdomain.DefaultField().
                    WithRequired(entdomain.ScopeCreate).
                    AsSearchable(),
            ),

        field.String("email").
            Optional().
            Annotations(
                entdomain.DefaultField().
                    WithFormat("email").
                    AsUniqueLookup(),
            ),

        field.Time("created_at").
            Default(time.Now).
            Immutable().
            Annotations(
                entdomain.OutputOnlyField().
                    AsRangeLookup(),
            ),
    }
}
```

## 生成的代码

为每个带注解的 schema 生成三个文件：

### 领域模型 (`user_domain_model.go`)

```go
type UserDomainModel struct { ... }     // 完整领域模型
type UserCreateRequest struct { ... }   // 创建请求 DTO
type UserUpdateRequest struct { ... }   // 更新请求 DTO（指针字段用于部分更新）
type UserResponse struct { ... }        // 响应 DTO
type UserQueryParams struct { ... }     // 搜索/过滤参数
type UserListResponse struct { ... }    // 分页响应
```

生成的转换方法：
- `ToDomainModel()` — 请求 → 领域模型
- `ApplyToDomainModel()` — 在现有模型上应用部分更新
- `ToResponse()` — 领域模型 → 响应（可选字段使用 `Ptr`/`PtrOrNil`/`PtrTimeOrNil` 辅助函数）
- `Clone()` — 深拷贝

### 仓储 (`user_domain_repository.go`)

```go
type UserRepository interface {
    Create(ctx, model) (*UserDomainModel, error)
    GetByID(ctx, id) (*UserDomainModel, error)
    Update(ctx, model) (*UserDomainModel, error)
    Delete(ctx, id) error
    List(ctx, req) ([]*UserDomainModel, int, error)
    Search(ctx, req) ([]*UserDomainModel, int, error)
    FindByEmail(ctx, email) (*UserDomainModel, error)  // 来自 AsUniqueLookup
    FindByCreatedAtRange(ctx, start, end) ([]*UserDomainModel, error)  // 来自 AsRangeLookup
    // ... 批量操作、计数、存在检查
}
```

### 服务 (`user_domain_service.go`)

服务模板只生成 **接口** — 不生成实现。你可以组合 `BaseGenericDomainService` 加上自定义业务逻辑：

```go
// 生成的接口（精简，无样板实现）：
type UserService interface {
    Create(ctx, model) (*UserDomainModel, error)
    GetByID(ctx, id) (*UserDomainModel, error)
    Update(ctx, model) (*UserDomainModel, error)
    Delete(ctx, id) error
    List(ctx, req) ([]*UserDomainModel, int, error)
    Search(ctx, req) ([]*UserDomainModel, int, error)
    Count(ctx, req) (int, error)
    Exists(ctx, id) (bool, error)
    CreateBatch(ctx, models) ([]*UserDomainModel, error)
    UpdateBatch(ctx, models) ([]*UserDomainModel, error)
    DeleteBatch(ctx, ids) error
    FindByEmail(ctx, email) (*UserDomainModel, error)       // 来自 AsUniqueLookup
    FindByCreatedAtRange(ctx, start, end) ([]*UserDomainModel, error)  // 来自 AsRangeLookup
}
```

同时生成 `DefaultUserService` — 嵌入它，只覆盖需要自定义逻辑的方法：

```go
type userServiceImpl struct {
    ent.DefaultUserService  // 嵌入默认 CRUD 委托
}

// 覆盖 Create 添加自定义验证：
func (s *userServiceImpl) Create(ctx context.Context, model *ent.UserDomainModel) (*ent.UserDomainModel, error) {
    if model.Email == "" {
        return nil, entdomain.ErrValidation
    }
    return s.DefaultUserService.Create(ctx, model)
}
```

## 类型化错误

仓储层将 Ent 错误包装为标准哨兵值：

```go
var (
    entdomain.ErrNotFound      // 实体未找到
    entdomain.ErrAlreadyExists // 唯一约束冲突
    entdomain.ErrValidation    // 验证失败
)

// 使用 errors.Is 或辅助函数检查：
if entdomain.IsNotFound(err) { ... }
if entdomain.IsAlreadyExists(err) { ... }
```

## 查询参数绑定

生成的 `QueryParams` 结构体包含 `form` tag，支持 gin 自动绑定：

```go
var params ent.UserQueryParams
if err := c.ShouldBindQuery(&params); err != nil { ... }
// params.Query, params.Page, params.Size 自动解析
```

## 游标分页

对于大数据集，使用基于游标（keyset）的分页替代偏移量分页：

```go
// 第一页：
items, pageInfo, err := repo.ListByCursor(ctx, &entdomain.ListRequest{
    Size: 20, SortBy: "created_at", Order: "desc",
})
// pageInfo.HasNextPage == true
// pageInfo.EndCursor == "eyJpZ..." (不透明)

// 下一页：
items, pageInfo, err = repo.ListByCursor(ctx, &entdomain.ListRequest{
    Cursor: pageInfo.EndCursor,
    Size: 20, SortBy: "created_at", Order: "desc",
})
```

游标分页是 O(log n) 复杂度，而偏移量分页是 O(n) — 它通过索引直接定位。
偏移量分页（`List()`）仍然可用于需要跳页的场景。

## ID 系统

```go
// 创建 ID
id := entdomain.NewIDFromInt64(12345)
id := entdomain.NewIDFromString("user-abc")

// 使用 ID
str := id.String()       // "12345"
num, err := id.Int64()   // 12345, nil
zero := id.IsZero()      // false
```

## 通用基础服务

用于处理层的 CRUD，自动验证和转换：

```go
service := entdomain.NewBaseGenericDomainService[
    *UserDomainModel,
    *UserCreateRequest,
    *UserUpdateRequest,
    *UserResponse,
    *UserListResponse,
    *UserQueryParams,
](repo, converters)

resp, err := service.Create(ctx, createReq)
resp, err := service.GetByID(ctx, id)
resp, err := service.Update(ctx, id, updateReq)
err := service.Delete(ctx, id)
```

## 字段作用域

作用域控制处理层 DTO 中包含哪些字段。它们**不会**限制服务层或仓储层的访问。

| 作用域 | 说明 |
|--------|------|
| `ScopeCreate` | 字段出现在创建请求中 |
| `ScopeUpdate` | 字段出现在更新请求中 |
| `ScopeQuery` | 字段出现在查询参数中 |
| `ScopeResponse` | 字段出现在响应中 |

## 扩展选项

```go
entdomain.WithOutputDir("./ent")              // 输出目录
entdomain.WithPackageName("ent")              // Go 包名
entdomain.WithRepository(true)                // 生成仓储（默认：true）
entdomain.WithService(true)                   // 生成服务（默认：true）
entdomain.WithEntDomainPackage("custom/path") // 覆盖 entdomain 导入路径
```

## 贡献

请参阅 [CONTRIBUTING.md](CONTRIBUTING.md) 了解开发配置和指南。

## 许可证

[MIT](LICENSE)
