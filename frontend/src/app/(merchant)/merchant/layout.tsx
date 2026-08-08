export default function MerchantLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="merchant-layout">
      {/* Future: Merchant header/sidebar */}
      <main>{children}</main>
    </div>
  );
}
