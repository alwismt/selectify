import MerchantPageShell from "@/components/MerchantDashboard/MerchantPageShell";
import VariantEdit from "@/components/Merchant/Products/VariantEdit";
import {
  getServerMerchantProduct,
  getServerMerchantProductVariants,
} from "@/lib/api/getServerMerchantProducts";
import {
  canonicalProductPath,
  merchantProductHref,
  parseProductIdFromPath,
} from "@/lib/productPath";
import { Metadata } from "next";
import { notFound, redirect } from "next/navigation";

export const metadata: Metadata = {
  title: "Edit Variant | Selectify Seller",
  description: "Edit a product variant",
};

type MerchantVariantEditPageProps = {
  params: Promise<{ slug: string; variant_id: string }>;
};

export default async function MerchantVariantEditPage({
  params,
}: MerchantVariantEditPageProps) {
  const { slug, variant_id } = await params;
  const productId = parseProductIdFromPath(slug);
  if (productId == null) {
    notFound();
  }

  const variantId = Number(variant_id);
  if (!Number.isInteger(variantId) || variantId <= 0) {
    notFound();
  }

  const fallbackProduct = await getServerMerchantProduct(productId);
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
      merchantProductHref(
        fallbackProduct.productId,
        productSlug,
        `/variants/${variantId}/edit`
      )
    );
  }

  const fallbackVariants = await getServerMerchantProductVariants(
    fallbackProduct.productId
  );
  if (fallbackVariants === null) {
    notFound();
  }

  const fallbackVariant = fallbackVariants.find(
    (item) => item.id === variantId
  );
  if (!fallbackVariant) {
    notFound();
  }

  return (
    <main>
      <MerchantPageShell>
        <VariantEdit
          productId={fallbackProduct.productId}
          variantId={variantId}
          fallbackProduct={fallbackProduct}
          fallbackVariants={fallbackVariants}
          fallbackVariant={fallbackVariant}
        />
      </MerchantPageShell>
    </main>
  );
}
