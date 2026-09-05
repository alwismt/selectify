import { cache } from "react";
import { apiUrl, API_PATHS } from "./config";
import type { SiteConfig, SiteCurrency } from "@/types/siteConfig";

function isSiteCurrency(value: unknown): value is SiteCurrency {
  if (!value || typeof value !== "object") return false;
  const item = value as Record<string, unknown>;
  return (
    typeof item.code === "string" &&
    typeof item.name === "string" &&
    typeof item.minorUnit === "number" &&
    typeof item.isActive === "boolean" &&
    typeof item.isDefault === "boolean"
  );
}

/**
 * Server-only: fetches site config (default currency) from GET /config.
 * Returns null on failure.
 *
 * Next.js fetch cache (1h) avoids frequent backend hits for slow-changing config.
 * React cache() dedupes within a single server render/request.
 */
export const getServerSiteConfig = cache(async function getServerSiteConfig(): Promise<SiteConfig | null> {
  const url = apiUrl(API_PATHS.config);

  try {
    const res = await fetch(url, {
      method: "GET",
      next: { revalidate: 3600 },
    });
    if (!res.ok) return null;

    const data: unknown = await res.json();
    if (!data || typeof data !== "object") return null;

    const currency = (data as Record<string, unknown>).currency;
    if (currency == null) {
      return { currency: null };
    }
    if (!isSiteCurrency(currency)) return null;
    return { currency };
  } catch {
    return null;
  }
});

/** Convenience: default site currency, or null if unavailable. */
export async function getServerSiteCurrency(): Promise<SiteCurrency | null> {
  const config = await getServerSiteConfig();
  return config?.currency ?? null;
}
