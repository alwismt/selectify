"use client";

import { userFileUrl } from "@/lib/api/config";
import type { UserFile } from "@/types/api/userFile";

type UserAvatarProps = {
  userFile: UserFile | null;
  size?: number;
  className?: string;
  alt?: string;
  fallback?: React.ReactNode;
};

export default function UserAvatar({
  userFile,
  size = 64,
  className = "",
  alt = "Profile photo",
  fallback = null,
}: UserAvatarProps) {
  const src = userFileUrl(userFile?.id);

  if (!src) {
    return fallback ? <>{fallback}</> : null;
  }

  return (
    <img
      src={src}
      alt={alt}
      width={size}
      height={size}
      className={`rounded-full object-cover ${className}`}
      style={{ width: size, height: size }}
    />
  );
}
