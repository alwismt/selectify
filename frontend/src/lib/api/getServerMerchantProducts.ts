import { cache } from "react";
import { cookies } from "next/headers";
import { apiUrl, API_PATHS, merchantProduct, merchantProductVariants } from "./config";
import type {
  MerchantApiProduct,
  MerchantApiProductVariant,
  MerchantApiVariantAttribute,
  MerchantApiVariantFile,
} from "@/types/merchantProduct";

function isMerchantApiProduct(value: unknown): value is MerchantApiProduct {
  if (!value || typeof value !== "object") return false;
  const item = value as Record<string, unknown>;
  return (
    typeof item.productId === "number" &&
    typeof item.sku === "string" &&
    typeof item.name === "string" &&
    typeof item.price === "number" &&
    typeof item.isActive === "boolean" &&
    typeof item.inStock === "boolean"
  );
}

/**
 * Server-only: fetches merchant products using request cookies.
 * Returns [] on successful empty list; null on request failure.
 */
export async function getServerMerchantProducts(): Promise<
  MerchantApiProduct[] | null
> {
  const cookieStore = await cookies();
  const cookieHeader = cookieStore.toString();
  const url = apiUrl(API_PATHS.merchantProducts);

  try {
    const res = await fetch(url, {
      method: "GET",
      headers: cookieHeader ? { Cookie: cookieHeader } : undefined,
      cache: "no-store",
    });
    if (!res.ok) return null;

    const data: unknown = await res.json();
    if (!Array.isArray(data)) return null;
    return data.filter(isMerchantApiProduct);
  } catch {
    return null;
  }
}

/**
 * Server-only: fetches a single merchant product by ID using request cookies.
 * Returns null on invalid ID, 400/404, or any error.
 */
export const getServerMerchantProduct = cache(async function getServerMerchantProduct(
  productId: number
): Promise<MerchantApiProduct | null> {
  if (!Number.isInteger(productId) || productId <= 0) return null;

  const cookieStore = await cookies();
  const cookieHeader = cookieStore.toString();
  const url = apiUrl(merchantProduct(productId));

  try {
    const res = await fetch(url, {
      method: "GET",
      headers: cookieHeader ? { Cookie: cookieHeader } : undefined,
      cache: "no-store",
    });
    if (!res.ok) return null;

    const data: unknown = await res.json();
    if (!isMerchantApiProduct(data)) return null;
    return data;
  } catch {
    return null;
  }
});

function isMerchantApiVariantAttribute(
  value: unknown
): value is MerchantApiVariantAttribute {
  if (!value || typeof value !== "object") return false;
  const item = value as Record<string, unknown>;
  return (
    typeof item.id === "number" &&
    typeof item.variant_id === "number" &&
    typeof item.name === "string" &&
    typeof item.value === "string"
  );
}

function isMerchantApiVariantFile(
  value: unknown
): value is MerchantApiVariantFile {
  if (!value || typeof value !== "object") return false;
  const item = value as Record<string, unknown>;
  return (
    typeof item.file_id === "string" &&
    typeof item.product_id === "number" &&
    typeof item.variant_id === "number" &&
    typeof item.content_type === "string" &&
    typeof item.position === "number" &&
    typeof item.is_primary === "boolean"
  );
}

function isMerchantApiProductVariant(
  value: unknown
): value is MerchantApiProductVariant {
  if (!value || typeof value !== "object") return false;
  const item = value as Record<string, unknown>;
  const attrs = item.product_variant_attributes;
  const files = item.files;
  return (
    typeof item.id === "number" &&
    typeof item.product_id === "number" &&
    typeof item.sku === "string" &&
    typeof item.price_amount === "number" &&
    typeof item.created_at === "string" &&
    typeof item.updated_at === "string" &&
    typeof item.stock_quantity === "number" &&
    typeof item.reserved_quantity === "number" &&
    Array.isArray(attrs) &&
    attrs.every(isMerchantApiVariantAttribute) &&
    Array.isArray(files) &&
    files.every(isMerchantApiVariantFile)
  );
}

/**
 * Server-only: fetches merchant product variants by product ID using request cookies.
 * Returns [] on successful empty list; null on request failure.
 */
export const getServerMerchantProductVariants = cache(async function getServerMerchantProductVariants(
  productId: number
): Promise<MerchantApiProductVariant[] | null> {
  if (!Number.isInteger(productId) || productId <= 0) return null;

  const cookieStore = await cookies();
  const cookieHeader = cookieStore.toString();
  const url = apiUrl(merchantProductVariants(productId));

  try {
    const res = await fetch(url, {
      method: "GET",
      headers: cookieHeader ? { Cookie: cookieHeader } : undefined,
      cache: "no-store",
    });
    if (!res.ok) return null;

    const data: unknown = await res.json();
    if (!Array.isArray(data)) return null;
    return data.filter(isMerchantApiProductVariant);
  } catch {
    return null;
  }
});
