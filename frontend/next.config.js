/** @type {import('next').NextConfig} */

function productFilesRemotePatterns() {
  const base = process.env.NEXT_PUBLIC_PRODUCT_FILES_BASE_URL;
  if (!base) return [];
  try {
    const url = new URL(base);
    const protocol = url.protocol === "http:" ? "http" : "https";
    return [
      {
        protocol,
        hostname: url.hostname,
        ...(url.port ? { port: url.port } : {}),
        pathname: "/products/**",
      },
    ];
  } catch {
    return [];
  }
}

const nextConfig = {
  allowedDevOrigins: [
    "selectify.alwis.dev",
    "localhost:3000",
  ],
  images: {
    remotePatterns: productFilesRemotePatterns(),
  },
};


module.exports = nextConfig;
