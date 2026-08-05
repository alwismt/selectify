"use client";

import React from "react";
import type { OrderShippingAddressInput } from "@/types/api/order";

export type DeliveryFormValues = {
  fullName: string;
  email: string;
  phone: string;
  line1: string;
  line2: string;
  city: string;
  region: string;
  postal_code: string;
  country_code: string;
};

export type DeliveryFormErrors = Partial<Record<keyof DeliveryFormValues, string>>;

const COUNTRIES: { code: string; name: string }[] = [
  { code: "AU", name: "Australia" },
  { code: "CA", name: "Canada" },
  { code: "DE", name: "Germany" },
  { code: "FR", name: "France" },
  { code: "GB", name: "United Kingdom" },
  { code: "IE", name: "Ireland" },
  { code: "NL", name: "Netherlands" },
  { code: "NZ", name: "New Zealand" },
  { code: "US", name: "United States" },
];

const inputClass =
  "rounded-md border border-gray-3 bg-gray-1 placeholder:text-dark-5 w-full py-1.5 px-3 text-custom-sm outline-none duration-200 focus:border-transparent focus:shadow-input focus:ring-2 focus:ring-blue/20";

type DeliveryFormProps = {
  values: DeliveryFormValues;
  errors: DeliveryFormErrors;
  onChange: (field: keyof DeliveryFormValues, value: string) => void;
};

export function validateDeliveryForm(
  values: DeliveryFormValues
): DeliveryFormErrors {
  const errors: DeliveryFormErrors = {};
  if (!values.fullName.trim()) errors.fullName = "Name is required";
  if (!values.email.trim()) errors.email = "Email is required";
  if (!values.line1.trim()) errors.line1 = "Address is required";
  if (!values.city.trim()) errors.city = "City is required";
  if (!values.postal_code.trim()) errors.postal_code = "Postal code is required";
  if (!values.country_code.trim() || values.country_code.length !== 2) {
    errors.country_code = "Country is required";
  }
  const country = values.country_code.trim().toUpperCase();
  if (["US", "CA", "AU"].includes(country) && !values.region.trim()) {
    errors.region = "State / Region is required";
  }
  return errors;
}

export function toShippingAddressInput(
  values: DeliveryFormValues
): OrderShippingAddressInput {
  return {
    line1: values.line1.trim(),
    line2: values.line2.trim() || undefined,
    city: values.city.trim(),
    region: values.region.trim() || undefined,
    postal_code: values.postal_code.trim(),
    country_code: values.country_code.trim().toUpperCase(),
    phone: values.phone.trim() || undefined,
  };
}

