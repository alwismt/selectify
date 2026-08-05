"use client";
import React, { createContext, useContext, useState } from "react";

interface PreviewSliderType {
  isModalPreviewOpen: boolean;
  previewIndex: number;
  openPreviewModal: (index?: number) => void;
  closePreviewModal: () => void;
}

const PreviewSlider = createContext<PreviewSliderType | undefined>(undefined);

export const usePreviewSlider = () => {
  const context = useContext(PreviewSlider);
  if (!context) {
    throw new Error("usePreviewSlider must be used within a ModalProvider");
  }
  return context;
};

export const PreviewSliderProvider = ({ children }) => {
  const [isModalPreviewOpen, setIsModalOpen] = useState(false);
  const [previewIndex, setPreviewIndex] = useState(0);

  const openPreviewModal = (index = 0) => {
    setPreviewIndex(index);
    setIsModalOpen(true);
  };

  const closePreviewModal = () => {
    setIsModalOpen(false);
  };

  return (
    <PreviewSlider.Provider
      value={{
        isModalPreviewOpen,
        previewIndex,
        openPreviewModal,
        closePreviewModal,
      }}
    >
      {children}
    </PreviewSlider.Provider>
  );
};
