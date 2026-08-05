import { cookies } from "next/headers";
import { apiUrl, API_PATHS } from "./config";
import type { UserAddress } from "@/types/api/userAddress";

/**
 * Server-only: fetches the current user's default (or latest) address.
 * Returns null when none exists or on any error.
 */
export async function getServerDefaultAddress(): Promise<UserAddress | null> {
  const cookieStore = await cookies();
  const cookieHeader = cookieStore.toString();
  const url = apiUrl(API_PATHS.userDefaultAddress);
  try {
    const res = await fetch(url, {
      method: "GET",
      headers: cookieHeader ? { Cookie: cookieHeader } : undefined,
      cache: "no-store",
    });
    if (!res.ok) return null;
    const data = await res.json();
    if (!data || typeof data !== "object" || !("id" in data) || !data.id) {
      return null;
    }
    return data as UserAddress;
  } catch {
    return null;
  }
}
