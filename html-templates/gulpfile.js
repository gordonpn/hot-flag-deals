'use strict';
const { src, dest, watch, series } = require('gulp');
const del = require('del');
const mjml = require('gulp-mjml');
const mjmlEngine = require('mjml');
const server = require('browser-sync').create();

const SRC = './src/**/*.mjml';
const DEST = './dist';

const reload = (done) => {
  server.reload();
  done();
};

const serve = (done) => {
  server.init({
    server: {
      baseDir: './dist'
    },
    directory: 'dist'
  });
  done();
};

const build = (done) => {
  return src(SRC)
    .pipe(mjml(mjmlEngine, {minify: false, validationLevel: 'strict'}))
    .on('error', err => {
      console.log(err.toString());
      done();
    })
    .pipe(dest(DEST));
};

const monitor = () => {
  watch(SRC, series(build, reload));
};

const clean = () => {
  return del(['./dist/**/*']);
};

exports.build = series(clean, build);
exports.default = series(clean, build, serve, monitor);
