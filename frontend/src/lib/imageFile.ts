export const ACCEPTED_IMAGE_TYPES = ["image/jpeg", "image/png", "image/webp"] as const;
export const MAX_IMAGE_FILE_SIZE = 5 * 1024 * 1024;

const acceptedTypeSet = new Set<string>(ACCEPTED_IMAGE_TYPES);

export function validateImageFile(file: File): string | null {
  if (!file.type.startsWith("image/")) {
    return "Please upload an image file.";
  }
  if (!acceptedTypeSet.has(file.type)) {
    return "Please upload a JPEG, PNG, or WebP image.";
  }
  if (file.size > MAX_IMAGE_FILE_SIZE) {
    return "Image must be 5 MB or smaller.";
  }
  return null;
}

export function partitionImageFiles(files: File[]): {
  accepted: File[];
  error: string | null;
} {
  const accepted: File[] = [];
  let error: string | null = null;

  for (const file of files) {
    const fileError = validateImageFile(file);
    if (fileError) {
      if (!error) error = fileError;
      continue;
    }
    accepted.push(file);
  }

  return { accepted, error };
}
