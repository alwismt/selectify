import ResetPassword from "@/components/Auth/ResetPassword";
import { getServerUser } from "@/lib/api/getServerUser";
import { redirect } from "next/navigation";
import React from "react";
import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Reset Password | Selectify",
  description: "Choose a new password for your Selectify account",
};

type ResetPasswordPageProps = {
  searchParams: Promise<{ token?: string }>;
};

async function validateResetToken(token: string): Promise<boolean> {
  try {
    const site =
      process.env.NEXT_PUBLIC_SITE_URL?.replace(/\/$/, "") ??
      "http://localhost:3000";
    const url = `${site}/api/auth/reset-password/validate?token=${encodeURIComponent(token)}`;
    const res = await fetch(url, {
      method: "GET",
      cache: "no-store",
    });
    if (!res.ok) {
      return false;
    }
    const data = (await res.json()) as { status?: string };
    return data.status === "ok";
  } catch {
    return false;
  }
}

export default async function ResetPasswordPage({
  searchParams,
}: ResetPasswordPageProps) {
  const user = await getServerUser();
  if (user) {
    redirect("/");
  }

  const params = await searchParams;
  const token = typeof params.token === "string" ? params.token : "";
  const tokenValid = token !== "" ? await validateResetToken(token) : false;

  return (
    <main>
      <ResetPassword token={token} tokenValid={tokenValid} />
    </main>
  );
}
