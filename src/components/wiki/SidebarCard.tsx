interface SidebarCardProps {
  label: string;
  children: React.ReactNode;
}

const SidebarCard = ({ label, children }: SidebarCardProps) => (
  <div className="bg-card rounded-lg p-4 px-[18px]" style={{ border: '0.5px solid hsl(220, 13%, 91%)' }}>
    <div className="text-[11px] font-semibold tracking-[0.08em] uppercase text-text-muted mb-3">{label}</div>
    {children}
  </div>
);

export default SidebarCard;
