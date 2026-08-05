import React from "react";
import BlogGrid from "@/components/BlogGrid";

import { Metadata } from "next";
export const metadata: Metadata = {
  title: "Blog Grid Page | NextCommerce Nextjs E-commerce template",
  description: "This is Blog Grid Page for NextCommerce Template",
  openGraph: {
    title: "Blog Grid Page | NextCommerce Nextjs E-commerce template",
    description: "This is Blog Grid Page for NextCommerce Template",
    url: "/blogs/blog-grid",
  },
  twitter: {
    card: "summary_large_image",
    title: "Blog Grid Page | NextCommerce Nextjs E-commerce template",
    description: "This is Blog Grid Page for NextCommerce Template",
  },
};

const BlogGridPage = () => {
  return (
    <main>
      <BlogGrid />
    </main>
  );
};

export default BlogGridPage;
