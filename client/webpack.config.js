const path = require("path")
const ExtractTextPlugin = require("extract-text-webpack-plugin")
const HtmlPlugin = require("html-webpack-plugin")
const CleanPlugin = require("clean-webpack-plugin")
const UglifyJsPlugin = require("uglifyjs-webpack-plugin")
const OptimizeCSSAssetsPlugin = require("optimize-css-assets-webpack-plugin")

var config = {
    context: __dirname,
    devtool: "inline source-map",
    entry: "./index.jsx",
    output: {
        filename: "bundle.js",
        path: path.resolve(__dirname, "dist"),
        publicPath: "/"
    },
    resolve: {
        alias: {
            Assets: path.resolve(__dirname, "assets"), 
            Components: path.resolve(__dirname, "components"), 
            Util: path.resolve(__dirname, "util")
        },
        extensions: [".jsx", ".js", ".json", "scss", "css", "*"]
    },
    target: "web",
    devServer: {
        compress: true,
        contentBase: "./dist",
        inline: true,
        hot: true,
        historyApiFallback: true,
        port: 9200,
        proxy: {
            "/websocket": {changeOrigin: true, target: "ws://localhost:8080", secure: false, ws:true},
            "/pause": {changeOrigin: true, target: "http://localhost:8080", secure: false},
            "/activate": {changeOrigin: true, target: "http://localhost:8080", secure: false},
            "/deactivate": {changeOrigin: true, target: "http://localhost:8080", secure: false},
            "/interval": {changeOrigin: true, target: "http://localhost:8080", secure: false},
            "/reset": {changeOrigin: true, target: "http://localhost:8080", secure: false}
        }
    },
    module: {
        rules: [
            {
                test: /\.jsx?$/,
                loader: "babel-loader",
                options: {presets: ["@babel/react"], babelrc: false, cacheDirectory: true}
            },
            {
                test: /\.(scss|css)$/,
                use: ExtractTextPlugin.extract({
                    use: [
                        {loader: "css-loader", options: {camelCase: true, importLoaders: 1, modules: true}},
                        {loader: "sass-loader", options: {sourceMap: false}}
                    ],
                    fallback: "style-loader"
                })
            },
            {
                test: /\.(png|jpg|gif)$/,
                use: [
                    // {loader: "url-loader", options: {limit: 8192}},
                    {loader: "file-loader", options: {limit: 8192, name: "[name].[ext]", outputPath: "assets"}}
                ]
            }
        ]
    },
    plugins: [
        new ExtractTextPlugin({filename: "bundle.css"}),
        new HtmlPlugin({filename: "./index.html", favicon: path.resolve(__dirname, "assets/favicon.ico"), template: path.resolve(__dirname, "templates/index.html")})
    ],
    optimization: {
        minimizer: [
            new UglifyJsPlugin({ cache: true, parallel: true, uglifyOptions: {compress: {drop_console: true}}}),
            new OptimizeCSSAssetsPlugin({})
        ]
    }
}

module.exports = (env, argv) => {
    if (argv.NODE_ENV === "production") {
        config.plugins.push(new CleanPlugin(["dist"], {}))
    }

    return config
}
