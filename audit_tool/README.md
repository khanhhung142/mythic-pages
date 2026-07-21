# vietmyth-auditor

CLI tool tự động audit entries vietmyth.vn.

## Setup

Set API keys for the providers you use:

| Provider | Env var | Default model |
|---|---|---|
| Claude (default LLM) | `ANTHROPIC_API_KEY_API_PLATFORM` | `claude-sonnet-4-6` |
| OpenAI | `OPENAI_API_KEY` | `gpt-4o` |
| DeepSeek | `DEEPSEEK_API_KEY` | `deepseek-chat` |
| Gemini | `GEMINI_API_KEY` or `GOOGLE_API_KEY` | `gemini-2.0-flash` |
| Perplexity (search) | `PERPLEXITY_API_KEY` | `sonar` |

Optional env: `AUDIT_LLM`, `AUDIT_LLM_MODEL`, `AUDIT_SEARCH`.

```bash
go mod tidy
go build -o vietmyth-auditor .
```

## Dùng

```bash
# Audit 1 entry theo slug
./vietmyth-auditor audit thanh-giong

# Dùng OpenAI làm LLM judge/extract
./vietmyth-auditor audit thanh-giong --llm openai --llm-model gpt-4o

# DeepSeek / Gemini
./vietmyth-auditor audit thanh-giong --llm deepseek
./vietmyth-auditor audit thanh-giong --llm gemini --llm-model gemini-2.0-flash

# Hoặc qua env
export AUDIT_LLM=deepseek
./vietmyth-auditor audit thanh-giong

# Đổi tên file trong audit/ (vẫn luôn ghi vào folder audit/)
./vietmyth-auditor audit thanh-giong -o thanh-giong-v2.md

# Verbose (xem progress)
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
| 1 | Extract claims | LLM (`--llm`) |
| 3 | Verify từng claim | Perplexity + LLM |
| 5 | Stance check (sovereignty, framing) | Regex |
| 6 | AI writing patterns | Regex |
| 7 | Source URLs + unlinked citations | HTTP + Regex |

Pass 7 kiểm tra:

- Mỗi mục `sources[]` có `url` không
- URL có reachable (HTTP), domain không nằm danh sách cấm (blog, wiki, …)
- Link markdown inline `[text](url)` trong thân bài
- Trích dẫn analysis thiếu link (tên + năm / *tựa sách* không bọc link)

## Kiến trúc AI

LLM và search tách qua interface — thêm provider mới bằng cách implement:

- `LLM.Complete(ctx, system, prompt, maxTokens)` — extract claims + judge verdicts
- `SearchProvider.Search(ctx, query)` — grounded search (hiện chỉ Perplexity)

Factory: `NewRuntime(AIConfig{...})` trong `ai.go`.

## Giới hạn

- Max 60 claims/entry (tránh API cost)
- Verify chạy concurrent 5 goroutines, rate limit 200ms/claim
- Perplexity `sonar` model — web search grounded
- Claim "not_found" không có nghĩa là sai, chỉ là không verify được — cần check tay
