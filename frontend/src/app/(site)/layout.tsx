import { getServerUser } from "@/lib/api/getServerUser";
import { getServerUserFile } from "@/lib/api/getServerUserFile";
import { getServerCart } from "@/lib/api/getServerCart";
import ClientSiteLayout from "./ClientSiteLayout";

export default async function SiteLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const [initialUser, initialUserFile] = await Promise.all([
    getServerUser(),
    getServerUserFile(),
  ]);

  const initialCart = initialUser ? await getServerCart() : null;

  return (
    <ClientSiteLayout
      initialUser={initialUser}
      initialUserFile={initialUser ? initialUserFile : null}
      initialCart={initialCart}
    >
      {children}
    </ClientSiteLayout>
  );
}
