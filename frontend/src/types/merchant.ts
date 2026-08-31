export type Merchant = {
  merchantId: number;
  name: string;
  slug: string;
  description?: string | null;
  logo?: string | null;
  countryCode: string;
  status: string;
  verificationStatus: string;
  createdAt: string;
  updatedAt: string;
};

export type UpdateMerchantInput = {
  name?: string;
  description?: string;
  countryCode?: string;
};

export type MerchantCountry = {
  code: string;
  name: string;
};
