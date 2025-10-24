const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
  // 多入口配置：键（如page1）是入口名称，值是JS文件路径
  entry: {
    page1: './src/page1/page1.js',  // 页面1的JS入口
    page2: './src/page2/page2.js'   // 页面2的JS入口
  },
  output: {
    // 输出的JS文件名：[name]会被替换为entry的键（如page1、page2）
    // 如果要控制客户端缓存，最好还要加上[chunkhash]，因为每个chunk所产生的[chunkhash]只与自身内容有关，单个chunk内容的改变不会影响其他资源，可以最精确地让客户端缓存得到更新。
    filename: '[name]@[chunkhash].bundle.js',  
    path: path.resolve(__dirname, 'dist'), // 所有文件输出到dist目录
    clean: true // 每次打包前清空dist目录（可选，推荐）
  },
  mode: 'development',
  devServer: {
    static: path.resolve(__dirname, 'dist'),
    port: 10240,
    open: '/page1.html'  // 自动打开默认页面（可指定打开page1.html）
  },
  plugins: [
    // 页面1的HTML配置：关联page1的JS
    new HtmlWebpackPlugin({
      template: './src/page1/page1.html', // 页面1的HTML模板路径
      filename: 'page1.html', // 输出到dist的文件名（如page1.html）
      chunks: ['page1','common'] // 只引入page1的JS（对应entry的键）
    }),
    // 页面2的HTML配置：关联page2的JS
    new HtmlWebpackPlugin({
      template: './src/page2/page2.html', // 页面2的HTML模板路径
      filename: 'page2.html', // 输出到dist的文件名
      chunks: ['page2'] // 只引入page2的JS
    })
  ]
};