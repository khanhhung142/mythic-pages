import { defineCollection, z } from 'astro:content';
import { glob } from 'astro/loaders';

const entries = defineCollection({
  loader: glob({ pattern: '**/*.md', base: './src/content/entries' }),
  schema: z.object({
    ten_vn: z.string(),
    ten_han: z.string().optional(),
    ten_khac: z.array(z.string()).optional(),
    ten_en: z.string().optional(),
    loai: z.string(),
    phan_loai_phu: z.array(z.string()).optional(),
    gioi_tinh: z.string().optional(),
    thoi_dai: z.string().optional(),
    nam_uoc_luong: z.number().optional(),
    vung_mien: z.string().optional(),
    dia_diem_chinh: z.array(z.string()).optional(),
    toa_do: z.array(z.number()).optional(),
    lien_quan: z.object({
      gia_dinh: z.array(z.string()).optional(),
      dong_minh: z.array(z.string()).optional(),
      ke_thu: z.array(z.string()).optional(),
      vat_pham: z.array(z.string()).optional(),
    }).optional(),
    nguon_co: z.array(z.object({
      ten: z.string(),
      chuong: z.string().optional(),
      ban_dich: z.string().optional(),
    })).optional(),
    chu_de: z.array(z.string()).optional(),
    do_pho_bien: z.number().default(1),
    trang_thai: z.string().default('draft'),
    tac_gia: z.string().optional(),
    cap_nhat: z.coerce.string().optional(),
  }),
});

export const collections = { entries };
