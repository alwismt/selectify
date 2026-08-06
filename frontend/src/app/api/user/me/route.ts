import { apiUrl, API_PATHS } from "@/lib/api/config";
import { NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  try {
    const url = apiUrl(API_PATHS.userMe);
    const cookie = request.headers.get("cookie") ?? undefined;
    const res = await fetch(url, {
      method: "GET",
      headers: cookie ? { Cookie: cookie } : undefined,
      cache: "no-store",
    });
    const data = await res.json();
    return NextResponse.json(data, { status: res.status });
  } catch (err) {
    console.error("User avatar proxy GET error:", err);
    return NextResponse.json(
      { error: "Avatar service unavailable" },
      { status: 503 }
    );
  }
}

export async function POST(request: NextRequest) {
  try {
    const url = apiUrl(API_PATHS.userMe);
    const cookie = request.headers.get("cookie") ?? undefined;
    const formData = await request.formData();
    const res = await fetch(url, {
      method: "POST",
      headers: cookie ? { Cookie: cookie } : undefined,
      body: formData,
    });
    const text = await res.text();
    return new NextResponse(text, {
      status: res.status,
      headers: { "Content-Type": "application/json" },
    });
  } catch (err) {
    console.error("User avatar proxy POST error:", err);
    return NextResponse.json(
      { error: "Avatar service unavailable" },
      { status: 503 }
    );
  }
}

export async function DELETE(request: NextRequest) {
  try {
    const url = apiUrl(API_PATHS.userMe);
    const cookie = request.headers.get("cookie") ?? undefined;
    const res = await fetch(url, {
      method: "DELETE",
      headers: cookie ? { Cookie: cookie } : undefined,
    });
    return new NextResponse(null, { status: res.status });
  } catch (err) {
    console.error("User avatar proxy DELETE error:", err);
    return NextResponse.json(
      { error: "Avatar service unavailable" },
      { status: 503 }
    );
  }
}
