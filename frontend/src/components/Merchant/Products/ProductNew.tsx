"use client";

import {
  ChangeEvent,
  FormEvent,
  useEffect,
  useRef,
  useState,
} from "react";
import { useRouter } from "next/navigation";
import { parseMajorAmountToMinor } from "@/lib/format";
import { ACCEPTED_IMAGE_TYPES, validateImageFile } from "@/lib/imageFile";
import { merchantProductHref } from "@/lib/productPath";
import type { ApiCategory } from "@/types/category";
import type { MerchantApiProduct } from "@/types/merchantProduct";
import type { SiteCurrency } from "@/types/siteConfig";
import CategoryTreeSelector from "./CategoryTreeSelector";

type ProductNewFormState = {
  name: string;
  description: string;
  sku: string;
  price: string;
};

const inputClassName =
  "rounded-md border border-gray-3 bg-gray-1 placeholder:text-dark-5 w-full py-2.5 px-5 outline-none duration-200 focus:border-transparent focus:shadow-input focus:ring-2 focus:ring-blue/20";

const KNOWN_FIELD_KEYS = [
  "name",
  "description",
  "sku",
  "price",
  "image",
  "categoryIds",
] as const;

type ProductNewProps = {
  currency: SiteCurrency | null;
  categories: ApiCategory[];
};

