import Link from "next/link";
import { formatMoney } from "@/lib/format";
import { merchantProductHref } from "@/lib/productPath";
import type { MerchantApiProduct } from "@/types/merchantProduct";
import type { SiteCurrency } from "@/types/siteConfig";
import ProductImagePlaceholder from "./ProductImagePlaceholder";
import ProductStatusPill from "./ProductStatusPill";

type MerchantProductsListProps = {
  products: MerchantApiProduct[] | null;
  currency: SiteCurrency | null;
};

function stockLabel(inStock: boolean): string {
  return inStock ? "In stock" : "Out of stock";
}

export default function MerchantProductsList({
  products,
  currency,
}: MerchantProductsListProps) {
  return (
    <div className="w-full min-w-0 bg-white rounded-xl shadow-1 py-9.5 px-4 sm:px-7.5 xl:px-10">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-7.5">
        <h2 className="font-medium text-xl sm:text-2xl text-dark">Products</h2>
        <Link
          href="/merchant/products/new"
          className="inline-flex items-center justify-center font-medium text-white bg-blue py-3 px-7 rounded-md ease-out duration-200 hover:bg-blue-dark"
        >
          + Add Product
        </Link>
      </div>

      {products === null ? (
        <p className="text-custom-sm text-red">
          Unable to load products. Please try again later.
        </p>
      ) : products.length === 0 ? (
        <p className="text-custom-sm text-dark-4">No products yet.</p>
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full min-w-[640px] text-left">
            <thead>
              <tr className="border-b border-gray-3">
                <th className="pb-4 pr-4 font-medium text-custom-sm text-dark-4">
                  Product
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
                  Status
                </th>
                <th className="pb-4 font-medium text-custom-sm text-dark-4 text-right">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody>
              {products.map((product) => {
                const slug =
                  typeof product.slug === "string" ? product.slug.trim() : "";
                const href = slug
                  ? merchantProductHref(product.productId, slug)
                  : null;

                return (
                  <tr
                    key={product.productId}
                    className="relative border-b border-gray-3 last:border-0"
                  >
                    <td className="py-4 pr-4">
                      <div className="flex items-center gap-3">
                        <ProductImagePlaceholder
                          productFile={product.productFile}
                          name={product.name}
                          size="sm"
                        />
                        {href ? (
                          <Link
                            href={href}
                            className="text-custom-sm font-medium text-dark before:absolute before:inset-0 before:content-['']"
                          >
                            {product.name}
                          </Link>
                        ) : (
                          <span className="text-custom-sm font-medium text-dark">
                            {product.name}
                          </span>
                        )}
                      </div>
                    </td>
                    <td className="relative z-10 py-4 pr-4 text-custom-sm text-dark">
                      {product.sku}
                    </td>
                    <td className="relative z-10 py-4 pr-4 text-custom-sm text-dark">
                      {currency ? formatMoney(product.price, currency) : "—"}
                    </td>
                    <td className="relative z-10 py-4 pr-4 text-custom-sm text-dark">
                      {stockLabel(product.inStock)}
                    </td>
                    <td className="relative z-10 py-4 pr-4">
                      <ProductStatusPill isActive={product.isActive} />
                    </td>
                    <td className="relative z-10 py-4 text-right">
                      {href ? (
                        <Link
                          href={href}
                          className="relative z-10 text-custom-sm font-medium text-blue hover:underline"
                        >
                          View
                        </Link>
                      ) : (
                        <span className="text-custom-sm text-dark-4">—</span>
                      )}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
