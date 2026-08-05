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
  price_amount: number;
  currency: string;
  attributes: Record<string, string>;
  product: CartProduct;
  available_qty: number;
};

/** Single cart line item */
export type CartItem = {
  id: number;
  quantity: number;
  variant: CartVariant;
};

/** GET /cart response */
export type CartResponse = {
  items: CartItem[];
  currency: string;
  subtotal: number;
  item_count: number;
};
