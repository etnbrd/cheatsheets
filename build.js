var path = require('path');
var metalsmith = require('metalsmith');
var template = require('metalsmith-in-place');
var partial = require('metalsmith-partial');
var layout = require('metalsmith-layouts');
var webpack = require('metalsmith-webpack');
var wp = require('webpack');
var sass = require('metalsmith-sass');
var ignore = require('metalsmith-ignore');
var watch = require('metalsmith-watch');
var serve = require('metalsmith-serve');
var filenames = require('metalsmith-filenames');
var sitemap = require('metalsmith-mapsite');

var metadata = require('./metadata');

var ms = metalsmith(__dirname)
  .metadata(metadata)
  .source('./source')
  .destination('./public');

ms.use(filenames())
  .use(sass())
  .use(template({
    engine: 'pug',
    pattern: '**/*.pug',
    partials: './partials',
    rename: true,
    basedir: __dirname
  }))
  .use(ignore([
    'layouts/**',
    '**/.csscomb.json',
    '**/.csslintrc'
  ]))

if (process.argv[2] == "serve") {
  ms.use(serve({
      host: '0.0.0.0',
      port: '8080'
    }))
    .use(watch({
      livereload: true,
      paths: {
        "${source}/**/*": '**/*.pug',
        "${source}/styles/**/*": '**/*.scss',
        "content/**/*": '**/*',
      }
    }))
}

ms.build(function(err) {
    if (err) console.log(err);
  });
