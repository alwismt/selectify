"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { useCart } from "@/app/context/CartContext";
import { useSiteConfig } from "@/app/context/SiteConfigContext";
import { clientOrdersPost } from "@/lib/api/client";
import { setOrderClientSecret } from "@/lib/stripe/clientSecret";
import { formatMoney } from "@/lib/format";

const OrderSummary = () => {
  const router = useRouter();
  const { cart } = useCart();
  const { currency } = useSiteConfig();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const cartItems = cart?.items ?? [];
  const subtotal = cart?.subtotal ?? 0;
  const hasItems = cartItems.length > 0;

  const handleProceedToCheckout = async () => {
    if (!hasItems) {
      setError("Your cart is empty.");
      return;
    }
    setError(null);
    setLoading(true);
    try {
      const data = await clientOrdersPost({});
      if (!data.client_secret) {
        setError("Payment could not be started. Please try again.");
        return;
      }
      setOrderClientSecret(data.id, data.client_secret);
      router.push(`/checkout?orderId=${data.id}`);
      router.refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create order. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="lg:max-w-[455px] w-full">
      <div className="bg-white shadow-1 rounded-[10px]">
        <div className="border-b border-gray-3 py-5 px-4 sm:px-8.5">
          <h3 className="font-medium text-xl text-dark">Order Summary</h3>
        </div>

        <div className="pt-2.5 pb-8.5 px-4 sm:px-8.5">
          <div className="flex items-center justify-between py-5 border-b border-gray-3">
            <div>
              <h4 className="font-medium text-dark">Product</h4>
            </div>
            <div>
              <h4 className="font-medium text-dark text-right">Subtotal</h4>
            </div>
          </div>

          {cartItems.map((item) => (
            <div key={item.id} className="flex items-center justify-between py-5 border-b border-gray-3">
              <div>
                <p className="text-dark">{item.variant.product.name}</p>
              </div>
              <div>
                <p className="text-dark text-right">
                  {currency
                    ? formatMoney(
                        item.variant.price_amount * item.quantity,
                        currency
                      )
                    : "—"}
                </p>
              </div>
            </div>
          ))}

          <div className="flex items-center justify-between pt-5">
            <div>
              <p className="font-medium text-lg text-dark">Total</p>
            </div>
            <div>
              <p className="font-medium text-lg text-dark text-right">
                {currency ? formatMoney(subtotal, currency) : "—"}
              </p>
            </div>
          </div>

          {error && (
            <p className="text-red text-custom-sm mt-2 mb-2">{error}</p>
          )}
          <button
            type="button"
            onClick={handleProceedToCheckout}
            disabled={!hasItems || loading}
            className="w-full flex justify-center font-medium text-white bg-blue py-3 px-6 rounded-md ease-out duration-200 hover:bg-blue-dark mt-7.5 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? "Creating order…" : "Process to Checkout"}
          </button>
        </div>
      </div>
    </div>
  );
};

export default OrderSummary;
