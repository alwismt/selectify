import MerchantPageShell from "@/components/MerchantDashboard/MerchantPageShell";
import ProductNew from "@/components/Merchant/Products/ProductNew";
import { getServerCategories } from "@/lib/api/getServerCategories";
import { getServerSiteCurrency } from "@/lib/api/getServerSiteConfig";
import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Add Product | Selectify Seller",
  description: "Add a new product to your store",
};

export default async function MerchantProductNewPage() {
  const [currency, categories] = await Promise.all([
    getServerSiteCurrency(),
    getServerCategories(),
  ]);

  return (
    <main>
      <MerchantPageShell>
        <ProductNew currency={currency} categories={categories ?? []} />
      </MerchantPageShell>
    </main>
  );
}
