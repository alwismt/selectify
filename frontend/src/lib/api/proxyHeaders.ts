import { NextRequest } from "next/server";

/**
 * Headers to forward from a browser request through a Next.js API proxy
 * so the Go backend can resolve the real client IP and user agent.
 */
export function backendProxyHeaders(
  request: NextRequest,
  extra?: HeadersInit
): Headers {
  const headers = new Headers(extra);

  const xff = request.headers.get("x-forwarded-for");
  const realIp = request.headers.get("x-real-ip");
  const ua = request.headers.get("user-agent");
  const nextIp = (request as NextRequest & { ip?: string }).ip;

  if (xff) {
    headers.set("X-Forwarded-For", xff);
  } else if (realIp) {
    headers.set("X-Forwarded-For", realIp);
  } else if (nextIp) {
    headers.set("X-Forwarded-For", nextIp);
  }

  if (realIp) {
    headers.set("X-Real-IP", realIp);
  } else if (xff) {
    headers.set("X-Real-IP", xff.split(",")[0]!.trim());
  } else if (nextIp) {
    headers.set("X-Real-IP", nextIp);
  }

  if (ua) {
    headers.set("User-Agent", ua);
  }

  return headers;
}
