import { apiUrl, API_PATHS } from "@/lib/api/config";
import { NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  try {
    const url = apiUrl(API_PATHS.userDefaultAddress);
    const cookie = request.headers.get("cookie") ?? undefined;
    const res = await fetch(url, {
      method: "GET",
      headers: cookie ? { Cookie: cookie } : undefined,
      cache: "no-store",
    });
    const data = await res.json();
    return NextResponse.json(data, { status: res.status });
  } catch (err) {
    console.error("User default address proxy GET error:", err);
    return NextResponse.json(
      { error: "Address service unavailable" },
      { status: 503 }
    );
  }
}
