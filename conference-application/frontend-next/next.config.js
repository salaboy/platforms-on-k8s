module.exports = {
    async rewrites() {
      return [
        {
          source: '/api/agenda',
          destination: 'http://frontend-go.default.74.220.17.238.sslip.io/agenda/',
        },
        {
            source: '/api/agenda/:path*',
            destination: 'http://frontend-go.default.74.220.17.238.sslip.io/agenda/:path*',
        },
        {
            source: '/api/c4p',
            destination: 'http://frontend-go.default.74.220.17.238.sslip.io/c4p/',
        },
        {
            source: '/api/c4p/:path*',
            destination: 'http://frontend-go.default.74.220.17.238.sslip.io/c4p/:path*',
        },
      ]
    },
  }