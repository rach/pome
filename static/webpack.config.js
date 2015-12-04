const autoprefixer = require('autoprefixer');
const ExtractTextPlugin = require('extract-text-webpack-plugin');
const path = require('path');
const webpack = require('webpack');


const sassIncludePaths = [
    path.resolve(__dirname, './node_modules/bootstrap/scss'),
    path.resolve(__dirname, './scss')
]

const sassLoaderIncludePathsArg = sassIncludePaths.map(
    function (v){
        return "includePaths[]=" + v;
    }
).join('&')

const sassLoaders = [
    'css-loader',
    'postcss-loader',
    'sass-loader?' + sassLoaderIncludePathsArg 
]

const config = {
    entry: {
        app: ['./js/index.js']
    },
    module: {
        loaders: [
            {
                test: /\.js$/,
                exclude: /node_modules/,
                loaders: ['babel-loader']
            },
            {
                test: /\.scss$/,
                loader: ExtractTextPlugin.extract('style-loader', sassLoaders.join('!'))
            }
        ]
    },
    output: {
        filename: '[name].js',
        path: path.join(__dirname, './build'),
        publicPath: '/build'
    },
    plugins: [
        new webpack.ProvidePlugin({
            'fetch': 'imports?this=>global!exports?global.fetch!whatwg-fetch'
        }),
        new ExtractTextPlugin('[name].css')
    ],
    postcss: [
        autoprefixer({
            browsers: ['last 2 versions']
        })
    ],
    resolve: {
        extensions: ['', '.js', '.scss', '.css'],
        modulesDirectories: ['node_modules']
    }
}

module.exports = config
