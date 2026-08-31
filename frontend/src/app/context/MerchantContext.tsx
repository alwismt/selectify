"use client";

import React, { createContext, useContext, useMemo, useState } from "react";
import type { Merchant } from "@/types/merchant";

interface MerchantContextType {
  merchant: Merchant | null;
  setMerchant: (merchant: Merchant | null) => void;
}

const MerchantContext = createContext<MerchantContextType | undefined>(
  undefined
);

export function useMerchant() {
  const context = useContext(MerchantContext);
  if (context === undefined) {
    throw new Error("useMerchant must be used within a MerchantProvider");
  }
  return context;
}

interface MerchantProviderProps {
  initialMerchant: Merchant | null;
  children: React.ReactNode;
}

export function MerchantProvider({
  initialMerchant,
  children,
}: MerchantProviderProps) {
  const [merchant, setMerchant] = useState<Merchant | null>(initialMerchant);

  const value = useMemo(
    () => ({ merchant, setMerchant }),
    [merchant]
  );

  return (
    <MerchantContext.Provider value={value}>{children}</MerchantContext.Provider>
  );
}
