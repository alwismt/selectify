import { getServerUser } from "@/lib/api/getServerUser";
import { getServerUserFile } from "@/lib/api/getServerUserFile";
import { getServerMerchant } from "@/lib/api/getServerMerchant";
import { getServerCart } from "@/lib/api/getServerCart";
import { redirect } from "next/navigation";
import ClientMerchantLayout from "./ClientMerchantLayout";

export default async function MerchantLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const [initialUser, initialUserFile] = await Promise.all([
    getServerUser(),
    getServerUserFile(),
  ]);

  if (!initialUser) {
    redirect("/signin");
  }

  if (initialUser.user_role?.role !== "merchant") {
    redirect("/");
  }

  const [initialMerchant, initialCart] = await Promise.all([
    getServerMerchant(),
    getServerCart(),
  ]);
  if (!initialMerchant) {
    redirect("/");
  }

  return (
    <ClientMerchantLayout
      initialUser={initialUser}
      initialUserFile={initialUserFile}
      initialMerchant={initialMerchant}
      initialCart={initialCart}
    >
      {children}
    </ClientMerchantLayout>
  );
}
