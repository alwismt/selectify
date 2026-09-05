"use client";

import React, { createContext, useContext, useMemo } from "react";
import type { SiteConfig, SiteCurrency } from "@/types/siteConfig";

type SiteConfigStatus = "ready" | "unavailable";

type SiteConfigContextValue = {
  config: SiteConfig | null;
  currency: SiteCurrency | null;
  status: SiteConfigStatus;
};

const SiteConfigContext = createContext<SiteConfigContextValue | undefined>(
  undefined
);

export function useSiteConfig(): SiteConfigContextValue {
  const context = useContext(SiteConfigContext);
  if (context === undefined) {
    throw new Error("useSiteConfig must be used within a SiteConfigProvider");
  }
  return context;
}

interface SiteConfigProviderProps {
  initialConfig: SiteConfig | null;
  children: React.ReactNode;
}

export function SiteConfigProvider({
  initialConfig,
  children,
}: SiteConfigProviderProps) {
  const value = useMemo<SiteConfigContextValue>(() => {
    const currency = initialConfig?.currency ?? null;
    return {
      config: initialConfig,
      currency,
      status: currency ? "ready" : "unavailable",
    };
  }, [initialConfig]);

  return (
    <SiteConfigContext.Provider value={value}>
      {children}
    </SiteConfigContext.Provider>
  );
}
