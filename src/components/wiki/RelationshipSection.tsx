interface RelGroup {
  label: string;
  items: { name: string; linked?: boolean }[];
}

const groups: RelGroup[] = [
  { label: "GIA ĐÌNH", items: [{ name: "Mẹ Gióng", linked: true }] },
  { label: "ĐỒNG MINH", items: [{ name: "Hùng Vương thứ 6", linked: true }] },
  { label: "KẺ THÙ", items: [{ name: "Giặc Ân" }] },
  {
    label: "VẬT PHẨM",
    items: [
      { name: "Ngựa sắt", linked: true },
      { name: "Roi sắt", linked: true },
      { name: "Tre đằng ngà", linked: true },
    ],
  },
];

const RelationshipSection = () => (
  <div className="space-y-3">
    {groups.map((g) => (
      <div key={g.label}>
        <div className="text-[10px] font-semibold tracking-[0.08em] uppercase text-text-light mb-1">{g.label}</div>
        <div className="text-[13px] leading-relaxed">
          {g.items.map((item, i) => (
            <span key={item.name}>
              {i > 0 && <span className="text-text-light"> · </span>}
              {item.linked ? (
                <a href="#" className="text-link-blue hover:underline">{item.name}</a>
              ) : (
                <span className="text-foreground">{item.name}</span>
              )}
            </span>
          ))}
        </div>
      </div>
    ))}
  </div>
);

export default RelationshipSection;
