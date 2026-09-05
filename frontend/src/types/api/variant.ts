import type { ApiProductFile } from "@/types/product";

/** Backend product_variant_attributes item */
export type ApiVariantAttribute = {
  id: number;
  variant_id: number;
  name: string;
  value: string;
};

/** Backend variant shape (GET /products/{product_id}/variants) */
export type ApiVariant = {
  id: number;
  product_id: number;
  sku: string;
  /** Price in integer minor units. */
  price_amount: number;
  created_at: string;
  updated_at: string;
  stock_quantity: number;
  reserved_quantity: number;
  product_variant_attributes: ApiVariantAttribute[];
  files: ApiProductFile[];
};

/** Available sellable quantity for a variant. */
export function availableVariantQuantity(variant: {
  stock_quantity: number;
  reserved_quantity: number;
}): number {
  return Math.max(0, variant.stock_quantity - variant.reserved_quantity);
}
