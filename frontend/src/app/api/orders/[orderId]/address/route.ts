import { apiUrl, orderAddress } from "@/lib/api/config";
import { NextRequest, NextResponse } from "next/server";

export async function PUT(
  request: NextRequest,
  context: { params: Promise<{ orderId: string }> }
) {
  try {
    const { orderId } = await context.params;
    const id = Number(orderId);
    if (!orderId || Number.isNaN(id) || id <= 0) {
      return NextResponse.json({ error: "Invalid order id" }, { status: 400 });
    }

    const cookie = request.headers.get("cookie") ?? undefined;
    let body: unknown = {};
    try {
      const text = await request.text();
      if (text) body = JSON.parse(text);
    } catch {
      return NextResponse.json({ error: "Invalid JSON body" }, { status: 400 });
    }

    const res = await fetch(apiUrl(orderAddress(id)), {
      method: "PUT",
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
    console.error("Order address proxy PUT error:", err);
    return NextResponse.json(
      { error: "Orders service unavailable" },
      { status: 503 }
    );
  }
}
