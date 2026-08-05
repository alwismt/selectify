"use client";

import { usePathname } from "next/navigation";
import Header from "../../components/Header";
import Footer from "../../components/Footer";
import { ModalProvider } from "../context/QuickViewModalContext";
import { CartModalProvider } from "../context/CartSidebarModalContext";
import { CartProvider } from "../context/CartContext";
import { UserProvider } from "../context/UserContext";
import { ReduxProvider } from "@/redux/provider";
import QuickViewModal from "@/components/Common/QuickViewModal";
import CartSidebarModal from "@/components/Common/CartSidebarModal";
import { PreviewSliderProvider } from "../context/PreviewSliderContext";
import PreviewSliderModal from "@/components/Common/PreviewSlider";
import ScrollToTop from "@/components/Common/ScrollToTop";
import type { User } from "@/types/user";

interface ClientSiteLayoutProps {
  initialUser: User | null;
  children: React.ReactNode;
}

export default function ClientSiteLayout({
  initialUser,
  children,
}: ClientSiteLayoutProps) {
  const pathname = usePathname();
  const isCheckout = pathname?.startsWith("/checkout") ?? false;

  return (
    <>
      <ReduxProvider>
        <UserProvider initialUser={initialUser}>
          <CartProvider>
            <CartModalProvider>
              <ModalProvider>
              <PreviewSliderProvider>
                {!isCheckout && <Header />}
                {children}
                <QuickViewModal />
                <CartSidebarModal />
                <PreviewSliderModal />
              </PreviewSliderProvider>
              </ModalProvider>
            </CartModalProvider>
          </CartProvider>
        </UserProvider>
      </ReduxProvider>
      {!isCheckout && <ScrollToTop />}
      {!isCheckout && <Footer />}
    </>
  );
}
