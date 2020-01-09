const CompressionWebpackPlugin = require('compression-webpack-plugin')

const productionGzipExtensions = ['js', 'css']
const isProduction = process.env.NODE_ENV === 'production'

module.exports = {
  publicPath: '/',
  outputDir: 'dist',
  // eslint
  lintOnSave: true,
  filenameHashing: true,
  productionSourceMap: false,
  transpileDependencies: [
    'vuetify',
  ],
  // 开启gzip
  configureWebpack: (config) => {
    if (isProduction) {
      // gzip
      config.plugins.push(new CompressionWebpackPlugin({
        algorithm: 'gzip',
        test: new RegExp(`\\.(${productionGzipExtensions.join('|')})$`),
        threshold: 10240,
        minRatio: 0.8,
      }))
    }
  },
}
