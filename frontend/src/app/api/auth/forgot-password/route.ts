import { apiUrl, API_PATHS } from "@/lib/api/config";
import { backendProxyHeaders } from "@/lib/api/proxyHeaders";
import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  try {
    const body = await request.text();
    const res = await fetch(apiUrl(API_PATHS.authForgotPassword), {
      method: "POST",
      headers: backendProxyHeaders(request, {
        "Content-Type": "application/json",
      }),
      body,
      cache: "no-store",
    });
    const text = await res.text();
    return new NextResponse(text, {
      status: res.status,
      headers: { "Content-Type": "application/json" },
    });
  } catch (err) {
    console.error("Forgot password proxy error:", err);
    return NextResponse.json(
      { status: "error", message: "Something went wrong. Please try again." },
      { status: 503 }
    );
  }
}
