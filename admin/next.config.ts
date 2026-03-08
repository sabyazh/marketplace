import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // 'standalone' output is only needed for Docker deployments
  // Vercel handles builds natively without it
  ...(process.env.VERCEL ? {} : { output: 'standalone' }),
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: '**',
      },
    ],
  },
};

export default nextConfig;
