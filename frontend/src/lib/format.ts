/**
 * Formats a date string for "Member Since MMM YYYY" display.
 * Returns a fallback string if the date is invalid.
 */
export function formatMemberSince(createdAt: string): string {
  const date = new Date(createdAt);
  if (Number.isNaN(date.getTime())) {
    return "Member";
  }
  const formatted = date.toLocaleDateString("en-US", {
    month: "short",
    year: "numeric",
  });
  return `Member Since ${formatted}`;
}

/** Format order date from item created_at or ISO string; returns "MMM d, yyyy" or fallback */
export function formatOrderDate(createdAt: string): string {
  const date = new Date(createdAt);
  if (Number.isNaN(date.getTime())) return "—";
  return date.toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
}

/** Format money with currency (e.g. EUR -> €, USD -> $) */
export function formatOrderTotal(total: number, currency: string): string {
  const symbol = currency === "EUR" ? "€" : currency === "USD" ? "$" : currency + " ";
  return `${symbol}${total.toFixed(2)}`;
}

/** Format a product/shop price (defaults to USD if currency missing). */
export function formatMoney(amount: number, currency?: string): string {
  return formatOrderTotal(amount, currency ?? "USD");
}
