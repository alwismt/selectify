import { apiUrl, API_PATHS } from "@/lib/api/config";
import { backendProxyHeaders } from "@/lib/api/proxyHeaders";
import { NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  try {
    const token = request.nextUrl.searchParams.get("token") ?? "";
    const url =
      apiUrl(API_PATHS.authValidateResetPassword) +
      `?token=${encodeURIComponent(token)}`;
    const res = await fetch(url, {
      method: "GET",
      headers: backendProxyHeaders(request),
      cache: "no-store",
    });
    const text = await res.text();
    return new NextResponse(text, {
      status: res.status,
      headers: { "Content-Type": "application/json" },
    });
  } catch (err) {
    console.error("Validate reset password proxy error:", err);
    return NextResponse.json(
      { status: "error", message: "Something went wrong. Please try again." },
      { status: 503 }
    );
  }
}
