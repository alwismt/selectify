"use client";

import React, { FormEvent, useEffect, useState } from "react";
import Link from "next/link";
import { Elements, useElements, useStripe } from "@stripe/react-stripe-js";
import { useUser } from "@/app/context/UserContext";
import { clientOrderAddressPut } from "@/lib/api/client";
import { formatOrderTotal } from "@/lib/format";
import { getOrderClientSecret } from "@/lib/stripe/clientSecret";
import { stripePromise } from "@/lib/stripe/stripePromise";
import type { Order } from "@/types/api/order";
import type { UserAddress } from "@/types/api/userAddress";
import type { User } from "@/types/user";
import CheckoutOrderSummary from "./CheckoutOrderSummary";
import DeliveryForm, {
  type DeliveryFormErrors,
  type DeliveryFormValues,
  toShippingAddressInput,
  validateDeliveryForm,
} from "./DeliveryForm";
import PaymentMethod from "./PaymentMethod";

function buildInitialDelivery(
  user: User | null,
  address: UserAddress | null
): DeliveryFormValues {
  return {
    fullName: [user?.first_name, user?.last_name].filter(Boolean).join(" "),
    email: user?.email ?? "",
    phone: address?.phone || user?.phone || "",
    line1: address?.line1 ?? "",
    line2: address?.line2 ?? "",
    city: address?.city ?? "",
    region: address?.region ?? "",
    postal_code: address?.postal_code ?? "",
    country_code: address?.country_code ?? "",
  };
}

function PaymentSection({
  order,
  delivery,
  onPaymentError,
  submitting,
  setSubmitting,
  onBeforeConfirm,
}: {
  order: Order | null;
  delivery: DeliveryFormValues;
  onPaymentError: (message: string | null) => void;
  submitting: boolean;
  setSubmitting: (v: boolean) => void;
  onBeforeConfirm: () => Promise<boolean>;
}) {
  const stripe = useStripe();
  const elements = useElements();

  const billingDetails = {
    name: delivery.fullName.trim(),
    email: delivery.email.trim(),
    phone: delivery.phone.trim(),
    address: {
      line1: delivery.line1.trim(),
      line2: delivery.line2.trim(),
      city: delivery.city.trim(),
      state: delivery.region.trim(),
      postal_code: delivery.postal_code.trim(),
      country: delivery.country_code.trim().toUpperCase(),
    },
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!stripe || !elements) return;

    setSubmitting(true);
    onPaymentError(null);

    const ready = await onBeforeConfirm();
    if (!ready) {
      setSubmitting(false);
      return;
    }

    const { error } = await stripe.confirmPayment({
      elements,
      confirmParams: {
        return_url: `${window.location.origin}/my-account`,
        payment_method_data: {
          billing_details: billingDetails,
        },
      },
    });

    if (error) {
      onPaymentError(error.message ?? "Payment failed. Please try again.");
      setSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <PaymentMethod
        onPaymentError={onPaymentError}
        onSubmitting={setSubmitting}
        onBeforeConfirm={onBeforeConfirm}
        billingDetails={billingDetails}
      />

      <button
        type="submit"
        disabled={!stripe || !elements || submitting || !order}
        className="w-full flex justify-center font-medium text-white bg-blue py-3.5 px-6 rounded-md ease-out duration-200 hover:bg-blue-dark disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {submitting
          ? "Processing…"
          : order
            ? `Pay ${formatOrderTotal(order.total, order.currency)}`
            : "Pay"}
      </button>
    </form>
  );
}

type CheckoutProps = {
  orderId: number;
  initialOrder: Order | null;
  initialAddress: UserAddress | null;
};

