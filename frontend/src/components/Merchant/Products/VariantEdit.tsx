"use client";

import { useMerchantProductEditor } from "@/app/context/MerchantProductEditorContext";
import type {
  MerchantApiProduct,
  MerchantApiProductVariant,
} from "@/types/merchantProduct";
import HydrateMerchantProductEditor from "./HydrateMerchantProductEditor";
import VariantForm from "./VariantForm";

type VariantEditProps = {
  productId: number;
  variantId: number;
  fallbackProduct: MerchantApiProduct;
  fallbackVariants: MerchantApiProductVariant[];
  fallbackVariant: MerchantApiProductVariant;
};

export default function VariantEdit({
  productId,
  variantId,
  fallbackProduct,
  fallbackVariants,
  fallbackVariant,
}: VariantEditProps) {
  const { product, variants, loadedProductId } = useMerchantProductEditor();
  const contextMatch = loadedProductId === productId && product != null;
  const contextVariant = contextMatch
    ? variants?.find((item) => item.id === variantId)
    : undefined;
  const resolvedVariant = contextVariant ?? fallbackVariant;
  const needsHydrate = !contextMatch || !contextVariant;
  const productSlug =
    typeof fallbackProduct.slug === "string"
      ? fallbackProduct.slug.trim()
      : "";

  return (
    <>
      {needsHydrate ? (
        <HydrateMerchantProductEditor
          product={fallbackProduct}
          variants={fallbackVariants}
        />
      ) : null}
      <VariantForm
        mode="edit"
        productId={productId}
        productSlug={productSlug}
        variant={resolvedVariant}
      />
    </>
  );
}
