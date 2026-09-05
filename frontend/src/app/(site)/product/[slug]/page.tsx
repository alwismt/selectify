import { Metadata } from "next";
import { notFound, redirect } from "next/navigation";
import ShopDetails from "@/components/ShopDetails";
import { getProductByIdFromApi, getVariants } from "@/lib/products";
import {
  canonicalProductPath,
  parseProductIdFromPath,
  productHref,
} from "@/lib/productPath";

type ProductPageProps = {
  params: Promise<{ slug: string }>;
};

export async function generateMetadata({
  params,
}: ProductPageProps): Promise<Metadata> {
  const { slug } = await params;
  const productId = parseProductIdFromPath(slug);
  if (productId == null) {
    return {
      title: "Product Not Found",
    };
  }

  const product = await getProductByIdFromApi(productId);
  if (!product) {
    return {
      title: "Product Not Found",
    };
  }
  const description = product.description ?? product.title;
  const image =
    product.imgs?.previews?.[0] != null ? product.imgs.previews[0] : undefined;
  const url = productHref(product.id, product.slug);
  return {
    title: product.title,
    description,
    openGraph: {
      title: product.title,
      description,
      images: image ? [image] : [],
      url,
    },
    twitter: {
      card: "summary_large_image",
      title: product.title,
      description,
    },
    alternates: {
      canonical: url,
    },
  };
}

export default async function ProductPage({ params }: ProductPageProps) {
  const { slug } = await params;
  const productId = parseProductIdFromPath(slug);
  if (productId == null) notFound();

  const product = await getProductByIdFromApi(productId);
  if (!product) notFound();

  const productSlug =
    typeof product.slug === "string" ? product.slug.trim() : "";
  if (!productSlug) notFound();

  const canonical = canonicalProductPath(product.id, productSlug);
  if (slug !== canonical) {
    redirect(productHref(product.id, productSlug));
  }

  const variants = await getVariants(product.id);
  return (
    <main>
      <ShopDetails initialProduct={product} variants={variants} />
    </main>
  );
}
