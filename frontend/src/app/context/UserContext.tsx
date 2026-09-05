"use client";

import React, { createContext, useContext, useState, useMemo, useEffect } from "react";
import type { User } from "@/types/user";
import type { UserFile } from "@/types/api/userFile";
import { SESSION_EXPIRED_EVENT } from "@/lib/api/session";

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
  const [prevInitialUser, setPrevInitialUser] = useState(initialUser);
  const [prevInitialUserFile, setPrevInitialUserFile] =
    useState(initialUserFile);

  // Sync after router.refresh() when layout re-fetches SSR user.
  if (initialUser !== prevInitialUser) {
    setPrevInitialUser(initialUser);
    setUser(initialUser);
  }
  if (initialUserFile !== prevInitialUserFile) {
    setPrevInitialUserFile(initialUserFile);
    setUserFile(initialUserFile);
  }

  useEffect(() => {
    const onExpired = () => {
      setUser(null);
      setUserFile(null);
    };
    window.addEventListener(SESSION_EXPIRED_EVENT, onExpired);
    return () => window.removeEventListener(SESSION_EXPIRED_EVENT, onExpired);
  }, []);

  const value = useMemo(
    () => ({ user, setUser, userFile, setUserFile }),
    [user, userFile]
  );
  return (
    <UserContext.Provider value={value}>{children}</UserContext.Provider>
  );
}
