import MerchantSidebar from "./MerchantSidebar";

type MerchantPageShellProps = {
  children: React.ReactNode;
};

/**
 * Shared merchant dashboard chrome: sidebar + main content area.
 * Server component; composes the client MerchantSidebar.
 */
export default function MerchantPageShell({ children }: MerchantPageShellProps) {
  return (
    <section className="overflow-hidden pt-[100px] sm:pt-[100px] lg:pt-[110px] xl:pt-[120px] pb-20 bg-gray-2">
      <div className="mx-auto max-w-[1500px] w-full px-4 sm:px-8 xl:px-0">
        <div className="flex flex-col gap-7.5 xl:flex-row">
          <MerchantSidebar />
          <main className="min-w-0 flex-1">{children}</main>
        </div>
      </div>
    </section>
  );
}
