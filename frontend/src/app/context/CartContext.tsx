"use client";

import React, { createContext, useContext, useState } from "react";
import { useUser } from "@/app/context/UserContext";
import type { CartResponse } from "@/types/api/cart";

type CartContextValue = {
  cart: CartResponse | null;
  loading: boolean;
  error: Error | null;
};

const CartContext = createContext<CartContextValue | null>(null);

interface CartProviderProps {
  initialCart?: CartResponse | null;
  children: React.ReactNode;
}

/**
 * Hydrate-only cart state from SSR (`getServerCart` in site/merchant layouts).
 * After mutations, callers should `router.refresh()` so the layout re-fetches.
 */
export function CartProvider({
  initialCart = null,
  children,
}: CartProviderProps) {
  const { user } = useUser();
  const userId = user?.id;
  const [cart, setCart] = useState<CartResponse | null>(
    userId ? initialCart : null
  );
  const [prevUserId, setPrevUserId] = useState(userId);
  const [prevInitialCart, setPrevInitialCart] = useState(initialCart);

  if (userId !== prevUserId) {
    setPrevUserId(userId);
    if (!userId) {
      setCart(null);
    } else {
      setCart(initialCart);
    }
  }

  if (initialCart !== prevInitialCart) {
    setPrevInitialCart(initialCart);
    if (userId) {
      setCart(initialCart);
    }
  }

  const value: CartContextValue = {
    cart: userId ? cart : null,
    loading: false,
    error: null,
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
