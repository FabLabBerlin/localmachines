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
					'dev/global/assets/css/global.css': ['dev/global/assets/less/global.less'], 
					'dev/machines/assets/css/main.css': ['dev/machines/assets/less/main.less'],
					'dev/admin/assets/css/main.css': ['dev/admin/assets/less/main.less']
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
					'prod/machines/assets/css/app.min.css': [
						'bower_components/bootstrap/less/bootstrap.less',
						'dev/global/assets/less/global.less',
						'dev/machines/assets/less/main.less'
					],
					'prod/admin/assets/css/app.min.css': [
						'bower_components/bootstrap/less/bootstrap.less',
						'dev/global/assets/less/global.less',
						'dev/admin/assets/less/main.less'
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
						'dev/global/assets/js/ie10-viewport-bug-workaround.js',
						'dev/global/assets/js/md5.js',
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
						'dev/machines/modules/login/login.js',
						'dev/machines/modules/machines/machines.js',
						'dev/machines/modules/logout/logout.js'
					],
					'tmp/admin-app.js': [
						'tmp/dependencies.js',
						'dev/admin/main.js',
						'dev/admin/modules/mainmenu/mainmenu.js',
						'dev/admin/modules/login/login.js',
						'dev/admin/modules/dashboard/dashboard.js'
					]
				}
			}
		}, // concat

		uglify: {
			prod: {
				files: {
					'prod/machines/assets/js/modernizr.min.js': 'bower_components/modernizr/modernizr.js',
					'prod/admin/assets/js/modernizr.min.js': 'bower_components/modernizr/modernizr.js',
					'prod/machines/assets/js/app.min.js': 'tmp/machines-app.js',
					'prod/admin/assets/js/app.min.js': 'tmp/admin-app.js'
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
					dest: 'prod/machines/assets/fonts/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/machines/modules/login/login.html',
					dest: 'prod/machines/modules/login/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/machines/modules/logout/logout.html',
					dest: 'prod/machines/modules/logout/',
					expand: true,
					flatten: true
				}, {
					src: [
						'dev/machines/modules/machines/machines.html',
						'dev/machines/modules/machines/machine-body-available.html',
						'dev/machines/modules/machines/machine-body-occupied.html',
						'dev/machines/modules/machines/machine-body-unavailable.html',
						'dev/machines/modules/machines/machine-body-used.html',
						'dev/machines/modules/machines/machine-item.html',
						'dev/machines/modules/machines/deactivate-modal.html'
					],
					dest: 'prod/machines/modules/machines/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/global/assets/img/*',
					dest: 'prod/machines/assets/img/',
					filter: 'isFile',
					expand: true,
					flatten: true
				}, 

				// Copy admin files
				{
					src: 'bower_components/bootstrap/fonts/*',
					dest: 'prod/admin/assets/fonts/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/admin/modules/dashboard/dashboard.html', 
					dest: 'prod/admin/modules/dashboard/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/admin/modules/login/login.html', 
					dest: 'prod/admin/modules/login/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/admin/modules/mainmenu/mainmenu.html', 
					dest: 'prod/admin/modules/mainmenu/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/global/assets/img/*',
					dest: 'prod/admin/assets/img/',
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
