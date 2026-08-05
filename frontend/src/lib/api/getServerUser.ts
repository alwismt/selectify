import { cookies } from "next/headers";
import { apiUrl } from "./config";
import { API_PATHS } from "./config";
import type { User } from "@/types/user";

/**
 * Server-only: fetches current user from backend using request cookies.
 * Returns null on 401, 404, or any error.
 */
export async function getServerUser(): Promise<User | null> {
  const cookieStore = await cookies();
  const cookieHeader = cookieStore.toString();
  const url = apiUrl(API_PATHS.userInfo);
  try {
    const res = await fetch(url, {
      method: "GET",
      headers: cookieHeader ? { Cookie: cookieHeader } : undefined,
      cache: "no-store",
    });
    if (!res.ok) return null;
    const data = (await res.json()) as User;
    return data;
  } catch {
    return null;
  }
}
