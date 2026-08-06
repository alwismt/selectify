import { apiUrl, orderAddress } from "./config";
import type { Order, OrderAddress, OrderShippingAddressInput } from "@/types/api/order";
import type { UserAddress } from "@/types/api/userAddress";
import { isUserFile, type UserFile } from "@/types/api/userFile";

export type ClientRequestInit = Omit<RequestInit, "method">;

/** Same-origin GET /api/orders (proxy). Hides backend URL. */
export async function clientOrdersGet(): Promise<Order[]> {
  const res = await fetch("/api/orders", { method: "GET", credentials: "include" });
  if (!res.ok) {
    throw new Error(`API error ${res.status}: ${res.statusText}`);
  }
  const data = await res.json();
  return Array.isArray(data) ? data : [];
}

/** Same-origin POST /api/orders (proxy). Hides backend URL. */
export async function clientOrdersPost(body?: unknown): Promise<Order> {
  const res = await fetch("/api/orders", {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body ?? {}),
  });
  if (!res.ok) {
    throw new Error(`API error ${res.status}: ${res.statusText}`);
  }
  return res.json() as Promise<Order>;
}

/** Same-origin PUT /api/orders/{id}/address (proxy). */
export async function clientOrderAddressPut(
  orderId: number,
  body: OrderShippingAddressInput
): Promise<OrderAddress> {
  const res = await fetch(`/api/orders/${orderId}/address`, {
    method: "PUT",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    throw new Error(`API error ${res.status}: ${res.statusText}`);
  }
  return res.json() as Promise<OrderAddress>;
}

/** Same-origin GET /api/user/addresses/default (proxy). Returns null if none. */
export async function clientUserDefaultAddressGet(): Promise<UserAddress | null> {
  const res = await fetch("/api/user/addresses/default", {
    method: "GET",
    credentials: "include",
  });
  if (!res.ok) {
    throw new Error(`API error ${res.status}: ${res.statusText}`);
  }
  const data = await res.json();
  if (!data || typeof data !== "object" || !data.id) {
    return null;
  }
  return data as UserAddress;
}

/** Same-origin GET /api/user/me (proxy). Returns null if no avatar. */
export async function clientUserFileGet(): Promise<UserFile | null> {
  const res = await fetch("/api/user/me", {
    method: "GET",
    credentials: "include",
  });
  if (!res.ok) {
    throw new Error(`API error ${res.status}: ${res.statusText}`);
  }
  const data: unknown = await res.json();
  return isUserFile(data) ? data : null;
}

/** Same-origin POST /api/user/me (proxy) with multipart image upload. */
export async function clientUserFileUpload(file: File): Promise<UserFile> {
  const formData = new FormData();
  formData.append("image", file);

  const res = await fetch("/api/user/me", {
    method: "POST",
    credentials: "include",
    body: formData,
  });
  if (!res.ok) {
    throw new Error(`API error ${res.status}: ${res.statusText}`);
  }
  const data: unknown = await res.json();
  if (!isUserFile(data)) {
    throw new Error("Invalid avatar response");
  }
  return data;
}

/** Same-origin DELETE /api/user/me (proxy). */
export async function clientUserFileDelete(): Promise<void> {
  const res = await fetch("/api/user/me", {
    method: "DELETE",
    credentials: "include",
  });
  if (!res.ok) {
    throw new Error(`API error ${res.status}: ${res.statusText}`);
  }
}

function clientFetch(path: string, init?: RequestInit): Promise<Response> {
  const url = apiUrl(path);
  return fetch(url, { ...init, credentials: "include" });
}

/**
 * GET request that returns parsed JSON. Throws on !res.ok.
 */
export async function apiClientGet<T = unknown>(
  path: string,
  init?: ClientRequestInit
): Promise<T> {
  const res = await clientFetch(path, { ...init, method: "GET" });
  if (!res.ok) {
    throw new Error(`API error ${res.status}: ${res.statusText}`);
  }
  return res.json() as Promise<T>;
}

/**
 * POST request with JSON body. Throws on !res.ok.
 */
export async function apiClientPost<T = unknown>(
  path: string,
  body: unknown,
  init?: Omit<RequestInit, "method" | "body">
): Promise<T> {
  const res = await clientFetch(path, {
    ...init,
    method: "POST",
    headers: { "Content-Type": "application/json", ...(init?.headers as object) },
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    throw new Error(`API error ${res.status}: ${res.statusText}`);
  }
  const text = await res.text();
  return (text ? JSON.parse(text) : undefined) as T;
}

/**
 * PATCH request with JSON body. Throws on !res.ok.
 */
export async function apiClientPatch<T = unknown>(
  path: string,
  body: unknown,
  init?: Omit<RequestInit, "method" | "body">
): Promise<T> {
  const res = await clientFetch(path, {
    ...init,
    method: "PATCH",
    headers: { "Content-Type": "application/json", ...(init?.headers as object) },
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    throw new Error(`API error ${res.status}: ${res.statusText}`);
  }
  const text = await res.text();
  return (text ? JSON.parse(text) : undefined) as T;
}

/**
 * DELETE request. Throws on !res.ok.
 */
export async function apiClientDelete(path: string): Promise<void> {
  const res = await clientFetch(path, { method: "DELETE" });
  if (!res.ok) {
    throw new Error(`API error ${res.status}: ${res.statusText}`);
  }
}

/**
 * Raw client fetch for other methods (POST, PUT, etc.). Use apiClientGet for GET + JSON.
 */
export { clientFetch, orderAddress };
