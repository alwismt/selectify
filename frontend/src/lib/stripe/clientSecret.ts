const storageKey = (orderId: number) => `order:${orderId}:client_secret`;

export function setOrderClientSecret(orderId: number, secret: string): void {
  if (typeof window === "undefined") return;
  sessionStorage.setItem(storageKey(orderId), secret);
}

export function getOrderClientSecret(orderId: number): string | null {
  if (typeof window === "undefined") return null;
  return sessionStorage.getItem(storageKey(orderId));
}

export function clearOrderClientSecret(orderId: number): void {
  if (typeof window === "undefined") return;
  sessionStorage.removeItem(storageKey(orderId));
}
