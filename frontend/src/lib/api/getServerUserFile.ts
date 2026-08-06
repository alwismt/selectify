import { cookies } from "next/headers";
import { apiUrl, API_PATHS } from "./config";
import { isUserFile, type UserFile } from "@/types/api/userFile";

/**
 * Server-only: fetches current user profile image metadata from backend using request cookies.
 * Returns null when no image exists or on auth/error.
 */
export async function getServerUserFile(): Promise<UserFile | null> {
  const cookieStore = await cookies();
  const cookieHeader = cookieStore.toString();
  const url = apiUrl(API_PATHS.userMe);

  try {
    const res = await fetch(url, {
      method: "GET",
      headers: cookieHeader ? { Cookie: cookieHeader } : undefined,
      cache: "no-store",
    });
    if (!res.ok) return null;

    const data: unknown = await res.json();
    return isUserFile(data) ? data : null;
  } catch {
    return null;
  }
}
