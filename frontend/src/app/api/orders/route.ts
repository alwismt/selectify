import { apiUrl } from "@/lib/api/config";
import { API_PATHS } from "@/lib/api/config";
import { NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  try {
    const url = apiUrl(API_PATHS.orders);
    const cookie = request.headers.get("cookie") ?? undefined;
    const res = await fetch(url, {
      method: "GET",
      headers: cookie ? { Cookie: cookie } : undefined,
      cache: "no-store",
    });
    const data = await res.json();
    return NextResponse.json(data, { status: res.status });
  } catch (err) {
    console.error("Orders proxy GET error:", err);
    return NextResponse.json(
      { error: "Orders service unavailable" },
      { status: 503 }
    );
  }
}

export async function POST(request: NextRequest) {
  try {
    const url = apiUrl(API_PATHS.orders);
    const cookie = request.headers.get("cookie") ?? undefined;
    let body: unknown = {};
    try {
      const text = await request.text();
      if (text) body = JSON.parse(text);
    } catch {
      // use empty body
    }
    const res = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        ...(cookie ? { Cookie: cookie } : {}),
      },
      body: JSON.stringify(body),
      cache: "no-store",
    });
    const data = await res.json();
    return NextResponse.json(data, { status: res.status });
  } catch (err) {
    console.error("Orders proxy POST error:", err);
    return NextResponse.json(
      { error: "Orders service unavailable" },
      { status: 503 }
    );
  }
}
