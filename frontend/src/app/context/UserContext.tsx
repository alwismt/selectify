"use client";

import React, { createContext, useContext, useState, useMemo } from "react";
import type { User } from "@/types/user";

interface UserContextType {
  user: User | null;
  setUser: (user: User | null) => void;
}

const UserContext = createContext<UserContextType | undefined>(undefined);

export function useUser() {
  const context = useContext(UserContext);
  if (context === undefined) {
    throw new Error("useUser must be used within a UserProvider");
  }
  return context;
}

interface UserProviderProps {
  initialUser: User | null;
  children: React.ReactNode;
}

export function UserProvider({ initialUser, children }: UserProviderProps) {
  const [user, setUser] = useState<User | null>(initialUser);
  const value = useMemo(() => ({ user, setUser }), [user]);
  return (
    <UserContext.Provider value={value}>{children}</UserContext.Provider>
  );
}
