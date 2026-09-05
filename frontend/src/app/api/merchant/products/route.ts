import { apiUrl, API_PATHS } from "@/lib/api/config";
import { NextRequest, NextResponse } from "next/server";

/**
 * Same-origin proxy: POST multipart create-product to Go POST /merchant/products.
 * Keeps the backend host server-only.
 */
export async function POST(request: NextRequest) {
  try {
    const url = apiUrl(API_PATHS.merchantProducts);
    const cookie = request.headers.get("cookie") ?? undefined;
    const formData = await request.formData();

    const res = await fetch(url, {
      method: "POST",
      headers: cookie ? { Cookie: cookie } : undefined,
      body: formData,
      cache: "no-store",
    });

    const text = await res.text();
    return new NextResponse(text, {
      status: res.status,
      headers: { "Content-Type": "application/json" },
    });
  } catch (err) {
    console.error("Merchant products proxy POST error:", err);
    return NextResponse.json(
      { error: "Product service unavailable" },
      { status: 503 }
    );
  }
}
