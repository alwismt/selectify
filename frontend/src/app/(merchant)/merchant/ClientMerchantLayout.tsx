"use client";

import Header from "@/components/Header";
import { CartModalProvider } from "@/app/context/CartSidebarModalContext";
import { CartProvider } from "@/app/context/CartContext";
import { UserProvider } from "@/app/context/UserContext";
import type { User } from "@/types/user";
import type { UserFile } from "@/types/api/userFile";

interface ClientMerchantLayoutProps {
  initialUser: User | null;
  initialUserFile: UserFile | null;
  children: React.ReactNode;
}

export default function ClientMerchantLayout({
  initialUser,
  initialUserFile,
  children,
}: ClientMerchantLayoutProps) {
  return (
    <UserProvider initialUser={initialUser} initialUserFile={initialUserFile}>
      <CartProvider>
        <CartModalProvider>
          <Header variant="merchant" />
          {children}
        </CartModalProvider>
      </CartProvider>
    </UserProvider>
  );
}
