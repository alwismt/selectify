import { productFileUrl } from "@/lib/api/config";
import type { SiteCurrency } from "@/types/siteConfig";

export type Product = {
  title: string;
  slug: string;
  reviews?: number;
  /** Price in minor units (e.g. cents). */
  price: number;
  /** Discounted price in minor units (e.g. cents). */
  discountedPrice: number;
  id: number;
  description?: string;
  currency?: string;
  minorUnit?: number;
  imgs?: {
    thumbnails: string[];
    previews: string[];
  };
};

/** Backend product_file / variant file item */
export type ApiProductFile = {
  file_id: string;
  product_id: number;
  variant_id: number | null;
  content_type: string;
  position: number;
  is_primary: boolean;
};

/** Backend API product shape (GET /products) */
export type ApiProduct = {
  productId: number;
  sku: string;
  name: string;
  description: string;
  slug: string;
  price: number;
  isActive: boolean;
  inStock: boolean;
  productFile?: ApiProductFile | null;
};

export function mapApiProductToProduct(
  api: ApiProduct,
  currency?: SiteCurrency | null
): Product {
  const url = productFileUrl(api.productFile?.file_id);
  return {
    id: api.productId,
    title: api.name,
    slug: api.slug,
    description: api.description,
    price: api.price,
    discountedPrice: api.price,
    currency: currency?.code,
    minorUnit: currency?.minorUnit,
    reviews: 0,
    imgs: url
      ? {
          thumbnails: [url],
          previews: [url],
        }
      : undefined,
  };
}
