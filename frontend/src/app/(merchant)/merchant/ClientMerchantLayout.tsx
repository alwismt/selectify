"use client";

import Header from "@/components/Header";
import { CartModalProvider } from "@/app/context/CartSidebarModalContext";
import { CartProvider } from "@/app/context/CartContext";
import { UserProvider } from "@/app/context/UserContext";
import { MerchantProvider } from "@/app/context/MerchantContext";
import type { User } from "@/types/user";
import type { UserFile } from "@/types/api/userFile";
import type { Merchant } from "@/types/merchant";

interface ClientMerchantLayoutProps {
  initialUser: User | null;
  initialUserFile: UserFile | null;
  initialMerchant: Merchant | null;
  children: React.ReactNode;
}

export default function ClientMerchantLayout({
  initialUser,
  initialUserFile,
  initialMerchant,
  children,
}: ClientMerchantLayoutProps) {
  return (
    <UserProvider initialUser={initialUser} initialUserFile={initialUserFile}>
      <MerchantProvider initialMerchant={initialMerchant}>
        <CartProvider>
          <CartModalProvider>
            <Header variant="merchant" />
            {children}
          </CartModalProvider>
        </CartProvider>
      </MerchantProvider>
    </UserProvider>
  );
}
