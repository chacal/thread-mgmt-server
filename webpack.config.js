const path = require('path')
const CopyWebpackPlugin = require('copy-webpack-plugin')
const ForkTsCheckerWebpackPlugin = require('fork-ts-checker-webpack-plugin')
const ForkTsCheckerNotifierWebpackPlugin = require('fork-ts-checker-notifier-webpack-plugin')

module.exports = {
  mode: 'development',
  entry: './public/index.tsx',
  output: {
    filename: 'bundle.js',
    path: path.resolve(__dirname, 'dist'),
  },
  devServer: {
    contentBase: './dist',
    port: 9000,
    hot: true,
    proxy: {
      '/v1': {
        target: 'http://localhost:8080',
        secure: false
      }
    }
  },
  module: {
    rules: [
      {
        test: /.tsx?$/,
        use: [
          {
            loader: 'ts-loader',
            options: {
              transpileOnly: true
            }
          }
        ],
        include: /public/,
      }
    ],
  },
  plugins: [
    new CopyWebpackPlugin({
      patterns: [
        { from: '*.html', context: 'public' },
      ]
    }),
    new ForkTsCheckerWebpackPlugin(),
    new ForkTsCheckerNotifierWebpackPlugin({ title: 'Webpack' })
  ],
  resolve: {
    extensions: ['.tsx', '.ts', '.js']
  },
}