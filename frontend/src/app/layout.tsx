import type { Metadata } from "next";
import "./css/euclid-circular-a-font.css";
import "./css/style.css";
import { getServerSiteConfig } from "@/lib/api/getServerSiteConfig";
import { SiteConfigProvider } from "@/app/context/SiteConfigContext";

const siteName = "NextCommerce";
const defaultTitle = "NextCommerce | Nextjs E-commerce template";
const defaultDescription =
  "NextCommerce - Nextjs E-commerce template for modern online stores";

export const metadata: Metadata = {
  metadataBase: new URL(
    process.env.NEXT_PUBLIC_SITE_URL ?? "http://localhost:3000"
  ),
  title: {
    default: defaultTitle,
    template: `%s | ${siteName}`,
  },
  description: defaultDescription,
  openGraph: {
    type: "website",
    siteName,
    title: defaultTitle,
    description: defaultDescription,
  },
  twitter: {
    card: "summary_large_image",
    title: defaultTitle,
    description: defaultDescription,
  },
};

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const initialConfig = await getServerSiteConfig();

  return (
    <html lang="en" data-scroll-behavior="smooth" suppressHydrationWarning>
      <body suppressHydrationWarning>
        <SiteConfigProvider initialConfig={initialConfig}>
          {children}
        </SiteConfigProvider>
      </body>
    </html>
  );
}
