import React from "react";
import MailSuccess from "@/components/MailSuccess";

import { Metadata } from "next";
export const metadata: Metadata = {
  title: "Mail Success Page | NextCommerce Nextjs E-commerce template",
  description: "This is Mail Success Page for NextCommerce Template",
  openGraph: {
    title: "Mail Success Page | NextCommerce Nextjs E-commerce template",
    description: "This is Mail Success Page for NextCommerce Template",
    url: "/mail-success",
  },
  twitter: {
    card: "summary_large_image",
    title: "Mail Success Page | NextCommerce Nextjs E-commerce template",
    description: "This is Mail Success Page for NextCommerce Template",
  },
};

const MailSuccessPage = () => {
  return (
    <main>
      <MailSuccess />
    </main>
  );
};

export default MailSuccessPage;
