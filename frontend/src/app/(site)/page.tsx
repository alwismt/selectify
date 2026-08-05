import Home from "@/components/Home";
import { Metadata } from "next";

export const metadata: Metadata = {
  title: "NextCommerce | Nextjs E-commerce template",
  description: "This is Home for NextCommerce Template",
  openGraph: {
    title: "NextCommerce | Nextjs E-commerce template",
    description: "This is Home for NextCommerce Template",
    url: "/",
  },
  twitter: {
    card: "summary_large_image",
    title: "NextCommerce | Nextjs E-commerce template",
    description: "This is Home for NextCommerce Template",
  },
};

export default function HomePage() {
  return (
    <>
      <Home />
    </>
  );
}
