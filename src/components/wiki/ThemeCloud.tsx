const themes = ["chiến tranh", "bảo vệ tổ quốc", "thần đồng", "hy sinh"];

const ThemeCloud = () => (
  <div className="flex flex-wrap gap-[6px]">
    {themes.map((t) => (
      <span key={t} className="bg-pill-slate-bg text-pill-slate-text text-xs px-2.5 py-1 rounded-md">
        {t}
      </span>
    ))}
  </div>
);

export default ThemeCloud;
