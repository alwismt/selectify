// Path constants only (no base URL)
export const API_PATHS = {
  order: "/order",
  orders: "/orders",
  products: "/products",
  cart: "/cart",
  cartItems: "/cart/items",
  authLogin: "/auth/login",
  userInfo: "/user/info",
  userDefaultAddress: "/user/addresses/default",
  logout: "/logout",
} as const;

export function productBySlug(slug: string): string {
  return `/products/${slug}`;
}

export function productVariants(productId: number): string {
  return `/products/${productId}/variants`;
}

export function cartItem(cartItemId: number): string {
  return `/cart/items/${cartItemId}`;
}

export function orderAddress(orderId: number): string {
  return `/orders/${orderId}/address`;
}

// Base URL: empty for relative, or from env
export function getApiBaseUrl(): string {
  if (typeof window !== "undefined") {
    return process.env.NEXT_PUBLIC_API_URL ?? "";
  }
  return process.env.API_URL ?? process.env.NEXT_PUBLIC_API_URL ?? "";
}

/** Build full API URL from a path (e.g. API_PATHS.order). */
export function apiUrl(path: string): string {
  const base = getApiBaseUrl();
  if (!base) return path;
  return base.replace(/\/$/, "") + (path.startsWith("/") ? path : `/${path}`);
}

/** Base URL for product/variant image files (no trailing slash). */
export function getProductFilesBaseUrl(): string {
  return (process.env.NEXT_PUBLIC_PRODUCT_FILES_BASE_URL ?? "").replace(
    /\/$/,
    ""
  );
}

/**
 * Public image URL for a product file id.
 * Shape: `{base}/products/{fileId}` (no extension).
 */
export function productFileUrl(fileId: string | null | undefined): string {
  const base = getProductFilesBaseUrl();
  if (!base || !fileId) return "";
  return `${base}/products/${fileId}`;
}
