import React, { useEffect, useState } from "react";
import { clientOrdersGet } from "@/lib/api/client";
import { formatOrderDate, formatOrderTotal } from "@/lib/format";
import type { Order } from "@/types/api/order";
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

function orderToDisplay(order: Order): DisplayOrder {
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
    total: formatOrderTotal(order.total, order.currency),
    raw: order,
  };
}

interface OrdersProps {
  initialOrders?: Order[];
}

const Orders = ({ initialOrders }: OrdersProps) => {
  const [orders, setOrders] = useState<DisplayOrder[]>(() =>
    initialOrders ? initialOrders.map(orderToDisplay) : []
  );

  useEffect(() => {
    if (initialOrders !== undefined) return;
    clientOrdersGet()
      .then((data) => setOrders(data.map(orderToDisplay)))
      .catch((err) => {
        console.log(err.message);
      });
  }, [initialOrders]);

  return (
    <>
      <div className="w-full overflow-x-auto">
        <div className="min-w-[770px]">
          {/* <!-- order item --> */}
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
            orders.map((orderItem, key) => (
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
