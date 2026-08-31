"use server";

import { cookies } from "next/headers";
import { apiUrl, API_PATHS } from "../config";

const ACCEPTED_TYPES = ["image/jpeg", "image/png", "image/webp"];
const MAX_FILE_SIZE = 5 * 1024 * 1024;

export type UploadMerchantLogoResult =
  | { ok: true; logo: string }
  | { ok: false; error: string };

/**
 * Server Action: validates and forwards a merchant logo upload to the backend.
 */
export async function uploadMerchantLogo(
  formData: FormData
): Promise<UploadMerchantLogoResult> {
  const image = formData.get("image");

  if (!(image instanceof File) || image.size === 0) {
    return { ok: false, error: "Image is required." };
  }

  if (!ACCEPTED_TYPES.includes(image.type)) {
    return { ok: false, error: "Please upload a JPEG, PNG, or WebP image." };
  }

  if (image.size > MAX_FILE_SIZE) {
    return { ok: false, error: "Image must be 5 MB or smaller." };
  }

  const cookieStore = await cookies();
  const cookieHeader = cookieStore.toString();
  const url = apiUrl(API_PATHS.merchantLogo);

  const forward = new FormData();
  forward.append("image", image);

  try {
    const res = await fetch(url, {
      method: "POST",
      headers: cookieHeader ? { Cookie: cookieHeader } : undefined,
      body: forward,
      cache: "no-store",
    });

    if (!res.ok) {
      return { ok: false, error: "Failed to upload logo. Please try again." };
    }

    const data: unknown = await res.json();
    if (typeof data !== "string" || !data) {
      return { ok: false, error: "Invalid logo response." };
    }

    return { ok: true, logo: data };
  } catch {
    return { ok: false, error: "Failed to upload logo. Please try again." };
  }
}
