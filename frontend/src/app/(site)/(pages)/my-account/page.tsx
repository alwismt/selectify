import MyAccount from "@/components/MyAccount";
import { getServerOrders } from "@/lib/api/getServerOrders";
import { getServerUser } from "@/lib/api/getServerUser";
import { Metadata } from "next";
import { redirect } from "next/navigation";
import React from "react";

export const metadata: Metadata = {
  title: "My Account | Selectify",
  description: "Manage your account details and view your orders",
  openGraph: {
    title: "My Account | Selectify",
    description: "Manage your account details and view your orders",
    url: "/my-account",
  },
  twitter: {
    card: "summary_large_image",
    title: "My Account | Selectify",
    description: "Manage your account details and view your orders",
  },
};

export default async function MyAccountPage() {
  const [user, orders] = await Promise.all([
    getServerUser(),
    getServerOrders(),
  ]);
  if (!user) {
    redirect("/signin");
  }
  return (
    <main>
      <MyAccount initialOrders={orders} />
    </main>
  );
}
