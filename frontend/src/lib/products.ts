import { cache } from "react";
import { Product, ApiProduct, mapApiProductToProduct } from "@/types/product";
import type { ApiVariant } from "@/types/api/variant";
import shopData from "@/components/Shop/shopData";
import { serverApiGet, serverFetch } from "@/lib/api/server";
import { API_PATHS, productById, productVariants } from "@/lib/api/config";
import { getServerSiteCurrency } from "@/lib/api/getServerSiteConfig";

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
  const [raw, currency] = await Promise.all([
    serverApiGet<ApiProduct[]>(API_PATHS.products),
    getServerSiteCurrency(),
  ]);
  return raw.map((api) => mapApiProductToProduct(api, currency));
}

/**
 * Fetch single product by ID from backend (server-only). Returns null on 404.
 * Cached so generateMetadata and the page share one request.
 */
export const getProductByIdFromApi = cache(async function getProductByIdFromApi(
  productId: number
): Promise<Product | null> {
  if (!Number.isInteger(productId) || productId <= 0) return null;

  const [res, currency] = await Promise.all([
    serverFetch(productById(productId), { method: "GET" }),
    getServerSiteCurrency(),
  ]);
  if (!res.ok) {
    if (res.status === 404) return null;
    throw new Error(`API error ${res.status}: ${res.statusText}`);
  }
  const api = (await res.json()) as ApiProduct;
  return mapApiProductToProduct(api, currency);
});

/**
 * Fetch variants for a product (server-only). Use in SSR. Returns [] on 404 or error.
 */
export async function getVariants(productId: number): Promise<ApiVariant[]> {
  const res = await serverFetch(productVariants(productId), { method: "GET" });
  if (!res.ok) return [];
  const raw = (await res.json()) as ApiVariant[];
  if (!Array.isArray(raw)) return [];
  return raw.map((v) => ({
    ...v,
    product_variant_attributes: v.product_variant_attributes ?? [],
    files: v.files ?? [],
  }));
}
