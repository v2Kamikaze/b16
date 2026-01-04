# B16 - Sistema de Autenticação e Autorização

Sistema de autenticação e autorização em Go com suporte a múltiplos métodos de autenticação (Basic Auth e Token JWT) e políticas de autorização flexíveis.

## 📋 Índice

- [Arquitetura](#arquitetura)
- [Conceitos Fundamentais](#conceitos-fundamentais)
- [Como Usar](#como-usar)
  - [Criando um Principal](#criando-um-principal)
  - [Criando um Manager](#criando-um-manager)
  - [Criando uma Policy](#criando-uma-policy)
  - [Criando um Handler](#criando-um-handler)
- [Exemplos Práticos](#exemplos-práticos)

## 🏗️ Arquitetura

O projeto segue uma arquitetura em camadas com separação clara de responsabilidades:

```
b16/
├── internal/
│   ├── auth/                    # Módulo de autenticação e autorização
│   │   ├── manager/             # Implementações de AuthManager
│   │   │   ├── basic_auth_manager.go
│   │   │   └── token_auth_manager.go
│   │   ├── middleware/          # Middlewares HTTP
│   │   │   ├── with_auth.go     # Middleware de autenticação
│   │   │   └── with_policy.go   # Middleware de autorização
│   │   ├── policy/              # Implementações de Policy
│   │   │   ├── role_policy.go   # Policy baseada em roles
│   │   │   ├── any_policy.go    # Policy OR (qualquer uma)
│   │   │   └── composite_policy.go # Policy AND (todas)
│   │   └── errors.go            # Erros do módulo auth
│   ├── domain/                  # Interfaces e modelos de domínio
│   │   ├── manager.go           # Interface AuthManager e Principal
│   │   ├── policy.go            # Interface Policy
│   │   ├── security.go          # Interfaces de segurança
│   │   └── models.go            # Modelos de domínio
│   ├── security/                # Implementações de segurança
│   │   ├── token_issuer.go      # Emissor de tokens JWT
│   │   └── password_hasher.go   # Hash de senhas
│   └── config/                  # Configuração da aplicação
│       └── config.go
└── main.go                      # Aplicação principal
```

### Fluxo de Autenticação e Autorização

```
Request → WithAuth Middleware → AuthManager.Authenticate() → Principal
                                                                    ↓
Response ← Handler ← WithPolicy Middleware ← Policy.Check() ←──────┘
```

1. **WithAuth Middleware**: Intercepta a requisição e chama o `AuthManager` para autenticar
2. **AuthManager**: Valida as credenciais e retorna um `Principal`
3. **WithPolicy Middleware**: (Opcional) Valida se o `Principal` tem permissão
4. **Policy**: Verifica se o `Principal` atende aos critérios de autorização
5. **Handler**: Processa a requisição autenticada e autorizada

## 🎯 Conceitos Fundamentais

### Principal

Um `Principal` representa a identidade autenticada do usuário. É uma interface genérica que encapsula os dados do usuário autenticado.

```go
type Principal[T any] interface {
    Principal() T
}
```

### AuthManager

Um `AuthManager` é responsável por autenticar requisições HTTP e retornar um `Principal`. Cada método de autenticação (Basic Auth, JWT, etc.) tem sua própria implementação.

```go
type AuthManager[T any] interface {
    Authenticate(req *http.Request) (Principal[T], error)
}
```

### Policy

Uma `Policy` define regras de autorização que verificam se um `Principal` tem permissão para acessar um recurso.

```go
type Policy[T any] interface {
    Check(credentials Principal[T]) error
}
```

### Handler

Um `Handler` é uma função que processa requisições HTTP autenticadas. Recebe o `Principal` como parâmetro.

```go
type AuthHandler[T any] func(w http.ResponseWriter, r *http.Request, credentials Principal[T])
```

## 🚀 Como Usar

### Criando um Principal

Um `Principal` é uma struct que implementa a interface `domain.Principal[T]`. O tipo `T` é tipicamente a própria struct do Principal.

**Exemplo:**

```go
package manager

import "github.com/v2code/b16/internal/domain"

type MyPrincipal struct {
    UserID   string
    Username string
    Email    string
}

// Implementa domain.Principal[*MyPrincipal]
func (p *MyPrincipal) Principal() *MyPrincipal {
    return p
}
```

### Criando um Manager

Um `Manager` implementa a interface `domain.AuthManager[T]` e é responsável por extrair e validar credenciais da requisição HTTP.

**Exemplo - Custom Manager:**

```go
package manager

import (
    "net/http"
    "github.com/v2code/b16/internal/auth"
    "github.com/v2code/b16/internal/domain"
)

type MyAuthManager struct {
    // Campos necessários para autenticação
    apiKey string
}

func NewMyAuthManager(apiKey string) domain.AuthManager[*MyPrincipal] {
    return &MyAuthManager{apiKey: apiKey}
}

func (m *MyAuthManager) Authenticate(req *http.Request) (domain.Principal[*MyPrincipal], error) {
    // Extrai credenciais do header, query params, etc.
    apiKey := req.Header.Get("X-API-Key")
    if apiKey == "" {
        return nil, auth.ErrUnauthorized
    }

    // Valida as credenciais
    if apiKey != m.apiKey {
        return nil, auth.ErrUnauthorized
    }

    // Retorna o Principal autenticado
    return &MyPrincipal{
        UserID:   "123",
        Username: "john",
        Email:    "john@example.com",
    }, nil
}
```

**Managers Disponíveis:**

- **BasicAuthManager**: Autenticação via HTTP Basic Auth
  ```go
  basicAuthManager := manager.NewBasicAuthManager("username", "password")
  ```

- **TokenAuthManager**: Autenticação via JWT Bearer Token
  ```go
  jwtIssuer := security.NewJwtIssuer(security.JwtIssuerParams{
      SecretKey: []byte("secret"),
      ExpireAt:  time.Hour * 2,
      Issuer:    "b16",
  })
  tokenAuthManager := manager.NewTokenAuthManager(jwtIssuer)
  ```

### Criando uma Policy

Uma `Policy` implementa a interface `domain.Policy[T]` e verifica se um `Principal` atende aos critérios de autorização.

**Exemplo - Custom Policy:**

```go
package policy

import (
    "github.com/v2code/b16/internal/auth"
    "github.com/v2code/b16/internal/auth/manager"
    "github.com/v2code/b16/internal/domain"
)

type RequireEmailDomain struct {
    domain string
}

func NewRequireEmailDomainPolicy(domain string) domain.Policy[*manager.TokenPrincipal] {
    return &RequireEmailDomain{domain: domain}
}

func (p *RequireEmailDomain) Check(credentials domain.Principal[*manager.TokenPrincipal]) error {
    email := credentials.Principal().Email
    // Verifica se o email termina com o domínio requerido
    if !strings.HasSuffix(email, "@"+p.domain) {
        return auth.ErrForbidden
    }
    return nil
}
```

**Policies Disponíveis:**

- **RequireRolePolicy**: Requer que o usuário tenha uma ou mais roles
  ```go
  policy := policy.RequireRolePolicy("ADMIN", "USER")
  ```

- **NewAnyPolicy**: Retorna sucesso se qualquer uma das policies passar (OR)
  ```go
  policy := policy.NewAnyPolicy(
      policy.RequireRolePolicy("ADMIN"),
      policy.RequireRolePolicy("SUPER_ADMIN"),
  )
  ```

- **NewCompositePolicy**: Retorna sucesso apenas se todas as policies passarem (AND)
  ```go
  policy := policy.NewCompositePolicy(
      policy.RequireRolePolicy("ADMIN"),
      NewRequireEmailDomainPolicy("company.com"),
  )
  ```

### Criando um Handler

Um `Handler` é uma função que processa requisições autenticadas. Recebe o `Principal` como terceiro parâmetro.

**Exemplo:**

```go
func MyHandler(w http.ResponseWriter, r *http.Request, credentials domain.Principal[*manager.MyPrincipal]) {
    principal := credentials.Principal()

    // Acessa os dados do principal
    userID := principal.UserID
    username := principal.Username

    // Processa a requisição
    fmt.Fprintf(w, "Hello %s (ID: %s)", username, userID)
}
```

## 📝 Exemplos Práticos

### Exemplo 1: Endpoint com Basic Auth

```go
func BasicAuthHandler(w http.ResponseWriter, r *http.Request, credentials domain.Principal[*manager.BasicAuthPrincipal]) {
    principal := credentials.Principal()
    fmt.Fprintf(w, "Hello %s", principal.Username)
}

func main() {
    basicAuthManager := manager.NewBasicAuthManager("admin", "password")

    mux := http.NewServeMux()
    mux.HandleFunc(
        "GET /basic-auth",
        middleware.WithAuth(basicAuthManager, BasicAuthHandler),
    )

    http.ListenAndServe(":8000", mux)
}
```

**Uso:**
```bash
curl -u admin:password http://localhost:8000/basic-auth
```

### Exemplo 2: Endpoint com Token Auth e Policy de Role

```go
func TokenAuthHandler(w http.ResponseWriter, r *http.Request, credentials domain.Principal[*manager.TokenPrincipal]) {
    principal := credentials.Principal()
    fmt.Fprintf(w, "Hello %s", principal.Email)
}

func main() {
    jwtIssuer := security.NewJwtIssuer(security.JwtIssuerParams{
        SecretKey: []byte("secret"),
        ExpireAt:  time.Hour * 2,
        Issuer:    "b16",
    })

    tokenAuthManager := manager.NewTokenAuthManager(jwtIssuer)

    mux := http.NewServeMux()
    mux.HandleFunc(
        "GET /token-auth",
        middleware.WithAuth(
            tokenAuthManager,
            middleware.WithPolicy(
                TokenAuthHandler,
                policy.RequireRolePolicy("ADMIN", "USER"),
            ),
        ),
    )

    http.ListenAndServe(":8000", mux)
}
```

**Uso:**
```bash
# Primeiro, gere um token JWT (exemplo)
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

curl -H "Authorization: Bearer $TOKEN" http://localhost:8000/token-auth
```

### Exemplo 3: Endpoint com Policy Composta (AND)

```go
mux.HandleFunc(
    "GET /admin-only",
    middleware.WithAuth(
        tokenAuthManager,
        middleware.WithPolicy(
            AdminHandler,
            policy.NewCompositePolicy(
                policy.RequireRolePolicy("ADMIN"),
                NewRequireEmailDomainPolicy("company.com"),
            ),
        ),
    ),
)
```

### Exemplo 4: Endpoint com Policy Any (OR)

```go
mux.HandleFunc(
    "GET /special-access",
    middleware.WithAuth(
        tokenAuthManager,
        middleware.WithPolicy(
            SpecialHandler,
            policy.NewAnyPolicy(
                policy.RequireRolePolicy("ADMIN"),
                policy.RequireRolePolicy("SUPER_ADMIN"),
            ),
        ),
    ),
)
```

### Exemplo 5: Criando um Manager Customizado

```go
package manager

import (
    "net/http"
    "github.com/v2code/b16/internal/auth"
    "github.com/v2code/b16/internal/domain"
)

type APIKeyPrincipal struct {
    APIKey string
    UserID string
}

func (p *APIKeyPrincipal) Principal() *APIKeyPrincipal {
    return p
}

type APIKeyManager struct {
    validKeys map[string]string // API Key -> User ID
}

func NewAPIKeyManager(validKeys map[string]string) domain.AuthManager[*APIKeyPrincipal] {
    return &APIKeyManager{validKeys: validKeys}
}

func (m *APIKeyManager) Authenticate(req *http.Request) (domain.Principal[*APIKeyPrincipal], error) {
    apiKey := req.Header.Get("X-API-Key")
    if apiKey == "" {
        return nil, auth.ErrUnauthorized
    }

    userID, ok := m.validKeys[apiKey]
    if !ok {
        return nil, auth.ErrUnauthorized
    }

    return &APIKeyPrincipal{
        APIKey: apiKey,
        UserID: userID,
    }, nil
}
```

**Uso:**
```go
apiKeyManager := manager.NewAPIKeyManager(map[string]string{
    "key123": "user1",
    "key456": "user2",
})

mux.HandleFunc(
    "GET /api/protected",
    middleware.WithAuth(apiKeyManager, APIKeyHandler),
)
```
