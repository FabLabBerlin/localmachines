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
				'!dev/**/md5.js'
			]
		}, // jshint

		less: {
			dev: {
				files: {
					'dev/global/global.css': ['dev/global/global.less'], 
					'dev/machines/main.css': ['dev/machines/main.less'],
					'dev/admin/main.css': ['dev/admin/main.less']
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
					'prod/machines/css/app.min.css': [
						'bower_components/bootstrap/less/bootstrap.less',
						'dev/global/global.less',
						'dev/machines/main.less'
					],
					'prod/admin/css/app.min.css': [
						'bower_components/bootstrap/less/bootstrap.less',
						'dev/global/global.less',
						'dev/admin/main.less'
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
						'dev/global/js/ie10-viewport-bug-workaround.js',
						'dev/global/js/md5.js',
						'bower_components/jquery/dist/jquery.js',
						'bower_components/angular/angular.js',
						'bower_components/angular-route/angular-route.js',
						'bower_components/angular-timer/dist/angular-timer.min.js',
						'bower_components/angular-cookies/angular-cookies.min.js',
						'bower_components/angular-ui-bootstrap/src/modal/modal.js',
						'bower_components/angular-ui-bootstrap/src/transition/transition.js',
						'dev/global/angular/components/version/version.js',
						'dev/global/angular/components/version/version-directive.js',
						'dev/global/angular/components/version/interpolate-filter.js',
						'dev/global/angular/filters.js',
						'bower_components/webfontloader/webfontloader.js'
					],
					'tmp/machines-app.js': [
						'tmp/dependencies.js',
						'dev/machines/main.js',
						'dev/machines/login/login.js',
						'dev/machines/machines/machines.js',
						'dev/machines/logout/logout.js'
					],
					'tmp/admin-app.js': [
						'tmp/dependencies.js',
						'dev/admin/main.js',
						'dev/admin/mainmenu/mainmenu.js',
						'dev/admin/login/login.js',
						'dev/admin/dashboard/dashboard.js'
					]
				}
			}
		}, // concat

		uglify: {
			prod: {
				files: {
					'prod/machines/modernizr.min.js': 'bower_components/modernizr/modernizr.js',
					'prod/admin/modernizr.min.js': 'bower_components/modernizr/modernizr.js',
					'prod/machines/app.min.js': 'tmp/machines-app.js',
					'prod/admin/app.min.js': 'tmp/admin-app.js'
				}
			}
		}, // uglify

		processhtml: {
			prod: {
				options: {
					process: true
				}, 
				files: {
					'prod/machines/index.html': ['dev/machines/index.html'],
					'prod/admin/index.html': ['dev/admin/index.html']
				}
			}
		}, // processhtml

		copy: {
			prod: { // copy Angular template files to production dir
				files: [

				// Copy machines files
				{
					src: 'bower_components/bootstrap/fonts/*',
					dest: 'prod/machines/fonts/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/machines/login/login.html',
					dest: 'prod/machines/login/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/machines/logout/logout.html',
					dest: 'prod/machines/logout/',
					expand: true,
					flatten: true
				}, {
					src: [
						'dev/machines/machines/machines.html',
						'dev/machines/machines/machine-body-available.html',
						'dev/machines/machines/machine-body-occupied.html',
						'dev/machines/machines/machine-body-unavailable.html',
						'dev/machines/machines/machine-body-used.html',
						'dev/machines/machines/machine-item.html',
					],
					dest: 'prod/machines/machines/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/global/img/*',
					dest: 'prod/machines/img/',
					filter: 'isFile',
					expand: true,
					flatten: true
				}, 

				// Copy admin files
				{
					src: 'bower_components/bootstrap/fonts/*',
					dest: 'prod/admin/fonts/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/admin/dashboard/dashboard.html', 
					dest: 'prod/admin/dashboard/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/admin/login/login.html', 
					dest: 'prod/admin/login/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/admin/mainmenu/mainmenu.html', 
					dest: 'prod/admin/mainmenu/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/global/img/*',
					dest: 'prod/admin/img/',
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
