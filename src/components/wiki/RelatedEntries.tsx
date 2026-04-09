const entries = [
  {
    name: "Sơn Tinh Thủy Tinh",
    preview: "Cuộc chiến giữa thần núi và thần nước tranh giành công chúa Mỵ Nương.",
    category: "Thần thoại",
  },
  {
    name: "Lạc Long Quân",
    preview: "Vị thần rồng khai sinh dân tộc, cha của trăm trứng trăm con.",
    category: "Thần",
  },
  {
    name: "Chử Đồng Tử",
    preview: "Chàng trai nghèo kết duyên cùng công chúa Tiên Dung, thành viên Tứ Bất Tử.",
    category: "Anh hùng",
  },
];

const categoryColors: Record<string, { bg: string; text: string }> = {
  "Thần thoại": { bg: "bg-pill-violet-bg", text: "text-pill-violet-text" },
  "Thần": { bg: "bg-pill-teal-bg", text: "text-pill-teal-text" },
  "Anh hùng": { bg: "bg-pill-coral-bg", text: "text-pill-coral-text" },
};

const RelatedEntries = () => (
  <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
    {entries.map((e) => {
      const colors = categoryColors[e.category] || { bg: "bg-pill-slate-bg", text: "text-pill-slate-text" };
      return (
        <a key={e.name} href="#" className="block rounded-lg overflow-hidden group" style={{ border: '0.5px solid hsl(220, 13%, 91%)' }}>
          <div className="bg-secondary h-32" />
          <div className="p-4">
            <div className="font-sans text-[15px] font-medium text-foreground mb-1 group-hover:text-link-blue">{e.name}</div>
            <p className="text-[13px] text-text-muted leading-snug mb-2">{e.preview}</p>
            <span className={`inline-block text-[11px] px-2 py-0.5 rounded ${colors.bg} ${colors.text}`}>{e.category}</span>
          </div>
        </a>
      );
    })}
  </div>
);

export default RelatedEntries;
