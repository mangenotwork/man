const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin'); // 引入插件

module.exports = {
    entry: './src/index.js',
    output: {
        filename: 'main.js',
        path: path.resolve(__dirname, 'dist'), // 明确指定打包到 dist 文件夹（必须是绝对路径）
    },
    mode: 'development',
    devServer: {
        static: path.resolve(__dirname, 'dist'), 
        port: 10240, // 其他配置（如端口、自动打开等）保留
        open: true
    },
    plugins: [
        // 配置 HTML 插件：根据模板生成 HTML 并注入 JS
        new HtmlWebpackPlugin({
            template: './src/index.html', // 你的 HTML 模板（需手动创建）
            filename: 'index.html' // 生成到 dist 目录的文件名
        })
    ]
}