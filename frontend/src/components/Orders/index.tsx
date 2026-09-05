"use client";

import React, { useState } from "react";
import { formatOrderDate, formatOrderMoney } from "@/lib/format";
import { useSiteConfig } from "@/app/context/SiteConfigContext";
import type { Order } from "@/types/api/order";
import type { SiteCurrency } from "@/types/siteConfig";
import SingleOrder from "./SingleOrder";

/** Display shape for SingleOrder / OrderModal (orderId, createdAt, status, title, total) */
type DisplayOrder = {
  orderId: string;
  createdAt: string;
  status: string;
  title: string;
  total: string;
  /** Full order for modal details */
  raw?: Order;
};

function orderToDisplay(
  order: Order,
  siteCurrency: SiteCurrency | null
): DisplayOrder {
  const dateStr = order.items[0]?.created_at
    ? formatOrderDate(order.items[0].created_at)
    : "—";
  const title =
    order.items.length === 0
      ? "No items"
      : order.items.length === 1
        ? order.items[0].sku
        : `${order.items.length} items`;
  return {
    orderId: String(order.id),
    createdAt: dateStr,
    status: order.status,
    title,
    total: formatOrderMoney(order.total, order.currency, siteCurrency),
    raw: order,
  };
}

interface OrdersProps {
  initialOrders?: Order[];
}

const Orders = ({ initialOrders = [] }: OrdersProps) => {
  const { currency: siteCurrency } = useSiteConfig();
  const [orders, setOrders] = useState<DisplayOrder[]>(() =>
    initialOrders.map((o) => orderToDisplay(o, siteCurrency))
  );
  const [prevInitialOrders, setPrevInitialOrders] = useState(initialOrders);

  if (initialOrders !== prevInitialOrders) {
    setPrevInitialOrders(initialOrders);
    setOrders(initialOrders.map((o) => orderToDisplay(o, siteCurrency)));
  }
  return (
    <>
      <div className="w-full overflow-x-auto">
        <div className="min-w-[770px]">
          {orders.length > 0 && (
            <div className="items-center justify-between py-4.5 px-7.5 hidden md:flex ">
              <div className="min-w-[111px]">
                <p className="text-custom-sm text-dark">Order</p>
              </div>
              <div className="min-w-[175px]">
                <p className="text-custom-sm text-dark">Date</p>
              </div>

              <div className="min-w-[128px]">
                <p className="text-custom-sm text-dark">Status</p>
              </div>

              <div className="min-w-[213px]">
                <p className="text-custom-sm text-dark">Title</p>
              </div>

              <div className="min-w-[113px]">
                <p className="text-custom-sm text-dark">Total</p>
              </div>

              <div className="min-w-[113px]">
                <p className="text-custom-sm text-dark">Action</p>
              </div>
            </div>
          )}
          {orders.length > 0 ? (
            orders.map((orderItem) => (
              <SingleOrder
                key={orderItem.orderId}
                orderItem={orderItem}
                smallView={false}
              />
            ))
          ) : (
            <p className="py-9.5 px-4 sm:px-7.5 xl:px-10">
              You don&apos;t have any orders!
            </p>
          )}
        </div>

        {orders.length > 0 &&
          orders.map((orderItem) => (
            <SingleOrder
              key={orderItem.orderId}
              orderItem={orderItem}
              smallView={true}
            />
          ))}
      </div>
    </>
  );
};

export default Orders;
