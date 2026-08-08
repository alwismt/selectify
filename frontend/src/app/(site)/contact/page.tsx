import Contact from "@/components/Contact";

import { Metadata } from "next";
export const metadata: Metadata = {
  title: "Contact Page | NextCommerce Nextjs E-commerce template",
  description: "This is Contact Page for NextCommerce Template",
  openGraph: {
    title: "Contact Page | NextCommerce Nextjs E-commerce template",
    description: "This is Contact Page for NextCommerce Template",
    url: "/contact",
  },
  twitter: {
    card: "summary_large_image",
    title: "Contact Page | NextCommerce Nextjs E-commerce template",
    description: "This is Contact Page for NextCommerce Template",
  },
};

const ContactPage = () => {
  return (
    <main>
      <Contact />
    </main>
  );
};

export default ContactPage;
