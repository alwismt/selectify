import { getServerUser } from "@/lib/api/getServerUser";
import { getServerUserFile } from "@/lib/api/getServerUserFile";
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

  return (
    <ClientMerchantLayout
      initialUser={initialUser}
      initialUserFile={initialUser ? initialUserFile : null}
    >
      {children}
    </ClientMerchantLayout>
  );
}
