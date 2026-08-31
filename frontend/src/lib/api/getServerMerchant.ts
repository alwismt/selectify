import { cookies } from "next/headers";
import { apiUrl, API_PATHS } from "./config";
import type { Merchant } from "@/types/merchant";

/**
 * Server-only: fetches current merchant from backend using request cookies.
 * Returns null on 401, 404, or any error.
 */
export async function getServerMerchant(): Promise<Merchant | null> {
  const cookieStore = await cookies();
  const cookieHeader = cookieStore.toString();
  const url = apiUrl(API_PATHS.merchant);

  try {
    const res = await fetch(url, {
      method: "GET",
      headers: cookieHeader ? { Cookie: cookieHeader } : undefined,
      cache: "no-store",
    });
    if (!res.ok) return null;
    return (await res.json()) as Merchant;
  } catch {
    return null;
  }
}
