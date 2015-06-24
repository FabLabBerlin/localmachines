var webpack = require('webpack');

module.exports = {
    entry: './src/js/main.js',
    output: {
        path: __dirname + '/dev/',
        filename: 'bundle.js'
    },
    module: {
        loaders: [
            { test: /\.js$/, exclude: /node_module/,  loader: 'babel' },
            { test: /\.less$/, loader: "style!css!less" }
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
