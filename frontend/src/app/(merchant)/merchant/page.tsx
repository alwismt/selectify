import MerchantDashboard from "@/components/MerchantDashboard";
import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Merchant Dashboard | Selectify",
  description: "Manage your products, orders, and inventory",
};

export default function MerchantDashboardPage() {
  return (
    <main>
      <MerchantDashboard />
    </main>
  );
}
