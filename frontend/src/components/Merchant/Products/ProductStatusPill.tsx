type ProductStatusPillProps = {
  isActive: boolean;
};

export function productStatusLabel(isActive: boolean): string {
  return isActive ? "Active" : "Inactive";
}

export function productStatusClass(isActive: boolean): string {
  return isActive
    ? "bg-green-light-6 text-green"
    : "bg-gray-1 text-dark-2";
}

export default function ProductStatusPill({ isActive }: ProductStatusPillProps) {
  return (
    <span
      className={`inline-flex rounded-full px-2.5 py-1 text-2xs font-medium ${productStatusClass(
        isActive
      )}`}
    >
      {productStatusLabel(isActive)}
    </span>
  );
}
