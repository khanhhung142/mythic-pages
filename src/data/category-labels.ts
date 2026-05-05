import type { Locale } from '../i18n/config';

export const CATEGORY_LABELS: Record<string, Record<Locale, string>> = {
  "than-linh": { vi: "Thần linh", en: "Deities" },
  "anh-hung": { vi: "Anh hùng", en: "Heroes" },
  "yeu-quai": { vi: "Yêu quái", en: "Demons" },
  "linh-vat": { vi: "Linh vật", en: "Sacred Beasts" },
  "dia-danh": { vi: "Địa danh", en: "Places" },
  "vat-pham": { vi: "Vật phẩm", en: "Artifacts" },
  "le-hoi": { vi: "Lễ hội", en: "Festivals" },
  "tich-co": { vi: "Truyện cổ", en: "Folktales" },
  "nhan-vat": { vi: "Nhân vật", en: "Figures" },
};

export const CATEGORY_SLUGS = Object.keys(CATEGORY_LABELS);

export function getCategoryLabel(slug: string, locale: Locale): string {
  return CATEGORY_LABELS[slug]?.[locale] ?? CATEGORY_LABELS[slug]?.vi ?? slug;
}
