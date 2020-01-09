const pluginsTransform = ['@vue/babel-plugin-transform-vue-jsx']
// 生产环境移除console
if (process.env.NODE_ENV === 'production') {
  pluginsTransform.push('transform-remove-console')
}
module.exports = {
  plugins: pluginsTransform,
  presets: [
    '@vue/app',
  ],
}
