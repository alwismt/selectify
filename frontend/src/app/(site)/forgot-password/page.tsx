import ForgotPassword from "@/components/Auth/ForgotPassword";
import { getServerUser } from "@/lib/api/getServerUser";
import { redirect } from "next/navigation";
import React from "react";
import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Forgot Password | Selectify",
  description: "Request a password reset link for your Selectify account",
};

export default async function ForgotPasswordPage() {
  const user = await getServerUser();
  if (user) {
    redirect("/");
  }
  return (
    <main>
      <ForgotPassword />
    </main>
  );
}
