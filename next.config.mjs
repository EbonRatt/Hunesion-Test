/** @type {import('next').NextConfig} */
const nextConfig = {
  /* config options here */
    async rewrites() {
        return [
            {
                source: "/tunnel",
                destination: "http://localhost:8080/tunnel", // <-- your Spring tunnel endpoint
            },
        ];
    },
};

export default nextConfig;
