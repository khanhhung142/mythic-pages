# vietmyth-auditor

CLI tool tự động audit entries vietmyth.vn.

## Setup

```bash
export ANTHROPIC_API_KEY_API_PLATFORM=sk-ant-...
export PERPLEXITY_API_KEY=pplx-...
```

```bash
go mod tidy
go build -o vietmyth-auditor .
```

## Dùng

```bash
# Audit 1 entry theo slug, tự đọc từ src/content/vi/entries/<slug>.md
# Output luôn vào audit/<slug>-audit.md
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
| 1 | Extract claims | Claude |
| 3 | Verify từng claim | Perplexity + Claude |
| 5 | Stance check (sovereignty, framing) | Regex |
| 6 | AI writing patterns | Regex |

## Giới hạn

- Max 60 claims/entry (tránh API cost)
- Verify chạy concurrent 5 goroutines, rate limit 200ms/claim
- Perplexity `sonar` model — web search grounded
- Claim "not_found" không có nghĩa là sai, chỉ là không verify được — cần check tay
