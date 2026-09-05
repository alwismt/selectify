import { cache } from "react";
import type { ApiCategory } from "@/types/category";
import { apiUrl, API_PATHS } from "./config";

function isApiCategory(value: unknown): value is ApiCategory {
  if (!value || typeof value !== "object") return false;
  const item = value as Record<string, unknown>;

  if (
    typeof item.categoryId !== "number" ||
    typeof item.name !== "string" ||
    typeof item.slug !== "string" ||
    typeof item.isActive !== "boolean"
  ) {
    return false;
  }

  if (
    item.parentId !== undefined &&
    typeof item.parentId !== "number"
  ) {
    return false;
  }

  if (item.children !== undefined) {
    if (!Array.isArray(item.children)) return false;
    if (!item.children.every(isApiCategory)) return false;
  }

  return true;
}

/**
 * Server-only: fetches the public category tree from GET /categories.
 * Returns null on request failure or invalid payload.
 */
export const getServerCategories = cache(async function getServerCategories(): Promise<
  ApiCategory[] | null
> {
  const url = apiUrl(API_PATHS.categories);

  try {
    const res = await fetch(url, {
      method: "GET",
      cache: "no-store",
    });
    if (!res.ok) return null;

    const data: unknown = await res.json();
    if (!Array.isArray(data)) return null;
    if (!data.every(isApiCategory)) return null;
    return data;
  } catch {
    return null;
  }
});
