module.exports = {
  entry: './src/js/app',
  output: {
    path: __dirname + '/static',
    filename: 'bundle.js'
  },
  module: {
    rules: [
      {
        test: /\.vue$/,
        loader: 'vue-loader',
      },
      {
        test: /\.css$/,
        use: [
          'style-loader',
          'css-loader'
        ]
      }
    ]
  },
  resolve: {
    alias: {
      'vue$': 'vue/dist/vue.esm.js'
    }
  },
  devServer: {
    host: 'localhost',
    port: 7778,
    contentBase: __dirname + '/static',
    proxy: {
      '/api': {
        target: 'http://localhost:7777'
      }
    }
  },
  mode: 'production'
};
