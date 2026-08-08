import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Merchant Dashboard | Selectify",
  description: "Manage your products, orders, and inventory",
};

export default function MerchantDashboard() {
  return (
    <div className="p-8">
      <h1 className="text-3xl font-bold mb-4">Merchant Dashboard</h1>
      <p>Coming soon: Manage your store</p>
    </div>
  );
}
