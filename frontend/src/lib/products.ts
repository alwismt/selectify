import { Product, ApiProduct, mapApiProductToProduct } from "@/types/product";
import type { ApiVariant } from "@/types/api/variant";
import shopData from "@/components/Shop/shopData";
import { serverApiGet, serverFetch } from "@/lib/api/server";
import { API_PATHS, productBySlug, productVariants } from "@/lib/api/config";

export function getProductBySlug(slug: string): Product | null {
  return shopData.find((p) => p.slug === slug) ?? null;
}

export function getAllProducts(): Product[] {
  return shopData;
}

/**
 * Fetch products from backend (server-only). Use in SSR pages.
 */
export async function getProducts(): Promise<Product[]> {
  const raw = await serverApiGet<ApiProduct[]>(API_PATHS.products);
  return raw.map(mapApiProductToProduct);
}

/**
 * Fetch single product by slug from backend (server-only). Returns null on 404.
 */
export async function getProductBySlugFromApi(slug: string): Promise<Product | null> {
  const res = await serverFetch(productBySlug(slug), { method: "GET" });
  if (!res.ok) {
    if (res.status === 404) return null;
    throw new Error(`API error ${res.status}: ${res.statusText}`);
  }
  const api = (await res.json()) as ApiProduct;
  return mapApiProductToProduct(api);
}

/**
 * Fetch variants for a product (server-only). Use in SSR. Returns [] on 404 or error.
 */
export async function getVariants(productId: number): Promise<ApiVariant[]> {
  const res = await serverFetch(productVariants(productId), { method: "GET" });
  if (!res.ok) return [];
  const raw = (await res.json()) as ApiVariant[];
  return Array.isArray(raw) ? raw : [];
}
