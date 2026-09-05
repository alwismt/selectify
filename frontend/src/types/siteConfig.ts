/** Site default currency from GET /api/v1/config */
export type SiteCurrency = {
  code: string;
  name: string;
  minorUnit: number;
  isActive: boolean;
  isDefault: boolean;
};

/** Matches GET /api/v1/config JSON (`currency` field). */
export type SiteConfig = {
  currency: SiteCurrency | null;
};
