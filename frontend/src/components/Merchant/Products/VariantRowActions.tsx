"use client";

import { useRouter } from "next/navigation";
import { useMerchantProductEditor } from "@/app/context/MerchantProductEditorContext";
import { merchantProductHref } from "@/lib/productPath";
import type {
  MerchantApiProduct,
  MerchantApiProductVariant,
} from "@/types/merchantProduct";

type VariantRowActionsProps = {
  product: MerchantApiProduct;
  variants: MerchantApiProductVariant[] | null;
  variantId: number;
};

export default function VariantRowActions({
  product,
  variants,
  variantId,
}: VariantRowActionsProps) {
  const router = useRouter();
  const { setProductEditorData } = useMerchantProductEditor();
  const slug = typeof product.slug === "string" ? product.slug.trim() : "";

  const handleEdit = () => {
    if (!slug) return;
    setProductEditorData(product, variants);
    router.push(
      merchantProductHref(
        product.productId,
        slug,
        `/variants/${variantId}/edit`
      )
    );
  };

  return (
    <div className="flex items-center gap-3">
      <button
        type="button"
        onClick={handleEdit}
        disabled={!slug}
        className="text-custom-sm font-medium text-blue hover:text-blue-dark disabled:text-gray-4 disabled:cursor-not-allowed"
      >
        Edit
      </button>
      <button
        type="button"
        disabled
        title="Variant deletion is not available yet"
        className="text-custom-sm font-medium text-red opacity-50 cursor-not-allowed"
      >
        Delete
      </button>
    </div>
  );
}
