"use client";

import React, { createContext, useContext, useState, useMemo } from "react";
import type { User } from "@/types/user";
import type { UserFile } from "@/types/api/userFile";

interface UserContextType {
  user: User | null;
  setUser: (user: User | null) => void;
  userFile: UserFile | null;
  setUserFile: (userFile: UserFile | null) => void;
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
  initialUserFile: UserFile | null;
  children: React.ReactNode;
}

export function UserProvider({
  initialUser,
  initialUserFile,
  children,
}: UserProviderProps) {
  const [user, setUser] = useState<User | null>(initialUser);
  const [userFile, setUserFile] = useState<UserFile | null>(initialUserFile);
  const value = useMemo(
    () => ({ user, setUser, userFile, setUserFile }),
    [user, userFile]
  );
  return (
    <UserContext.Provider value={value}>{children}</UserContext.Provider>
  );
}
