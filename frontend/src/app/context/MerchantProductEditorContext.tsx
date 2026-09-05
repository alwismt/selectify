"use client";

import React, {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useState,
} from "react";
import type { ApiProductFile } from "@/types/product";
import type {
  MerchantApiProduct,
  MerchantApiProductVariant,
} from "@/types/merchantProduct";

type MerchantProductEditorContextType = {
  product: MerchantApiProduct | null;
  variants: MerchantApiProductVariant[] | null;
  loadedProductId: number | null;
  setProductEditorData: (
    product: MerchantApiProduct,
    variants: MerchantApiProductVariant[] | null
  ) => void;
  updateProduct: (patch: Partial<MerchantApiProduct>) => void;
  updateVariant: (variant: MerchantApiProductVariant) => void;
  addVariant: (variant: MerchantApiProductVariant) => void;
  removeVariant: (variantId: number) => void;
  setProductImage: (file: ApiProductFile | null) => void;
};

const MerchantProductEditorContext = createContext<
  MerchantProductEditorContextType | undefined
>(undefined);

export function useMerchantProductEditor() {
  const context = useContext(MerchantProductEditorContext);
  if (context === undefined) {
    throw new Error(
      "useMerchantProductEditor must be used within a MerchantProductEditorProvider"
    );
  }
  return context;
}

export function MerchantProductEditorProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const [product, setProduct] = useState<MerchantApiProduct | null>(null);
  const [variants, setVariants] = useState<MerchantApiProductVariant[] | null>(
    null
  );
  const [loadedProductId, setLoadedProductId] = useState<number | null>(null);

  const setProductEditorData = useCallback(
    (
      nextProduct: MerchantApiProduct,
      nextVariants: MerchantApiProductVariant[] | null
    ) => {
      setProduct(nextProduct);
      setVariants(nextVariants);
      setLoadedProductId(nextProduct.productId);
    },
    []
  );

  const updateProduct = useCallback((patch: Partial<MerchantApiProduct>) => {
    setProduct((prev) => (prev ? { ...prev, ...patch } : prev));
  }, []);

  const updateVariant = useCallback((variant: MerchantApiProductVariant) => {
    setVariants((prev) =>
      prev
        ? prev.map((item) => (item.id === variant.id ? variant : item))
        : prev
    );
  }, []);

  const addVariant = useCallback((variant: MerchantApiProductVariant) => {
    setVariants((prev) => (prev ? [...prev, variant] : [variant]));
  }, []);

  const removeVariant = useCallback((variantId: number) => {
    setVariants((prev) =>
      prev ? prev.filter((item) => item.id !== variantId) : prev
    );
  }, []);

  const setProductImage = useCallback((file: ApiProductFile | null) => {
    setProduct((prev) => (prev ? { ...prev, productFile: file } : prev));
  }, []);

  const value = useMemo(
    () => ({
      product,
      variants,
      loadedProductId,
      setProductEditorData,
      updateProduct,
      updateVariant,
      addVariant,
      removeVariant,
      setProductImage,
    }),
    [
      product,
      variants,
      loadedProductId,
      setProductEditorData,
      updateProduct,
      updateVariant,
      addVariant,
      removeVariant,
      setProductImage,
    ]
  );

  return (
    <MerchantProductEditorContext.Provider value={value}>
      {children}
    </MerchantProductEditorContext.Provider>
  );
}
