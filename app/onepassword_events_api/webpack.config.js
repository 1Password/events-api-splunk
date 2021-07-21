const path = require("path");
const { WebpackManifestPlugin } = require('webpack-manifest-plugin');
const MiniCssExtractPlugin = require("mini-css-extract-plugin");

module.exports = {
    entry: "./appserver/static/javascript/setup_page.js",
    output: {
        filename: "main_[contenthash].js",
        publicPath: "",
        path: path.resolve(__dirname, "appserver/static/build"),
    },
    plugins: [
        new WebpackManifestPlugin(),
        new MiniCssExtractPlugin({
            filename: "main_[contenthash].css"
        }),
    ],
    module: {
        rules: [
            {
                test: /\.css$/i,
                use: [MiniCssExtractPlugin.loader, "css-loader"],
            }
        ],
    }
};
