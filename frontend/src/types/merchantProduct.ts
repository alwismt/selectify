import type { ApiProductFile } from "@/types/product";

/**
 * Merchant API product shape (GET /merchant/products, GET /merchant/products/{product_id}).
 * Separate from public ApiProduct: description/slug are nullable; merchantId is present.
 * Currency is site-level (GET /config), not on the product.
 */
export type MerchantApiProduct = {
  productId: number;
  sku: string;
  name: string;
  description: string | null;
  slug?: string | null;
  /** Price in minor units (e.g. cents). */
  price: number;
  isActive: boolean;
  inStock: boolean;
  merchantId?: number | null;
  productFile?: ApiProductFile | null;
};

export type MerchantApiVariantAttribute = {
  id: number;
  variant_id: number;
  name: string;
  value: string;
};

export type MerchantApiVariantFile = {
  file_id: string;
  product_id: number;
  variant_id: number;
  content_type: string;
  position: number;
  is_primary: boolean;
};

/** Merchant API variant shape (GET /merchant/products/{product_id}/variants). */
export type MerchantApiProductVariant = {
  id: number;
  product_id: number;
  sku: string;
  /** Price in integer minor units. */
  price_amount: number;
  created_at: string;
  updated_at: string;
  stock_quantity: number;
  reserved_quantity: number;
  product_variant_attributes: MerchantApiVariantAttribute[];
  files: MerchantApiVariantFile[];
};
