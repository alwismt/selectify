import Signup from "@/components/Auth/Signup";
import React from "react";

import { Metadata } from "next";
export const metadata: Metadata = {
  title: "Signup Page | NextCommerce Nextjs E-commerce template",
  description: "This is Signup Page for NextCommerce Template",
  openGraph: {
    title: "Signup Page | NextCommerce Nextjs E-commerce template",
    description: "This is Signup Page for NextCommerce Template",
    url: "/signup",
  },
  twitter: {
    card: "summary_large_image",
    title: "Signup Page | NextCommerce Nextjs E-commerce template",
    description: "This is Signup Page for NextCommerce Template",
  },
};

const SignupPage = () => {
  return (
    <main>
      <Signup />
    </main>
  );
};

export default SignupPage;
