"use client";

import React, { createContext, useCallback, useContext, useEffect, useState } from "react";
import { useUser } from "@/app/context/UserContext";
import { apiClientGet } from "@/lib/api/client";
import { API_PATHS } from "@/lib/api/config";
import type { CartResponse } from "@/types/api/cart";

type CartContextValue = {
  cart: CartResponse | null;
  loading: boolean;
  error: Error | null;
  refetch: () => Promise<void>;
};

const CartContext = createContext<CartContextValue | null>(null);

export function CartProvider({ children }: { children: React.ReactNode }) {
  const { user } = useUser();
  const [cart, setCart] = useState<CartResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  const refetch = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await apiClientGet<CartResponse>(API_PATHS.cart);
      setCart(data);
    } catch (err) {
      setError(err instanceof Error ? err : new Error(String(err)));
      setCart(null);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    refetch();
  }, [user?.id, refetch]);

  const value: CartContextValue = {
    cart,
    loading,
    error,
    refetch,
  };

  return (
    <CartContext.Provider value={value}>{children}</CartContext.Provider>
  );
}

export function useCart(): CartContextValue {
  const ctx = useContext(CartContext);
  if (!ctx) {
    throw new Error("useCart must be used within CartProvider");
  }
  return ctx;
}
