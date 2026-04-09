import { useState } from "react";
import SidebarCard from "@/components/wiki/SidebarCard";
import InfoTable from "@/components/wiki/InfoTable";
import RelationshipSection from "@/components/wiki/RelationshipSection";
import ThemeCloud from "@/components/wiki/ThemeCloud";
import RelatedEntries from "@/components/wiki/RelatedEntries";

const Index = () => {
  const [lang, setLang] = useState<"vi" | "en">("vi");

  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-[1200px] mx-auto px-4 sm:px-6 py-6">
        {/* Breadcrumb + Language switcher */}
        <div className="flex items-center justify-between mb-4">
          <nav className="text-[13px] text-text-muted">
            <span className="hover:text-foreground cursor-pointer">Wiki</span>
            <span className="mx-1.5">›</span>
            <span className="hover:text-foreground cursor-pointer">Anh hùng</span>
            <span className="mx-1.5">›</span>
            <span className="text-foreground">Thánh Gióng</span>
          </nav>
          <div className="flex rounded-full overflow-hidden" style={{ border: '1px solid hsl(220, 13%, 91%)' }}>
            <button
              onClick={() => setLang("vi")}
              className={`px-3 py-1 text-xs font-medium transition-colors ${lang === "vi" ? "bg-foreground text-background" : "text-text-muted hover:text-foreground"}`}
            >
              VI
            </button>
            <button
              onClick={() => setLang("en")}
              className={`px-3 py-1 text-xs font-medium transition-colors ${lang === "en" ? "bg-foreground text-background" : "text-text-muted hover:text-foreground"}`}
            >
              EN
            </button>
          </div>
        </div>

        {/* Category pills */}
        <div className="flex flex-wrap gap-2 mb-4">
          <span className="text-[12px] px-2.5 py-1 rounded-full bg-pill-coral-bg text-pill-coral-text">Anh hùng</span>
          <span className="text-[12px] px-2.5 py-1 rounded-full bg-pill-teal-bg text-pill-teal-text">Bắc Bộ</span>
          <span className="text-[12px] px-2.5 py-1 rounded-full bg-pill-amber-bg text-pill-amber-text">Hùng Vương VI</span>
        </div>

        {/* H1 */}
        <h1 className="text-[32px] font-medium text-foreground leading-tight mb-1">Thánh Gióng</h1>

        {/* Alternative names */}
        <p className="font-serif italic text-[14px] text-text-muted mb-0.5">
          Phù Đổng Thiên Vương · Sóc Thiên Vương · Ông Gióng
        </p>
        <p className="text-[13px] text-text-light mb-8">
          聖揀 · Saint Gióng
        </p>

        {/* Two-column layout */}
        <div className="flex flex-col lg:flex-row gap-8">
          {/* LEFT — main content */}
          <article className="flex-1 min-w-0">
            {/* Summary */}
            <div className="bg-secondary rounded-lg py-4 px-5 mb-8">
              <p className="font-serif italic text-foreground leading-relaxed">
                Cậu bé ba tuổi không nói không cười, nghe tin giặc Ân xâm lược liền vươn vai thành tráng sĩ khổng lồ, cưỡi ngựa sắt, cầm roi sắt đánh tan quân thù, rồi bay về trời.
              </p>
            </div>

            {/* Câu chuyện */}
            <h2 className="text-[18px] font-medium text-foreground mb-4">Câu chuyện</h2>
            <div className="font-serif text-[16px] leading-[1.7] text-foreground space-y-4 mb-6">
              <p>
                Vào đời Hùng Vương thứ sáu, tại làng Phù Đổng, bộ Vũ Ninh (nay thuộc Gia Lâm, Hà Nội), có hai vợ chồng nông dân hiền lành, tuổi đã cao mà chưa có con. Một hôm bà vợ ra đồng, thấy một vết chân rất lớn in trên mặt đất, bèn ướm thử chân mình vào. Về nhà bà thụ thai, mười hai tháng sau sinh được một cậu bé khôi ngô. Lạ thay, đứa trẻ lên ba tuổi vẫn không biết nói, biết cười, đặt đâu nằm đấy.
              </p>
              <p>
                Bấy giờ giặc Ân từ phương Bắc kéo sang xâm lược, thế giặc rất mạnh. Vua Hùng lo lắng, sai sứ giả đi khắp nơi tìm người tài giỏi cứu nước. Khi sứ giả đi qua làng Phù Đổng, cậu bé bỗng cất tiếng nói, nhờ mẹ gọi sứ giả vào. Cậu bảo sứ giả về tâu vua sắm cho một con ngựa sắt, một cái roi sắt và một áo giáp sắt, rồi cậu sẽ đánh tan giặc. Sứ giả vừa mừng vừa ngạc nhiên, vội về triều tâu lại.
              </p>
              <p>
                Từ ngày gặp sứ giả, cậu bé lớn nhanh như thổi, cơm ăn mấy cũng không no, áo mặc mấy cũng không vừa. Dân làng góp gạo, góp vải nuôi cậu. Khi ngựa sắt, roi sắt, áo giáp sắt được đem đến, cậu vươn vai biến thành một tráng sĩ khổng lồ, mình cao hơn trượng, oai phong lẫm liệt. Tráng sĩ mặc giáp, cầm roi, nhảy lên ngựa. Ngựa sắt phun lửa, phi thẳng vào trận giặc. Roi sắt gãy, tráng sĩ nhổ từng bụi tre đằng ngà bên đường quật vào quân thù. Giặc Ân tan tác.
              </p>
              <p>
                Đánh tan giặc, tráng sĩ không trở về triều nhận thưởng mà thúc ngựa lên đỉnh núi Sóc Sơn, cởi áo giáp, rồi cả người lẫn ngựa bay lên trời. Vua Hùng nhớ ơn, phong là Phù Đổng Thiên Vương, lập đền thờ tại quê nhà. Ngày nay tại Sóc Sơn vẫn còn dấu tre đằng ngà bị cháy vàng và những ao hồ — tương truyền là dấu chân ngựa sắt.
              </p>
            </div>

            {/* Dị bản */}
            <h3 className="text-[16px] font-medium text-foreground mb-3">Dị bản</h3>
            <div className="font-serif text-[16px] leading-[1.7] text-foreground space-y-4 mb-8">
              <p>
                Trong tập <em>Việt Điện U Linh Tập</em> của Lý Tế Xuyên (thế kỷ XIV), nhân vật này được gọi là Sóc Thiên Vương và câu chuyện có thêm chi tiết về việc vua Hùng lập miếu thờ ngay tại chân núi Sóc. Một số dị bản vùng Bắc Ninh kể rằng ngựa sắt do chính dân làng Phù Đổng rèn đúc, nhấn mạnh tinh thần cộng đồng trong cuộc kháng chiến.
              </p>
            </div>

            {/* Ý nghĩa văn hóa */}
            <h2 className="text-[18px] font-medium text-foreground mb-4">Ý nghĩa văn hóa</h2>
            <div className="font-serif text-[16px] leading-[1.7] text-foreground space-y-4 mb-8">
              <p>
                Thánh Gióng là một trong Tứ Bất Tử của tín ngưỡng dân gian Việt Nam, cùng với Tản Viên Sơn Thánh, Chử Đồng Tử và Liễu Hạnh. Bốn vị thánh đại diện cho bốn khát vọng lớn của dân tộc: chống ngoại xâm (Thánh Gióng), chinh phục thiên nhiên (Tản Viên), xây dựng cuộc sống ấm no (Chử Đồng Tử) và tự do cá nhân (Liễu Hạnh).
              </p>
              <p>
                Hình tượng cậu bé ba tuổi vươn vai thành người khổng lồ mang ý nghĩa sâu sắc về sức mạnh tiềm ẩn của dân tộc: khi Tổ quốc lâm nguy, ngay cả những người nhỏ bé nhất cũng có thể vùng lên. Việc tráng sĩ bay về trời sau chiến thắng — không nhận công danh phú quý — thể hiện tinh thần vô tư, hy sinh vì đại nghĩa, một giá trị cốt lõi trong truyền thống Việt Nam.
              </p>
            </div>

            {/* Trong văn hóa hiện đại */}
            <h2 className="text-[18px] font-medium text-foreground mb-4">Trong văn hóa hiện đại</h2>
            <div className="font-serif text-[16px] leading-[1.7] text-foreground space-y-4 mb-8">
              <p>
                Hội Gióng tại đền Phù Đổng và đền Sóc được UNESCO công nhận là Di sản văn hóa phi vật thể đại diện của nhân loại vào năm 2010. Lễ hội diễn ra vào mùng 9 tháng 4 âm lịch hàng năm, thu hút hàng vạn người tham gia với các nghi thức tái hiện trận đánh giặc Ân. Hình tượng Thánh Gióng cũng xuất hiện trong sách giáo khoa, phim hoạt hình, truyện tranh và nhiều tác phẩm nghệ thuật đương đại Việt Nam.
              </p>
            </div>

            {/* Divider + Sources */}
            <div className="border-t border-border pt-6 mb-8">
              <div className="text-[11px] font-semibold tracking-[0.08em] uppercase text-text-muted mb-3">Nguồn tham khảo</div>
              <ol className="list-decimal list-inside text-[13px] text-text-muted space-y-1.5">
                <li>Trần Thế Pháp, <em className="text-foreground">Lĩnh Nam chích quái</em>, thế kỷ XV.</li>
                <li>Lý Tế Xuyên, <em className="text-foreground">Việt Điện U Linh Tập</em>, thế kỷ XIV.</li>
                <li>Nguyễn Đổng Chi, <em className="text-foreground">Kho tàng truyện cổ tích Việt Nam</em>, NXB Giáo dục, 1957.</li>
                <li>Trần Quốc Vượng, "Thánh Gióng và biểu tượng anh hùng chống ngoại xâm," <em className="text-foreground">Tạp chí Nghiên cứu Lịch sử</em>, 1985.</li>
              </ol>
            </div>

            {/* Related entries */}
            <h2 className="text-[18px] font-medium text-foreground mb-4">Đọc thêm</h2>
            <RelatedEntries />
          </article>

          {/* RIGHT — sidebar */}
          <aside className="w-full lg:w-[280px] lg:flex-shrink-0">
            <div className="lg:sticky lg:top-6 space-y-4">
              <SidebarCard label="Thông tin">
                <InfoTable />
              </SidebarCard>

              <SidebarCard label="Liên quan">
                <RelationshipSection />
              </SidebarCard>

              <SidebarCard label="Chủ đề">
                <ThemeCloud />
              </SidebarCard>
            </div>
          </aside>
        </div>
      </div>
    </div>
  );
};

export default Index;
