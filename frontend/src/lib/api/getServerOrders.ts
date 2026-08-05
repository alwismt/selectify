import { cookies } from "next/headers";
import { apiUrl } from "./config";
import { API_PATHS } from "./config";
import type { Order } from "@/types/api/order";

/**
 * Server-only: fetches current user's orders from backend using request cookies.
 * Returns [] on non-OK or any error.
 */
export async function getServerOrders(): Promise<Order[]> {
  const cookieStore = await cookies();
  const cookieHeader = cookieStore.toString();
  const url = apiUrl(API_PATHS.orders);
  try {
    const res = await fetch(url, {
      method: "GET",
      headers: cookieHeader ? { Cookie: cookieHeader } : undefined,
      cache: "no-store",
    });
    if (!res.ok) return [];
    const data = await res.json();
    return Array.isArray(data) ? data : [];
  } catch {
    return [];
  }
}
