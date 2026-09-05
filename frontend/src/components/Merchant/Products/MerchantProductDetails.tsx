import Link from "next/link";
import { formatMoney } from "@/lib/format";
import type {
  MerchantApiProduct,
  MerchantApiProductVariant,
  MerchantApiVariantFile,
} from "@/types/merchantProduct";
import { availableVariantQuantity } from "@/types/api/variant";
import type { SiteCurrency } from "@/types/siteConfig";
import AddVariantButton from "./AddVariantButton";
import EditProductButton from "./EditProductButton";
import HydrateMerchantProductEditor from "./HydrateMerchantProductEditor";
import ProductImagePlaceholder from "./ProductImagePlaceholder";
import ProductStatusPill from "./ProductStatusPill";
import VariantRowActions from "./VariantRowActions";

type MerchantProductDetailsProps = {
  product: MerchantApiProduct;
  variants: MerchantApiProductVariant[] | null;
  currency: SiteCurrency | null;
};

function stockLabel(inStock: boolean): string {
  return inStock ? "In stock" : "Out of stock";
}

function capitalizeAttributeName(name: string): string {
  if (!name) return name;
  return name.charAt(0).toUpperCase() + name.slice(1);
}

function pickVariantFile(
  files: MerchantApiVariantFile[]
): MerchantApiVariantFile | null {
  if (files.length === 0) return null;
  const primary = files.find((file) => file.is_primary);
  if (primary) return primary;
  return files.reduce((lowest, file) =>
    file.position < lowest.position ? file : lowest
  );
}

function variantIdentityLabel(variant: MerchantApiProductVariant): string {
  return (variant.product_variant_attributes ?? [])
    .map((attr) => `${capitalizeAttributeName(attr.name)}: ${attr.value}`)
    .join(", ");
}

export default function MerchantProductDetails({
  product,
  variants,
  currency,
}: MerchantProductDetailsProps) {
  return (
    <div className="w-full min-w-0 bg-white rounded-xl shadow-1 py-9.5 px-4 sm:px-7.5 xl:px-10">
      <HydrateMerchantProductEditor product={product} variants={variants} />
      <Link
        href="/merchant/products"
        className="inline-flex text-custom-sm font-medium text-dark-4 hover:text-blue mb-5"
      >
        ← Back to Products
      </Link>

      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-7.5">
        <h2 className="font-medium text-xl sm:text-2xl text-dark">
          {product.name}
        </h2>
        <EditProductButton product={product} variants={variants} />
      </div>

      <div className="flex flex-col xl:flex-row gap-7.5">
        <div className="shrink-0">
          <ProductImagePlaceholder
            productFile={product.productFile}
            name={product.name}
            size="lg"
          />
        </div>

        <div className="min-w-0 flex-1 space-y-5">
          <div>
            <p className="text-custom-sm text-dark-4 mb-1">Product Name</p>
            <p className="text-custom-sm font-medium text-dark">{product.name}</p>
          </div>

          <div>
            <p className="text-custom-sm text-dark-4 mb-1.5">Status</p>
            <ProductStatusPill isActive={product.isActive} />
          </div>

          <div>
            <p className="text-custom-sm text-dark-4 mb-1">SKU</p>
            <p className="text-custom-sm text-dark">{product.sku}</p>
          </div>

          <div>
            <p className="text-custom-sm text-dark-4 mb-1">Price</p>
            <p className="text-custom-sm font-medium text-dark">
              {currency ? formatMoney(product.price, currency) : "—"}
            </p>
          </div>

          <div>
            <p className="text-custom-sm text-dark-4 mb-1">Stock</p>
            <p className="text-custom-sm text-dark">
              {stockLabel(product.inStock)}
            </p>
          </div>

          {product.description != null && product.description !== "" && (
            <div>
              <p className="text-custom-sm text-dark-4 mb-1">Description</p>
              <p className="text-custom-sm text-dark whitespace-pre-wrap">
                {product.description}
              </p>
            </div>
          )}
        </div>
      </div>

      <div className="mt-8 pt-7.5 border-t border-gray-3">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-5">
          <h3 className="font-medium text-lg text-dark">Variants</h3>
          <AddVariantButton product={product} variants={variants} />
        </div>

        {variants === null ? (
          <p className="text-custom-sm text-red">Unable to load variants.</p>
        ) : variants.length === 0 ? (
          <p className="text-custom-sm text-dark-4">
            No variants for this product.
          </p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full min-w-[840px] text-left">
              <thead>
                <tr className="border-b border-gray-3">
                  <th className="pb-4 pr-4 font-medium text-custom-sm text-dark-4">
                    Variant
                  </th>
                  <th className="pb-4 pr-4 font-medium text-custom-sm text-dark-4">
                    SKU
                  </th>
                  <th className="pb-4 pr-4 font-medium text-custom-sm text-dark-4">
                    Price
                  </th>
                  <th className="pb-4 pr-4 font-medium text-custom-sm text-dark-4">
                    Stock
                  </th>
                  <th className="pb-4 pr-4 font-medium text-custom-sm text-dark-4">
                    Reserved
                  </th>
                  <th className="pb-4 pr-4 font-medium text-custom-sm text-dark-4">
                    Available
                  </th>
                  <th className="pb-4 font-medium text-custom-sm text-dark-4">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody>
                {variants.map((variant) => {
                  const file = pickVariantFile(variant.files ?? []);
                  const available = availableVariantQuantity(variant);
                  const identity = variantIdentityLabel(variant);

                  return (
                    <tr
                      key={variant.id}
                      className="border-b border-gray-3 last:border-0"
                    >
                      <td className="py-4 pr-4">
                        <div className="flex items-center gap-3">
                          <ProductImagePlaceholder
                            productFile={file}
                            name={identity || variant.sku}
                            size="sm"
                          />
                          <div className="min-w-0">
                            {(variant.product_variant_attributes ?? [])
                              .length === 0 ? (
                              <span className="text-custom-sm text-dark-4">
                                —
                              </span>
                            ) : (
                              variant.product_variant_attributes.map((attr) => (
                                <p
                                  key={attr.id}
                                  className="text-custom-sm text-dark"
                                >
                                  {capitalizeAttributeName(attr.name)}:{" "}
                                  {attr.value}
                                </p>
                              ))
                            )}
                          </div>
                        </div>
                      </td>
                      <td className="py-4 pr-4 text-custom-sm text-dark">
                        {variant.sku}
                      </td>
                      <td className="py-4 pr-4 text-custom-sm text-dark">
                        {currency
                          ? formatMoney(variant.price_amount, currency)
                          : "—"}
                      </td>
                      <td className="py-4 pr-4 text-custom-sm text-dark">
                        {variant.stock_quantity}
                      </td>
                      <td className="py-4 pr-4 text-custom-sm text-dark">
                        {variant.reserved_quantity}
                      </td>
                      <td className="py-4 pr-4 text-custom-sm text-dark">
                        {available}
                      </td>
                      <td className="py-4">
                        <VariantRowActions
                          product={product}
                          variants={variants}
                          variantId={variant.id}
                        />
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}
