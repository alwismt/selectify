"use client";

import EditableImageGallery from "./EditableImageGallery";
import type { EditableImage } from "@/lib/editableImage";

type ProductImagesEditorProps = {
  images: EditableImage[];
  onAdd: (files: File[]) => void;
  onRemove: (id: string) => void;
  onSetPrimary: (id: string) => void;
  error?: string | null;
};

export default function ProductImagesEditor({
  images,
  onAdd,
  onRemove,
  onSetPrimary,
  error,
}: ProductImagesEditorProps) {
  return (
    <EditableImageGallery
      title="Product Images"
      size="lg"
      images={images}
      onAdd={onAdd}
      onRemove={onRemove}
      onSetPrimary={onSetPrimary}
      error={error}
    />
  );
}
