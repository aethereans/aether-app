const path = require('path')
const VueLoaderPlugin = require('vue-loader/lib/plugin')
const NodeExternals = require('webpack-node-externals');
// ^ This makes it so that we can exclude all of node_modules.
const webpack = require('webpack')

module.exports = {
    // target: 'node', // in order to ignore built-in modules like path, fs, etc.
    // mode: 'development',
    mode: 'production',
    target: 'electron-renderer',
    node: {
        __dirname: false,
        __filename: false,
    },
    externals: [NodeExternals()], // in order to ignore all modules in node_modules folder
    entry: './src/app/renderermain.ts',
    output: {
        path: path.resolve(__dirname, 'dist'),
        filename: 'bundle.js'
    },
    resolve: {
        // Add `.ts` as a resolvable extension.
        extensions: ['.ts', '.js', '.vue']
    },
    module: {
        rules: [{
                test: /\.ts$/,
                // exclude: /node_modules|vue\/src/,
                loader: "ts-loader",
                options: {
                    appendTsSuffixTo: [/\.vue$/]
                }
            }, {
                test: /\.vue$/,
                loader: 'vue-loader'
            },
            // this will apply to both plain `.scss` files
            // AND `<style>` blocks in `.vue` files, same for all below
            {
                test: /\.scss$/,
                use: [
                    'vue-style-loader',
                    'css-loader',
                    'resolve-url-loader',
                    'sass-loader?sourceMap'
                ]
            },
            {
                test: /\.(ttf|eot|woff|woff2)$/,
                use: {
                    loader: "file-loader",
                    options: {
                        name: "[name].[ext]",
                        outputPath: 'typefaces/',
                        publicPath: 'dist/typefaces/'
                    },
                },
            },
            {
                test: /\.(png|svg|jpg|gif)$/,
                use: {
                    loader: "file-loader",
                    options: {
                        name: "[name].[ext]",
                        outputPath: 'images/',
                        publicPath: 'dist/images/'
                    },
                }
            },
            {
                test: /\.css$/,
                use: [
                    'vue-style-loader',
                    'css-loader'
                ]
            }
        ]
    },
    plugins: [
        // make sure to include the plugin!
        new VueLoaderPlugin(),
        new webpack.DefinePlugin({
            $dirname: '__dirname',
        })
    ]
};