import Signin from "@/components/Auth/Signin";
import { getServerUser } from "@/lib/api/getServerUser";
import { redirect } from "next/navigation";
import React from "react";
import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Signin Page | NextCommerce Nextjs E-commerce template",
  description: "This is Signin Page for NextCommerce Template",
  openGraph: {
    title: "Signin Page | NextCommerce Nextjs E-commerce template",
    description: "This is Signin Page for NextCommerce Template",
    url: "/signin",
  },
  twitter: {
    card: "summary_large_image",
    title: "Signin Page | NextCommerce Nextjs E-commerce template",
    description: "This is Signin Page for NextCommerce Template",
  },
};

export default async function SigninPage() {
  const user = await getServerUser();
  if (user) {
    redirect("/");
  }
  return (
    <main>
      <Signin />
    </main>
  );
}
