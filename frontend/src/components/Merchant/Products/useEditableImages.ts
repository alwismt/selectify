"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import {
  appendImageFiles,
  imagesAreDirty,
  removeEditableImage,
  revokeNewImageUrls,
  setPrimaryImage,
  type EditableImage,
} from "@/lib/editableImage";
import { partitionImageFiles } from "@/lib/imageFile";

export function useEditableImages(initialImages: EditableImage[]) {
  const [images, setImages] = useState<EditableImage[]>(initialImages);
  const [error, setError] = useState<string | null>(null);
  const imagesRef = useRef(images);
  const initialRef = useRef(initialImages);

  imagesRef.current = images;

  useEffect(() => {
    return () => {
      revokeNewImageUrls(imagesRef.current);
    };
  }, []);

  const addFiles = useCallback((files: File[]) => {
    const { accepted, error: fileError } = partitionImageFiles(files);
    setError(fileError);
    if (accepted.length === 0) return;
    setImages((prev) => appendImageFiles(prev, accepted));
  }, []);

  const remove = useCallback((id: string) => {
    setError(null);
    setImages((prev) => removeEditableImage(prev, id));
  }, []);

  const setPrimary = useCallback((id: string) => {
    setError(null);
    setImages((prev) => setPrimaryImage(prev, id));
  }, []);

  const isDirty = imagesAreDirty(images, initialRef.current);

  return { images, error, addFiles, remove, setPrimary, isDirty };
}
