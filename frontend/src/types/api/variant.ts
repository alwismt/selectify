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
  price_amount: number;
  currency: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  stock_quantity: number;
  reserved_quantity: number;
  product_variant_attributes: ApiVariantAttribute[];
  files?: ApiProductFile[];
};
