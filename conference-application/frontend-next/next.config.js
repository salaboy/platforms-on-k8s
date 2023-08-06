/** @type {import('next').NextConfig} */
module.exports = {
    trailingSlash: true,
    output: 'export', //-> this is to get static files
    images: {
      loader: "custom",
      imageSizes: [16, 32, 48, 64, 96, 128, 256, 384],
      deviceSizes: [640, 750, 828, 1080, 1200, 1920, 2048, 3840],
    },
    transpilePackages: ["next-image-export-optimizer"],
    env: {
      nextImageExportOptimizer_imageFolderPath: "public/images",
      nextImageExportOptimizer_exportFolderPath: "out",
      nextImageExportOptimizer_quality: 75,
      nextImageExportOptimizer_storePicturesInWEBP: true,
      nextImageExportOptimizer_exportFolderName: "nextImageExportOptimizer",
  
      // If you do not want to use blurry placeholder images, then you can set
      // nextImageExportOptimizer_generateAndUseBlurImages to false and pass
      // `placeholder="empty"` to all <ExportedImage> components.
      nextImageExportOptimizer_generateAndUseBlurImages: true,
    },
    // output: 'standalone', -> this is to get a server for nodejs
    // async rewrites() {
    //   return [
    //     {
    //       source: '/api/agenda',
    //       destination: 'http://frontend.default.74.220.17.154.sslip.io/api/agenda/',
    //     },
    //     {
    //         source: '/api/agenda/:path*',
    //         destination: 'http://frontend.default.74.220.17.154.sslip.io/api/agenda/:path*',
    //     },
    //     {
    //         source: '/api/c4p',
    //         destination: 'http://frontend.default.74.220.17.154.sslip.io/api/c4p/',
    //     },
    //     {
    //         source: '/api/c4p/:path*',
    //         destination: 'http://frontend.default.74.220.17.154.sslip.io/api/c4p/:path*',
    //     },
    //     {
    //       source: '/api/events',
    //       destination: 'http://frontend.default.74.220.17.154.sslip.io/api/events/',
    //     },
    //     {
    //       source: '/api/events/:path*',
    //       destination: 'http://frontend.default.74.220.17.154.sslip.io/api/events/:path*',
    //     },
    //     {
    //       source: '/api/notifications',
    //       destination: 'http://frontend.default.74.220.17.154.sslip.io/api/notifications/',
    //     },
    //     {
    //       source: '/api/notifications/:path*',
    //       destination: 'http://frontend.default.74.220.17.154.sslip.io/api/notifications/:path*',
    //     },
    //     {
    //       source: '/service/info',
    //       destination: 'http://frontend.default.74.220.17.154.sslip.io/service/info',
    //     },
    //     {
    //       source: '/api/features',
    //       destination: 'http://frontend.default.74.220.17.154.sslip.io/api/features/',
    //     },
    //   ]
    // },
  }