const DeliveryForm = ({ values, errors, onChange }: DeliveryFormProps) => {
  return (
    <div>
      <h2 className="font-medium text-base text-dark mb-2">Shipping address</h2>

      <div className="space-y-2">
        <div>
          <label htmlFor="fullName" className="block mb-0.5 text-custom-xs text-dark">
            Full name <span className="text-red">*</span>
          </label>
          <input
            type="text"
            id="fullName"
            name="fullName"
            autoComplete="name"
            value={values.fullName}
            onChange={(e) => onChange("fullName", e.target.value)}
            className={inputClass}
          />
          {errors.fullName && (
            <p className="text-red text-custom-xs mt-1">{errors.fullName}</p>
          )}
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
          <div>
            <label htmlFor="email" className="block mb-0.5 text-custom-xs text-dark">
              Email <span className="text-red">*</span>
            </label>
            <input
              type="email"
              id="email"
              name="email"
              autoComplete="email"
              value={values.email}
              onChange={(e) => onChange("email", e.target.value)}
              className={inputClass}
            />
            {errors.email && (
              <p className="text-red text-custom-xs mt-1">{errors.email}</p>
            )}
          </div>
          <div>
            <label htmlFor="phone" className="block mb-0.5 text-custom-xs text-dark">
              Phone
            </label>
            <input
              type="tel"
              id="phone"
              name="phone"
              autoComplete="tel"
              value={values.phone}
              onChange={(e) => onChange("phone", e.target.value)}
              className={inputClass}
            />
          </div>
        </div>

        <div>
          <label htmlFor="line1" className="block mb-0.5 text-custom-xs text-dark">
            Address <span className="text-red">*</span>
          </label>
          <input
            type="text"
            id="line1"
            name="line1"
            autoComplete="address-line1"
            placeholder="Street address"
            value={values.line1}
            onChange={(e) => onChange("line1", e.target.value)}
            className={inputClass}
          />
          {errors.line1 && (
            <p className="text-red text-custom-xs mt-1">{errors.line1}</p>
          )}
        </div>

        <div>
          <label htmlFor="line2" className="block mb-0.5 text-custom-xs text-dark">
            Apartment, suite, etc.
          </label>
          <input
            type="text"
            id="line2"
            name="line2"
            autoComplete="address-line2"
            value={values.line2}
            onChange={(e) => onChange("line2", e.target.value)}
            className={inputClass}
          />
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
          <div>
            <label htmlFor="city" className="block mb-0.5 text-custom-xs text-dark">
              City <span className="text-red">*</span>
            </label>
            <input
              type="text"
              id="city"
              name="city"
              autoComplete="address-level2"
              value={values.city}
              onChange={(e) => onChange("city", e.target.value)}
              className={inputClass}
            />
            {errors.city && (
              <p className="text-red text-custom-xs mt-1">{errors.city}</p>
            )}
          </div>
          <div>
            <label htmlFor="postal_code" className="block mb-0.5 text-custom-xs text-dark">
              Postal code <span className="text-red">*</span>
            </label>
            <input
              type="text"
              id="postal_code"
              name="postal_code"
              autoComplete="postal-code"
              value={values.postal_code}
              onChange={(e) => onChange("postal_code", e.target.value)}
              className={inputClass}
            />
            {errors.postal_code && (
              <p className="text-red text-custom-xs mt-1">{errors.postal_code}</p>
            )}
          </div>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
          <div>
            <label htmlFor="region" className="block mb-0.5 text-custom-xs text-dark">
              State / Region
            </label>
            <input
              type="text"
              id="region"
              name="region"
              autoComplete="address-level1"
              value={values.region}
              onChange={(e) => onChange("region", e.target.value)}
              className={inputClass}
            />
            {errors.region && (
              <p className="text-red text-custom-xs mt-1">{errors.region}</p>
            )}
          </div>
          <div>
            <label htmlFor="country_code" className="block mb-0.5 text-custom-xs text-dark">
              Country <span className="text-red">*</span>
            </label>
            <div className="relative">
              <select
                id="country_code"
                name="country_code"
                autoComplete="country"
                value={values.country_code}
                onChange={(e) => onChange("country_code", e.target.value)}
                className={`${inputClass} appearance-none pr-9`}
              >
                <option value="">Select country</option>
                {COUNTRIES.map((c) => (
                  <option key={c.code} value={c.code}>
                    {c.name}
                  </option>
                ))}
              </select>
              <span className="pointer-events-none absolute right-4 top-1/2 -translate-y-1/2 text-dark-4">
                <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
                  <path
                    fillRule="evenodd"
                    clipRule="evenodd"
                    d="M2.29664 5.24176C2.56035 4.9679 2.99354 4.96047 3.2674 5.22418L8.00161 9.80261L12.7323 5.25705C13.0061 4.99333 13.4393 5.00076 13.703 5.27462C13.9668 5.54848 13.9593 5.98167 13.6855 6.24539L8.45019 11.2774C8.18291 11.5341 7.76339 11.5341 7.4961 11.2774L2.24895 6.24539C1.9751 5.98167 1.96767 5.54848 2.29664 5.24176Z"
                    fill="currentColor"
                  />
                </svg>
              </span>
            </div>
            {errors.country_code && (
              <p className="text-red text-custom-xs mt-1">{errors.country_code}</p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default DeliveryForm;
