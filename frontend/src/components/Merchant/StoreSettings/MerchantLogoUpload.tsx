"use client";

import { useEffect, useRef, useState } from "react";
import { useMerchant } from "@/app/context/MerchantContext";
import { uploadMerchantLogo } from "@/lib/api/actions/uploadMerchantLogo";
import { merchantLogoUrl } from "@/lib/api/config";

const ACCEPTED_TYPES = ["image/jpeg", "image/png", "image/webp"];
const MAX_FILE_SIZE = 5 * 1024 * 1024;

export default function MerchantLogoUpload() {
  const { merchant, setMerchant } = useMerchant();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [isUploading, setIsUploading] = useState(false);
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
    if (!merchant) return;

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
      const formData = new FormData();
      formData.append("image", file);
      const result = await uploadMerchantLogo(formData);

      if (result.ok === false) {
        setError(result.error);
        clearPreview();
      } else {
        setMerchant({ ...merchant, logo: result.logo });
        clearPreview();
      }
    } catch {
      setError("Failed to upload logo. Please try again.");
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

  const logoSrc = merchantLogoUrl(merchant?.logo);
  const displayPreview = previewUrl ? (
    <img
      src={previewUrl}
      alt="Selected store logo preview"
      width={96}
      height={96}
      className="h-24 w-24 rounded-full object-cover"
    />
  ) : logoSrc ? (
    <img
      src={logoSrc}
      alt={`${merchant?.name ?? "Store"} logo`}
      width={96}
      height={96}
      className="h-24 w-24 rounded-full object-cover"
    />
  ) : (
    <div className="flex h-24 w-24 items-center justify-center rounded-full bg-gray-1 text-custom-sm text-dark-4">
      No logo
    </div>
  );

  return (
    <div className="mb-8 border-b border-gray-3 pb-8">
      <p className="mb-4 font-medium text-xl text-dark">Store Logo</p>

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
            disabled={isUploading || !merchant}
            onClick={() => fileInputRef.current?.click()}
            className="inline-flex font-medium text-white bg-blue py-2.5 px-6 rounded-md ease-out duration-200 hover:bg-blue-dark disabled:cursor-not-allowed disabled:opacity-60"
          >
            {isUploading
              ? "Uploading..."
              : merchant?.logo
                ? "Replace Logo"
                : "Upload Logo"}
          </button>

          <p className="text-custom-xs text-dark-4">
            JPEG, PNG, or WebP. Max 5 MB.
          </p>

          {error && <p className="text-custom-sm text-red">{error}</p>}
        </div>
      </div>
    </div>
  );
}
