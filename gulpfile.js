const gulp = require('gulp');
const conf = require('./conf')
const sass = require('gulp-sass');
const connect = require('gulp-connect')
const del = require('del');
const concat = require('gulp-concat');
const autoprefixer = require('gulp-autoprefixer');
const cleanCSS = require('gulp-clean-css');
const browserify = require('gulp-browserify');
const rename = require("gulp-rename");
const htmlmin = require('gulp-htmlmin');
const jsminify = require('gulp-minify');
const babel = require("gulp-babel");
const inject = require('gulp-inject');
const rev = require('gulp-rev');
const gulpif = require('gulp-if');
const proxy = require('http-proxy-middleware');

var MIN_BUILD = false;

gulp.task('sass', function () {
  return gulp.src(conf.path.src('**/*.scss'))
    .pipe(sass().on('error', sass.logError))
    .pipe(autoprefixer())
    .pipe(concat('styles.css'))
    .pipe(cleanCSS())
    .pipe(gulpif(MIN_BUILD, rev()))
    .pipe(gulp.dest(conf.path.tmp()))
    .pipe(gulp.dest(conf.path.dist()))
    .pipe(connect.reload());
});
gulp.task('scripts', function() {
  return gulp.src(conf.path.src('index.js'))
      .pipe(browserify({
        insertGlobals : false,
        transform: ['hbsfy']
      }))
      .pipe(babel({ presets: ['@babel/env']}))
      .pipe(gulpif(MIN_BUILD, rev()))
      .pipe(gulpif(MIN_BUILD, jsminify({
        noSource: true
      })))
      .pipe(gulp.dest(conf.path.tmp()))
      .pipe(gulp.dest(conf.path.dist()))
      .pipe(connect.reload());
});

gulp.task('html', function() {

  var sources = gulp.src([
      conf.path.tmp('*.js'),
      conf.path.tmp('*.css')
    ], {read: false});

  return gulp.src(conf.path.src('**/*.html'))
      .pipe(inject(sources, {
        ignorePath: '/tmp/'
      }))
      .pipe(gulpif(MIN_BUILD, htmlmin({
        collapseWhitespace: true,
        removeComments: true,
        minifyCSS: true,
        minifyJS: true
      })))
      .pipe(gulp.dest(conf.path.tmp()))
      .pipe(gulp.dest(conf.path.dist()))
      .pipe(connect.reload());

});
gulp.task('assets', () => {
  return gulp.src(conf.path.src('assets/**'))
      .pipe(gulp.dest(conf.path.tmp('assets')))
      .pipe(gulp.dest(conf.path.dist('assets')))
      .pipe(connect.reload());  
});




gulp.task('watch', function (done) {
  gulp.watch(conf.path.src('**/*.scss'), gulp.series('sass'));
  gulp.watch(conf.path.src('**/*.html'), gulp.series('html'));
  gulp.watch(conf.path.src('**/*.js'), gulp.series('scripts'));
  gulp.watch(conf.path.src('**/*.hbs'), gulp.series('scripts'));
  gulp.watch(conf.path.src('assets/**'), gulp.series('assets'));
  done();
});

gulp.task('clean', function () {
  return del(conf.path.tmp())
      .then(del(conf.path.dist()));
});

gulp.task('connect', function (done) {
  connect.server({
      root: conf.path.dist(),
      port: conf.connect.port,
      host: conf.connect.host,
      livereload: conf.connect.livereload,
      middleware: function(connect, opt) {
        return [
            proxy('/', {
                target: 'http://localhost:3001',
                changeOrigin:true
            })
        ]
    }
  });
  done();
});

gulp.task('dev:prepare',(done) => {
  MIN_BUILD = false;
  done();
});
gulp.task('prod:prepare',(done) => {
  MIN_BUILD = true;
  done();
});

gulp.task('serve', gulp.series('dev:prepare', 'clean', 'assets', 'sass', 'scripts', 'html', 'connect', 'watch'));
gulp.task('dist', gulp.series('prod:prepare', 'clean', 'assets', 'sass', 'scripts', 'html'));
