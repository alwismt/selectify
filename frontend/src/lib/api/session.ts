export const SESSION_EXPIRED_EVENT = "selectify:session-expired";

const AUTH_EXEMPT_PREFIXES = ["/auth/login", "/auth/register", "/auth/forgot-password", "/auth/reset-password", "/logout"];

export function isAuthExemptPath(path: string): boolean {
  return AUTH_EXEMPT_PREFIXES.some(
    (prefix) => path === prefix || path.startsWith(`${prefix}?`)
  );
}

/** Clear client auth state and navigate to sign-in after a definitive 401. */
export function handleSessionExpired(): void {
  if (typeof window === "undefined") return;
  window.dispatchEvent(new Event(SESSION_EXPIRED_EVENT));
  const target = "/signin?reason=session-expired";
  if (!window.location.pathname.startsWith("/signin")) {
    // API client layer cannot use Next.js router hooks.
    // eslint-disable-next-line @next/next/no-location-assign-relative-destination -- non-React fetch helper
    window.location.assign(target);
  }
}
