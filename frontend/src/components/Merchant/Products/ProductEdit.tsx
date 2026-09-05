"use client";

import { FormEvent, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { useMerchantProductEditor } from "@/app/context/MerchantProductEditorContext";
import { imagesFromProductFile } from "@/lib/editableImage";
import { formatMinorToMajorInput } from "@/lib/format";
import { merchantProductHref } from "@/lib/productPath";
import type { MerchantApiProduct } from "@/types/merchantProduct";
import type { SiteCurrency } from "@/types/siteConfig";
import HydrateMerchantProductEditor from "./HydrateMerchantProductEditor";
import ProductImagesEditor from "./ProductImagesEditor";
import { useEditableImages } from "./useEditableImages";

type ProductFormState = {
  name: string;
  sku: string;
  description: string;
  /** Major-unit decimal string for the price input. */
  price: string;
  isActive: boolean;
  inStock: boolean;
};

const inputClassName =
  "rounded-md border border-gray-3 bg-gray-1 placeholder:text-dark-5 w-full py-2.5 px-5 outline-none duration-200 focus:border-transparent focus:shadow-input focus:ring-2 focus:ring-blue/20";

function formFromProduct(
  product: MerchantApiProduct,
  minorUnit: number
): ProductFormState {
  return {
    name: product.name,
    sku: product.sku,
    description: product.description ?? "",
    price: formatMinorToMajorInput(product.price, minorUnit),
    isActive: product.isActive,
    inStock: product.inStock,
  };
}

type ProductEditProps = {
  productId: number;
  fallbackProduct: MerchantApiProduct;
  currency: SiteCurrency | null;
};

export default function ProductEdit({
  productId,
  fallbackProduct,
  currency,
}: ProductEditProps) {
  const router = useRouter();
  const { product, variants, loadedProductId } =
    useMerchantProductEditor();
  const resolvedProduct =
    loadedProductId === productId && product != null ? product : fallbackProduct;
  const minorUnit = currency?.minorUnit ?? 2;

  const [form, setForm] = useState<ProductFormState>(() =>
    formFromProduct(resolvedProduct, minorUnit)
  );
  const {
    images,
    error: imageError,
    addFiles,
    remove,
    setPrimary,
    isDirty: imagesDirty,
  } = useEditableImages(imagesFromProductFile(resolvedProduct.productFile));

  const needsHydrate = loadedProductId !== productId || product == null;

  const isDirty = useMemo(() => {
    const initialPrice = formatMinorToMajorInput(
      resolvedProduct.price,
      minorUnit
    );
    return (
      form.name.trim() !== resolvedProduct.name.trim() ||
      form.sku.trim() !== resolvedProduct.sku.trim() ||
      form.description.trim() !== (resolvedProduct.description ?? "").trim() ||
      form.price !== initialPrice ||
      form.isActive !== resolvedProduct.isActive ||
      form.inStock !== resolvedProduct.inStock ||
      imagesDirty
    );
  }, [form, resolvedProduct, imagesDirty, minorUnit]);

  const productSlug =
    typeof resolvedProduct.slug === "string"
      ? resolvedProduct.slug.trim()
      : "";
  const detailsHref = productSlug
    ? merchantProductHref(resolvedProduct.productId, productSlug)
    : "/merchant/products";

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

  const slugValue =
    typeof resolvedProduct.slug === "string" ? resolvedProduct.slug : "";

  const currencyLabel = currency
    ? `${currency.code} (${currency.name})`
    : "Site default";
  const priceHelper = currency
    ? `Enter amount in ${currency.code} (e.g. 399.99).`
    : "Enter the price in major units.";

  return (
    <form
      onSubmit={handleSubmit}
      className="w-full min-w-0 bg-white rounded-xl shadow-1 py-9.5 px-4 sm:px-7.5 xl:px-10"
    >
      {needsHydrate ? (
        <HydrateMerchantProductEditor
          product={fallbackProduct}
          variants={variants}
        />
      ) : null}
      <button
        type="button"
        onClick={handleCancel}
        className="inline-flex text-custom-sm font-medium text-dark-4 hover:text-blue mb-5"
      >
        ← Back to Product
      </button>

      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-7.5">
        <h2 className="font-medium text-xl sm:text-2xl text-dark">
          Edit Product
        </h2>
      </div>

      <div className="mb-5">
        <label htmlFor="productName" className="block mb-2.5">
          Name
        </label>
        <input
          type="text"
          id="productName"
          name="productName"
          value={form.name}
          onChange={(e) =>
            setForm((prev) => ({ ...prev, name: e.target.value }))
          }
          className={inputClassName}
        />
      </div>

      <div className="mb-5">
        <label htmlFor="productSlug" className="block mb-2.5">
          Slug
        </label>
        <input
          type="text"
          id="productSlug"
          name="productSlug"
          value={slugValue}
          readOnly
          className={inputClassName}
        />
      </div>

      <div className="mb-5">
        <label htmlFor="productSku" className="block mb-2.5">
          SKU
        </label>
        <input
          type="text"
          id="productSku"
          name="productSku"
          value={form.sku}
          onChange={(e) =>
            setForm((prev) => ({ ...prev, sku: e.target.value }))
          }
          className={inputClassName}
        />
      </div>

      <div className="mb-5">
        <label htmlFor="productDescription" className="block mb-2.5">
          Description
        </label>
        <textarea
          id="productDescription"
          name="productDescription"
          rows={4}
          value={form.description}
          onChange={(e) =>
            setForm((prev) => ({ ...prev, description: e.target.value }))
          }
          className={inputClassName}
        />
      </div>

      <div className="flex flex-col lg:flex-row gap-5 sm:gap-8 mb-5">
        <div className="w-full">
          <label htmlFor="productPrice" className="block mb-2.5">
            Price
          </label>
          <input
            type="text"
            inputMode="decimal"
            id="productPrice"
            name="productPrice"
            value={form.price}
            onChange={(e) =>
              setForm((prev) => ({ ...prev, price: e.target.value }))
            }
            placeholder="399.99"
            className={inputClassName}
          />
          <p className="mt-1.5 text-custom-sm text-dark-4">{priceHelper}</p>
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

      <div className="flex flex-col sm:flex-row gap-5 mb-7.5">
        <label className="inline-flex items-center gap-2 text-custom-sm text-dark">
          <input
            type="checkbox"
            checked={form.isActive}
            onChange={(e) =>
              setForm((prev) => ({ ...prev, isActive: e.target.checked }))
            }
          />
          Active
        </label>
        <label className="inline-flex items-center gap-2 text-custom-sm text-dark">
          <input
            type="checkbox"
            checked={form.inStock}
            onChange={(e) =>
              setForm((prev) => ({ ...prev, inStock: e.target.checked }))
            }
          />
          In stock
        </label>
      </div>

      <ProductImagesEditor
        images={images}
        onAdd={addFiles}
        onRemove={remove}
        onSetPrimary={setPrimary}
        error={imageError}
      />

      <p className="mt-8 mb-4 text-custom-sm text-dark-4">
        Product updates are not available yet.
      </p>

      <div className="flex flex-wrap gap-3">
        <button
          type="submit"
          disabled
          className="inline-flex font-medium text-white py-3 px-7 rounded-md bg-gray-4 cursor-not-allowed"
        >
          Save Product
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
