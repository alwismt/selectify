"use client";

import Breadcrumb from "@/components/Common/Breadcrumb";
import Link from "next/link";
import { useRouter } from "next/navigation";
import React, { useState } from "react";

const invalidResetMessage =
  "This password reset link is invalid or has expired. Please request a new one.";

type ResetPasswordProps = {
  token: string;
  tokenValid: boolean;
};

const InvalidResetLink = () => (
  <>
    <Breadcrumb title={"Reset Password"} pages={["Reset Password"]} />
    <section className="overflow-hidden py-20 bg-gray-2">
      <div className="max-w-[1170px] w-full mx-auto px-4 sm:px-8 xl:px-0">
        <div className="max-w-[570px] w-full mx-auto rounded-xl bg-white shadow-1 p-4 sm:p-7.5 xl:p-11 text-center">
          <h2 className="font-semibold text-xl sm:text-2xl xl:text-heading-5 text-dark mb-1.5">
            Invalid reset link
          </h2>
          <p className="mb-6">{invalidResetMessage}</p>
          <Link
            href="/forgot-password"
            className="inline-flex justify-center font-medium text-white bg-dark py-3 px-6 rounded-lg ease-out duration-200 hover:bg-blue"
          >
            Request a new link
          </Link>
        </div>
      </div>
    </section>
  </>
);

const ResetPassword = ({ token, tokenValid }: ResetPasswordProps) => {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [invalid, setInvalid] = useState(!token || !tokenValid);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError(null);
    const form = e.currentTarget;
    const password = (form.elements.namedItem("password") as HTMLInputElement)
      .value;
    const confirmPassword = (
      form.elements.namedItem("confirm_password") as HTMLInputElement
    ).value;

    if (!password || !confirmPassword) {
      setError("Please enter and confirm your new password.");
      return;
    }
    if (password !== confirmPassword) {
      setError("Passwords do not match.");
      return;
    }
    if (password.length < 8) {
      setError("Password must be at least 8 characters.");
      return;
    }
    if (!token) {
      setInvalid(true);
      return;
    }

    setLoading(true);
    try {
      const res = await fetch("/api/auth/reset-password", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          token,
          password,
          confirm_password: confirmPassword,
        }),
      });
      if (!res.ok) {
        const data = (await res.json().catch(() => null)) as {
          message?: string;
        } | null;
        if (
          res.status === 400 &&
          typeof data?.message === "string" &&
          data.message.includes("invalid or has expired")
        ) {
          setInvalid(true);
          return;
        }
        setError("Something went wrong. Please try again.");
        return;
      }
      router.replace("/signin?reset=1");
    } catch {
      setError("Something went wrong. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  if (invalid) {
    return <InvalidResetLink />;
  }

  return (
    <>
      <Breadcrumb title={"Reset Password"} pages={["Reset Password"]} />
      <section className="overflow-hidden py-20 bg-gray-2">
        <div className="max-w-[1170px] w-full mx-auto px-4 sm:px-8 xl:px-0">
          <div className="max-w-[570px] w-full mx-auto rounded-xl bg-white shadow-1 p-4 sm:p-7.5 xl:p-11">
            <div className="text-center mb-6">
              <h2 className="font-semibold text-xl sm:text-2xl xl:text-heading-5 text-dark mb-1.5">
                Reset your password
              </h2>
              <p>Choose a new password for your account.</p>
            </div>

            <form onSubmit={handleSubmit}>
              {error && (
                <p className="mb-5 text-red-600 text-sm" role="alert">
                  {error}
                </p>
              )}
              <div className="mb-5">
                <label htmlFor="password" className="block mb-2.5">
                  New password
                </label>
                <input
                  type="password"
                  name="password"
                  id="password"
                  placeholder="Enter new password"
                  autoComplete="new-password"
                  required
                  minLength={8}
                  disabled={loading}
                  className="rounded-lg border border-gray-3 bg-gray-1 placeholder:text-dark-5 w-full py-3 px-5 outline-none duration-200 focus:border-transparent focus:shadow-input focus:ring-2 focus:ring-blue/20"
                />
              </div>
              <div className="mb-5">
                <label htmlFor="confirm_password" className="block mb-2.5">
                  Confirm password
                </label>
                <input
                  type="password"
                  name="confirm_password"
                  id="confirm_password"
                  placeholder="Confirm new password"
                  autoComplete="new-password"
                  required
                  minLength={8}
                  disabled={loading}
                  className="rounded-lg border border-gray-3 bg-gray-1 placeholder:text-dark-5 w-full py-3 px-5 outline-none duration-200 focus:border-transparent focus:shadow-input focus:ring-2 focus:ring-blue/20"
                />
              </div>

              <button
                type="submit"
                disabled={loading}
                className="w-full flex justify-center font-medium text-white bg-dark py-3 px-6 rounded-lg ease-out duration-200 hover:bg-blue mt-5 disabled:opacity-70"
              >
                {loading ? "Saving..." : "Reset password"}
              </button>

              <p className="text-center mt-6">
                <Link
                  href="/signin"
                  className="text-dark ease-out duration-200 hover:text-blue"
                >
                  Back to Sign In
                </Link>
              </p>
            </form>
          </div>
        </div>
      </section>
    </>
  );
};

export default ResetPassword;
