const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  outputDir: './../www',
  productionSourceMap: false,
  transpileDependencies: true,
  pages: {
    index: {
      entry: 'src/main.js',
      title: 'TRAF'
    }
  },
  devServer: {
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:8086',
        changeOrigin: true,
        //pathRewrite: { '^/api': '' },
        secure: false,
        // onProxyReq: (proxyReq, req, res) => {
        //   console.log('Proxying request:', req.url);
        // },
        // onProxyRes: (proxyRes, req, res) => {
        //   console.log('Received response from target:', req.url);
        // }
      }
    }
  }
})
