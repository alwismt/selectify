"use client";

import { FormEvent, useEffect, useMemo, useState } from "react";
import { useMerchant } from "@/app/context/MerchantContext";
import { clientUpdateMerchant } from "@/lib/api/client";
import type { MerchantCountry } from "@/types/merchant";
import MerchantLogoUpload from "./MerchantLogoUpload";

type FormState = {
  name: string;
  description: string;
  countryCode: string;
};

type StoreSettingsFormProps = {
  countries: MerchantCountry[];
};

export default function StoreSettingsForm({
  countries,
}: StoreSettingsFormProps) {
  const { merchant, setMerchant } = useMerchant();
  const [form, setForm] = useState<FormState>({
    name: "",
    description: "",
    countryCode: "",
  });
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  useEffect(() => {
    if (!merchant) return;
    setForm({
      name: merchant.name ?? "",
      description: merchant.description ?? "",
      countryCode: merchant.countryCode ?? "",
    });
  }, [merchant]);

  const isDirty = useMemo(() => {
    if (!merchant) return false;
    return (
      form.name.trim() !== (merchant.name ?? "").trim() ||
      form.description.trim() !== (merchant.description ?? "").trim() ||
      form.countryCode !== (merchant.countryCode ?? "")
    );
  }, [form, merchant]);

  const countryOptions = useMemo(() => {
    if (!form.countryCode) return countries;
    if (countries.some((c) => c.code === form.countryCode)) return countries;
    return [{ code: form.countryCode, name: form.countryCode }, ...countries];
  }, [countries, form.countryCode]);

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    if (!merchant || !isDirty) return;

    setError(null);
    setSuccess(null);
    setIsSaving(true);

    try {
      const payload: {
        name?: string;
        description?: string;
        countryCode?: string;
      } = {};

      const name = form.name.trim();
      const description = form.description.trim();
      const countryCode = form.countryCode.trim().toUpperCase();

      if (name !== (merchant.name ?? "").trim()) payload.name = name;
      if (description !== (merchant.description ?? "").trim()) {
        payload.description = description;
      }
      if (countryCode !== (merchant.countryCode ?? "")) {
        payload.countryCode = countryCode;
      }

      const updated = await clientUpdateMerchant(payload);
      setMerchant(updated);
      setSuccess("Store settings saved.");
    } catch {
      setError("Failed to save store settings. Please try again.");
    } finally {
      setIsSaving(false);
    }
  };

  if (!merchant) {
    return (
      <div className="bg-white shadow-1 rounded-xl p-4 sm:p-8.5">
        <p className="text-dark">Unable to load store settings.</p>
      </div>
    );
  }

  return (
    <form onSubmit={(e) => void handleSubmit(e)}>
      <div className="bg-white shadow-1 rounded-xl p-4 sm:p-8.5">
        <MerchantLogoUpload />

        <div className="mb-5">
          <label htmlFor="storeName" className="block mb-2.5">
            Store Name <span className="text-red">*</span>
          </label>
          <input
            type="text"
            name="storeName"
            id="storeName"
            placeholder="Store name"
            value={form.name}
            onChange={(e) => setForm((prev) => ({ ...prev, name: e.target.value }))}
            required
            maxLength={255}
            className="rounded-md border border-gray-3 bg-gray-1 placeholder:text-dark-5 w-full py-2.5 px-5 outline-none duration-200 focus:border-transparent focus:shadow-input focus:ring-2 focus:ring-blue/20"
          />
        </div>

        <div className="mb-5">
          <label htmlFor="storeSlug" className="block mb-2.5">
            Store Slug
          </label>
          <input
            type="text"
            name="storeSlug"
            id="storeSlug"
            value={merchant.slug}
            readOnly
            className="rounded-md border border-gray-3 bg-gray-1 w-full py-2.5 px-5 outline-none duration-200 focus:border-transparent focus:shadow-input focus:ring-2 focus:ring-blue/20"
          />
        </div>

        <div className="mb-5">
          <label htmlFor="storeDescription" className="block mb-2.5">
            Description
          </label>
          <textarea
            name="storeDescription"
            id="storeDescription"
            rows={4}
            placeholder="Describe your store"
            value={form.description}
            onChange={(e) =>
              setForm((prev) => ({ ...prev, description: e.target.value }))
            }
            className="rounded-md border border-gray-3 bg-gray-1 placeholder:text-dark-5 w-full py-2.5 px-5 outline-none duration-200 focus:border-transparent focus:shadow-input focus:ring-2 focus:ring-blue/20"
          />
        </div>

        <div className="mb-5">
          <label htmlFor="countryCode" className="block mb-2.5">
            Country / Region <span className="text-red">*</span>
          </label>
          <div className="relative">
            <select
              id="countryCode"
              name="countryCode"
              value={form.countryCode}
              onChange={(e) =>
                setForm((prev) => ({ ...prev, countryCode: e.target.value }))
              }
              required
              className="w-full bg-gray-1 rounded-md border border-gray-3 text-dark-4 py-3 pl-5 pr-9 duration-200 appearance-none outline-none focus:border-transparent focus:shadow-input focus:ring-2 focus:ring-blue/20"
            >
              <option value="" disabled>
                Select a country
              </option>
              {countryOptions.map((country) => (
                <option key={country.code} value={country.code}>
                  {country.name}
                </option>
              ))}
            </select>
            <span className="absolute right-4 top-1/2 -translate-y-1/2 text-dark-4 pointer-events-none">
              <svg
                className="fill-current"
                width="16"
                height="16"
                viewBox="0 0 16 16"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  d="M2.41469 5.03569L2.41467 5.03571L2.41749 5.03846L7.76749 10.2635L8.0015 10.492L8.23442 10.2623L13.5844 4.98735L13.5844 4.98735L13.5861 4.98569C13.6809 4.89086 13.8199 4.89087 13.9147 4.98569C14.0092 5.08024 14.0095 5.21864 13.9155 5.31345C13.9152 5.31373 13.915 5.31401 13.9147 5.31429L8.16676 10.9622L8.16676 10.9622L8.16469 10.9643C8.06838 11.0606 8.02352 11.0667 8.00039 11.0667C7.94147 11.0667 7.89042 11.0522 7.82064 10.9991L2.08526 5.36345C1.99127 5.26865 1.99154 5.13024 2.08609 5.03569C2.18092 4.94086 2.31986 4.94086 2.41469 5.03569Z"
                  fill=""
                  stroke=""
                  strokeWidth="0.666667"
                />
              </svg>
            </span>
          </div>
        </div>

        <div className="flex flex-col lg:flex-row gap-5 sm:gap-8 mb-5">
          <div className="w-full">
            <label htmlFor="storeStatus" className="block mb-2.5">
              Status
            </label>
            <input
              type="text"
              id="storeStatus"
              value={merchant.status}
              readOnly
              className="rounded-md border border-gray-3 bg-gray-1 w-full py-2.5 px-5 outline-none capitalize"
            />
          </div>
          <div className="w-full">
            <label htmlFor="verificationStatus" className="block mb-2.5">
              Verification Status
            </label>
            <input
              type="text"
              id="verificationStatus"
              value={merchant.verificationStatus}
              readOnly
              className="rounded-md border border-gray-3 bg-gray-1 w-full py-2.5 px-5 outline-none capitalize"
            />
          </div>
        </div>

        {error && <p className="mb-4 text-custom-sm text-red">{error}</p>}
        {success && (
          <p className="mb-4 text-custom-sm text-green">{success}</p>
        )}

        <button
          type="submit"
          disabled={!isDirty || isSaving}
          className={`inline-flex font-medium text-white py-3 px-7 rounded-md ease-out duration-200 ${
            !isDirty || isSaving
              ? "bg-gray-4 cursor-not-allowed"
              : "bg-blue hover:bg-blue-dark"
          }`}
        >
          {isSaving ? "Saving..." : "Save Changes"}
        </button>
      </div>
    </form>
  );
}
