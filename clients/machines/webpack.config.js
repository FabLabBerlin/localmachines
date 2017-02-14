var webpack = require('webpack');

const IS_PRODUCTION = process.env.NODE_ENV === 'production';
const IS_DEVELOPMENT = !IS_PRODUCTION;

function getOutputDevOrProd() {
  if (IS_DEVELOPMENT) {
    var outputDir = '/dev';
  } else {
    var outputDir = '/prod';
  }
  return outputDir;
}

var config = {
  devtool: IS_DEVELOPMENT ? 'source-map' : undefined,
  entry: './src/js/main.js',
  output: {
    path: __dirname + getOutputDevOrProd(),
    filename: 'bundle.js',
    publicPath: '/machines/assets/'
  },
  module: {
    loaders: [
      { test: /\.js$/, exclude: /node_modules/,  loader: 'babel-loader' },
      { test: /\.css/, loader: "style-loader!css-loader" },
      { test: /\.less$/, loader: "style-loader!css-loader!less-loader" },
      { test: /\.woff(\?v=\d+\.\d+\.\d+)?$/, loader: "url-loader?limit=10000&mimetype=application/font-woff" },
      { test: /\.woff2(\?v=\d+\.\d+\.\d+)?$/, loader: "url-loader?limit=10000&mimetype=application/font-woff" },
      { test: /\.ttf(\?v=\d+\.\d+\.\d+)?$/, loader: "url-loader?limit=10000&mimetype=application/octet-stream" },
      { test: /\.eot(\?v=\d+\.\d+\.\d+)?$/, loader: "file-loader" },
      { test: /\.svg(\?v=\d+\.\d+\.\d+)?$/, loader: "url-loader?limit=10000&mimetype=image/svg+xml" }
    ]
  },
  plugins: [
    new webpack.ProvidePlugin({
      $: "jquery",
      jQuery: "jquery",
      "window.jQuery": "jquery",
      "root.jQuery": "jquery"
    })
  ],
  resolve: {
    alias: {
      vex: 'vex-js'
    }
  }
};

if (IS_PRODUCTION) {
  config.resolve.alias['react'] = 'react/dist/react.min';
  config.resolve.alias['react-dom'] = 'react-dom/dist/react-dom.min';
  console.log('config.resolve.alias=');
  console.log(config.resolve.alias);
}

module.exports = config;
