"use client";

import Breadcrumb from "@/components/Common/Breadcrumb";
import Link from "next/link";
import React, { useState } from "react";

const defaultSuccessMessage =
  "If an account exists with this email address, a password reset link has been sent.";

const ForgotPassword = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);
    const form = e.currentTarget;
    const email = (form.elements.namedItem("email") as HTMLInputElement).value;
    if (!email) {
      setError("Please enter your email.");
      return;
    }
    setLoading(true);
    try {
      const res = await fetch("/api/auth/forgot-password", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email }),
      });
      if (!res.ok) {
        setError("Something went wrong. Please try again.");
        return;
      }
      const data = (await res.json()) as { message?: string };
      setSuccess(
        typeof data.message === "string" ? data.message : defaultSuccessMessage
      );
    } catch {
      setError("Something went wrong. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <Breadcrumb title={"Forgot Password"} pages={["Forgot Password"]} />
      <section className="overflow-hidden py-20 bg-gray-2">
        <div className="max-w-[1170px] w-full mx-auto px-4 sm:px-8 xl:px-0">
          <div className="max-w-[570px] w-full mx-auto rounded-xl bg-white shadow-1 p-4 sm:p-7.5 xl:p-11">
            {success ? (
              <div className="text-center">
                <h2 className="font-semibold text-xl sm:text-2xl xl:text-heading-5 text-dark mb-1.5">
                  Check your email
                </h2>
                <p className="mb-6" role="status">
                  {success}
                </p>
                <Link
                  href="/signin"
                  className="inline-flex justify-center font-medium text-white bg-dark py-3 px-6 rounded-lg ease-out duration-200 hover:bg-blue"
                >
                  Back to Sign In
                </Link>
              </div>
            ) : (
              <>
                <div className="text-center mb-6">
                  <h2 className="font-semibold text-xl sm:text-2xl xl:text-heading-5 text-dark mb-1.5">
                    Forgot your password?
                  </h2>
                  <p>Enter your email below to receive a password reset link.</p>
                </div>

                <form onSubmit={handleSubmit}>
                  {error && (
                    <p className="mb-5 text-red-600 text-sm" role="alert">
                      {error}
                    </p>
                  )}
                  <div className="mb-5">
                    <label htmlFor="email" className="block mb-2.5">
                      Email
                    </label>
                    <input
                      type="email"
                      name="email"
                      id="email"
                      placeholder="Enter your email"
                      required
                      disabled={loading}
                      className="rounded-lg border border-gray-3 bg-gray-1 placeholder:text-dark-5 w-full py-3 px-5 outline-none duration-200 focus:border-transparent focus:shadow-input focus:ring-2 focus:ring-blue/20"
                    />
                  </div>

                  <button
                    type="submit"
                    disabled={loading}
                    className="w-full flex justify-center font-medium text-white bg-dark py-3 px-6 rounded-lg ease-out duration-200 hover:bg-blue mt-5 disabled:opacity-70"
                  >
                    {loading ? "Sending..." : "Send reset link"}
                  </button>

                  <p className="text-center mt-6">
                    Remembered your password?
                    <Link
                      href="/signin"
                      className="text-dark ease-out duration-200 hover:text-blue pl-2"
                    >
                      Sign In
                    </Link>
                  </p>
                </form>
              </>
            )}
          </div>
        </div>
      </section>
    </>
  );
};

export default ForgotPassword;
