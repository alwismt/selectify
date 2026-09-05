"use client";

import { useEffect } from "react";
import { useMerchantProductEditor } from "@/app/context/MerchantProductEditorContext";
import type {
  MerchantApiProduct,
  MerchantApiProductVariant,
} from "@/types/merchantProduct";

type HydrateMerchantProductEditorProps = {
  product: MerchantApiProduct;
  variants: MerchantApiProductVariant[] | null;
};

/**
 * Optional Details-page hydrate. Edit Product still writes context
 * synchronously before navigation so it does not depend on this effect.
 */
export default function HydrateMerchantProductEditor({
  product,
  variants,
}: HydrateMerchantProductEditorProps) {
  const { setProductEditorData } = useMerchantProductEditor();

  useEffect(() => {
    setProductEditorData(product, variants);
  }, [product, variants, setProductEditorData]);

  return null;
}
