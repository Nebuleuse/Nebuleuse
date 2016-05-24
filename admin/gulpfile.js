var gulp = require('gulp'),
    usemin = require('gulp-usemin'),
    wrap = require('gulp-wrap'),
    connect = require('gulp-connect'),
    watch = require('gulp-watch'),
    minifyCss = require('gulp-minify-css'),
    minifyJs = require('gulp-uglify'),
    concat = require('gulp-concat'),
    sass = require('gulp-sass'),
    rename = require('gulp-rename'),
    minifyHTML = require('gulp-minify-html');

var paths = {
    scripts: 'src/js/**/*.*',
    styles: 'src/sass/**/*.scss',
    css: 'src/css/**/*.*',
    images: 'src/img/**/*.*',
    templates: 'src/templates/**/*.html',
    index: 'src/index.html',
    fonts: 'src/fonts/*.{ttf,woff,eof,svg}',
    bower_fonts: 'bower_components/**/*.{ttf,woff,woff2,eof,svg}',
};

/**
 * Handle bower components from index
 */
gulp.task('usemin', function() {
    return gulp.src(paths.index)
        .pipe(usemin({
            js: [minifyJs(), 'concat'],
            //css: [minifyCss({keepSpecialComments: 0}), 'concat'],
        }))
        .pipe(gulp.dest('dist/'));
});

/**
 * Copy assets
 */
gulp.task('build-assets', ['copy-bower_fonts', 'copy-fonts']);

gulp.task('copy-bower_fonts', function() {
    return gulp.src(paths.bower_fonts)
        .pipe(rename({
            dirname: '/fonts'
        }))
        .pipe(gulp.dest('dist/lib'));
});
gulp.task('copy-fonts', function() {
    return gulp.src(paths.fonts)
        .pipe(rename({
            dirname: '/fonts'
        }))
        .pipe(gulp.dest('dist/lib'));
});

/**
 * Handle custom files
 */
gulp.task('build-custom', ['custom-images', 'custom-js', 'custom-sass','custom-css', 'custom-templates']);
gulp.task('custom-images', function() {
    return gulp.src(paths.images)
        .pipe(gulp.dest('dist/img'));
});

gulp.task('custom-js', function() {
    return gulp.src(paths.scripts)
        .pipe(minifyJs())
        .pipe(concat('dashboard.min.js'))
        .pipe(gulp.dest('dist/js'));
});

gulp.task('custom-sass', function() {
    return gulp.src(paths.styles)
        .pipe(sass())
        .pipe(gulp.dest('dist/css'));
});
gulp.task('custom-css', function() {
    return gulp.src(paths.css)
        .pipe(gulp.dest('dist/css'));
});


gulp.task('custom-templates', function() {
    return gulp.src(paths.templates)
        .pipe(minifyHTML())
        .pipe(gulp.dest('dist/templates'));
});

/**
 * Watch custom files
 */
gulp.task('watch', function() {
    gulp.watch([paths.images], ['custom-images']);
    gulp.watch([paths.css], ['custom-css']);
    gulp.watch([paths.scripts], ['custom-js']);
    gulp.watch([paths.templates], ['custom-templates']);
    gulp.watch([paths.index], ['usemin']);
});

/**
 * Live reload server
 */
gulp.task('webserver', function() {
    connect.server({
        root: 'dist',
        livereload: true,
        port: 8888
    });
});

gulp.task('livereload', function() {
    gulp.src(['dist/**/*.*'])
        .pipe(watch())
        .pipe(connect.reload());
});

/**
 * Gulp tasks
 */
gulp.task('build', ['usemin', 'build-assets', 'build-custom']);
gulp.task('default', ['build', 'watch']);