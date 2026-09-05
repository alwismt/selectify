"use client";

import React, { useState } from "react";
import Link from "next/link";
import { formatOrderMoney } from "@/lib/format";
import { useSiteConfig } from "@/app/context/SiteConfigContext";
import type { Order } from "@/types/api/order";

type CheckoutOrderSummaryProps = {
  order: Order | null;
  orderLoading: boolean;
  orderError: string | null;
  /** Mobile collapsible strip vs full desktop panel. */
  variant?: "panel" | "mobile";
};

const CheckoutOrderSummary = ({
  order,
  orderLoading,
  orderError,
  variant = "panel",
}: CheckoutOrderSummaryProps) => {
  const [open, setOpen] = useState(false);
  const { currency: siteCurrency } = useSiteConfig();

  const totalLabel = order
    ? formatOrderMoney(order.total, order.currency, siteCurrency)
    : null;

  const lines = (
    <>
      {orderLoading && (
        <p className="py-4 text-white/70">Loading order…</p>
      )}
      {orderError && !orderLoading && (
        <p className="py-4 text-red-light-3">{orderError}</p>
      )}
      {order && !orderLoading && (
        <>
          <ul className="space-y-4 mb-6">
            {order.items.map((item) => (
              <li
                key={item.id}
                className="flex items-start justify-between gap-4"
              >
                <p className="text-white/90 text-custom-sm pr-2">
                  {item.sku}
                  {item.quantity > 1 ? ` × ${item.quantity}` : ""}
                </p>
                <p className="text-white text-custom-sm text-right shrink-0">
                  {formatOrderMoney(
                    item.unit_price * item.quantity,
                    order.currency,
                    siteCurrency
                  )}
                </p>
              </li>
            ))}
          </ul>

          <div className="border-t border-white/20 pt-4 space-y-3">
            <div className="flex items-center justify-between">
              <p className="text-white/80 text-custom-sm">Subtotal</p>
              <p className="text-white text-custom-sm">
                {formatOrderMoney(order.subtotal, order.currency, siteCurrency)}
              </p>
            </div>

            {order.shipping > 0 && (
              <div className="flex items-center justify-between">
                <p className="text-white/80 text-custom-sm">Shipping</p>
                <p className="text-white text-custom-sm">
                  {formatOrderMoney(
                    order.shipping,
                    order.currency,
                    siteCurrency
                  )}
                </p>
              </div>
            )}

            {order.discount > 0 && (
              <div className="flex items-center justify-between">
                <p className="text-white/80 text-custom-sm">Discount</p>
                <p className="text-white text-custom-sm">
                  −
                  {formatOrderMoney(
                    order.discount,
                    order.currency,
                    siteCurrency
                  )}
                </p>
              </div>
            )}

            <div className="flex items-center justify-between pt-2">
              <p className="font-medium text-white">Total due</p>
              <p className="font-medium text-white text-lg">
                {formatOrderMoney(order.total, order.currency, siteCurrency)}
              </p>
            </div>
          </div>
        </>
      )}
    </>
  );

  if (variant === "mobile") {
    return (
      <div className="bg-blue text-white lg:hidden">
        <button
          type="button"
          className="w-full flex items-center justify-between px-4 sm:px-6 py-4"
          onClick={() => setOpen((v) => !v)}
          aria-expanded={open}
        >
          <span className="text-custom-sm text-white/90">
            {open ? "Hide order summary" : "Show order summary"}
          </span>
          <span className="flex items-center gap-2">
            {totalLabel && (
              <span className="font-semibold text-white">{totalLabel}</span>
            )}
            <svg
              className={`transition-transform ${open ? "rotate-180" : ""}`}
              width="16"
              height="16"
              viewBox="0 0 16 16"
              fill="none"
              aria-hidden
            >
              <path
                fillRule="evenodd"
                clipRule="evenodd"
                d="M2.29664 5.24176C2.56035 4.9679 2.99354 4.96047 3.2674 5.22418L8.00161 9.80261L12.7323 5.25705C13.0061 4.99333 13.4393 5.00076 13.703 5.27462C13.9668 5.54848 13.9593 5.98167 13.6855 6.24539L8.45019 11.2774C8.18291 11.5341 7.76339 11.5341 7.4961 11.2774L2.24895 6.24539C1.9751 5.98167 1.96767 5.54848 2.29664 5.24176Z"
                fill="currentColor"
              />
            </svg>
          </span>
        </button>
        {open && (
          <div className="px-4 sm:px-6 pb-5 border-t border-white/15">{lines}</div>
        )}
      </div>
    );
  }

  return (
    <div className="h-full flex flex-col text-white">
      <div className="mb-8">
        <Link
          href="/cart"
          className="inline-flex items-center gap-2 text-custom-sm text-white/80 hover:text-white transition-colors mb-6"
        >
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden>
            <path
              d="M10 3L5 8L10 13"
              stroke="currentColor"
              strokeWidth="1.5"
              strokeLinecap="round"
              strokeLinejoin="round"
            />
          </svg>
          Back to cart
        </Link>

        <p className="text-white/80 text-custom-sm mb-1">Pay</p>
        <p className="font-semibold text-heading-6 sm:text-heading-5 text-white tracking-tight">
          {totalLabel ?? "—"}
        </p>
      </div>

      {lines}
    </div>
  );
};

export default CheckoutOrderSummary;
