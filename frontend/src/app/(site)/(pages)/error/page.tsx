import React from "react";
import Error from "@/components/Error";

import { Metadata } from "next";
export const metadata: Metadata = {
  title: "Error Page | NextCommerce Nextjs E-commerce template",
  description: "This is Error Page for NextCommerce Template",
  openGraph: {
    title: "Error Page | NextCommerce Nextjs E-commerce template",
    description: "This is Error Page for NextCommerce Template",
    url: "/error",
  },
  twitter: {
    card: "summary_large_image",
    title: "Error Page | NextCommerce Nextjs E-commerce template",
    description: "This is Error Page for NextCommerce Template",
  },
};

const ErrorPage = () => {
  return (
    <main>
      <Error />
    </main>
  );
};

export default ErrorPage;
