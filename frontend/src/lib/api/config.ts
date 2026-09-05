// Path constants only (no base URL)
export const API_PATHS = {
  order: "/order",
  orders: "/orders",
  products: "/products",
  categories: "/categories",
  config: "/config",
  cart: "/cart",
  cartItems: "/cart/items",
  authLogin: "/auth/login",
  authForgotPassword: "/auth/forgot-password",
  authResetPassword: "/auth/reset-password",
  authValidateResetPassword: "/auth/reset-password/validate",
  userInfo: "/user/info",
  userMe: "/user/me",
  userDefaultAddress: "/user/addresses/default",
  merchant: "/merchant",
  merchantLogo: "/merchant/logo",
  merchantCountries: "/merchant/countries",
  merchantProducts: "/merchant/products",
  logout: "/logout",
} as const;

export function merchantProduct(productId: number): string {
  return `/merchant/products/${productId}`;
}

export function merchantProductVariants(productId: number): string {
  return `/merchant/products/${productId}/variants`;
}

export function productById(productId: number): string {
  return `/products/${productId}`;
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

/** Public CDN URL for any path under the product-files base. */
export function filesUrl(path: string): string {
  const base = getProductFilesBaseUrl();
  const normalized = path.startsWith("/") ? path : `/${path}`;
  if (!base) return normalized;
  return `${base}${normalized}`;
}

/**
 * Public image URL for a product file id.
 * Shape: `{base}/products/{fileId}` (no extension).
 */
export function productFileUrl(fileId: string | null | undefined): string {
  if (!fileId || !getProductFilesBaseUrl()) return "";
  return filesUrl(`/products/${fileId}`);
}

/**
 * Public image URL for a user profile file id.
 * Shape: `{base}/users/{fileId}` (no extension).
 */
export function userFileUrl(fileId: string | null | undefined): string {
  if (!fileId || !getProductFilesBaseUrl()) return "";
  return filesUrl(`/users/${fileId}`);
}

/**
 * Public image URL for a merchant logo object key.
 * Logo is stored as a full object key (e.g. `merchant/logo/{uuid}-{slug}`).
 */
export function merchantLogoUrl(logo: string | null | undefined): string {
  if (!logo || !getProductFilesBaseUrl()) return "";
  return filesUrl(logo.startsWith("/") ? logo : `/${logo}`);
}
