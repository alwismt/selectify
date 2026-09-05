import MerchantPageShell from "@/components/MerchantDashboard/MerchantPageShell";
import MerchantProductDetails from "@/components/Merchant/Products/MerchantProductDetails";
import {
  getServerMerchantProduct,
  getServerMerchantProductVariants,
} from "@/lib/api/getServerMerchantProducts";
import { getServerSiteCurrency } from "@/lib/api/getServerSiteConfig";
import {
  canonicalProductPath,
  merchantProductHref,
  parseProductIdFromPath,
} from "@/lib/productPath";
import { Metadata } from "next";
import { notFound, redirect } from "next/navigation";

type MerchantProductPageProps = {
  params: Promise<{ slug: string }>;
};

export async function generateMetadata({
  params,
}: MerchantProductPageProps): Promise<Metadata> {
  const { slug } = await params;
  const productId = parseProductIdFromPath(slug);
  if (productId == null) {
    return { title: "Product Not Found | Selectify Seller" };
  }

  const product = await getServerMerchantProduct(productId);
  if (!product) {
    return { title: "Product Not Found | Selectify Seller" };
  }

  return {
    title: `${product.name} | Selectify Seller`,
    description: product.description ?? `Manage ${product.name}`,
  };
}

export default async function MerchantProductPage({
  params,
}: MerchantProductPageProps) {
  const { slug } = await params;
  const productId = parseProductIdFromPath(slug);
  if (productId == null) {
    notFound();
  }

  const [product, currency] = await Promise.all([
    getServerMerchantProduct(productId),
    getServerSiteCurrency(),
  ]);
  if (!product) {
    notFound();
  }

  const productSlug =
    typeof product.slug === "string" ? product.slug.trim() : "";
  if (!productSlug) {
    notFound();
  }

  const canonical = canonicalProductPath(product.productId, productSlug);
  if (slug !== canonical) {
    redirect(merchantProductHref(product.productId, productSlug));
  }

  const variants = await getServerMerchantProductVariants(product.productId);

  return (
    <main>
      <MerchantPageShell>
        <MerchantProductDetails
          product={product}
          variants={variants}
          currency={currency}
        />
      </MerchantPageShell>
    </main>
  );
}
