import { MerchantProductEditorProvider } from "@/app/context/MerchantProductEditorContext";

export default function MerchantProductSlugLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <MerchantProductEditorProvider>{children}</MerchantProductEditorProvider>
  );
}
