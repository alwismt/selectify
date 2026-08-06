"use client";

import { useEffect, useRef, useState } from "react";
import { useUser } from "@/app/context/UserContext";
import {
  clientUserFileDelete,
  clientUserFileUpload,
} from "@/lib/api/client";
import UserAvatar from "@/components/Common/UserAvatar";

const ACCEPTED_TYPES = ["image/jpeg", "image/png", "image/webp"];
const MAX_FILE_SIZE = 5 * 1024 * 1024;

export default function AvatarUpload() {
  const { userFile, setUserFile } = useUser();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    return () => {
      if (previewUrl) {
        URL.revokeObjectURL(previewUrl);
      }
    };
  }, [previewUrl]);

  const clearPreview = () => {
    if (previewUrl) {
      URL.revokeObjectURL(previewUrl);
      setPreviewUrl(null);
    }
  };

  const handleFileSelect = async (file: File) => {
    setError(null);

    if (!ACCEPTED_TYPES.includes(file.type)) {
      setError("Please upload a JPEG, PNG, or WebP image.");
      return;
    }

    if (file.size > MAX_FILE_SIZE) {
      setError("Image must be 5 MB or smaller.");
      return;
    }

    clearPreview();
    const localPreview = URL.createObjectURL(file);
    setPreviewUrl(localPreview);
    setIsUploading(true);

    try {
      const updatedFile = await clientUserFileUpload(file);
      setUserFile(updatedFile);
      clearPreview();
    } catch {
      setError("Failed to upload image. Please try again.");
      clearPreview();
    } finally {
      setIsUploading(false);
      if (fileInputRef.current) {
        fileInputRef.current.value = "";
      }
    }
  };

  const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      void handleFileSelect(file);
    }
  };

  const handleDelete = async () => {
    if (!userFile) return;

    setError(null);
    setIsDeleting(true);

    try {
      await clientUserFileDelete();
      setUserFile(null);
      clearPreview();
    } catch {
      setError("Failed to delete image. Please try again.");
    } finally {
      setIsDeleting(false);
    }
  };

  const displayPreview = previewUrl ? (
    <img
      src={previewUrl}
      alt="Selected profile photo preview"
      width={96}
      height={96}
      className="h-24 w-24 rounded-full object-cover"
    />
  ) : (
    <UserAvatar
      userFile={userFile}
      size={96}
      fallback={
        <div className="flex h-24 w-24 items-center justify-center rounded-full bg-gray-1 text-custom-sm text-dark-4">
          No photo
        </div>
      }
    />
  );

  return (
    <div className="mb-8 border-b border-gray-3 pb-8">
      <p className="mb-4 font-medium text-xl text-dark">Profile Photo</p>

      <div className="flex flex-col gap-5 sm:flex-row sm:items-center">
        {displayPreview}

        <div className="flex flex-col gap-3">
          <input
            ref={fileInputRef}
            type="file"
            accept={ACCEPTED_TYPES.join(",")}
            className="hidden"
            onChange={handleInputChange}
          />

          <button
            type="button"
            disabled={isUploading || isDeleting}
            onClick={() => fileInputRef.current?.click()}
            className="inline-flex font-medium text-white bg-blue py-2.5 px-6 rounded-md ease-out duration-200 hover:bg-blue-dark disabled:cursor-not-allowed disabled:opacity-60"
          >
            {isUploading
              ? "Uploading..."
              : userFile
                ? "Replace Photo"
                : "Upload Photo"}
          </button>

          {userFile && (
            <button
              type="button"
              disabled={isUploading || isDeleting}
              onClick={() => void handleDelete()}
              className="inline-flex font-medium text-red py-2.5 px-6 rounded-md border border-red ease-out duration-200 hover:bg-red/5 disabled:cursor-not-allowed disabled:opacity-60"
            >
              {isDeleting ? "Deleting..." : "Delete Photo"}
            </button>
          )}

          <p className="text-custom-xs text-dark-4">
            JPEG, PNG, or WebP. Max 5 MB.
          </p>

          {error && <p className="text-custom-sm text-red">{error}</p>}
        </div>
      </div>
    </div>
  );
}
