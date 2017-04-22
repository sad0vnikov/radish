var gulp = require('gulp');
var webpack = require('gulp-webpack');
var rimraf = require('rimraf')

gulp.task('default', ['watch', 'webpack'])

gulp.task('webpack', function() {
    return gulp.src("src/index.html")
        .pipe(webpack(require("./webpack.config.js")))
        .pipe(gulp.dest("dist/"));
})

gulp.task('pre-build-clean', function(cb) {
    rimraf('./dist/', cb);
})

gulp.task("watch", function() {
    gulp.watch(["./src/elm/**/*.elm", "src/index.html"], ['webpack'])
})