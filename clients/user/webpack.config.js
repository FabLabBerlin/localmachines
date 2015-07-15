var webpack = require('webpack');

function getOutputDevOrProd() {
  if(process.env.NODE_ENV !== 'production') {
    var outputDir = '/dev';
  } else {
    var outputDir = '/prod';
  }
  return outputDir;
}

module.exports = {
  entry: './src/js/main.js',
  output: {
    path: __dirname + getOutputDevOrProd(),
    filename: 'bundle.js'
  },
  module: {
    loaders: [
      { test: /\.js$/, exclude: /node_modules/,  loader: 'babel' },
      { test: /\.css/, loader: "style!css" },
      { test: /\.less$/, loader: "style!css!less" },
      { test: /\.woff(\?v=\d+\.\d+\.\d+)?$/, loader: "url-loader?limit=10000&minetype=application/font-woff" },
      { test: /\.woff2(\?v=\d+\.\d+\.\d+)?$/, loader: "url-loader?limit=10000&minetype=application/font-woff" },
      { test: /\.ttf(\?v=\d+\.\d+\.\d+)?$/, loader: "url-loader?limit=10000&minetype=application/octet-stream" },
      { test: /\.eot(\?v=\d+\.\d+\.\d+)?$/, loader: "file-loader" },
      { test: /\.svg(\?v=\d+\.\d+\.\d+)?$/, loader: "url-loader?limit=10000&minetype=image/svg+xml" }
    ]
  },
  plugins: [
    new webpack.ProvidePlugin({
      $: "jquery",
      jQuery: "jquery",
      "window.jQuery": "jquery",
      "root.jQuery": "jquery"
    })
  ]
};
