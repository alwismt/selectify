import React from "react";
import ShopWithSidebar from "@/components/ShopWithSidebar";
import { getProducts } from "@/lib/products";
import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Shop | NextCommerce Nextjs E-commerce template",
  description: "Browse all products",
  openGraph: {
    title: "Shop | NextCommerce Nextjs E-commerce template",
    description: "Browse all products",
    url: "/shop",
  },
  twitter: {
    card: "summary_large_image",
    title: "Shop | NextCommerce Nextjs E-commerce template",
    description: "Browse all products",
  },
};

export default async function ShopPage() {
  let products;
  try {
    products = await getProducts();
  } catch (err) {
    return (
      <main>
        <div className="max-w-[1170px] w-full mx-auto px-4 py-20 text-center">
          <p className="text-dark-4">
            Unable to load products. Please try again later.
          </p>
        </div>
      </main>
    );
  }

  return (
    <main>
      <ShopWithSidebar products={products} breadcrumbCurrent="Shop" />
    </main>
  );
}
