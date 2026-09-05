/**
 * Extract a positive product ID from the start of a browser path segment
 * like `2-sony-wh-1000xm5` or `2-urlsf`.
 */
export function parseProductIdFromPath(param: string): number | null {
  const trimmed = param.trim();
  if (!trimmed) return null;

  const match = /^(\d+)/.exec(trimmed);
  if (!match) return null;

  const id = Number(match[1]);
  if (!Number.isInteger(id) || id <= 0) return null;
  return id;
}

/** Canonical browser segment: `{productId}-{slug}`. */
export function canonicalProductPath(productId: number, slug: string): string {
  return `${productId}-${slug.trim()}`;
}

/** Customer product page href: `/product/{id}-{slug}`. */
export function productHref(productId: number, slug: string): string {
  return `/product/${canonicalProductPath(productId, slug)}`;
}

/**
 * Merchant product page href: `/merchant/products/{id}-{slug}` plus optional
 * suffix (e.g. `/edit`, `/variants/new`).
 */
export function merchantProductHref(
  productId: number,
  slug: string,
  suffix = ""
): string {
  const base = `/merchant/products/${canonicalProductPath(productId, slug)}`;
  if (!suffix) return base;
  return `${base}${suffix.startsWith("/") ? suffix : `/${suffix}`}`;
}
