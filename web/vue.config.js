const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  outputDir: './../www',
  productionSourceMap: false,
  transpileDependencies: true
})
