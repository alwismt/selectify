import { getServerUser } from "@/lib/api/getServerUser";
import { getServerUserFile } from "@/lib/api/getServerUserFile";
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

  return (
    <ClientSiteLayout
      initialUser={initialUser}
      initialUserFile={initialUser ? initialUserFile : null}
    >
      {children}
    </ClientSiteLayout>
  );
}