export default function ProductNew({ currency, categories }: ProductNewProps) {
  const router = useRouter();
  const fileInputRef = useRef<HTMLInputElement>(null);

  const [form, setForm] = useState<ProductNewFormState>({
    name: "",
    description: "",
    sku: "",
    price: "",
  });
  const [selectedCategoryIds, setSelectedCategoryIds] = useState<number[]>([]);
  const [imageFile, setImageFile] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});
  const [formError, setFormError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const minorUnit = currency?.minorUnit ?? 2;
  const currencyLabel = currency
    ? `${currency.code} (${currency.name})`
    : "Site default";
  const priceHelper = currency
    ? `Enter amount in ${currency.code} (e.g. 1299.99).`
    : "Enter amount in major units (e.g. 1299.99).";

  useEffect(() => {
    return () => {
      if (previewUrl) {
        URL.revokeObjectURL(previewUrl);
      }
    };
  }, [previewUrl]);

  const clearFieldError = (key: string) => {
    setFieldErrors((prev) => {
      if (!(key in prev)) return prev;
      const next = { ...prev };
      delete next[key];
      return next;
    });
  };

  const handleCategoriesChange = (ids: number[]) => {
    setSelectedCategoryIds(ids);
    clearFieldError("categoryIds");
    setFormError(null);
  };

  const handleImageChange = (event: ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0] ?? null;
    if (!file) return;

    const validationError = validateImageFile(file);
    if (validationError) {
      setFieldErrors((prev) => ({ ...prev, image: validationError }));
      event.target.value = "";
      return;
    }

    if (previewUrl) {
      URL.revokeObjectURL(previewUrl);
    }
    setImageFile(file);
    setPreviewUrl(URL.createObjectURL(file));
    clearFieldError("image");
    setFormError(null);
  };

  const handleRemoveImage = () => {
    if (previewUrl) {
      URL.revokeObjectURL(previewUrl);
    }
    setImageFile(null);
    setPreviewUrl(null);
    clearFieldError("image");
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  const handleCancel = () => {
    const hasValues =
      form.name.trim() !== "" ||
      form.description.trim() !== "" ||
      form.sku.trim() !== "" ||
      form.price.trim() !== "" ||
      selectedCategoryIds.length > 0 ||
      imageFile != null;
    if (
      hasValues &&
      !window.confirm("You have unsaved changes. Discard them and leave?")
    ) {
      return;
    }
    router.push("/merchant/products");
  };

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    if (isSubmitting) return;

    setFormError(null);

    const localErrors: Record<string, string> = {};
    if (!form.name.trim()) {
      localErrors.name = "This field is required";
    }
    if (!form.description.trim()) {
      localErrors.description = "This field is required";
    }
    if (!form.sku.trim()) {
      localErrors.sku = "This field is required";
    }

    const priceMinor = parseMajorAmountToMinor(form.price, minorUnit);
    if (!form.price.trim()) {
      localErrors.price = "This field is required";
    } else if (priceMinor == null) {
      localErrors.price = "Invalid value";
    }

    if (selectedCategoryIds.length === 0) {
      localErrors.categoryIds = "This field is required";
    }

    if (!imageFile) {
      localErrors.image = "Image is required";
    }

    if (Object.keys(localErrors).length > 0) {
      setFieldErrors(localErrors);
      return;
    }

    setIsSubmitting(true);
    setFieldErrors({});

    try {
      const formData = new FormData();
      formData.append("name", form.name.trim());
      formData.append("description", form.description.trim());
      formData.append("sku", form.sku.trim());
      formData.append("price", String(priceMinor));
      formData.append("image", imageFile as File);
      for (const categoryId of selectedCategoryIds) {
        formData.append("categoryIds", String(categoryId));
      }

      const res = await fetch("/api/merchant/products", {
        method: "POST",
        body: formData,
        credentials: "include",
      });

      const data: unknown = await res.json().catch(() => null);

      if (res.status === 400) {
        if (
          data &&
          typeof data === "object" &&
          !Array.isArray(data) &&
          Object.values(data).every((v) => typeof v === "string")
        ) {
          const fieldErrorsMap = data as Record<string, string>;
          setFieldErrors(fieldErrorsMap);
          const unknown = Object.entries(fieldErrorsMap).filter(
            ([key]) =>
              !(KNOWN_FIELD_KEYS as readonly string[]).includes(key)
          );
          if (unknown.length > 0) {
            setFormError(unknown.map(([, message]) => message).join(" "));
          }
          return;
        }
        setFormError("Validation failed. Please check your input.");
        return;
      }

      if (res.status !== 201) {
        setFormError("Failed to create product. Please try again.");
        return;
      }

      if (
        !data ||
        typeof data !== "object" ||
        typeof (data as MerchantApiProduct).productId !== "number" ||
        typeof (data as MerchantApiProduct).sku !== "string" ||
        typeof (data as MerchantApiProduct).name !== "string"
      ) {
        setFormError("Invalid product response.");
        return;
      }

      const product = data as MerchantApiProduct;
      const slug =
        typeof product.slug === "string" ? product.slug.trim() : "";
      if (!slug) {
        setFormError("Product was created but is missing a slug.");
        return;
      }

      router.push(merchantProductHref(product.productId, slug));
    } catch {
      setFormError("Failed to create product. Please try again.");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form
      onSubmit={handleSubmit}
      className="w-full min-w-0 bg-white rounded-xl shadow-1 py-9.5 px-4 sm:px-7.5 xl:px-10"
      noValidate
    >
      <button
        type="button"
        onClick={handleCancel}
        className="inline-flex text-custom-sm font-medium text-dark-4 hover:text-blue mb-5"
      >
        ← Back to Products
      </button>

      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-7.5">
        <h2 className="font-medium text-xl sm:text-2xl text-dark">
          Add Product
        </h2>
      </div>

      {formError ? (
        <p className="mb-4 text-custom-sm text-red">{formError}</p>
      ) : null}

      <div className="mb-5">
        <label htmlFor="productName" className="block mb-2.5">
          Name <span className="text-red">*</span>
        </label>
        <input
          type="text"
          id="productName"
          name="productName"
          value={form.name}
          onChange={(e) => {
            setForm((prev) => ({ ...prev, name: e.target.value }));
            clearFieldError("name");
          }}
          className={inputClassName}
          disabled={isSubmitting}
        />
        {fieldErrors.name ? (
          <p className="mt-2 text-custom-sm text-red">{fieldErrors.name}</p>
        ) : null}
      </div>

      <div className="mb-5">
        <label htmlFor="productDescription" className="block mb-2.5">
          Description <span className="text-red">*</span>
        </label>
        <textarea
          id="productDescription"
          name="productDescription"
          rows={4}
          value={form.description}
          onChange={(e) => {
            setForm((prev) => ({ ...prev, description: e.target.value }));
            clearFieldError("description");
          }}
          className={inputClassName}
          disabled={isSubmitting}
        />
        {fieldErrors.description ? (
          <p className="mt-2 text-custom-sm text-red">
            {fieldErrors.description}
          </p>
        ) : null}
      </div>

      <div className="mb-5">
        <label htmlFor="productSku" className="block mb-2.5">
          SKU <span className="text-red">*</span>
        </label>
        <input
          type="text"
          id="productSku"
          name="productSku"
          value={form.sku}
          onChange={(e) => {
            setForm((prev) => ({ ...prev, sku: e.target.value }));
            clearFieldError("sku");
          }}
          className={inputClassName}
          disabled={isSubmitting}
          autoCapitalize="characters"
          spellCheck={false}
        />
        <p className="mt-2 text-custom-sm text-dark-4">
          Uppercase letters, numbers, hyphens, and underscores.
        </p>
        {fieldErrors.sku ? (
          <p className="mt-2 text-custom-sm text-red">{fieldErrors.sku}</p>
        ) : null}
      </div>

      <div className="mb-5">
        <label className="block mb-2.5">
          Categories <span className="text-red">*</span>
        </label>
        <p className="mb-2.5 text-custom-sm text-dark-4">
          Select at least one category. Use search to find nested categories.
        </p>
        <CategoryTreeSelector
          categories={categories}
          selectedIds={selectedCategoryIds}
          onChange={handleCategoriesChange}
          disabled={isSubmitting}
          error={fieldErrors.categoryIds ?? null}
        />
      </div>

      <div className="flex flex-col lg:flex-row gap-5 sm:gap-8 mb-5">
        <div className="w-full">
          <label htmlFor="productPrice" className="block mb-2.5">
            Price <span className="text-red">*</span>
          </label>
          <input
            type="text"
            inputMode="decimal"
            id="productPrice"
            name="productPrice"
            value={form.price}
            onChange={(e) => {
              setForm((prev) => ({ ...prev, price: e.target.value }));
              clearFieldError("price");
            }}
            placeholder="1299.99"
            className={inputClassName}
            disabled={isSubmitting}
          />
          <p className="mt-2 text-custom-sm text-dark-4">{priceHelper}</p>
          {fieldErrors.price ? (
            <p className="mt-2 text-custom-sm text-red">{fieldErrors.price}</p>
          ) : null}
        </div>
        <div className="w-full">
          <label htmlFor="productCurrency" className="block mb-2.5">
            Currency
          </label>
          <input
            type="text"
            id="productCurrency"
            name="productCurrency"
            value={currencyLabel}
            readOnly
            className={inputClassName}
          />
        </div>
      </div>

      <section className="mt-8 pt-7.5 border-t border-gray-3">
        <h3 className="font-medium text-lg text-dark mb-5">
          Product Image <span className="text-red">*</span>
        </h3>
        <div className="flex flex-wrap gap-3">
          {previewUrl ? (
            <div className="flex flex-col gap-2">
              <div className="relative">
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img
                  src={previewUrl}
                  alt="Selected product preview"
                  className="h-40 w-40 sm:h-48 sm:w-48 rounded-md object-cover border border-gray-3"
                />
              </div>
              <button
                type="button"
                onClick={handleRemoveImage}
                disabled={isSubmitting}
                aria-label="Remove image"
                className="text-custom-sm font-medium text-red hover:underline disabled:opacity-50"
              >
                Remove
              </button>
            </div>
          ) : (
            <div
              className="h-40 w-40 sm:h-48 sm:w-48 rounded-md border border-gray-3 bg-gray-1"
              aria-hidden="true"
            />
          )}

          <button
            type="button"
            onClick={() => fileInputRef.current?.click()}
            disabled={isSubmitting}
            aria-label="Add product image"
            className="h-40 w-40 sm:h-48 sm:w-48 rounded-lg border border-dashed border-gray-3 bg-gray-1 flex flex-col items-center justify-center gap-1 px-3 text-center hover:border-blue disabled:opacity-50"
          >
            <span className="text-custom-sm font-medium text-blue">
              {previewUrl ? "Change image" : "+ Add image"}
            </span>
            <span className="text-custom-xs text-dark-4">PNG, JPG, WEBP</span>
          </button>
        </div>

        <input
          ref={fileInputRef}
          type="file"
          accept={ACCEPTED_IMAGE_TYPES.join(",")}
          className="sr-only"
          onChange={handleImageChange}
          disabled={isSubmitting}
        />

        {fieldErrors.image ? (
          <p className="mt-3 text-custom-sm text-red">{fieldErrors.image}</p>
        ) : null}
      </section>

      <div className="flex flex-wrap gap-3 mt-8">
        <button
          type="submit"
          disabled={isSubmitting}
          className={`inline-flex font-medium text-white py-3 px-7 rounded-md ease-out duration-200 ${
            isSubmitting
              ? "bg-gray-4 cursor-not-allowed"
              : "bg-blue hover:bg-blue-dark"
          }`}
        >
          {isSubmitting ? "Creating..." : "Create Product"}
        </button>
        <button
          type="button"
          onClick={handleCancel}
          disabled={isSubmitting}
          className="inline-flex font-medium text-dark bg-gray-2 py-3 px-7 rounded-md ease-out duration-200 hover:bg-gray-3 disabled:opacity-50"
        >
          Cancel
        </button>
      </div>
    </form>
  );
}
