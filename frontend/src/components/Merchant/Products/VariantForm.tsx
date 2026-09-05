"use client";

import { FormEvent, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { useSiteConfig } from "@/app/context/SiteConfigContext";
import { imagesFromVariantFiles } from "@/lib/editableImage";
import { formatMinorToMajorInput } from "@/lib/format";
import { merchantProductHref } from "@/lib/productPath";
import { availableVariantQuantity } from "@/types/api/variant";
import type { MerchantApiProductVariant } from "@/types/merchantProduct";
import EditableImageGallery from "./EditableImageGallery";
import VariantAttributesEditor, {
  type VariantAttributeRow,
} from "./VariantAttributesEditor";
import { useEditableImages } from "./useEditableImages";

export const variantFormInputClassName =
  "rounded-md border border-gray-3 bg-gray-1 placeholder:text-dark-5 w-full py-2.5 px-5 outline-none duration-200 focus:border-transparent focus:shadow-input focus:ring-2 focus:ring-blue/20";

type VariantFormState = {
  sku: string;
  /** Major-unit decimal string for the price input. */
  priceAmount: string;
  stockQuantity: string;
  attributes: VariantAttributeRow[];
};

type VariantFormProps = {
  mode: "create" | "edit";
  productId: number;
  productSlug: string;
  variant?: MerchantApiProductVariant;
};

let attributeKey = 0;

function nextAttributeKey(): string {
  attributeKey += 1;
  return `attr-${attributeKey}`;
}

function formFromVariant(
  variant: MerchantApiProductVariant | undefined,
  minorUnit: number
): VariantFormState {
  if (!variant) {
    return {
      sku: "",
      priceAmount: "",
      stockQuantity: "0",
      attributes: [{ key: nextAttributeKey(), name: "", value: "" }],
    };
  }

  return {
    sku: variant.sku,
    priceAmount: formatMinorToMajorInput(variant.price_amount, minorUnit),
    stockQuantity: String(variant.stock_quantity),
    attributes:
      (variant.product_variant_attributes ?? []).length === 0
        ? [{ key: nextAttributeKey(), name: "", value: "" }]
        : variant.product_variant_attributes.map((attr) => ({
            key: String(attr.id),
            name: attr.name,
            value: attr.value,
          })),
  };
}

export default function VariantForm({
  mode,
  productId,
  productSlug,
  variant,
}: VariantFormProps) {
  const router = useRouter();
  const { currency } = useSiteConfig();
  const minorUnit = currency?.minorUnit ?? 2;
  const [form, setForm] = useState<VariantFormState>(() =>
    formFromVariant(variant, minorUnit)
  );
  const {
    images,
    error: imageError,
    addFiles,
    remove,
    setPrimary,
    isDirty: imagesDirty,
  } = useEditableImages(imagesFromVariantFiles(variant?.files ?? []));

  const reservedQuantity = variant?.reserved_quantity ?? 0;
  const stockQuantity = Number(form.stockQuantity);
  const available = useMemo(() => {
    const stock = Number.isFinite(stockQuantity) ? stockQuantity : 0;
    return availableVariantQuantity({
      stock_quantity: stock,
      reserved_quantity: reservedQuantity,
    });
  }, [stockQuantity, reservedQuantity]);

  const detailsHref = productSlug
    ? merchantProductHref(productId, productSlug)
    : "/merchant/products";

  const isDirty = useMemo(() => {
    const initial = formFromVariant(variant, minorUnit);
    return (
      form.sku !== initial.sku ||
      form.priceAmount !== initial.priceAmount ||
      form.stockQuantity !== initial.stockQuantity ||
      form.attributes.length !== initial.attributes.length ||
      form.attributes.some(
        (row, index) =>
          row.name !== initial.attributes[index]?.name ||
          row.value !== initial.attributes[index]?.value
      ) ||
      imagesDirty
    );
  }, [form, variant, minorUnit, imagesDirty]);

  const handleCancel = () => {
    if (
      isDirty &&
      !window.confirm("You have unsaved changes. Discard them and leave?")
    ) {
      return;
    }
    router.push(detailsHref);
  };

  const handleSubmit = (event: FormEvent) => {
    event.preventDefault();
  };

  const title = mode === "create" ? "Add Variant" : "Edit Variant";
  const saveLabel = mode === "create" ? "Add Variant" : "Save Variant";
  const helperText =
    mode === "create"
      ? "Adding a variant is not available yet."
      : "Variant updates are not available yet.";
  const priceHelper = currency
    ? `Enter amount in ${currency.code} (e.g. 399.99).`
    : "Enter the price in major units.";

  return (
    <form
      onSubmit={handleSubmit}
      className="w-full min-w-0 bg-white rounded-xl shadow-1 py-9.5 px-4 sm:px-7.5 xl:px-10"
    >
      <button
        type="button"
        onClick={handleCancel}
        className="inline-flex text-custom-sm font-medium text-dark-4 hover:text-blue mb-5"
      >
        ← Back to Product
      </button>

      <h2 className="font-medium text-xl sm:text-2xl text-dark mb-7.5">
        {title}
      </h2>

      <div className="mb-5">
        <label htmlFor="variantSku" className="block mb-2.5">
          SKU
        </label>
        <input
          type="text"
          id="variantSku"
          name="variantSku"
          value={form.sku}
          onChange={(e) => setForm((prev) => ({ ...prev, sku: e.target.value }))}
          className={variantFormInputClassName}
        />
      </div>

      <div className="mb-5">
        <label htmlFor="variantPrice" className="block mb-2.5">
          Price
        </label>
        <input
          type="text"
          inputMode="decimal"
          id="variantPrice"
          name="variantPrice"
          value={form.priceAmount}
          onChange={(e) =>
            setForm((prev) => ({ ...prev, priceAmount: e.target.value }))
          }
          placeholder="399.99"
          className={variantFormInputClassName}
        />
        <p className="mt-1.5 text-custom-sm text-dark-4">{priceHelper}</p>
      </div>

      <div className="mb-5">
        <label htmlFor="variantStock" className="block mb-2.5">
          Stock Quantity
        </label>
        <input
          type="number"
          id="variantStock"
          name="variantStock"
          value={form.stockQuantity}
          onChange={(e) =>
            setForm((prev) => ({ ...prev, stockQuantity: e.target.value }))
          }
          className={variantFormInputClassName}
        />
      </div>

      {mode === "edit" && (
        <div className="flex flex-col lg:flex-row gap-5 sm:gap-8 mb-5">
          <div className="w-full">
            <label htmlFor="variantReserved" className="block mb-2.5">
              Reserved Quantity
            </label>
            <input
              type="text"
              id="variantReserved"
              name="variantReserved"
              value={String(reservedQuantity)}
              readOnly
              className={variantFormInputClassName}
            />
          </div>
          <div className="w-full">
            <label htmlFor="variantAvailable" className="block mb-2.5">
              Available
            </label>
            <input
              type="text"
              id="variantAvailable"
              name="variantAvailable"
              value={String(available)}
              readOnly
              className={variantFormInputClassName}
            />
          </div>
        </div>
      )}

      <VariantAttributesEditor
        attributes={form.attributes}
        onChange={(attributes) =>
          setForm((prev) => ({ ...prev, attributes }))
        }
        onAdd={() =>
          setForm((prev) => ({
            ...prev,
            attributes: [
              ...prev.attributes,
              { key: nextAttributeKey(), name: "", value: "" },
            ],
          }))
        }
      />

      <EditableImageGallery
        title="Variant Images"
        size="md"
        images={images}
        onAdd={addFiles}
        onRemove={remove}
        onSetPrimary={setPrimary}
        error={imageError}
      />

      <p className="mt-8 mb-4 text-custom-sm text-dark-4">{helperText}</p>

      <div className="flex flex-wrap gap-3">
        <button
          type="submit"
          disabled
          className="inline-flex font-medium text-white py-3 px-7 rounded-md bg-gray-4 cursor-not-allowed"
        >
          {saveLabel}
        </button>
        <button
          type="button"
          onClick={handleCancel}
          className="inline-flex font-medium text-dark bg-gray-2 py-3 px-7 rounded-md ease-out duration-200 hover:bg-gray-3"
        >
          Cancel
        </button>
      </div>
    </form>
  );
}
