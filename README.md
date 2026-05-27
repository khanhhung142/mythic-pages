# Mythic Pages — Kho Truyền Thuyết & Truyện Dân Gian Việt Nam

**[vietmyth.vn](https://vietmyth.vn)** — Tổng hợp truyền thuyết, thần thoại và truyện dân gian Việt Nam theo dạng wiki mở.

---

## Giới thiệu

Mythic Pages là kho lưu trữ mở về thần linh, anh hùng, yêu quái, linh vật và truyện cổ của người Việt. Mỗi entry là một bài viết có cấu trúc — gồm tóm tắt, nguồn tham khảo, địa danh, chủ đề — được xây dựng để vừa dễ đọc vừa tra cứu được.

Dự án được build bằng [Astro](https://astro.build/) (static site), deploy trên Cloudflare Pages, hỗ trợ song ngữ Việt–Anh.

---

## Đóng góp

Có hai cách đóng góp:

### 1. Trực tiếp trên website

Truy cập [vietmyth.vn](https://vietmyth.vn) và dùng form **Đóng góp** có sẵn trên UI. Không cần biết Git hay Markdown.

### 2. Pull Request trên GitHub

Viết entry theo định dạng Markdown, tạo file trong `src/content/vi/entries/`, sau đó mở Pull Request.

Xem **[CONTRIBUTING.md](./CONTRIBUTING.md)** để biết cách viết entry đúng định dạng.

---

## Tech Stack

| Thành phần | Công nghệ |
|---|---|
| Framework | Astro 6 (Static SSG) |
| Ngôn ngữ | TypeScript 5, Markdown |
| Styling | Tailwind CSS 3 + CSS variables |
| Content | Astro Content Collections + Zod |
| Backend | Cloudflare Pages Functions (`functions/api/*`) |
| Deploy | Cloudflare Pages |

---

## Cấu trúc thư mục

```
src/
├── content/
│   ├── vi/entries/   # Entries tiếng Việt (canonical)
│   └── en/entries/   # Bản dịch tiếng Anh (optional)
├── pages/            # Routes (Astro dynamic routing)
├── components/       # Astro / UI components
└── i18n/             # Config đa ngôn ngữ
functions/
└── api/              # Cloudflare Pages Functions
docs/                 # Tài liệu nội bộ cho dev / AI agents
```

---

## Dev

```bash
npm install
npm run dev      # localhost:4321
npm run build    # build static
npm run preview  # preview build
```

**Deploy (Cloudflare Pages):**
- Build command: `npm run build`
- Output directory: `dist`
- Functions: `functions/api/*` → `/api/*`

---

## Giấy phép

Nội dung (entries, truyện) thuộc cộng đồng đóng góp — sử dụng tự do với attribution.  
Code thuộc [MIT License](./LICENSE).
