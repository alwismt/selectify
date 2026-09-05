/** Single line item in an order (GET/POST response) */
export type OrderItem = {
  id: number;
  order_id: number;
  variant_id: number;
  sku: string;
  /** Unit price in integer minor units. */
  unit_price: number;
  currency: string;
  quantity: number;
  attributes: Record<string, string>;
  created_at: string;
};

/**
 * Order (GET array element or POST response).
 * Monetary fields (subtotal, shipping, discount, total, unit_price) are integer minor units.
 * `currency` is a financial snapshot — format with order.currency, not a later site currency.
 */
export type Order = {
  id: number;
  user_id: number;
  status: string;
  currency: string;
  subtotal: number;
  shipping: number;
  discount: number;
  total: number;
  items: OrderItem[];
  /** Present on POST create when a Stripe PaymentIntent is created */
  client_secret?: string;
};

/** Shipping address payload for PUT /orders/{id}/address */
export type OrderShippingAddressInput = {
  line1: string;
  line2?: string;
  city: string;
  region?: string;
  postal_code: string;
  country_code: string;
  phone?: string;
};

/** Saved order address (API response) */
export type OrderAddress = {
  id: number;
  order_id: number;
  type: string;
  phone?: string;
  line1: string;
  line2?: string;
  city: string;
  region?: string;
  postal_code: string;
  country_code: string;
  created_at: string;
};
