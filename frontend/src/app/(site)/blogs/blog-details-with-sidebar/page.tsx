import React from "react";
import BlogDetailsWithSidebar from "@/components/BlogDetailsWithSidebar";

import { Metadata } from "next";
export const metadata: Metadata = {
  title: "Blog Details Page | NextCommerce Nextjs E-commerce template",
  description: "This is Blog Details Page for NextCommerce Template",
  openGraph: {
    title: "Blog Details Page | NextCommerce Nextjs E-commerce template",
    description: "This is Blog Details Page for NextCommerce Template",
    url: "/blogs/blog-details-with-sidebar",
  },
  twitter: {
    card: "summary_large_image",
    title: "Blog Details Page | NextCommerce Nextjs E-commerce template",
    description: "This is Blog Details Page for NextCommerce Template",
  },
};

const BlogDetailsWithSidebarPage = () => {
  return (
    <main>
      <BlogDetailsWithSidebar />
    </main>
  );
};

export default BlogDetailsWithSidebarPage;
