import React from "react";
import ShopWithoutSidebar from "@/components/ShopWithoutSidebar";

import { Metadata } from "next";
export const metadata: Metadata = {
  title: "Shop Page | NextCommerce Nextjs E-commerce template",
  description: "This is Shop Page for NextCommerce Template",
  openGraph: {
    title: "Shop Page | NextCommerce Nextjs E-commerce template",
    description: "This is Shop Page for NextCommerce Template",
    url: "/shop-without-sidebar",
  },
  twitter: {
    card: "summary_large_image",
    title: "Shop Page | NextCommerce Nextjs E-commerce template",
    description: "This is Shop Page for NextCommerce Template",
  },
};

const ShopWithoutSidebarPage = () => {
  return (
    <main>
      <ShopWithoutSidebar />
    </main>
  );
};

export default ShopWithoutSidebarPage;
