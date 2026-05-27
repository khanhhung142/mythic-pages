# Hướng dẫn viết Entry

Entry là một file Markdown mô tả một nhân vật, truyện cổ, địa danh, yêu quái, linh vật hoặc hiện vật trong kho truyền thuyết và dân gian Việt Nam.

---

## 1. Tạo file

Tạo file mới tại: `src/content/vi/entries/<tên-entry>.md`

Tên file dùng **kebab-case**, khớp với nội dung `id` trong frontmatter.

**Ví dụ:** Thánh Gióng → `thanh-giong.md`

---

## 2. Cấu trúc frontmatter

Mỗi entry bắt đầu bằng khối YAML giữa hai dấu `---`. Dưới đây là toàn bộ các trường có thể dùng:

```yaml
---
id: ten-entry-kebab-case           # khớp với tên file (bắt buộc)
name_vi: Tên tiếng Việt            # bắt buộc
name_han: 漢字                     # chữ Hán nếu có
aliases:                           # tên gọi khác, tên địa phương, biệt danh
  - Tên khác 1
  - Tên khác 2
name_en: English name              # tên tiếng Anh nếu có
category: than-linh                # xem bảng danh mục bên dưới (bắt buộc)
subcategories:                     # phân loại chi tiết hơn (tùy chọn)
  - tu-bat-tu
  - anh-hung-chong-ngoai-xam
gender: nam                        # nam | nu | khong-xac-dinh
era: Hùng Vương thứ 6              # thời đại, ghi tự do
year_approx: -1718                 # năm gần đúng (âm = TCN), hoặc null nếu không xác định
region: bac                        # bac | trung | nam
locations:                         # địa danh liên quan
  - Phù Đổng, Gia Lâm, Hà Nội
  - Núi Sóc Sơn
coordinates:                       # [lat, lng] của địa điểm chính
  - 21.0833
  - 105.9500
relations:
  family:
    - Âu Cơ (mẹ)
    - Lạc Long Quân (cha)
  allies:
    - Hùng Vương thứ 6
  enemies:
    - Giặc Ân
  artifacts:
    - Ngựa sắt
    - Roi sắt
  related_sites:
    - Đền Phù Đổng
sources:                           # nguồn tham khảo
  - title: Lĩnh Nam Chích Quái
    author: Trần Thế Pháp
    chapter: Truyện Đổng Thiên Vương
    edition: Đinh Gia Khánh – Nguyễn Ngọc San, NXB Văn hóa 1960
summary: |
  Tóm tắt ngắn về nhân vật / truyện, 1–3 câu.
themes:                            # chủ đề, dùng slug kebab-case
  - chien-tranh
  - bao-ve-to-quoc
  - anh-hung-hoa-than
popularity: 5                      # 1–10, ảnh hưởng thứ tự hiển thị
status: published                  # published = hiện trên site
updated_at: 2026-05-27
---
```

**Các trường bắt buộc:** `id`, `name_vi`, `category`, `status`

---

## 3. Danh mục (category)

| Slug | Nhãn tiếng Việt | Dùng cho |
|---|---|---|
| `than-linh` | Thần linh | Thần, thánh, tiên, phật |
| `anh-hung` | Anh hùng | Anh hùng dân tộc, nhân vật lịch sử huyền thoại |
| `yeu-quai` | Yêu quái | Ma, quỷ, tinh, chồn, cáo, rắn thành tinh |
| `linh-vat` | Linh vật | Rồng, phượng hoàng, kỳ lân, rùa thần |
| `dia-danh` | Địa danh | Núi, sông, hồ, vùng đất có tích |
| `vat-pham` | Vật phẩm | Bảo vật, linh khí, vũ khí thần |
| `le-hoi` | Lễ hội | Tục lệ, lễ hội dân gian |
| `tich-co` | Truyện cổ | Truyện cổ tích, truyện ký, truyện kỳ ảo |
| `nhan-vat` | Nhân vật | Nhân vật lịch sử được huyền thoại hóa |

---

## 4. Phần thân (Markdown body)

Viết sau khối frontmatter. Dùng heading `##` để chia phần, và dùng _nghiêng_ cho phần tên/tiêu đề phụ theo quy ước của site:

```markdown
## Cốt _truyện_

Kể lại nội dung câu chuyện theo thứ tự tự nhiên...

## Phân _tích_

Ý nghĩa biểu tượng, motif dân gian, so sánh dị bản...

## Địa _danh_

Vị trí thực tế ngày nay, di tích liên quan...

## Tín _ngưỡng_

Tục thờ, lễ hội, sắc phong còn lưu...
```

**Bảng so sánh dị bản** (dùng GFM table):

```markdown
| Yếu tố | Bản A | Bản B |
|---|---|---|
| Tên nhân vật | Thôi Lượng | Thôi Lạng |
| Ngày mất | Mồng 3 tháng Giêng | 13 tháng Giêng |
```

---

## 5. Ví dụ entry đơn giản

```markdown
---
id: ngu-tinh
name_vi: Ngư Tinh
name_han: 魚精
aliases:
  - Cá tinh
  - Xác Cáo
name_en: Fish Demon
category: yeu-quai
era: Thời Lạc Long Quân
region: bac
locations:
  - Biển Đông
  - Vùng biển ngoài khơi đất Việt
relations:
  enemies:
    - Lạc Long Quân
sources:
  - title: Lĩnh Nam Chích Quái
    author: Trần Thế Pháp
    chapter: Ngư Tinh truyện
summary: |
  Cá tinh khổng lồ quấy nhiễu dân Việt, bị Lạc Long Quân giết bằng đuốc lửa. Xác cá hóa thành đảo.
themes:
  - tru-yeu-diet-ma
  - lap-quoc
popularity: 3
status: published
updated_at: 2026-05-27
---

## Cốt _truyện_

Thuở xưa, ngoài biển Đông có con cá tinh khổng lồ...

## Ý _nghĩa_

Motif trừ yêu giải thích nguồn gốc địa lý của các đảo ven biển...
```

---

## 6. Quy tắc chung

- **Một file = một entry.** Không gộp nhiều nhân vật vào cùng một file.
- **`status: published`** thì entry mới hiện trên site. Dùng `draft` nếu chưa xong.
- **`popularity`** từ 1–10. Nhân vật trung tâm như Lạc Long Quân, Thánh Gióng thì 8–10; nhân vật phụ thì 1–3.
- **Nguồn tham khảo:** ghi ít nhất một nguồn trong `sources`. Nếu không có nguồn thành văn, ghi rõ "truyền miệng" hoặc "thần tích địa phương".
- **`relations`:** viết tự do bằng tiếng Việt, có thể ghi thêm chú thích trong ngoặc đơn, ví dụ `Âu Cơ (mẹ)`.
- **Không bịa đặt.** Nếu không chắc, ghi rõ "theo một số bản kể" hoặc "không rõ nguồn gốc".

---

## 7. Gửi đóng góp

### Qua Pull Request

1. Fork repo
2. Tạo branch mới: `git checkout -b add/ten-entry`
3. Tạo file `src/content/vi/entries/ten-entry.md`
4. Commit và mở Pull Request về `main`

### Qua website

Truy cập [vietmyth.vn](https://vietmyth.vn) → nút **Đóng góp** → điền form. Không cần tài khoản GitHub.
