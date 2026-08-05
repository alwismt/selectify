import { productFileUrl } from "@/lib/api/config";

export type Product = {
  title: string;
  slug: string;
  reviews?: number;
  price: number;
  discountedPrice: number;
  id: number;
  description?: string;
  currency?: string;
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
  priceAmount: number;
  currency: string;
  isActive: boolean;
  inStock: boolean;
  productFile?: ApiProductFile | null;
};

export function mapApiProductToProduct(api: ApiProduct): Product {
  const url = productFileUrl(api.productFile?.file_id);
  return {
    id: api.productId,
    title: api.name,
    slug: api.slug,
    description: api.description,
    price: api.priceAmount,
    discountedPrice: api.priceAmount,
    currency: api.currency,
    reviews: 0,
    imgs: url
      ? {
          thumbnails: [url],
          previews: [url],
        }
      : undefined,
  };
}
