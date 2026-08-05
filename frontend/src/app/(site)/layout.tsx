import { getServerUser } from "@/lib/api/getServerUser";
import ClientSiteLayout from "./ClientSiteLayout";

export default async function SiteLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const initialUser = await getServerUser();
  return (
    <ClientSiteLayout initialUser={initialUser}>{children}</ClientSiteLayout>
  );
}
