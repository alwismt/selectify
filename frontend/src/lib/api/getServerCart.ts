import { cache } from "react";
import { cookies } from "next/headers";
import { apiUrl, API_PATHS, productById, productVariants } from "./config";
import type { CartResponse } from "@/types/api/cart";
import type { ApiProduct } from "@/types/product";
import type { ApiVariant } from "@/types/api/variant";
import {
  attachCartImageUrls,
  collectCartProductIds,
  indexProductAndVariantFiles,
} from "@/lib/cart/enrichCartImages";

function normalizeCart(data: unknown): CartResponse | null {
  if (!data || typeof data !== "object") return null;
  const raw = data as Record<string, unknown>;
  if (!Array.isArray(raw.items)) return null;
  return {
    items: raw.items as CartResponse["items"],
    subtotal: typeof raw.subtotal === "number" ? raw.subtotal : 0,
    item_count: typeof raw.item_count === "number" ? raw.item_count : 0,
  };
}

async function fetchProductForCart(
  productId: number,
  cookieHeader: string
): Promise<ApiProduct | null> {
  try {
    const res = await fetch(apiUrl(productById(productId)), {
      method: "GET",
      headers: cookieHeader ? { Cookie: cookieHeader } : undefined,
      cache: "no-store",
    });
    if (!res.ok) return null;
    return (await res.json()) as ApiProduct;
  } catch {
    return null;
  }
}

async function fetchVariantsForCart(
  productId: number,
  cookieHeader: string
): Promise<ApiVariant[]> {
  try {
    const res = await fetch(apiUrl(productVariants(productId)), {
      method: "GET",
      headers: cookieHeader ? { Cookie: cookieHeader } : undefined,
      cache: "no-store",
    });
    if (!res.ok) return [];
    const raw = (await res.json()) as ApiVariant[];
    if (!Array.isArray(raw)) return [];
    return raw.map((v) => ({
      ...v,
      product_variant_attributes: v.product_variant_attributes ?? [],
      files: v.files ?? [],
    }));
  } catch {
    return [];
  }
}

/**
 * Enrich cart items with imageUrl from product/variant files (frontend-only).
 */
export async function enrichCartImagesServer(
  cart: CartResponse
): Promise<CartResponse> {
  const productIds = collectCartProductIds(cart);
  if (productIds.length === 0) {
    return { ...cart, items: cart.items.map((i) => ({ ...i, imageUrl: null })) };
  }

  const cookieStore = await cookies();
  const cookieHeader = cookieStore.toString();

  const products = await Promise.all(
    productIds.map((id) => fetchProductForCart(id, cookieHeader))
  );
  const variantsByProduct = await Promise.all(
    productIds.map((id) => fetchVariantsForCart(id, cookieHeader))
  );

  const { productFiles, variantFiles } = indexProductAndVariantFiles(
    products,
    variantsByProduct
  );
  return attachCartImageUrls(cart, productFiles, variantFiles);
}

/**
 * Server-only: fetches cart using request cookies and enriches images.
 * Returns null when unauthenticated or on error.
 */
export const getServerCart = cache(async function getServerCart(): Promise<CartResponse | null> {
  const cookieStore = await cookies();
  const cookieHeader = cookieStore.toString();
  if (!cookieHeader) return null;

  const url = apiUrl(API_PATHS.cart);
  try {
    const res = await fetch(url, {
      method: "GET",
      headers: { Cookie: cookieHeader },
      cache: "no-store",
    });
    if (!res.ok) return null;
    const cart = normalizeCart(await res.json());
    if (!cart) return null;
    return enrichCartImagesServer(cart);
  } catch {
    return null;
  }
});
