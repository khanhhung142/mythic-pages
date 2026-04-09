const rows = [
  { label: "Loại", value: "Anh hùng" },
  { label: "Giới tính", value: "Nam" },
  { label: "Thời đại", value: "Hùng Vương VI" },
  { label: "Vùng", value: "Bắc Ninh" },
  { label: "Địa điểm", value: "Phù Đổng", isLink: true },
  { label: "Nhóm", value: "Tứ Bất Tử" },
];

const InfoTable = () => (
  <table className="w-full text-[13px]">
    <tbody>
      {rows.map((row) => (
        <tr key={row.label} className="border-b border-border last:border-b-0">
          <td className="py-2 text-text-muted pr-3">{row.label}</td>
          <td className="py-2 text-right text-foreground">
            {row.isLink ? (
              <a href="#" className="text-link-blue hover:underline">{row.value}</a>
            ) : row.value}
          </td>
        </tr>
      ))}
    </tbody>
  </table>
);

export default InfoTable;
