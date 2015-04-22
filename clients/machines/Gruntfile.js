'use strict';

module.exports = function (grunt) {
	grunt.initConfig({
		pkg: grunt.file.readJSON('package.json'),
		
		jshint: {
			options: {
				jshintrc: '.jshintrc'
			},
			all: [
				'Gruntfile.js',
				'dev/**/*.js',
				'!dev/bower_components/**'
			]
		}, // jshint

		less: {
			dev: {
				files: {
					'dev/assets/css/main.css': ['dev/assets/less/main.less']
				},
				options: {
					compress: false
				}
			},

			prod: {
				files: {
					// Here it would be possible to create a custom bootstrap 
					// compile script that compiles only the bootstrap components
					// we actually need
					'prod/assets/css/app.min.css': [
						'dev/bower_components/bootstrap/less/bootstrap.less',
						'dev/bower_components/font-awesome/less/font-awesome.less',
						'dev/bower_components/toastr/toastr.less',
						'dev/bower_components/vex/css/vex.css',
						'dev/bower_components/vex/css/vex-theme-plain.css',
						'dev/assets/less/main.less'
					]
				},
				options: {
					compress: true,
					cleancss: true
				}
			}
		}, // less

		concat: {
			prod: {
				files: {
					'tmp/dependencies.js': [
						'dev/assets/js/ie10-viewport-bug-workaround.js',
						'dev/bower_components/jquery/dist/jquery.js',
						'dev/bower_components/lodash/lodash.js',
						'dev/bower_components/vex/js/vex.combined.min.js',
						'dev/bower_components/angular/angular.js',
						'dev/bower_components/angular-route/angular-route.js',
						'dev/bower_components/angular-cookies/angular-cookies.min.js',
						'dev/bower_components/angular-ui-bootstrap/src/transition/transition.js',
						'dev/bower_components/toastr/toastr.js',
						'dev/ng-components/version/version.js',
						'dev/ng-components/version/version-directive.js',
						'dev/ng-components/version/interpolate-filter.js',
						'dev/ng-components/filters.js'
					],
					'tmp/app.js': [
						'tmp/dependencies.js',
						'dev/ng-main.js',
						'dev/ng-modules/login/login.js',
						'dev/ng-modules/machines/machines.js',
						'dev/ng-modules/logout/logout.js'
					]
				}
			}
		}, // concat

		uglify: {
			prod: {
				files: {
					'prod/assets/js/modernizr.min.js': 'dev/bower_components/modernizr/modernizr.js',
					'prod/assets/js/app.min.js': 'tmp/app.js'
				}
			}
		}, // uglify

		processhtml: {
			prod: {
				options: {
					process: true
				}, 
				files: {
					'prod/index.html': ['dev/index.html']
				}
			}
		}, // processhtml

		copy: {
			prod: { // copy Angular template files to production dir
				files: [

				// Copy machines files
				{
					src: 'dev/bower_components/bootstrap/fonts/*',
					dest: 'prod/assets/fonts/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/bower_components/font-awesome/fonts/*',
					dest: 'prod/assets/fonts/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/login/login.html',
					dest: 'prod/ng-modules/login/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/logout/logout.html',
					dest: 'prod/ng-modules/logout/',
					expand: true,
					flatten: true
				}, {
					src: [
						'dev/ng-modules/machines/machines.html',
						'dev/ng-modules/machines/machine-body-available.html',
						'dev/ng-modules/machines/machine-body-occupied.html',
						'dev/ng-modules/machines/machine-body-unavailable.html',
						'dev/ng-modules/machines/machine-body-used.html',
						'dev/ng-modules/machines/machine-item.html'
					],
					dest: 'prod/ng-modules/machines/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/assets/img/*',
					dest: 'prod/assets/img/',
					filter: 'isFile',
					expand: true,
					flatten: true
				}]
			}
		}, // copy

		clean: {
			preprod: ['prod/**'],
			postprod: ['tmp/**']
		}, // clean

		watch: {
			js: {
				files: ['dev/**/*.js'],
				tasks: ['jshint'],
				options: {livereload: true, atBegin: true}
			},
			css: {
				files: ['dev/**/*.less'], 
				tasks: ['less:dev'],
				options: {livereload: true, atBegin: true}
			}
		} // watch

	});

	// Load tasks
	grunt.loadNpmTasks('grunt-contrib-clean');
	grunt.loadNpmTasks('grunt-contrib-jshint');
	grunt.loadNpmTasks('grunt-contrib-concat');
	grunt.loadNpmTasks('grunt-contrib-uglify');
	grunt.loadNpmTasks('grunt-notify');
	grunt.loadNpmTasks('grunt-contrib-watch');
	grunt.loadNpmTasks('grunt-text-replace');
	grunt.loadNpmTasks('grunt-processhtml');
	grunt.loadNpmTasks('grunt-contrib-less');
	grunt.loadNpmTasks('grunt-contrib-copy');

	// Register tasks
	grunt.registerTask('default', ['jshint', 'less:dev']);
	grunt.registerTask('dev', ['watch']);
	grunt.registerTask('prod', ['clean:preprod', 'jshint', 'less:prod', 'concat:prod', 'uglify:prod', 'processhtml:prod', 'copy:prod', 'clean:postprod']);
};
