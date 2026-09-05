"use client";

import { useRouter } from "next/navigation";
import { useMerchantProductEditor } from "@/app/context/MerchantProductEditorContext";
import { merchantProductHref } from "@/lib/productPath";
import type {
  MerchantApiProduct,
  MerchantApiProductVariant,
} from "@/types/merchantProduct";

type AddVariantButtonProps = {
  product: MerchantApiProduct;
  variants: MerchantApiProductVariant[] | null;
};

export default function AddVariantButton({
  product,
  variants,
}: AddVariantButtonProps) {
  const router = useRouter();
  const { setProductEditorData } = useMerchantProductEditor();
  const slug = typeof product.slug === "string" ? product.slug.trim() : "";

  const handleAdd = () => {
    if (!slug) return;
    setProductEditorData(product, variants);
    router.push(
      merchantProductHref(product.productId, slug, "/variants/new")
    );
  };

  return (
    <button
      type="button"
      onClick={handleAdd}
      disabled={!slug}
      className={`inline-flex items-center justify-center font-medium text-white py-3 px-7 rounded-md ease-out duration-200 ${
        slug
          ? "bg-blue hover:bg-blue-dark"
          : "bg-gray-4 cursor-not-allowed"
      }`}
    >
      + Add Variant
    </button>
  );
}
