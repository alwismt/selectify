import React from "react";
import BlogGridWithSidebar from "@/components/BlogGridWithSidebar";

import { Metadata } from "next";
export const metadata: Metadata = {
  title: "Blog Grid Page | NextCommerce Nextjs E-commerce template",
  description: "This is Blog Grid Page for NextCommerce Template",
  openGraph: {
    title: "Blog Grid Page | NextCommerce Nextjs E-commerce template",
    description: "This is Blog Grid Page for NextCommerce Template",
    url: "/blogs/blog-grid-with-sidebar",
  },
  twitter: {
    card: "summary_large_image",
    title: "Blog Grid Page | NextCommerce Nextjs E-commerce template",
    description: "This is Blog Grid Page for NextCommerce Template",
  },
};

const BlogGridWithSidebarPage = () => {
  return (
    <>
      <BlogGridWithSidebar />
    </>
  );
};

export default BlogGridWithSidebarPage;
