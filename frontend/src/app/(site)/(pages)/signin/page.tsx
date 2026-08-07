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

type SigninPageProps = {
  searchParams: Promise<{ reset?: string }>;
};

export default async function SigninPage({ searchParams }: SigninPageProps) {
  const user = await getServerUser();
  if (user) {
    redirect("/");
  }
  const params = await searchParams;
  return (
    <main>
      <Signin resetSuccess={params.reset === "1"} />
    </main>
  );
}
