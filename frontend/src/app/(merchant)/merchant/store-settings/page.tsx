import MerchantSidebar from "@/components/MerchantDashboard/MerchantSidebar";
import StoreSettingsForm from "@/components/Merchant/StoreSettings";
import { getServerMerchantCountries } from "@/lib/api/getServerMerchantCountries";
import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Store Settings | Selectify Seller",
  description: "Manage your store details and logo",
};

export default async function StoreSettingsPage() {
  const countries = await getServerMerchantCountries();

  return (
    <main>
      <section className="overflow-hidden pt-[100px] sm:pt-[100px] lg:pt-[110px] xl:pt-[120px] pb-20 bg-gray-2">
        <div className="mx-auto max-w-[1500px] w-full px-4 sm:px-8 xl:px-0">
          <div className="flex flex-col gap-7.5 xl:flex-row">
            <MerchantSidebar />
            <div className="min-w-0 flex-1">
              <h1 className="mb-6 text-2xl font-medium text-dark">
                Store Settings
              </h1>
              <StoreSettingsForm countries={countries} />
            </div>
          </div>
        </div>
      </section>
    </main>
  );
}
