import MerchantPageShell from "@/components/MerchantDashboard/MerchantPageShell";
import MerchantProductsList from "@/components/Merchant/Products/MerchantProductsList";
import { getServerMerchantProducts } from "@/lib/api/getServerMerchantProducts";
import { getServerSiteCurrency } from "@/lib/api/getServerSiteConfig";
import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Products | Selectify Seller",
  description: "Manage your store products",
};

export default async function MerchantProductsPage() {
  const [products, currency] = await Promise.all([
    getServerMerchantProducts(),
    getServerSiteCurrency(),
  ]);

  return (
    <main>
      <MerchantPageShell>
        <MerchantProductsList products={products} currency={currency} />
      </MerchantPageShell>
    </main>
  );
}