const Checkout = ({
  orderId,
  initialOrder,
  initialAddress,
}: CheckoutProps) => {
  const { user } = useUser();
  const [clientSecret, setClientSecret] = useState<string | null>(null);
  const [secretError, setSecretError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [paymentError, setPaymentError] = useState<string | null>(null);
  const [deliveryErrors, setDeliveryErrors] = useState<DeliveryFormErrors>({});
  const [delivery, setDelivery] = useState<DeliveryFormValues>(() =>
    buildInitialDelivery(user, initialAddress)
  );

  const orderError = initialOrder ? null : "Order not found.";

  useEffect(() => {
    const secret = getOrderClientSecret(orderId);
    if (!secret) {
      setSecretError(
        "Payment session expired. Please return to your cart and try checkout again."
      );
      return;
    }
    setClientSecret(secret);
  }, [orderId]);

  const handleDeliveryChange = (
    field: keyof DeliveryFormValues,
    value: string
  ) => {
    setDelivery((prev) => ({ ...prev, [field]: value }));
    if (deliveryErrors[field]) {
      setDeliveryErrors((prev) => {
        const next = { ...prev };
        delete next[field];
        return next;
      });
    }
  };

  const saveDeliveryAddress = async (): Promise<boolean> => {
    const errors = validateDeliveryForm(delivery);
    setDeliveryErrors(errors);
    if (Object.keys(errors).length > 0) {
      setPaymentError("Please complete the delivery details.");
      return false;
    }
    if (!initialOrder) {
      setPaymentError("Order not found.");
      return false;
    }

    try {
      await clientOrderAddressPut(
        initialOrder.id,
        toShippingAddressInput(delivery)
      );
      setPaymentError(null);
      return true;
    } catch {
      setPaymentError("Failed to save delivery address. Please try again.");
      return false;
    }
  };

  return (
    <div className="min-h-screen bg-white lg:bg-transparent lg:grid lg:grid-cols-[minmax(0,0.9fr)_minmax(0,1.1fr)]">
      {/* Desktop summary panel */}
      <aside className="hidden lg:flex bg-blue min-h-screen">
        <div className="w-full max-w-md ml-auto px-10 xl:px-14 py-10">
          <CheckoutOrderSummary
            order={initialOrder}
            orderLoading={false}
            orderError={orderError}
            variant="panel"
          />
        </div>
      </aside>

      {/* Form panel */}
      <div className="bg-white min-h-screen flex flex-col">
        <CheckoutOrderSummary
          order={initialOrder}
          orderLoading={false}
          orderError={orderError}
          variant="mobile"
        />

        <div className="flex-1 w-full max-w-lg mx-auto px-4 sm:px-6 py-5 sm:py-6 space-y-4">
          <div className="lg:hidden">
            <Link
              href="/cart"
              className="inline-flex items-center gap-2 text-custom-sm text-dark-4 hover:text-blue transition-colors"
            >
              <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden>
                <path
                  d="M10 3L5 8L10 13"
                  stroke="currentColor"
                  strokeWidth="1.5"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
              </svg>
              Back to cart
            </Link>
          </div>

          <DeliveryForm
            values={delivery}
            errors={deliveryErrors}
            onChange={handleDeliveryChange}
          />

          {secretError && (
            <div className="rounded-md border border-red-light-3 bg-red-light-6 p-5 text-center">
              <p className="text-red mb-4">{secretError}</p>
              <Link
                href="/cart"
                className="inline-flex font-medium text-white bg-blue py-3 px-6 rounded-md hover:bg-blue-dark"
              >
                Back to cart
              </Link>
            </div>
          )}

          {!secretError && !stripePromise && (
            <div className="rounded-md border border-red-light-3 bg-red-light-6 p-5 text-center">
              <p className="text-red">
                Payment is not configured. Missing publishable key.
              </p>
            </div>
          )}

          {!secretError && clientSecret && stripePromise && (
            <Elements
              stripe={stripePromise}
              options={{
                clientSecret,
                appearance: {
                  theme: "stripe",
                  variables: {
                    colorPrimary: "#3C50E0",
                    borderRadius: "6px",
                  },
                },
              }}
            >
              <PaymentSection
                order={initialOrder}
                delivery={delivery}
                onPaymentError={setPaymentError}
                submitting={submitting}
                setSubmitting={setSubmitting}
                onBeforeConfirm={saveDeliveryAddress}
              />
            </Elements>
          )}

          {paymentError && (
            <p className="text-red text-custom-sm">{paymentError}</p>
          )}
        </div>
      </div>
    </div>
  );
};

export default Checkout;
