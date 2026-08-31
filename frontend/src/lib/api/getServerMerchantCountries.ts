import { cookies } from "next/headers";
import { apiUrl, API_PATHS } from "./config";
import type { MerchantCountry } from "@/types/merchant";

function isMerchantCountry(value: unknown): value is MerchantCountry {
  if (!value || typeof value !== "object") return false;
  const item = value as Record<string, unknown>;
  return typeof item.code === "string" && typeof item.name === "string";
}

/**
 * Server-only: fetches active merchant countries using request cookies.
 * Returns an empty array on auth/error.
 */
export async function getServerMerchantCountries(): Promise<MerchantCountry[]> {
  const cookieStore = await cookies();
  const cookieHeader = cookieStore.toString();
  const url = apiUrl(API_PATHS.merchantCountries);

  try {
    const res = await fetch(url, {
      method: "GET",
      headers: cookieHeader ? { Cookie: cookieHeader } : undefined,
      cache: "no-store",
    });
    if (!res.ok) return [];

    const data: unknown = await res.json();
    if (!Array.isArray(data)) return [];
    return data.filter(isMerchantCountry);
  } catch {
    return [];
  }
}
