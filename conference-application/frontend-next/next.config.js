/** @type {import('next').NextConfig} */
module.exports = {
    // output: 'export', //-> this is to get static files
    // output: 'standalone', -> this is to get a server for nodejs
    async rewrites() {
      return [
        {
          source: '/api/agenda',
          destination: 'http://frontend-go.default.74.220.17.238.sslip.io/api/agenda/',
        },
        {
            source: '/api/agenda/:path*',
            destination: 'http://frontend-go.default.74.220.17.238.sslip.io/api/agenda/:path*',
        },
        {
            source: '/api/c4p',
            destination: 'http://frontend-go.default.74.220.17.238.sslip.io/api/c4p/',
        },
        {
            source: '/api/c4p/:path*',
            destination: 'http://frontend-go.default.74.220.17.238.sslip.io/api/c4p/:path*',
        },
        {
          source: '/api/events',
          destination: 'http://frontend-go.default.74.220.17.238.sslip.io/api/events/',
        },
        {
          source: '/api/events/:path*',
          destination: 'http://frontend-go.default.74.220.17.238.sslip.io/api/events/:path*',
        },
        {
          source: '/api/notifications',
          destination: 'http://frontend-go.default.74.220.17.238.sslip.io/api/notifications/',
        },
        {
          source: '/api/notifications/:path*',
          destination: 'http://frontend-go.default.74.220.17.238.sslip.io/api/notifications/:path*',
        },
      ]
    },
  }