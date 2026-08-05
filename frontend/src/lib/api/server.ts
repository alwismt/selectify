import { apiUrl } from "./config";

function serverFetch(path: string, init?: RequestInit): Promise<Response> {
  const url = apiUrl(path);
  return fetch(url, init);
}

/**
 * GET request that returns parsed JSON. Throws on !res.ok.
 */
export async function serverApiGet<T = unknown>(
  path: string,
  init?: Omit<RequestInit, "method">
): Promise<T> {
  const res = await serverFetch(path, { ...init, method: "GET" });
  if (!res.ok) {
    throw new Error(`API error ${res.status}: ${res.statusText}`);
  }
  return res.json() as Promise<T>;
}

/**
 * Raw server fetch for other methods. Use serverApiGet for GET + JSON.
 */
export { serverFetch };
