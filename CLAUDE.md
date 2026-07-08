# vietmyth.vn

Bilingual (VI/EN) encyclopedia: Vietnamese mythology, folklore, historical
figures, places. Model: yokai.com.

## Repo layout

```
.src/content/
  vi/entries/*.md        # entry tiếng Việt (YAML frontmatter + body)
  en/entries/*.md        # entry tiếng Anh, cùng slug
research/
  [slug].md              # research report đã approve (input cho writing)
prompts/                 # prompt templates cho từng bước
.claude/skills/
  vietmyth-entry/        # skill chính — MỌI entry theo skill này
  humanizer/             # pass quét cuối
```

(Nếu đường dẫn content khác — Astro/Hugo/Next — giữ nguyên quy ước
vi/en song song cùng slug.)

## Quy tắc bất biến

1. **Mọi entry tuân theo skill `vietmyth-entry` v2.** Đọc SKILL.md trước
   khi viết bất kỳ entry nào. Không viết entry khi chưa có research report
   approved trong `research/[slug].md`.
2. **Calibration bắt buộc:** trước khi viết entry mới, đọc 2 entry
   `status: published` gần nhất CÙNG TYPE trong `content/vi/entries/`
   để khớp giọng. Nếu chưa có entry cùng type, dùng entry Bạch Trĩ.
3. **Hai file mỗi entry:** `content/vi/entries/[slug].md` và
   `content/en/entries/[slug].md`. Bản EN viết lại từ research report,
   KHÔNG dịch từ bản VI.
4. **Stance rules (skill §3) là non-negotiable:** Hoàng Sa Trường Sa là
   của Việt Nam; East Sea, không bao giờ "South China Sea" trong prose EN;
   tôn trọng nhân vật lịch sử; accuracy > pride.
5. **Zero em dash.** Sau khi viết, grep cả 2 file: `—`, `–`, ` -- `.
   Kết quả phải rỗng (trừ range số trong YAML/table như `938–944`).
6. **Self-audit 4 pass (skill §7) chạy TRƯỚC khi báo xong.** Pass 4
   stance check output từng câu nghi vấn, không output "đã kiểm tra, ổn".
7. **Hán tự:** verify từng ký tự bằng web search khi không chắc. Đây là
   lỗi phổ biến nhất.
8. **Citation:** mọi claim học thuật có tên + tác phẩm + năm + trang.
   Không chắc nguồn → bỏ claim, không bịa.
9. **humanizer:** dùng skill có sẵn đã cài đặt vào máy

## Không được làm

- Không sửa entry `status: published` khi chưa được yêu cầu rõ.
- Không viết cả 2 phase (research + writing) trong 1 lần chạy. Research
  report phải được người duyệt trước.
- Không thêm chi tiết vào Chuyện kể mà research report không có nguồn.
- Không dùng từ trong BANNED lists (skill §6) kể cả trong summary
  frontmatter và commit message... commit message viết thường, mô tả
  thay đổi, không marketing.

## Workflow chuẩn

```
1. research/[slug].md tồn tại + approved?
   Chưa → dừng, báo cần chạy research phase trước (prompts/01-research.md)
2. Đọc SKILL.md + 2 entry calibration cùng type
3. Viết bản VI → bản EN (từ report, không dịch)
4. Grep em dash cả 2 file
5. Self-audit 4 pass, output bảng vi phạm + fix
6. Báo xong kèm: type, word count Chuyện kể, số nguồn cite, audit summary
```

## Frontmatter status flow

`draft` → (người review) → `published`. Claude chỉ tạo `draft`.