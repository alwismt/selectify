"use client";

import { useMerchantProductEditor } from "@/app/context/MerchantProductEditorContext";
import type { MerchantApiProduct } from "@/types/merchantProduct";
import type { SiteCurrency } from "@/types/siteConfig";
import HydrateMerchantProductEditor from "./HydrateMerchantProductEditor";
import VariantForm from "./VariantForm";

type VariantNewProps = {
  productId: number;
  fallbackProduct: MerchantApiProduct;
  currency: SiteCurrency | null;
};

export default function VariantNew({
  productId,
  fallbackProduct,
  currency: _currency,
}: VariantNewProps) {
  const { product, variants, loadedProductId } = useMerchantProductEditor();
  const needsHydrate = loadedProductId !== productId || product == null;
  const productSlug =
    typeof fallbackProduct.slug === "string"
      ? fallbackProduct.slug.trim()
      : "";

  return (
    <>
      {needsHydrate ? (
        <HydrateMerchantProductEditor
          product={fallbackProduct}
          variants={variants}
        />
      ) : null}
      <VariantForm
        mode="create"
        productId={productId}
        productSlug={productSlug}
      />
    </>
  );
}
