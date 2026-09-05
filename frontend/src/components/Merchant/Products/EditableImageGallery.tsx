"use client";

import { ChangeEvent, useRef } from "react";
import type { EditableImage } from "@/lib/editableImage";

type EditableImageGalleryProps = {
  title: string;
  images: EditableImage[];
  onAdd: (files: File[]) => void;
  onRemove: (id: string) => void;
  onSetPrimary: (id: string) => void;
  error?: string | null;
  size?: "lg" | "md";
};

const sizeClass = {
  lg: "h-40 w-40 sm:h-48 sm:w-48",
  md: "h-24 w-24 sm:h-28 sm:w-28",
};

export default function EditableImageGallery({
  title,
  images,
  onAdd,
  onRemove,
  onSetPrimary,
  error,
  size = "md",
}: EditableImageGalleryProps) {
  const inputRef = useRef<HTMLInputElement>(null);
  const cardSize = sizeClass[size];

  const handleInputChange = (event: ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files ? Array.from(event.target.files) : [];
    if (files.length > 0) {
      onAdd(files);
    }
    event.target.value = "";
  };

  return (
    <section className="mt-8 pt-7.5 border-t border-gray-3">
      <h3 className="font-medium text-lg text-dark mb-5">{title}</h3>
      <div className="flex flex-wrap gap-3">
        {images.map((image) => (
          <div key={image.id} className="flex flex-col gap-2">
            <div className="relative">
              {image.previewUrl ? (
                // eslint-disable-next-line @next/next/no-img-element
                <img
                  src={image.previewUrl}
                  alt=""
                  className={`${cardSize} rounded-md object-cover border border-gray-3`}
                />
              ) : (
                <div
                  className={`${cardSize} rounded-md border border-gray-3 bg-gray-1`}
                  aria-hidden="true"
                />
              )}
              {image.isPrimary ? (
                <span className="absolute top-2 left-2 rounded bg-blue px-2 py-0.5 text-xs font-medium text-white">
                  Primary
                </span>
              ) : null}
            </div>
            <div className="flex flex-col items-start gap-1">
              {!image.isPrimary ? (
                <button
                  type="button"
                  onClick={() => onSetPrimary(image.id)}
                  aria-label="Set image as primary"
                  className="text-custom-sm font-medium text-blue hover:text-blue-dark"
                >
                  Set as Primary
                </button>
              ) : null}
              <button
                type="button"
                onClick={() => onRemove(image.id)}
                aria-label="Remove image"
                className="text-custom-sm font-medium text-red hover:underline"
              >
                Remove
              </button>
            </div>
          </div>
        ))}

        {images.length === 0 ? (
          <div
            className={`${cardSize} rounded-md border border-gray-3 bg-gray-1`}
            aria-hidden="true"
          />
        ) : null}

        <button
          type="button"
          onClick={() => inputRef.current?.click()}
          aria-label="Add images"
          className={`${cardSize} rounded-lg border border-dashed border-gray-3 bg-gray-1 flex flex-col items-center justify-center gap-1 px-3 text-center hover:border-blue`}
        >
          <span className="text-custom-sm font-medium text-blue">+ Add image</span>
          <span className="text-custom-xs text-dark-4">PNG, JPG, WEBP</span>
        </button>
      </div>

      <input
        ref={inputRef}
        type="file"
        accept="image/*"
        multiple
        className="sr-only"
        onChange={handleInputChange}
      />

      {error ? <p className="mt-3 text-custom-sm text-red">{error}</p> : null}
    </section>
  );
}
