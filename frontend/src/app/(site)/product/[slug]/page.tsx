import { Metadata } from "next";
import { notFound } from "next/navigation";
import ShopDetails from "@/components/ShopDetails";
import { getProductBySlugFromApi, getVariants } from "@/lib/products";

type ProductPageProps = {
  params: Promise<{ slug: string }>;
};

export async function generateMetadata({
  params,
}: ProductPageProps): Promise<Metadata> {
  const { slug } = await params;
  const product = await getProductBySlugFromApi(slug);
  if (!product) {
    return {
      title: "Product Not Found",
    };
  }
  const description = product.description ?? product.title;
  const image =
    product.imgs?.previews?.[0] != null ? product.imgs.previews[0] : undefined;
  const url = `/product/${product.slug}`;
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
  const product = await getProductBySlugFromApi(slug);
  if (!product) notFound();
  const variants = await getVariants(product.id);
  return (
    <main>
      <ShopDetails initialProduct={product} variants={variants} />
    </main>
  );
}
