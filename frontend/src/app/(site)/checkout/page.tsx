import Checkout from "@/components/Checkout";
import { getServerDefaultAddress } from "@/lib/api/getServerDefaultAddress";
import { getServerOrders } from "@/lib/api/getServerOrders";
import { getServerUser } from "@/lib/api/getServerUser";
import { Metadata } from "next";
import { redirect } from "next/navigation";
import React from "react";

export const metadata: Metadata = {
  title: "Checkout | Selectify",
  description: "Complete your order checkout",
  openGraph: {
    title: "Checkout | Selectify",
    description: "Complete your order checkout",
    url: "/checkout",
  },
  twitter: {
    card: "summary_large_image",
    title: "Checkout | Selectify",
    description: "Complete your order checkout",
  },
};

export default async function CheckoutPage({
  searchParams,
}: {
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>;
}) {
  const params = await searchParams;
  const orderIdParam =
    typeof params.orderId === "string" ? params.orderId : undefined;
  if (!orderIdParam?.trim()) {
    redirect("/cart");
  }

  const orderId = Number(orderIdParam);
  if (Number.isNaN(orderId) || orderId <= 0) {
    redirect("/cart");
  }

  const [user, orders, defaultAddress] = await Promise.all([
    getServerUser(),
    getServerOrders(),
    getServerDefaultAddress(),
  ]);

  if (!user) {
    redirect("/signin");
  }

  const order = orders.find((o) => o.id === orderId) ?? null;

  return (
    <main>
      <Checkout
        orderId={orderId}
        initialOrder={order}
        initialAddress={defaultAddress}
      />
    </main>
  );
}
