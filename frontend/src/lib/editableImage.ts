import { productFileUrl } from "@/lib/api/config";
import type { ApiProductFile } from "@/types/product";
import type { MerchantApiVariantFile } from "@/types/merchantProduct";

export type EditableImage = {
  id: string;
  fileId?: string;
  file?: File;
  previewUrl: string;
  isPrimary: boolean;
  isNew: boolean;
};

export function imagesFromProductFile(
  file: ApiProductFile | null | undefined
): EditableImage[] {
  if (!file) return [];
  return [
    {
      id: file.file_id,
      fileId: file.file_id,
      previewUrl: productFileUrl(file.file_id),
      isPrimary: file.is_primary ?? true,
      isNew: false,
    },
  ];
}

export function imagesFromVariantFiles(
  files: MerchantApiVariantFile[]
): EditableImage[] {
  const mapped = [...files]
    .sort((a, b) => a.position - b.position)
    .map((file) => ({
      id: file.file_id,
      fileId: file.file_id,
      previewUrl: productFileUrl(file.file_id),
      isPrimary: file.is_primary,
      isNew: false,
    }));

  if (mapped.length > 0 && !mapped.some((image) => image.isPrimary)) {
    mapped[0] = { ...mapped[0], isPrimary: true };
  }

  return mapped;
}

export function appendImageFiles(
  existing: EditableImage[],
  files: File[]
): EditableImage[] {
  const added = files.map((file, index) => ({
    id: crypto.randomUUID(),
    file,
    previewUrl: URL.createObjectURL(file),
    isPrimary: existing.length === 0 && index === 0,
    isNew: true,
  }));

  return [...existing, ...added];
}

export function setPrimaryImage(
  images: EditableImage[],
  id: string
): EditableImage[] {
  return images.map((image) => ({
    ...image,
    isPrimary: image.id === id,
  }));
}

export function removeEditableImage(
  images: EditableImage[],
  id: string
): EditableImage[] {
  const target = images.find((image) => image.id === id);
  if (target?.isNew) {
    URL.revokeObjectURL(target.previewUrl);
  }

  const remaining = images.filter((image) => image.id !== id);
  if (remaining.length > 0 && !remaining.some((image) => image.isPrimary)) {
    remaining[0] = { ...remaining[0], isPrimary: true };
  }
  return remaining;
}

export function revokeNewImageUrls(images: EditableImage[]): void {
  for (const image of images) {
    if (image.isNew) {
      URL.revokeObjectURL(image.previewUrl);
    }
  }
}

export function imagesAreDirty(
  current: EditableImage[],
  initial: EditableImage[]
): boolean {
  if (current.length !== initial.length) return true;

  return current.some((image, index) => {
    const other = initial[index];
    return (
      image.id !== other.id ||
      image.fileId !== other.fileId ||
      image.isPrimary !== other.isPrimary ||
      image.isNew !== other.isNew
    );
  });
}
