/** Product summary in cart variant */
export type CartProduct = {
  id: number;
  name: string;
  description: string;
};

/** Variant in cart item */
export type CartVariant = {
  id: number;
  sku: string;
  /** Unit price in integer minor units. */
  price_amount: number;
  attributes: Record<string, string>;
  product: CartProduct;
  /** Server-computed available quantity (stock - reserved). */
  available_qty: number;
};

/** Single cart line item */
export type CartItem = {
  id: number;
  quantity: number;
  variant: CartVariant;
  /**
   * Frontend-enriched image URL (SSR/client), not returned by GET /cart.
   */
  imageUrl?: string | null;
};

/** GET /cart response — monetary fields are integer minor units. */
export type CartResponse = {
  items: CartItem[];
  /** Subtotal in integer minor units. */
  subtotal: number;
  item_count: number;
};
