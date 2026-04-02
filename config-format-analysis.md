# Config Format Analysis for Go Projects

## JSON (Standard Library)

**Pros:**
- Built into Go standard library (`encoding/json`)
- No external dependencies
- Universally understood and supported
- Good for simple, flat configurations
- Excellent for programmatic generation

**Cons:**
- No comments (workarounds with `//` in Go struct tags)
- Less human-readable for nested structures
- Verbose syntax
- Strict formatting (commas, quotes)

**Example:**
```json
{
  "agent": {
    "endpoint": "http://localhost:8080/chat/completions",
    "timeout": "30s",
    "model": "custom-model"
  }
}
```

## YAML (Third-party)

**Pros:**
- Human-readable with clean syntax
- Supports comments
- Better for complex, nested configurations
- Widely used in DevOps/Kubernetes ecosystems
- Multi-line strings without escaping

**Cons:**
- Requires external dependency (`gopkg.in/yaml.v3`, `sigs.k8s.io/yaml`)
- Whitespace-sensitive (can be error-prone)
- More complex parsing
- Potential security issues with anchors/aliases

**Example:**
```yaml
agent:
  endpoint: "http://localhost:8080/chat/completions"
  timeout: "30s"
  model: "custom-model"
  # Optional parameters
  temperature: 0.7
  max_tokens: 1000
```

## TOML (Emerging Standard)

**Pros:**
- Explicit, less ambiguous than YAML
- Good mix of readability and structure
- Used by major Go projects (Hugo, Caddy)
- Comments supported
- Type-aware (dates, times, arrays)

**Cons:**
- Less common than JSON/YAML
- Still requires external dependency
- Can be verbose for deeply nested data

**Example:**
```toml
[agent]
endpoint = "http://localhost:8080/chat/completions"
timeout = "30s"
model = "custom-model"

# Optional parameters
temperature = 0.7
max_tokens = 1000
```

## Viper (Multi-format Solution)

**Common Pattern:**
```go
import "github.com/spf13/viper"
```
- Supports JSON, YAML, TOML, HCL, env vars, flags
- Widely used in CLI applications
- Single interface for multiple config sources
- Automatic type conversion

## Recommendations for Your TUI Agent Harness

### **Option 1: JSON + Struct Tags (Recommended for MVP)**
```go
type Config struct {
    Agent struct {
        Endpoint string `json:"endpoint"`
        Timeout  string `json:"timeout"`
        Model    string `json:"model"`
    } `json:"agent"`
}
```
**Why:** Simplest, no dependencies, standard library, easy to extend.

### **Option 2: YAML with Viper (Professional Setup)**
```go
// With viper, supports multiple formats
viper.SetConfigName("config")
viper.SetConfigType("yaml")
viper.AddConfigPath(".")
```
**Why:** Flexible, supports multiple formats, comments, industry standard for complex configs.

### **Option 3: TOML (Clean Middle Ground)**
```go
import "github.com/pelletier/go-toml/v2"
```
**Why:** Good balance, type-safe, used by respected Go projects.

## Decision Matrix

| Criteria | JSON | YAML | TOML |
|----------|------|------|------|
| **No Dependencies** | ✅ | ❌ | ❌ |
| **Human Readable** | ❌ | ✅ | ✅ |
| **Comments** | ❌ | ✅ | ✅ |
| **Std Library** | ✅ | ❌ | ❌ |
| **Simple Syntax** | ❌ | ✅ | ✅ |
| **Type Safety** | ❌ | ❌ | ✅ |
| **K8s Ecosystem** | ❌ | ✅ | ❌ |

## Recommendation

**For your MVP: Use JSON** with the standard library.
- Simplest implementation
- No external dependencies
- Easy to parse and validate
- Can switch to YAML/Viper later if needed

**If you anticipate complex config needs: Use YAML with Viper**
- More flexible for future features
- Supports multiple config sources
- Industry standard for agent/DevOps tools

**Implementation path:**
1. Start with JSON for MVP (faster, simpler)
2. If config becomes complex, switch to Viper with YAML support
3. Use struct tags compatible with both formats