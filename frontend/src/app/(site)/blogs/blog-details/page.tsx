import BlogDetails from "@/components/BlogDetails";
import React from "react";

import { Metadata } from "next";
export const metadata: Metadata = {
  title: "Blog Details Page | NextCommerce Nextjs E-commerce template",
  description: "This is Blog Details Page for NextCommerce Template",
  openGraph: {
    title: "Blog Details Page | NextCommerce Nextjs E-commerce template",
    description: "This is Blog Details Page for NextCommerce Template",
    url: "/blogs/blog-details",
  },
  twitter: {
    card: "summary_large_image",
    title: "Blog Details Page | NextCommerce Nextjs E-commerce template",
    description: "This is Blog Details Page for NextCommerce Template",
  },
};

const BlogDetailsPage = () => {
  return (
    <main>
      <BlogDetails />
    </main>
  );
};

export default BlogDetailsPage;
