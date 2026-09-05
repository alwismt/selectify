import MerchantPageShell from "@/components/MerchantDashboard/MerchantPageShell";
import ProductEdit from "@/components/Merchant/Products/ProductEdit";
import { getServerMerchantProduct } from "@/lib/api/getServerMerchantProducts";
import { getServerSiteCurrency } from "@/lib/api/getServerSiteConfig";
import {
  canonicalProductPath,
  merchantProductHref,
  parseProductIdFromPath,
} from "@/lib/productPath";
import { Metadata } from "next";
import { notFound, redirect } from "next/navigation";

export const metadata: Metadata = {
  title: "Edit Product | Selectify Seller",
  description: "Edit your store product",
};

type MerchantProductEditPageProps = {
  params: Promise<{ slug: string }>;
};

export default async function MerchantProductEditPage({
  params,
}: MerchantProductEditPageProps) {
  const { slug } = await params;
  const productId = parseProductIdFromPath(slug);
  if (productId == null) {
    notFound();
  }

  const [fallbackProduct, currency] = await Promise.all([
    getServerMerchantProduct(productId),
    getServerSiteCurrency(),
  ]);
  if (!fallbackProduct) {
    notFound();
  }

  const productSlug =
    typeof fallbackProduct.slug === "string"
      ? fallbackProduct.slug.trim()
      : "";
  if (!productSlug) {
    notFound();
  }

  const canonical = canonicalProductPath(
    fallbackProduct.productId,
    productSlug
  );
  if (slug !== canonical) {
    redirect(
      merchantProductHref(fallbackProduct.productId, productSlug, "/edit")
    );
  }

  return (
    <main>
      <MerchantPageShell>
        <ProductEdit
          productId={fallbackProduct.productId}
          fallbackProduct={fallbackProduct}
          currency={currency}
        />
      </MerchantPageShell>
    </main>
  );
}
