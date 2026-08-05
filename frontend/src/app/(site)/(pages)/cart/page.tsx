import React from "react";
import Cart from "@/components/Cart";

import { Metadata } from "next";
export const metadata: Metadata = {
  title: "Cart Page | NextCommerce Nextjs E-commerce template",
  description: "This is Cart Page for NextCommerce Template",
  openGraph: {
    title: "Cart Page | NextCommerce Nextjs E-commerce template",
    description: "This is Cart Page for NextCommerce Template",
    url: "/cart",
  },
  twitter: {
    card: "summary_large_image",
    title: "Cart Page | NextCommerce Nextjs E-commerce template",
    description: "This is Cart Page for NextCommerce Template",
  },
};

const CartPage = () => {
  return (
    <>
      <Cart />
    </>
  );
};

export default CartPage;
