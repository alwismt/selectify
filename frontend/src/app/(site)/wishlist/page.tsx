import React from "react";
import { Wishlist } from "@/components/Wishlist";
import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Wishlist Page | NextCommerce Nextjs E-commerce template",
  description: "This is Wishlist Page for NextCommerce Template",
  openGraph: {
    title: "Wishlist Page | NextCommerce Nextjs E-commerce template",
    description: "This is Wishlist Page for NextCommerce Template",
    url: "/wishlist",
  },
  twitter: {
    card: "summary_large_image",
    title: "Wishlist Page | NextCommerce Nextjs E-commerce template",
    description: "This is Wishlist Page for NextCommerce Template",
  },
};

const WishlistPage = () => {
  return (
    <main>
      <Wishlist />
    </main>
  );
};

export default WishlistPage;
