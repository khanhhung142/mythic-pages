# vietmyth-auditor

CLI tool tự động audit entries vietmyth.vn.

## Kiến trúc AI (mặc định)

| Job | Provider | Lý do |
|---|---|---|
| **Search** | Perplexity `sonar` | Grounded web search + citations |
| **Extract** | DeepSeek | Rẻ, JSON extraction |
| **Judge** | Claude Sonnet | Bám hard rules, ít false invented |

Judge **không nên** dùng Gemini/DeepSeek cho encyclopedia Việt/Hán Nôm — tool sẽ cảnh báo nếu bạn override.

## Setup

Set API keys:

| Job | Env var | Default model |
|---|---|---|
| Extract (DeepSeek) | `DEEPSEEK_API_KEY` | `deepseek-chat` |
| Judge (Claude) | `ANTHROPIC_API_KEY_API_PLATFORM` | `claude-sonnet-4-6` |
| Search | `PERPLEXITY_API_KEY` | `sonar` |

Optional overrides: `AUDIT_EXTRACT_LLM`, `AUDIT_EXTRACT_MODEL`, `AUDIT_JUDGE_LLM`, `AUDIT_JUDGE_MODEL`, `AUDIT_SEARCH`.

```bash
go mod tidy
go build -o vietmyth-auditor .
```

## Dùng

```bash
# Mặc định: extract=deepseek, judge=claude, search=perplexity
./vietmyth-auditor audit thanh-giong

# Đổi extract sang Gemini
./vietmyth-auditor audit thanh-giong --extract-llm gemini

# Chỉ định model cụ thể
./vietmyth-auditor audit thanh-giong --extract-llm gemini --extract-model gemini-2.0-flash
./vietmyth-auditor audit thanh-giong --judge-llm claude --judge-model claude-sonnet-4-6

# Batch thử nghiệm: một provider cho cả extract + judge (không khuyến nghị)
./vietmyth-auditor audit thanh-giong --llm deepseek

# Verbose
./vietmyth-auditor audit thanh-giong -v
```

## Output

File `audit/<entry>-audit.md` gồm:

- **Verdict:** PASS / REVISE / REJECT
- **Claim table:** từng claim → verified / wrong / suspicious / not_found
- **Wrong claims:** highlight riêng, fix ngay
- **Suspicious:** cần human review
- **Pattern issues:** AI writing patterns theo line number

## Verdict logic

| Condition | Verdict |
|---|---|
| wrong >= 3 hoặc wrong% >= 15% | REJECT |
| wrong >= 1 hoặc suspicious% >= 20% hoặc stance issue hoặc patterns >= 5 | REVISE |
| else | PASS |

## Passes

| Pass | Nội dung | Tool |
|---|---|---|
| 1 | Extract claims | Extract LLM (`--extract-llm`) |
| 3 | Verify từng claim | Perplexity + Judge LLM |
| 5 | Stance check (sovereignty, framing) | Regex |
| 6 | AI writing patterns | Regex |
| 7 | Source URLs + unlinked citations | HTTP + Regex |

## Interface

- `LLM.Complete(ctx, system, prompt, maxTokens)` — extract + judge
- `SearchProvider.Search(ctx, query)` — Perplexity

Factory: `NewRuntime(AIConfig{...})` trong `ai.go`.

## Giới hạn

- Max 60 claims/entry (tránh API cost)
- Verify chạy concurrent 5 goroutines, rate limit 200ms/claim
- Claim "not_found" không có nghĩa là sai, chỉ là không verify được — cần check tay
