import type { CartItem, CartResponse } from "@/types/api/cart";
import type { ApiProduct, ApiProductFile } from "@/types/product";
import type { ApiVariant } from "@/types/api/variant";
import { productFileUrl } from "@/lib/api/config";

type FileLike = Pick<
  ApiProductFile,
  "file_id" | "position" | "is_primary"
>;

function pickBestFile(files: FileLike[] | null | undefined): FileLike | null {
  if (!files || files.length === 0) return null;
  const primary = files.find((f) => f.is_primary);
  if (primary) return primary;
  return [...files].sort((a, b) => a.position - b.position)[0] ?? null;
}

/**
 * Attach imageUrl to cart items using product/variant file maps.
 * Shared by SSR and client enrichment.
 */
export function attachCartImageUrls(
  cart: CartResponse,
  productFiles: Map<number, ApiProductFile | null | undefined>,
  variantFiles: Map<number, FileLike[]>
): CartResponse {
  const items: CartItem[] = cart.items.map((item) => {
    const variantFile = pickBestFile(variantFiles.get(item.variant.id));
    const productFile = productFiles.get(item.variant.product.id);
    const fileId = variantFile?.file_id ?? productFile?.file_id ?? null;
    const imageUrl = productFileUrl(fileId) || null;
    return { ...item, imageUrl };
  });
  return { ...cart, items };
}

export function collectCartProductIds(cart: CartResponse): number[] {
  const ids = new Set<number>();
  for (const item of cart.items) {
    if (item.variant?.product?.id) {
      ids.add(item.variant.product.id);
    }
  }
  return Array.from(ids);
}

export function indexProductAndVariantFiles(
  products: Array<ApiProduct | null>,
  variantsByProduct: Array<ApiVariant[] | null>
): {
  productFiles: Map<number, ApiProductFile | null | undefined>;
  variantFiles: Map<number, FileLike[]>;
} {
  const productFiles = new Map<number, ApiProductFile | null | undefined>();
  const variantFiles = new Map<number, FileLike[]>();

  for (const product of products) {
    if (!product) continue;
    productFiles.set(product.productId, product.productFile ?? null);
  }

  for (const variants of variantsByProduct) {
    if (!variants) continue;
    for (const variant of variants) {
      variantFiles.set(variant.id, variant.files ?? []);
    }
  }

  return { productFiles, variantFiles };
}
