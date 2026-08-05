"use client";

import React, { useState } from "react";
import {
  ExpressCheckoutElement,
  PaymentElement,
  useElements,
  useStripe,
} from "@stripe/react-stripe-js";
import type { StripeExpressCheckoutElementConfirmEvent } from "@stripe/stripe-js";

type PaymentMethodProps = {
  onPaymentError: (message: string | null) => void;
  onSubmitting: (submitting: boolean) => void;
  /** Called before confirmPayment; return false to abort. */
  onBeforeConfirm: () => Promise<boolean>;
  billingDetails: {
    name: string;
    email: string;
    phone: string;
    address: {
      line1: string;
      line2: string;
      city: string;
      state: string;
      postal_code: string;
      country: string;
    };
  };
};

function hasAvailableWallets(
  methods: Record<string, boolean> | null | undefined
): boolean {
  if (!methods) return false;
  return Object.values(methods).some(Boolean);
}

const PaymentMethod = ({
  onPaymentError,
  onSubmitting,
  onBeforeConfirm,
  billingDetails,
}: PaymentMethodProps) => {
  const stripe = useStripe();
  const elements = useElements();
  const [expressVisible, setExpressVisible] = useState(false);

  const handleExpressConfirm = async (
    _event: StripeExpressCheckoutElementConfirmEvent
  ) => {
    if (!stripe || !elements) return;

    onSubmitting(true);
    onPaymentError(null);

    const ready = await onBeforeConfirm();
    if (!ready) {
      onSubmitting(false);
      return;
    }

    const { error: submitError } = await elements.submit();
    if (submitError) {
      onPaymentError(submitError.message ?? "Payment failed. Please try again.");
      onSubmitting(false);
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
      onSubmitting(false);
    }
  };

  return (
    <div>
      <h3 className="font-medium text-base text-dark mb-2">Payment</h3>

      <div>
        <div className={expressVisible ? "mb-0" : "hidden"}>
          <ExpressCheckoutElement
              options={{
                paymentMethods: {
                  applePay: "always",
                  googlePay: "always",
                  link: "auto",
                  paypal: "never",
                },
              }}
              onConfirm={handleExpressConfirm}
              onReady={({ availablePaymentMethods }) => {
                setExpressVisible(hasAvailableWallets(availablePaymentMethods));
              }}
          />
        </div>

        {expressVisible && (
          <div className="relative my-5 flex items-center">
            <div className="flex-grow border-t border-gray-3" />
            <span className="mx-3 text-custom-sm text-dark-4 shrink-0 uppercase tracking-wide">
              or
            </span>
            <div className="flex-grow border-t border-gray-3" />
          </div>
        )}

        <PaymentElement
          options={{
            layout: {
              type: "accordion",
              defaultCollapsed: false,
              radios: "always",
              spacedAccordionItems: false,
            },
            // layout: "tabs",
            fields: {
              billingDetails: "never",
            },
          }}
        />
      </div>
    </div>
  );
};

export default PaymentMethod;
