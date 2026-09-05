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

export type MoneyCurrency = {
  code: string;
  minorUnit: number;
};

/**
 * Format an integer minor-unit amount for display.
 * Converts to major units only at the presentation boundary.
 */
export function formatMoney(
  amountMinor: number,
  currency: MoneyCurrency
): string {
  const units =
    typeof currency.minorUnit === "number" &&
    Number.isFinite(currency.minorUnit) &&
    currency.minorUnit >= 0
      ? currency.minorUnit
      : 0;
  const amount = Number.isFinite(amountMinor)
    ? amountMinor / 10 ** units
    : 0;

  return new Intl.NumberFormat("en", {
    style: "currency",
    currency: currency.code,
    minimumFractionDigits: units,
    maximumFractionDigits: units,
  }).format(amount);
}

/**
 * Resolve currency metadata for an order financial snapshot.
 * `order.currency` always controls the display code.
 * Prefer site minorUnit when codes match; otherwise use Intl ISO fraction digits.
 */
export function resolveOrderMoneyCurrency(
  orderCurrencyCode: string,
  siteCurrency: MoneyCurrency | null | undefined
): MoneyCurrency {
  const code = orderCurrencyCode || siteCurrency?.code || "XXX";
  if (siteCurrency && code === siteCurrency.code) {
    return { code, minorUnit: siteCurrency.minorUnit };
  }
  const minorUnit = intlCurrencyMinorUnit(code);
  return { code, minorUnit };
}

/** ISO/Intl fraction digits for a currency code (display-only fallback). */
export function intlCurrencyMinorUnit(currencyCode: string): number {
  try {
    const digits = new Intl.NumberFormat("en", {
      style: "currency",
      currency: currencyCode,
    }).resolvedOptions().maximumFractionDigits;
    return typeof digits === "number" && digits >= 0 ? digits : 2;
  } catch {
    return 2;
  }
}

/**
 * Format order minor-unit amounts using order.currency as the snapshot code.
 */
export function formatOrderMoney(
  amountMinor: number,
  orderCurrencyCode: string,
  siteCurrency: MoneyCurrency | null | undefined
): string {
  return formatMoney(
    amountMinor,
    resolveOrderMoneyCurrency(orderCurrencyCode, siteCurrency)
  );
}

/**
 * Format minor units as a decimal string for merchant price inputs (no symbol).
 * Example: 39999 + minorUnit 2 → "399.99"
 */
export function formatMinorToMajorInput(
  amountMinor: number,
  minorUnit: number
): string {
  const units =
    typeof minorUnit === "number" &&
    Number.isInteger(minorUnit) &&
    minorUnit >= 0
      ? minorUnit
      : 0;
  if (!Number.isFinite(amountMinor) || !Number.isInteger(amountMinor)) {
    return "";
  }
  if (units === 0) return String(amountMinor);
  const sign = amountMinor < 0 ? "-" : "";
  const abs = Math.abs(amountMinor);
  const scale = 10 ** units;
  const whole = Math.floor(abs / scale);
  const fraction = String(abs % scale).padStart(units, "0");
  return `${sign}${whole}.${fraction}`;
}

/**
 * Parse a human-readable major-unit amount (e.g. "1299.99") into integer
 * minor units using integer arithmetic only. Returns null if invalid or <= 0.
 */
export function parseMajorAmountToMinor(
  input: string,
  minorUnit = 2
): number | null {
  const trimmed = input.trim();
  if (!trimmed) return null;

  const units =
    typeof minorUnit === "number" &&
    Number.isInteger(minorUnit) &&
    minorUnit >= 0
      ? minorUnit
      : 2;

  if (!/^\d+(\.\d+)?$/.test(trimmed)) return null;

  const [wholePart, fractionPart = ""] = trimmed.split(".");
  if (fractionPart.length > units) return null;

  const whole = Number(wholePart);
  if (!Number.isInteger(whole) || whole < 0) return null;

  const paddedFraction = fractionPart.padEnd(units, "0");
  const fraction = units === 0 ? 0 : Number(paddedFraction);
  if (!Number.isInteger(fraction) || fraction < 0) return null;

  const scale = 10 ** units;
  const amountMinor = whole * scale + fraction;
  if (!Number.isInteger(amountMinor) || amountMinor <= 0) return null;

  return amountMinor;
}
