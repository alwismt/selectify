import { productFileUrl } from "@/lib/api/config";
import type { ApiProductFile } from "@/types/product";

type ProductImagePlaceholderProps = {
  productFile?: ApiProductFile | null;
  name: string;
  size?: "sm" | "lg";
  className?: string;
};

/**
 * Product image with gray placeholder as the normal fallback when productFile is absent.
 */
export default function ProductImagePlaceholder({
  productFile,
  name,
  size = "sm",
  className = "",
}: ProductImagePlaceholderProps) {
  const url = productFileUrl(productFile?.file_id);
  const sizeClass =
    size === "lg"
      ? "h-48 w-full sm:h-56 sm:w-56 xl:h-64 xl:w-64"
      : "h-12 w-12 shrink-0";

  if (url) {
    return (
      // eslint-disable-next-line @next/next/no-img-element
      <img
        src={url}
        alt={name}
        className={`${sizeClass} rounded-md object-cover border border-gray-3 ${className}`}
      />
    );
  }

  return (
    <div
      className={`${sizeClass} rounded-md border border-gray-3 bg-gray-1 ${className}`}
      aria-hidden="true"
    />
  );
}
